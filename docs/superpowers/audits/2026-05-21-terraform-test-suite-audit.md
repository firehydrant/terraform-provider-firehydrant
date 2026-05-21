# Terraform Acceptance Test Suite Audit

**Date:** 2026-05-21
**Scope:** `provider/*_test.go` (47 files, ~87 test steps)
**Goal:** Identify causes of flakiness (lifecycle ordering, name collisions, cascades) and runtime cost.

---

## TL;DR — Top Fixes by Impact

| # | Severity | Fix | Effort |
|---|----------|-----|--------|
| 1 | **Critical** | Fix `TF_ACC=1` vs `TF_ACC=="true"` mismatch — shared resource init never runs today | Trivial |
| 2 | **Critical** | Remove `testAccCheckTeamResourceDestroy` from tests that use `getSharedTeamID` — they're destroying the shared team mid-suite | Small |
| 3 | **High** | Randomize hardcoded slugs in `custom_event_source_resource_test.go`, `lifecycle_milestone_resource_test.go` | Small |
| 4 | **High** | Convert in-test team/service creation (~25 tests) to use `getSharedTeamID/getSharedServiceID` | Medium |
| 5 | **Medium** | Remove `time.Sleep(500ms)` per provider init in `provider.go:106` once rate limiting is otherwise mitigated | Small |
| 6 | **Medium** | Add `t.Parallel()` to data source tests + parent-resource-isolated tests, set `-parallel 4` in Makefile | Medium |
| 7 | **Medium** | Fix CheckDestroy ordering in `escalation_policy_data_test.go` (child → parent, not parent first) | Small |

---

## Critical: `TF_ACC` env var mismatch breaks shared resource init

**Files:** `Makefile:38`, `provider/provider_test.go:20`, `provider/provider.go:105,121`, `provider/test_helpers.go:58`

The Makefile sets `TF_ACC=1`, but every Go-side check is `os.Getenv("TF_ACC") == "true"`. Consequences:

- **`TestMain` skips `InitializeSharedResources` entirely** — the "shared resources" mechanism in `test_resources.go` is dead code under `make testacc`. Tests calling `getSharedTeamID(t)` fall through to `LoadFromAPI()` (which only loads teams by `tf-test-shared-` prefix) or `t.Fatalf`. If a stale shared team from a prior run happens to exist, tests "work"; otherwise they fail.
- **Cleanup at end of run also skipped** — `DestroyCreatedResources` is wrapped in the same `if` block, so any shared resources created get leaked.
- **The 500ms provider-init sleep also gated on `"true"`** — actually means tests skip the sleep today. Removing the sleep is a no-op for the current runs but matters once #1 is fixed.

**Fix:** Pick one (`true` or `1`) and use it consistently. `TF_ACC=1` is the Terraform SDK convention (`resource.TestCase` itself only checks if it's non-empty), so changing the Go checks to `os.Getenv("TF_ACC") != ""` is the safer choice.

---

## Critical: Tests destroy the shared team they depend on

**File:** `provider/on_call_schedule_resource_test.go:658,668` and `:736,746`

```go
sharedTeamID := getSharedTeamID(t)
...
CheckDestroy: testAccCheckTeamResourceDestroy(),
```

The test acquires the shared team ID, then asserts the team is destroyed at end-of-test. If `TF_ACC` were correctly wired (#1 above), this would actively delete the shared team that other tests in the same run depend on. Today it appears benign only because shared init isn't happening — fixing #1 turns this into a hard failure.

**Fix:** Drop `CheckDestroy: testAccCheckTeamResourceDestroy()` for these tests. The schedule's own destroy check is what matters. Audit other tests that combine `getShared*ID` with destroy checks for the parent type.

---

## High: Hardcoded globally-unique fields

These will collide across parallel CI jobs and on retried runs:

| File:line | Field | Value |
|-----------|-------|-------|
| `custom_event_source_resource_test.go:153,163` | `slug` | `"foo"` |
| `lifecycle_milestone_resource_test.go:209` | `slug` | `"test-milestone"` |
| `severity_data_test.go:57,69` | `slug` | `"TESTSEVERITYBASIC"`, `"TESTSEVERITYALL"` |
| `incident_type_resource_test.go:266-267` | `severity_slug`, `priority_slug` | `"SEV1"`, `"TESTPRIORITY"` |
| `incident_type_data_test.go:116-117` | same as above | |

The severity/priority slugs in `incident_type_*` are actually a different problem: they reference fixtures the test account *must already have*. That's a hidden environment dependency, not a collision risk — but it should be documented in `TESTS.md` or moved into `InitializeSharedResources`.

**Fix:** `acctest.RandStringFromCharSet` the user-created slugs. For platform fixtures (SEV1, TESTPRIORITY), promote them to shared resources or document the prerequisite.

---

## High: In-test team/service creation defeats the shared-resource design

**Roughly 25 test config functions** embed `resource "firehydrant_team"` or `resource "firehydrant_service"` instead of calling `getSharedTeamID/getSharedServiceID`. Examples:

- `escalation_policy_data_test.go` (×2)
- `inbound_email_resource_test.go`
- `incident_type_data_test.go`, `incident_type_resource_test.go`
- `on_call_schedule_data_test.go`
- `rotation_data_test.go`, `rotation_resource_test.go`
- `runbook_data_test.go`
- `functionality_data_test.go`, `functionality_resource_test.go`
- `service_data_test.go`, `service_resource_test.go`
- ...and others

Each pays API create + destroy cost for the parent resource (typically ~1-3s of API time), and each is a chance for a name collision or cascade leak. Some legitimately need *fresh* teams (mutation tests), but most just need a team to reference.

**Fix:** Audit each and either (a) inject `getSharedTeamID(t)` into the config template via `fmt.Sprintf` or (b) document why a fresh parent is required.

---

## Medium: `time.Sleep(500ms)` in provider init

**File:** `provider/provider.go:106`

Comment says it's to dodge 429s. With shared provider factories (`test_helpers.go:74`) most tests reuse one configured client, so this fires rarely. But for the 8 files using `mockProviderFactories`/`defaultProviderFactories`, every test step triggers it.

Once #1 is fixed and shared init runs, the rate-limit pressure goes down. Consider:
- Remove the sleep entirely and rely on the SDK's retry/backoff.
- Or replace with exponential backoff at the API call sites.

This becomes relevant after #1 — today the gate is `=="true"`, so the sleep doesn't fire under `make testacc`.

---

## Medium: No parallelism

**Files:** zero `t.Parallel()` calls across 47 test files. `Makefile:38` doesn't pass `-parallel`.

Go test default is `GOMAXPROCS` (often 4-10), but `t.Parallel()` is opt-in. Today the suite runs serially.

**Caveats before enabling:**
- Shared resources must be safe for concurrent reads (they are — they're just IDs).
- Tests that create their own parent resources (#4 above) must have collision-free names (#3 above) before this is safe.
- The provider's 500ms init sleep (#5) becomes more painful in parallel because more inits happen concurrently.

**Fix:**
1. Resolve #3 and #4 first.
2. Add `t.Parallel()` to data source tests and resource tests that don't mutate shared state.
3. Set `-parallel 4` in `Makefile:38`.

Realistic 2-4× speedup on independent tests.

---

## Medium: CheckDestroy ordering in multi-resource tests

**File:** `provider/escalation_policy_data_test.go:142-145` and similar

Multi-resource configs (team → schedule → escalation_policy) have CheckDestroy walking the state but no guarantee about order. When the parent (team) is destroyed first, the API may cascade-delete the children, causing the child's own destroy check to see 404 (which it interprets as "destroyed correctly") — masking real bugs in child cleanup.

**Files at risk:**
- `escalation_policy_data_test.go:142-145`
- `rotation_data_test.go:22-79` (team → schedule → rotation in one step)
- `functionality_resource_test.go:308-349` (3 teams + 2 services in one step)

**Fix:** Use `resource.ComposeAggregateTestCheckFunc` with explicit per-resource destroy checks in child-first order, OR add a delay between `terraform destroy` and the API check.

---

## Low / Informational

- **`testAccCheckResourceDestroy` only handles 4 resource types** (`test_helpers.go:135`) and prints "Warning: Generic destroy check not implemented" for the rest — meaning many tests have effectively no destroy check. Worth a follow-up to extend coverage.
- **`LoadFromAPI` has stubs** for users and incident roles (`test_resources.go:307,313`) that return nil silently. Either implement or remove the calls.
- **`time.Now().Unix()` for shared resource names** (`test_resources.go:104,114,etc.`) collides if two test runs start in the same second. Use `acctest.RandString` or `time.Now().UnixNano()` + random suffix.

---

## Suggested execution order

1. **Fix `TF_ACC` mismatch** — unblocks shared resource design entirely.
2. **Fix on-call-schedule team destroy check** — prevents shared team annihilation once #1 lands.
3. **Randomize hardcoded slugs** — required before parallelism.
4. **Audit & migrate in-test parent resource creation** — biggest single runtime win, also reduces collision surface.
5. **Add `t.Parallel()` and `-parallel 4`** — payoff phase.
6. **Tighten destroy checks** — improves signal quality from now on.
