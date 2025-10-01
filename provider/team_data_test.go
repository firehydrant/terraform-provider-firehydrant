package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testTeamDataSourceConfig_basic() string {
	return fmt.Sprintln(`
data "firehydrant_team" "test_team" {
  id = "123"
}`)
}

func TestTeamDataSource_OneMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/ping") && !strings.HasPrefix(r.URL.Path, "/v1/ping") && !strings.HasPrefix(r.URL.Path, "/teams") {
			t.Errorf("Expected to request '/ping' or '/teams', got: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"id": "123",
			"name": "Test Team",
			"slug": "test_team",
			"description": "Test team description",
			"memberships": [
				{
					"user": {"id": "user-123"},
					"schedule": {"id": "schedule-456"},
					"default_incident_role": {"id": "role-789"}
				}
			]
		}`))
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
				Config: testTeamDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_team.test_team", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_team.test_team", "id", "123"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_team.test_team", "name", "Test Team"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_team.test_team", "slug", "test_team"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_team.test_team", "memberships.#", "1"),
					resource.TestCheckResourceAttrSet(
						"data.firehydrant_team.test_team", "memberships.0.user_id"),
					resource.TestCheckResourceAttrSet(
						"data.firehydrant_team.test_team", "memberships.0.schedule_id"),
					resource.TestCheckResourceAttrSet(
						"data.firehydrant_team.test_team", "memberships.0.default_incident_role_id"),
				),
			},
		},
	})
}

func TestTeamDataSource_NoMatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.HasPrefix(r.URL.Path, "/ping") && !strings.HasPrefix(r.URL.Path, "/v1/ping") && !strings.HasPrefix(r.URL.Path, "/teams") {
			t.Errorf("Expected to request '/ping' or '/teams', got: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"detail":"record not found"}`))
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
				Config:      testTeamDataSourceConfig_basic(),
				ExpectError: regexp.MustCompile(`resource not found`),
			},
		},
	})
}
