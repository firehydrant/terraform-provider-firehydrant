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
		if !strings.HasPrefix(r.URL.Path, "/ping") && !strings.HasPrefix(r.URL.Path, "/teams") {
			t.Errorf("Expected to request '/ping' or '/teams', got: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": "123", "name": "Test Team", "slug": "test_team", memberships: [{user:{"id": "456", "name": "Test User", "email": "user@example.com"}}]}`))
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
						"data.firehydrant_team.test_team", "memberships.0.id", "456"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_team.test_team", "memberships.0.name", "Test User"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_team.test_team", "memberships.0.email", "user@example.com"),

					resource.TestCheckResourceAttr(
						"data.firehydrant_team.test_team", "slug", "test_team"),
				),
			},
		},
	})
}

func TestTeamDataSource_NoMatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.HasPrefix(r.URL.Path, "/ping") && !strings.HasPrefix(r.URL.Path, "/teams") {
			t.Errorf("Expected to request '/ping' or '/teams', got: %s", r.URL.Path)
		}

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
