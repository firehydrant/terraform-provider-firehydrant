package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testEscalationPolicyDataSourceConfig_basic() string {
	return fmt.Sprintln(`
data "firehydrant_escalation_policy" "test_escalation_policy" {
  team_id = "test-team-id"
  name = "My Escalation Policy"
}`)
}

func testAccEscalationPolicyDataSourceConfig_basic() string {
	return `
resource "firehydrant_team" "test_team" {
	name = "test-team-acc"
}

resource "firehydrant_escalation_policy" "test_escalation_policy" {
	name = "My Test Escalation Policy"
	description = "Test escalation policy for acceptance testing"
	team_id = firehydrant_team.test_team.id
	repetitions = 1
	
	step {
		timeout = "PT5M"
		targets {
			type = "user"
			id = "test-user-id"
		}
	}
}

data "firehydrant_escalation_policy" "test_escalation_policy" {
	team_id = firehydrant_team.test_team.id
	name = firehydrant_escalation_policy.test_escalation_policy.name
}`
}

func testAccEscalationPolicyDataSourceConfig_dynamic() string {
	return `
resource "firehydrant_team" "test_team" {
	name = "test-team-acc-dynamic"
}

resource "firehydrant_on_call_schedule" "test_on_call_schedule" {
	team_id = firehydrant_team.test_team.id
	name = "test-on-call-schedule-dynamic"
	time_zone = "America/New_York"
	slack_user_group_id = "test-group-dynamic"

	strategy {
		type         = "weekly"
		handoff_time = "10:00:00"
		handoff_day  = "monday"
	}
}

resource "firehydrant_escalation_policy" "test_escalation_policy" {
	name = "My Dynamic Escalation Policy"
	description = "Test dynamic escalation policy for acceptance testing"
	team_id = firehydrant_team.test_team.id
	repetitions = 1
	step_strategy = "dynamic_by_priority"

	step {
		timeout = "PT1M"
		targets {
			type = "OnCallSchedule"
			id   = firehydrant_on_call_schedule.test_on_call_schedule.id
		}
	}

	notification_priority_policies {
		priority = "HIGH"
		repetitions = 2
		
		step {
			timeout = "PT2M"
			targets {
				type = "OnCallSchedule"
				id   = firehydrant_on_call_schedule.test_on_call_schedule.id
			}
		}
		
		handoff_step {
			target_type = "Team"
			target_id   = firehydrant_team.test_team.id
		}
	}

	notification_priority_policies {
		priority = "LOW"
		repetitions = 1
		
		step {
			timeout = "PT5M"
			targets {
				type = "OnCallSchedule"
				id   = firehydrant_on_call_schedule.test_on_call_schedule.id
			}
		}
	}
}

data "firehydrant_escalation_policy" "test_escalation_policy" {
	team_id = firehydrant_team.test_team.id
	name = firehydrant_escalation_policy.test_escalation_policy.name
}`
}

func testEscalationPolicyDataSourceConfig_exactMatch() string {
	return fmt.Sprintln(`
data "firehydrant_escalation_policy" "test_escalation_policy" {
  team_id = "test-team-id"
  name = "My Escalation Policy"
}`)
}

func TestAccEscalationPolicyDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccEscalationPolicyDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_escalation_policy.test_escalation_policy", "id"),
					resource.TestCheckResourceAttrSet("data.firehydrant_escalation_policy.test_escalation_policy", "team_id"),
					resource.TestCheckResourceAttrSet("data.firehydrant_escalation_policy.test_escalation_policy", "name"),
					resource.TestCheckResourceAttrSet("data.firehydrant_escalation_policy.test_escalation_policy", "description"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "name", "My Test Escalation Policy"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "description", "Test escalation policy for acceptance testing"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "step.#", "1"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "step.0.timeout", "PT5M"),
				),
			},
		},
	})
}

func TestAccEscalationPolicyDataSource_dynamic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccEscalationPolicyDataSourceConfig_dynamic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_escalation_policy.test_escalation_policy", "id"),
					resource.TestCheckResourceAttrSet("data.firehydrant_escalation_policy.test_escalation_policy", "team_id"),
					resource.TestCheckResourceAttrSet("data.firehydrant_escalation_policy.test_escalation_policy", "name"),
					resource.TestCheckResourceAttrSet("data.firehydrant_escalation_policy.test_escalation_policy", "description"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "name", "My Dynamic Escalation Policy"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "description", "Test dynamic escalation policy for acceptance testing"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "step_strategy", "dynamic_by_priority"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "repetitions", "1"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "step.#", "1"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "step.0.timeout", "PT1M"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "step.0.targets.0.type", "OnCallSchedule"),
					resource.TestCheckResourceAttrSet("data.firehydrant_escalation_policy.test_escalation_policy", "step.0.targets.0.id"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.#", "2"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.priority", "HIGH"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.repetitions", "2"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.step.#", "1"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.step.0.timeout", "PT2M"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.step.0.targets.0.type", "OnCallSchedule"),
					resource.TestCheckResourceAttrSet("data.firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.step.0.targets.0.id"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.handoff_step.0.target_type", "Team"),
					resource.TestCheckResourceAttrSet("data.firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.handoff_step.0.target_id"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.1.priority", "LOW"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.1.repetitions", "1"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.1.step.#", "1"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.1.step.0.timeout", "PT5M"),
					resource.TestCheckResourceAttr("data.firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.1.step.0.targets.0.type", "OnCallSchedule"),
					resource.TestCheckResourceAttrSet("data.firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.1.step.0.targets.0.id"),
				),
			},
		},
	})
}

func TestEscalationPolicyDataSource_OneMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/v1/ping" && r.URL.Path != "/v1/teams/test-team-id/escalation_policies" {
			t.Errorf("Expected to request '/ping' or '/v1/teams/test-team-id/escalation_policies', got: %s", r.URL.Path)
		}

		if r.URL.Path == "/v1/teams/test-team-id/escalation_policies" && r.URL.Query().Get("query") != "My Escalation Policy" {
			t.Errorf("Expected query param 'query' to be 'My Escalation Policy', got: %s", r.URL.Query().Get("query"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"id": "123", "name":"My Escalation Policy", "description": "Test policy", "default": false, "repetitions": 3, "steps": [], "handoff_step": null, "step_strategy": "static", "notification_priority_policies": []}]}`))
	}))

	defer server.Close()

	orig := os.Getenv("FIREHYDRANT_BASE_URL")
	os.Setenv("FIREHYDRANT_BASE_URL", server.URL)
	t.Cleanup(func() { os.Setenv("FIREHYDRANT_BASE_URL", orig) })

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testEscalationPolicyDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_escalation_policy.test_escalation_policy", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_escalation_policy.test_escalation_policy", "name", "My Escalation Policy"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_escalation_policy.test_escalation_policy", "id", "123"),
				),
			},
		},
	})
}

func TestEscalationPolicyDataSource_MultipleMatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/v1/ping" && r.URL.Path != "/v1/teams/test-team-id/escalation_policies" {
			t.Errorf("Expected to request '/ping' or '/v1/teams/test-team-id/escalation_policies', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/v1/teams/test-team-id/escalation_policies" && r.URL.Query().Get("query") != "My Escalation Policy" {
			t.Errorf("Expected query param 'query' to be 'My Escalation Policy', got: %s", r.URL.Query().Get("query"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"id": "456", "name":"My Escalation Policy 2", "description": "Test policy 2", "default": false, "repetitions": 3, "steps": [], "handoff_step": null, "step_strategy": "static", "notification_priority_policies": []}, {"id": "123", "name":"My Escalation Policy", "description": "Test policy", "default": false, "repetitions": 3, "steps": [], "handoff_step": null, "step_strategy": "static", "notification_priority_policies": []}, {"id": "789", "name":"My Escalation Policy 3", "description": "Test policy 3", "default": false, "repetitions": 3, "steps": [], "handoff_step": null, "step_strategy": "static", "notification_priority_policies": []}]}`))
	}))

	defer server.Close()

	orig := os.Getenv("FIREHYDRANT_BASE_URL")
	os.Setenv("FIREHYDRANT_BASE_URL", server.URL)
	t.Cleanup(func() { os.Setenv("FIREHYDRANT_BASE_URL", orig) })

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testEscalationPolicyDataSourceConfig_exactMatch(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_escalation_policy.test_escalation_policy", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_escalation_policy.test_escalation_policy", "name", "My Escalation Policy"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_escalation_policy.test_escalation_policy", "id", "123"),
				),
			},
		},
	})
}

func TestEscalationPolicyDataSource_MultipleMatchesNoExactMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/v1/ping" && r.URL.Path != "/v1/teams/test-team-id/escalation_policies" {
			t.Errorf("Expected to request '/ping' or '/v1/teams/test-team-id/escalation_policies', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/v1/teams/test-team-id/escalation_policies" && r.URL.Query().Get("query") != "My Escalation Policy" {
			t.Errorf("Expected query param 'query' to be 'My Escalation Policy', got: %s", r.URL.Query().Get("query"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"id": "456", "name":"My Escalation Policy 2", "description": "Test policy 2", "default": false, "repetitions": 3, "steps": [], "handoff_step": null, "step_strategy": "static", "notification_priority_policies": []}, {"id": "789", "name":"My Escalation Policy 3", "description": "Test policy 3", "default": false, "repetitions": 3, "steps": [], "handoff_step": null, "step_strategy": "static", "notification_priority_policies": []}]}`))
	}))

	defer server.Close()

	orig := os.Getenv("FIREHYDRANT_BASE_URL")
	os.Setenv("FIREHYDRANT_BASE_URL", server.URL)
	t.Cleanup(func() { os.Setenv("FIREHYDRANT_BASE_URL", orig) })

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testEscalationPolicyDataSourceConfig_exactMatch(),
				ExpectError: regexp.MustCompile(`Did not find escalation policy matching 'My Escalation Policy'`),
			},
		},
	})
}

func TestEscalationPolicyDataSource_NoMatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/v1/ping" && r.URL.Path != "/v1/teams/test-team-id/escalation_policies" {
			t.Errorf("Expected to request '/ping' or '/v1/teams/test-team-id/escalation_policies', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/v1/teams/test-team-id/escalation_policies" && r.URL.Query().Get("query") != "My Escalation Policy" {
			t.Errorf("Expected query param 'query' to be 'My Escalation Policy', got: %s", r.URL.Query().Get("query"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[]}`))
	}))

	defer server.Close()

	orig := os.Getenv("FIREHYDRANT_BASE_URL")
	os.Setenv("FIREHYDRANT_BASE_URL", server.URL)
	t.Cleanup(func() { os.Setenv("FIREHYDRANT_BASE_URL", orig) })

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testEscalationPolicyDataSourceConfig_basic(),
				ExpectError: regexp.MustCompile(`Did not find escalation policy matching 'My Escalation Policy'`),
			},
		},
	})
}
