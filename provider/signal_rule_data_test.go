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

func testSignalRuleDataSourceConfig_basic() string {
	return fmt.Sprintln(`
data "firehydrant_signal_rule" "test_signal_rule" {
  team_id = "test-team-id"
  name = "My Signal Rule"
}`)
}

func testSignalRuleDataSourceConfig_exactMatch() string {
	return fmt.Sprintln(`
data "firehydrant_signal_rule" "test_signal_rule" {
  team_id = "test-team-id"
  name = "My Signal Rule"
}`)
}

func TestSignalRuleDataSource_OneMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/v1/ping" && r.URL.Path != "/v1/teams/test-team-id/signal_rules" {
			t.Errorf("Expected to request '/ping' or '/v1/teams/test-team-id/signal_rules', got: %s", r.URL.Path)
		}

		if r.URL.Path == "/v1/teams/test-team-id/signal_rules" && r.URL.Query().Get("query") != "My Signal Rule" {
			t.Errorf("Expected query param 'query' to be 'My Signal Rule', got: %s", r.URL.Query().Get("query"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"id": "123", "name":"My Signal Rule", "expression": "severity == 'high'", "target": {"type": "escalation_policy", "id": "ep-123", "name": "Test Target", "team_id": "t-123", "is_pageable": true}, "incident_type": {"id": "it-123"}, "notification_priority_override": "HIGH", "create_incident_condition_when": "WHEN_ALWAYS", "deduplication_expiry": "PT30M"}]}`))
	}))

	defer server.Close()

	orig := os.Getenv("FIREHYDRANT_BASE_URL")
	os.Setenv("FIREHYDRANT_BASE_URL", server.URL)
	t.Cleanup(func() { os.Setenv("FIREHYDRANT_BASE_URL", orig) })

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: mockProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testSignalRuleDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_signal_rule.test_signal_rule", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_signal_rule.test_signal_rule", "name", "My Signal Rule"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_signal_rule.test_signal_rule", "id", "123"),
				),
			},
		},
	})
}

func TestSignalRuleDataSource_MultipleMatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/v1/ping" && r.URL.Path != "/v1/teams/test-team-id/signal_rules" {
			t.Errorf("Expected to request '/ping' or '/v1/teams/test-team-id/signal_rules', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/v1/teams/test-team-id/signal_rules" && r.URL.Query().Get("query") != "My Signal Rule" {
			t.Errorf("Expected query param 'query' to be 'My Signal Rule', got: %s", r.URL.Query().Get("query"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"id": "456", "name":"My Signal Rule 2", "expression": "severity == 'medium'", "target": {"type": "escalation_policy", "id": "ep-456", "name": "Test Target 2", "team_id": "t-456", "is_pageable": false}, "incident_type": {"id": "it-456"}, "notification_priority_override": "MEDIUM", "create_incident_condition_when": "WHEN_UNSPECIFIED", "deduplication_expiry": "PT1H"}, {"id": "123", "name":"My Signal Rule", "expression": "severity == 'high'", "target": {"type": "escalation_policy", "id": "ep-123", "name": "Test Target", "team_id": "t-123", "is_pageable": true}, "incident_type": {"id": "it-123"}, "notification_priority_override": "HIGH", "create_incident_condition_when": "WHEN_ALWAYS", "deduplication_expiry": "PT30M"}, {"id": "789", "name":"My Signal Rule 3", "expression": "severity == 'low'", "target": {"type": "escalation_policy", "id": "ep-789", "name": "Test Target 3", "team_id": "t-789", "is_pageable": true}, "incident_type": {"id": "it-789"}, "notification_priority_override": "LOW", "create_incident_condition_when": "WHEN_UNSPECIFIED", "deduplication_expiry": "PT2H"}]}`))
	}))

	defer server.Close()

	orig := os.Getenv("FIREHYDRANT_BASE_URL")
	os.Setenv("FIREHYDRANT_BASE_URL", server.URL)
	t.Cleanup(func() { os.Setenv("FIREHYDRANT_BASE_URL", orig) })

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: mockProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testSignalRuleDataSourceConfig_exactMatch(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_signal_rule.test_signal_rule", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_signal_rule.test_signal_rule", "name", "My Signal Rule"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_signal_rule.test_signal_rule", "id", "123"),
				),
			},
		},
	})
}

func TestSignalRuleDataSource_MultipleMatchesNoExactMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/v1/ping" && r.URL.Path != "/v1/teams/test-team-id/signal_rules" {
			t.Errorf("Expected to request '/ping' or '/v1/teams/test-team-id/signal_rules', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/v1/teams/test-team-id/signal_rules" && r.URL.Query().Get("query") != "My Signal Rule" {
			t.Errorf("Expected query param 'query' to be 'My Signal Rule', got: %s", r.URL.Query().Get("query"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"id": "456", "name":"My Signal Rule 2", "expression": "severity == 'medium'", "target": {"type": "escalation_policy", "id": "ep-456", "name": "Test Target 2", "team_id": "t-456", "is_pageable": false}, "incident_type": {"id": "it-456"}, "notification_priority_override": "MEDIUM", "create_incident_condition_when": "WHEN_UNSPECIFIED", "deduplication_expiry": "PT1H"}, {"id": "789", "name":"My Signal Rule 3", "expression": "severity == 'low'", "target": {"type": "escalation_policy", "id": "ep-789", "name": "Test Target 3", "team_id": "t-789", "is_pageable": true}, "incident_type": {"id": "it-789"}, "notification_priority_override": "LOW", "create_incident_condition_when": "WHEN_UNSPECIFIED", "deduplication_expiry": "PT2H"}]}`))
	}))

	defer server.Close()

	orig := os.Getenv("FIREHYDRANT_BASE_URL")
	os.Setenv("FIREHYDRANT_BASE_URL", server.URL)
	t.Cleanup(func() { os.Setenv("FIREHYDRANT_BASE_URL", orig) })

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: mockProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testSignalRuleDataSourceConfig_exactMatch(),
				ExpectError: regexp.MustCompile(`Did not find signal rule matching 'My Signal Rule'`),
			},
		},
	})
}

func TestSignalRuleDataSource_NoMatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/v1/ping" && r.URL.Path != "/v1/teams/test-team-id/signal_rules" {
			t.Errorf("Expected to request '/ping' or '/v1/teams/test-team-id/signal_rules', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/v1/teams/test-team-id/signal_rules" && r.URL.Query().Get("query") != "My Signal Rule" {
			t.Errorf("Expected query param 'query' to be 'My Signal Rule', got: %s", r.URL.Query().Get("query"))
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
		ProviderFactories: mockProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testSignalRuleDataSourceConfig_basic(),
				ExpectError: regexp.MustCompile(`Did not find signal rule matching 'My Signal Rule'`),
			},
		},
	})
}
