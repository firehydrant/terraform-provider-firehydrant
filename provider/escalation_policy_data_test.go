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

func testEscalationPolicyDataSourceConfig_exactMatch() string {
	return fmt.Sprintln(`
data "firehydrant_escalation_policy" "test_escalation_policy" {
  team_id = "test-team-id"
  name = "My Escalation Policy"
}`)
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
