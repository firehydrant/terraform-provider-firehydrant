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

func testTeamDataSourceConfig_basic() string {
	return fmt.Sprintln(`
data "firehydrant_team" "test_team" {
  name = "Test Team"
}`)
}

func TestTeamDataSource_OneMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/teams" {
			t.Errorf("Expected to request '/ping' or '/teams', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/teams" && r.URL.Query().Get("query") != "Test Team" {
			t.Errorf("Expected query param 'query' to be 'Test Team', got: %s", r.URL.Query().Get("query"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"id": "123", "name": "Test Team", "slug":"test-team"}]}`))
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
						"data.firehydrant_team.test_team", "slug", "test-team"),
				),
			},
		},
	})
}

func TestTeamDataSource_MultipleMatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/teams" {
			t.Errorf("Expected to request '/ping' or '/teams', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/teams" && r.URL.Query().Get("query") != "Test Team" {
			t.Errorf("Expected query param 'query' to be 'Test Team', got: %s", r.URL.Query().Get("query"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"id": "123", "name": "Test Team", "slug": "test-team"}]}`))
		w.Write([]byte(`{"data":[{"id": "123", "name": "Test Team", "slug": "test-team"},{"id": "456", "name": "Test Team 2", "slug": "test-team-2"}]}`))
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
				ExpectError: regexp.MustCompile(`Found multiple matching teams for 'Test Team'`),
			},
		},
	})
}

func TestTeamDataSource_NoMatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/teams" {
			t.Errorf("Expected to request '/ping' or '/teams', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/teams" && r.URL.Query().Get("query") != "Test Team" {
			t.Errorf("Expected query param 'query' to be 'Test Team', got: %s", r.URL.Query().Get("query"))
		}

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
				Config:      testTeamDataSourceConfig_basic(),
				ExpectError: regexp.MustCompile(`Did not find team matching 'Test Team'`),
			},
		},
	})
}
