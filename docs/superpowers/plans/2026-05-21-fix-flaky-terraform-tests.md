# Fix Flaky Terraform Acceptance Tests — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix the critical and high-severity flakiness sources in the Terraform acceptance test suite identified in `docs/superpowers/audits/2026-05-21-terraform-test-suite-audit.md`.

**Architecture:** Fix the four highest-impact issues in order: (1) repair `TF_ACC` env-var mismatch so shared-resource init actually runs, (2) stop tests from destroying the shared team they depend on, (3) randomize hardcoded globally-unique slugs, (4) document the platform-fixture prerequisites that currently exist as hidden assumptions. Lower-priority items (parallelism, sleep removal, in-test parent resource migration) are out of scope here and tracked as a follow-up.

**Tech Stack:** Go, Terraform Plugin SDK v2, FireHydrant Go SDK.

**Test runner compatibility (CI + local):**

This plan must keep both runners green:
- **GitHub Actions** (`.github/workflows/ci.yml:40-46`) — runs `go test -v -timeout=25m ./...` directly with `TF_ACC: 'true'` and `FIREHYDRANT_API_KEY` from secrets. Does NOT use `make testacc`.
- **Local** (`Makefile:38`) — runs `TF_ACC=1 go test ...`.

Today the values disagree (`"true"` vs `"1"`) and Go checks compare against `"true"`. The fix is to make Go-side checks treat any non-empty value as truthy, so both runners work without changing either entry-point's env var. The Makefile and CI workflow do **not** need to change.

**Out of scope (separate follow-up plan):**
- Migrating ~25 tests that create their own teams/services to use `getSharedTeamID`/`getSharedServiceID` helpers.
- Enabling `t.Parallel()` and `-parallel 4` in Makefile.
- Removing `time.Sleep(500ms)` from `provider.go:106`.
- Tightening per-resource `CheckDestroy` ordering in multi-resource configs.
- Extending `testAccCheckResourceDestroy` to cover more resource types.

These build on the fixes here — do them after this plan lands.

**Verification approach:** Each task has its own verification. The end-to-end check is `make testacc` runs without the previously-observed flakes. Because acceptance tests hit a real FireHydrant account, full-suite verification requires `FIREHYDRANT_API_KEY` and is run by the user, not automated.

---

## File Structure

**Modified files:**
- `provider/provider.go` — accept any non-empty `TF_ACC`
- `provider/provider_test.go` — accept any non-empty `TF_ACC`
- `provider/test_helpers.go` — accept any non-empty `TF_ACC`
- `provider/on_call_schedule_resource_test.go` — drop incorrect team-destroy checks (2 sites)
- `provider/custom_event_source_resource_test.go` — randomize slug
- `provider/lifecycle_milestone_resource_test.go` — randomize slug
- `provider/test_resources.go` — replace `time.Now().Unix()` with `acctest.RandString` for shared resource names

**New files:**
- None.

**Docs:**
- `TESTS.md` — add prerequisite section documenting platform fixtures (`SEV1`, `TESTPRIORITY`) required by `incident_type_*` tests.

---

## Task 1: Fix `TF_ACC` env var mismatch

**Why:** `Makefile:38` sets `TF_ACC=1` but Go-side checks use `os.Getenv("TF_ACC") == "true"`. The shared-resource init in `TestMain` never runs under `make testacc`. The Terraform SDK convention is `TF_ACC=1`; align Go checks to accept any non-empty value.

**Files:**
- Modify: `provider/provider_test.go:20`
- Modify: `provider/provider.go:105`
- Modify: `provider/provider.go:121`
- Modify: `provider/test_helpers.go:58`

- [ ] **Step 1: Confirm current state of the four call sites**

Run: `grep -n 'TF_ACC' provider/provider.go provider/provider_test.go provider/test_helpers.go Makefile`

Expected output (exact):
```
provider/provider.go:105:	if os.Getenv("TF_ACC") == "true" {
provider/provider.go:121:	if os.Getenv("TF_ACC") != "true" {
provider/provider_test.go:20:	if os.Getenv("TF_ACC") == "true" {
provider/test_helpers.go:58:		if os.Getenv("TF_ACC") != "true" {
Makefile:38:	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 15m
```

If the lines differ from this output, stop and ask.

- [ ] **Step 2: Update `provider/provider_test.go:20`**

Change:
```go
if os.Getenv("TF_ACC") == "true" {
```
To:
```go
if os.Getenv("TF_ACC") != "" {
```

- [ ] **Step 3: Update `provider/provider.go:105`**

Change:
```go
if os.Getenv("TF_ACC") == "true" {
```
To:
```go
if os.Getenv("TF_ACC") != "" {
```

- [ ] **Step 4: Update `provider/provider.go:121`**

Change:
```go
if os.Getenv("TF_ACC") != "true" {
```
To:
```go
if os.Getenv("TF_ACC") == "" {
```

- [ ] **Step 5: Update `provider/test_helpers.go:58`**

Change:
```go
if os.Getenv("TF_ACC") != "true" {
```
To:
```go
if os.Getenv("TF_ACC") == "" {
```

- [ ] **Step 6: Verify build passes**

Run: `go build ./...`
Expected: no output, exit 0.

- [ ] **Step 7: Verify the unit-test (non-acceptance) suite still works**

Run: `make test`
Expected: tests pass. (This runs without `TF_ACC` set, so the acceptance code paths stay off.)

- [ ] **Step 8: Commit**

```bash
git add provider/provider.go provider/provider_test.go provider/test_helpers.go
git commit -m "$(cat <<'EOF'
Fix TF_ACC env check to accept any non-empty value

The Makefile sets TF_ACC=1 (the Terraform SDK convention) but Go-side
checks were comparing against "true", so shared-resource init in
TestMain never ran under make testacc. Align all four call sites to
the SDK convention.

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

---

## Task 2: Stop on-call schedule tests from destroying the shared team

**Why:** `TestAccOnCallScheduleResource_effectiveAt` and `TestAccOnCallScheduleResourceImport_basic` both call `getSharedTeamID(t)` (lines 658, 736) and then include `testAccCheckTeamResourceDestroy()` in CheckDestroy (lines 668, 746). After Task 1, shared init will actually run — and these tests will then actively destroy the shared team mid-suite, breaking every other test that uses it.

The on-call schedule's own destroy check is what these tests need. The team destroy check is wrong for tests using a shared team.

**Files:**
- Modify: `provider/on_call_schedule_resource_test.go:666-669`
- Modify: `provider/on_call_schedule_resource_test.go:744-747`

- [ ] **Step 1: Confirm the two CheckDestroy blocks**

Run: `sed -n '665,670p;743,748p' provider/on_call_schedule_resource_test.go`

Expected (exact):
```
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckOnCallScheduleResourceDestroy(),
			testAccCheckTeamResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckOnCallScheduleResourceDestroy(),
			testAccCheckTeamResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
```

- [ ] **Step 2: Update the first CheckDestroy (around line 666-669)**

Replace:
```go
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckOnCallScheduleResourceDestroy(),
			testAccCheckTeamResourceDestroy(),
		),
```
With:
```go
		CheckDestroy: testAccCheckOnCallScheduleResourceDestroy(),
```

Note: there are two occurrences of this exact block. Use the line-anchored edit (lines 666-669) for the first one; do not use `replace_all` for this step.

- [ ] **Step 3: Update the second CheckDestroy (around line 744-747)**

Replace the second occurrence of the same block with `CheckDestroy: testAccCheckOnCallScheduleResourceDestroy(),`.

- [ ] **Step 4: Verify both blocks are updated**

Run: `grep -n 'testAccCheckTeamResourceDestroy' provider/on_call_schedule_resource_test.go`

Expected: no output (the references are gone).

- [ ] **Step 5: Verify build**

Run: `go build ./...`
Expected: no output.

- [ ] **Step 6: Verify `testAccCheckTeamResourceDestroy` is still defined**

Run: `grep -rn 'func testAccCheckTeamResourceDestroy' provider/`
Expected: at least one match (the function may still be used by team-specific tests). Confirm it's still defined and that other call sites are tests that do create+destroy their own teams.

If it's no longer used anywhere, that's a follow-up cleanup, not part of this task.

- [ ] **Step 7: Commit**

```bash
git add provider/on_call_schedule_resource_test.go
git commit -m "$(cat <<'EOF'
Drop team destroy check from on-call schedule tests using shared team

TestAccOnCallScheduleResource_effectiveAt and TestAccOnCallScheduleResourceImport_basic
both use getSharedTeamID(t) but also asserted CheckDestroy on the team
itself, which would destroy the shared team that other tests in the
suite depend on. The schedule's own destroy check is the correct
assertion for these tests.

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

---

## Task 3: Randomize hardcoded slug in custom_event_source_resource_test.go

**Why:** `slug = "foo"` is hardcoded in both basic and update configs (lines 153, 163). Custom event source slugs are globally unique in FireHydrant. This collides between parallel CI jobs and on test retries.

**Files:**
- Modify: `provider/custom_event_source_resource_test.go` (imports, both tests, both config helpers)

- [ ] **Step 1: Read the current file structure**

Run: `sed -n '1,50p' provider/custom_event_source_resource_test.go`

Confirm imports include `"testing"`, `"context"`, `"fmt"`, and `terraform-plugin-sdk/v2/helper/resource`. We need to add `"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"` to imports.

- [ ] **Step 2: Add `acctest` import**

Modify imports block at lines 1-10. Add `"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"`.

Final imports block:
```go
package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)
```

- [ ] **Step 3: Change config helper signatures to accept `slug`**

Replace `testAccCustomEventSourceResourceConfig_basic` (lines 149-157):
```go
func testAccCustomEventSourceResourceConfig_basic(slug string) string {
	return fmt.Sprintf(`
resource "firehydrant_custom_event_source" "foo_transposer" {
  name = "The Foo Transposer"
	slug = "%s"
	description = "This is the foo transposer"
	javascript = "function transpose(input) {\n  return input.data;\n}"
}`, slug)
}
```

Replace `testAccCustomEventSourceResourceConfig_update` (lines 159-167):
```go
func testAccCustomEventSourceResourceConfig_update(slug string) string {
	return fmt.Sprintf(`
resource "firehydrant_custom_event_source" "foo_transposer" {
  name = "The Foo Transposer"
	slug = "%s"
	description = "A new foo transposer description"
	javascript = "function transpose(input) {\n  return input.foo;\n}"
}`, slug)
}
```

- [ ] **Step 4: Update `TestAccCustomEventSourceResource_basic`**

Replace the function body (lines 12-36) so that it generates a random slug and passes it to the config helper and the attribute checks:

```go
func TestAccCustomEventSourceResource_basic(t *testing.T) {
	slug := "tf-acc-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckCustomEventSourceResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomEventSourceResourceConfig_basic(slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCustomEventSourceResourceExistsWithAttributes_basic("firehydrant_custom_event_source.foo_transposer"),
					resource.TestCheckResourceAttrSet("firehydrant_custom_event_source.foo_transposer", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "name", "The Foo Transposer"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "slug", slug),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "description", "This is the foo transposer"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "javascript", "function transpose(input) {\n  return input.data;\n}"),
					resource.TestCheckResourceAttrSet("firehydrant_custom_event_source.foo_transposer", "ingest_url"),
				),
			},
		},
	})
}
```

- [ ] **Step 5: Update `TestAccCustomEventSourceResource_update`**

Replace the function body (lines 38-78) so it generates one random slug and reuses it across both steps:

```go
func TestAccCustomEventSourceResource_update(t *testing.T) {
	slug := "tf-acc-" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckCustomEventSourceResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomEventSourceResourceConfig_basic(slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCustomEventSourceResourceExistsWithAttributes_basic("firehydrant_custom_event_source.foo_transposer"),
					resource.TestCheckResourceAttrSet("firehydrant_custom_event_source.foo_transposer", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "name", "The Foo Transposer"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "slug", slug),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "description", "This is the foo transposer"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "javascript", "function transpose(input) {\n  return input.data;\n}"),
					resource.TestCheckResourceAttrSet("firehydrant_custom_event_source.foo_transposer", "ingest_url"),
				),
			},
			{
				Config: testAccCustomEventSourceResourceConfig_update(slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCustomEventSourceResourceExistsWithAttributes_basic("firehydrant_custom_event_source.foo_transposer"),
					resource.TestCheckResourceAttrSet("firehydrant_custom_event_source.foo_transposer", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "name", "The Foo Transposer"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "slug", slug),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "description", "A new foo transposer description"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "javascript", "function transpose(input) {\n  return input.foo;\n}"),
					resource.TestCheckResourceAttrSet("firehydrant_custom_event_source.foo_transposer", "ingest_url"),
				),
			},
		},
	})
}
```

- [ ] **Step 6: Update `TestAccCustomEventSourceResourceImport_basic`**

Locate the import test (around line 169). Apply the same pattern: generate a `slug` once with `acctest.RandStringFromCharSet`, pass it to `testAccCustomEventSourceResourceConfig_basic(slug)`.

Read the existing test:
```bash
sed -n '169,200p' provider/custom_event_source_resource_test.go
```

Update it to declare the slug at the top of the function and pass it to the config helper(s) it calls.

- [ ] **Step 7: Verify build**

Run: `go build ./...`
Expected: no output, exit 0.

- [ ] **Step 8: Verify no remaining hardcoded `slug = "foo"`**

Run: `grep -n 'slug = "foo"' provider/custom_event_source_resource_test.go`
Expected: no matches.

- [ ] **Step 9: Commit**

```bash
git add provider/custom_event_source_resource_test.go
git commit -m "$(cat <<'EOF'
Randomize slug in custom_event_source acceptance tests

The hardcoded slug "foo" collided between parallel CI jobs and on
test retries. Generate a random slug per test invocation so each run
uses a unique value.

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

---

## Task 4: Randomize hardcoded slug in lifecycle_milestone_resource_test.go

**Why:** `slug = "test-milestone"` is hardcoded at line 209 in the update config. Lifecycle milestone slugs are unique per phase; this collides between concurrent runs.

**Files:**
- Modify: `provider/lifecycle_milestone_resource_test.go`

- [ ] **Step 1: Read the affected function**

Run: `sed -n '199,213p' provider/lifecycle_milestone_resource_test.go`

Expected (last lines):
```go
func testAccLifecycleMilestoneResourceConfig_update(rName string) string {
	return fmt.Sprintf(`
data "firehydrant_lifecycle_phase" "started" {
  name = "started"
}

resource "firehydrant_lifecycle_milestone" "new_milestone" {
  name        = "Test Milestone %s"
  description = "test description %s"
	phase_id    = data.firehydrant_lifecycle_phase.started.id
	slug        = "test-milestone"
	position    = 2
	auto_assign_timestamp_on_create = "never_set_on_create"
}`, rName, rName)
}
```

- [ ] **Step 2: Update the config function to derive slug from `rName`**

Replace the function (lines 199-213):

```go
func testAccLifecycleMilestoneResourceConfig_update(rName string) string {
	return fmt.Sprintf(`
data "firehydrant_lifecycle_phase" "started" {
  name = "started"
}

resource "firehydrant_lifecycle_milestone" "new_milestone" {
  name        = "Test Milestone %s"
  description = "test description %s"
	phase_id    = data.firehydrant_lifecycle_phase.started.id
	slug        = "test-milestone-%s"
	position    = 2
	auto_assign_timestamp_on_create = "never_set_on_create"
}`, rName, rName, rName)
}
```

- [ ] **Step 3: Check whether any test asserts the slug attribute literally**

Run: `grep -n 'slug.*"test-milestone"' provider/lifecycle_milestone_resource_test.go`
Expected: only the line we just modified (or no matches if the modification already happened). If a `resource.TestCheckResourceAttr` references `"test-milestone"` literally, that assertion must be updated too — either dropped or rewritten to match `"test-milestone-" + rName`.

If you find such an assertion, read the surrounding test and update it. Otherwise proceed.

- [ ] **Step 4: Verify build**

Run: `go build ./...`
Expected: no output.

- [ ] **Step 5: Commit**

```bash
git add provider/lifecycle_milestone_resource_test.go
git commit -m "$(cat <<'EOF'
Randomize slug in lifecycle_milestone update test config

The hardcoded slug "test-milestone" could collide between concurrent
runs. Derive it from the per-test random name.

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

---

## Task 5: Document platform-fixture prerequisites and randomize shared resource names

**Why:** Two problems combined into one task because both touch the test-setup story:

(a) `incident_type_resource_test.go:266-267` and `incident_type_data_test.go:116-117` reference `severity_slug = "SEV1"` and `priority_slug = "TESTPRIORITY"`. These must already exist in the test FireHydrant account — there's no setup code that creates them. This is currently an undocumented prerequisite.

(b) `provider/test_resources.go` uses `time.Now().Unix()` for shared resource names (lines 104, 114, 124, 134, 144, 154). Two test runs starting in the same second collide. With Task 1 fixed, shared init runs every time, so this becomes load-bearing.

Also `severity_data_test.go` uses hardcoded slugs `TESTSEVERITYBASIC` and `TESTSEVERITYALL` — but these are *resources the test creates*, not platform fixtures. They have the same collision problem as Task 3/4. We'll randomize them here too.

**Files:**
- Modify: `provider/test_resources.go:104,114,124,134,144,154`
- Modify: `provider/severity_data_test.go` (both config helpers + test assertions)
- Modify: `TESTS.md` (add prerequisites section)

- [ ] **Step 1: Update shared resource name generation in `test_resources.go`**

Read the current name-generation pattern:
```bash
grep -n 'time.Now().Unix()' provider/test_resources.go
```

Expected: 6 matches at the lines noted above.

For each `fmt.Sprintf("tf-test-shared-...-%d", time.Now().Unix())` call, replace `time.Now().Unix()` with a random suffix. At the top of `InitializeSharedResources` (line 101), add a unique run identifier:

```go
func (r *SharedTestResources) InitializeSharedResources(ctx context.Context, client *firehydrant.APIClient) error {
	runID := acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
```

Then update each call site to use `runID` instead of `time.Now().Unix()`. For example:
```go
teamID, err := r.createSharedTeam(ctx, client, fmt.Sprintf("tf-test-shared-default-%s", runID))
```

Apply this transform to all 6 call sites (default team, default schedule, default incident role, default service, service2, team2).

Add the import: `"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"` to the imports block at the top of the file.

Remove the `"time"` import if no other code in the file uses it. (Check with `grep -n 'time\\.' provider/test_resources.go` — if anything other than the removed `time.Now()` calls references the package, keep the import. There IS a `time.RFC3339` reference at line 188, so keep the import.)

- [ ] **Step 2: Verify build**

Run: `go build ./...`
Expected: no output.

- [ ] **Step 3: Randomize the slug in `severity_data_test.go`**

The two configs at lines 54-64 and 66-77 both hardcode unique slugs. Update each to accept a slug parameter and pass a randomized value from the test.

Read first:
```bash
sed -n '1,80p' provider/severity_data_test.go
```

Add `acctest` to imports. Update each test function to generate `slug := "TESTSEV" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)`, pass it to the config helper, and update the `TestCheckResourceAttr` assertion to compare against the same generated slug.

Final state of `testAccSeverityDataSourceConfig_basic`:
```go
func testAccSeverityDataSourceConfig_basic(slug string) string {
	return fmt.Sprintf(`
resource "firehydrant_severity" "test_severity" {
  slug        = "%s"
  description = "test-description"
}

data "firehydrant_severity" "test_severity" {
  slug = firehydrant_severity.test_severity.id
}`, slug)
}
```

Final state of `testAccSeverityDataSourceConfig_allAttributes`:
```go
func testAccSeverityDataSourceConfig_allAttributes(slug string) string {
	return fmt.Sprintf(`
resource "firehydrant_severity" "test_severity" {
  slug        = "%s"
  description = "test-description"
  type        = "gameday"
}

data "firehydrant_severity" "test_severity" {
  slug = firehydrant_severity.test_severity.id
}`, slug)
}
```

Both tests at the top of the file need the corresponding update: generate `slug` at the start, pass it to the config helper, and update the slug assertion to use the same variable.

- [ ] **Step 4: Verify build**

Run: `go build ./...`
Expected: no output.

- [ ] **Step 5: Verify no hardcoded uniqueness-critical slugs remain in the modified files**

Run:
```bash
grep -n 'TESTSEVERITY\|test-milestone[^"]*"$\|slug = "foo"' provider/*_test.go
```

Expected: matches only inside `incident_type_*_test.go` and `incident_type_data_test.go` (which reference platform fixtures `SEV1` and `TESTPRIORITY`, not user-created slugs — those stay).

- [ ] **Step 6: Document the platform fixture prerequisites in `TESTS.md`**

Read the current file:
```bash
cat TESTS.md
```

Add a new section at the end (after section "2. Run the tests"):

```markdown
## 3. Required platform fixtures

Some acceptance tests depend on FireHydrant resources that the suite does NOT create automatically. Before running `make testacc` against a fresh account, ensure the following exist:

| Resource type | Slug / identifier | Used by |
|---------------|------------------|---------|
| Severity      | `SEV1`           | `incident_type_resource_test.go`, `incident_type_data_test.go` |
| Priority      | `TESTPRIORITY`   | `incident_type_resource_test.go`, `incident_type_data_test.go` |

These can be created via the FireHydrant UI or API. They are read by tests but never mutated; one-time setup per test account is sufficient.

> **Note:** These prerequisites are an artifact of the current test design — incident-type tests reference fixed slugs rather than creating their own severity/priority. A future refactor may move these into the shared-resource initialization in `provider/test_resources.go`.
```

- [ ] **Step 7: Commit**

```bash
git add provider/test_resources.go provider/severity_data_test.go TESTS.md
git commit -m "$(cat <<'EOF'
Randomize shared resource and severity slugs; document platform fixtures

Replace time.Now().Unix() with acctest.RandString for shared test
resource names so back-to-back test runs don't collide. Randomize the
hardcoded TESTSEVERITY* slugs in severity_data_test.go for the same
reason. Document the SEV1 / TESTPRIORITY platform fixtures that
incident_type tests assume exist in the target account.

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
EOF
)"
```

---

## Task 6: End-to-end verification

**Why:** Confirm the suite still compiles, unit tests pass, and the audit's headline issues are resolved.

- [ ] **Step 1: Build check**

Run: `go build ./...`
Expected: no output.

- [ ] **Step 2: `go vet` check**

Run: `go vet ./...`
Expected: no output.

- [ ] **Step 3: Run non-acceptance tests**

Run: `make test`
Expected: all pass.

- [ ] **Step 4: Confirm `TF_ACC` checks are now consistent**

Run: `grep -n 'TF_ACC' provider/provider.go provider/provider_test.go provider/test_helpers.go`

Expected: every Go check is `!= ""` or `== ""` — no `"true"` references.

- [ ] **Step 4a: Confirm CI workflow still uses `TF_ACC: 'true'` and we didn't accidentally touch it**

Run: `grep -n 'TF_ACC' .github/workflows/ci.yml Makefile`

Expected (exact):
```
.github/workflows/ci.yml:44:          TF_ACC: 'true'
Makefile:38:	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 15m
```

Both entry-points are unchanged; the Go-side fix accommodates both values.

- [ ] **Step 4b: Simulate CI invocation locally (without API key, just to confirm the binary loads)**

Run: `TF_ACC=true go test -run NoSuchTest ./provider/...`

Expected: tests find no matches and exit cleanly. This validates that `TF_ACC=true` (CI's value) does not crash the test binary during package init — important because TestMain now runs shared-resource init when `TF_ACC != ""`.

Note: this command WILL fail with "FIREHYDRANT_API_KEY not set" inside TestMain if the env var is missing — that's expected and confirms the gating logic now fires on `TF_ACC=true`. The test binary itself should still exit (with a non-zero code from the fatal log). If it segfaults or panics in an unexpected way, investigate before proceeding.

- [ ] **Step 5: Confirm no remaining `testAccCheckTeamResourceDestroy` in shared-team tests**

Run: `grep -B2 'testAccCheckTeamResourceDestroy' provider/on_call_schedule_resource_test.go`
Expected: no matches.

- [ ] **Step 6: Confirm no hardcoded unique slugs remain in the files we touched**

Run:
```bash
grep -n 'slug = "foo"\|slug = "test-milestone"$\|TESTSEVERITYBASIC\|TESTSEVERITYALL' provider/*_test.go
```
Expected: no matches.

- [ ] **Step 7: Manual acceptance run (user, not automated)**

This step requires `FIREHYDRANT_API_KEY` and access to a test account. The user runs:

```bash
envchain terraform-provider-firehydrant make testacc
```

Watch for:
- Shared-resource creation log lines (`tf-test-shared-default-<random>`, etc.) — should appear once per run.
- Cleanup log lines at the end.
- No "still exists" errors in CheckDestroy.
- No 4xx errors complaining about duplicate slugs.

If a flake reappears that's NOT addressed by this plan, capture which test and add to the audit doc's follow-up section.

---

## Self-Review Notes

- **Spec coverage:** Tasks 1-5 map to audit items #1-#4 plus the `time.Now().Unix()` and `TESTSEVERITY*` issues. Out-of-scope items are explicitly documented.
- **Placeholder scan:** No TBDs, no "implement later" — every step shows the exact code or command.
- **Type consistency:** `acctest.RandStringFromCharSet(N, acctest.CharSetAlphaNum)` is used consistently across Tasks 3, 4, 5. Slug prefixes (`tf-acc-`, `TESTSEV`) are deliberately different per resource type.
- **Ambiguity:** Task 5 step 3 says "both tests at the top need the corresponding update" — that's two tests (`TestAccSeverityDataSource_basic` and `TestAccSeverityDataSource_allAttributes`), each updated to mirror Task 3's pattern.
