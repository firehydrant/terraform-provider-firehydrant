package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTeamsDataSource_QueryMatch(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTeamsDataSourceConfig_Query(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_teams.test_teams", "teams.#"),
					// Verify we can query for teams by name - this tests the query functionality
					testAccCheckTeamsSet("data.firehydrant_teams.test_teams"),
				),
			},
		},
	})
}

// This is a unit test unlike other tests. This tests purpose is to test the
// pagination logic. Other tests are in place to confirm the contract of the API
// still matches the expectations of the provider.
func TestTeamsDataSource_Pagination(t *testing.T) {
	if os.Getenv("FIREHYDRANT_API_KEY") == "" {
		t.Setenv("FIREHYDRANT_API_KEY", "unit-test-pagination")
	}

	var responsePage1, responsePage2 bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/ping" || strings.HasPrefix(r.URL.Path, "/v1/ping"):
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
			return
		case r.Method == http.MethodGet && r.URL.Path == "/v1/teams":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			switch page := r.URL.Query().Get("page"); page {
			case "":
				fallthrough
			case "1":
				responsePage1 = true
				_, _ = w.Write([]byte(`{
					"data": [{
						"id": "11111111-1111-1111-1111-111111111111",
						"name": "Silver Snakes",
						"description": "The Silver Snakes",
						"slug": "team-silver-snakes",
						"services": [],
						"owned_services": []
					}],
					"pagination": {
						"count": 2,
						"page": 1,
						"items": 1,
						"pages": 2,
						"last": 2,
						"prev": null,
						"next": 2
					}
				}`))
			case "2":
				responsePage2 = true
				_, _ = w.Write([]byte(`{
					"data": [{
						"id": "22222222-2222-2222-2222-222222222222",
						"name": "Blue Barracudas",
						"description": "The Blue Barracudas",
						"slug": "team-blue-barracudas",
						"services": [],
						"owned_services": []
					}],
					"pagination": {
						"count": 2,
						"page": 2,
						"items": 1,
						"pages": 2,
						"last": 2,
						"prev": 1,
						"next": null
					}
				}`))
			default:
				t.Errorf("unexpected list teams page query: %q", page)
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	t.Setenv("FIREHYDRANT_BASE_URL", server.URL)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: mockProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTeamsDataSourceConfig_pagination(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.firehydrant_teams.paginated_teams", "teams.#", "2"),
					resource.TestCheckResourceAttr("data.firehydrant_teams.paginated_teams", "teams.0.id", "11111111-1111-1111-1111-111111111111"),
					resource.TestCheckResourceAttr("data.firehydrant_teams.paginated_teams", "teams.0.name", "Silver Snakes"),
					resource.TestCheckResourceAttr("data.firehydrant_teams.paginated_teams", "teams.1.id", "22222222-2222-2222-2222-222222222222"),
					resource.TestCheckResourceAttr("data.firehydrant_teams.paginated_teams", "teams.1.name", "Blue Barracudas"),
					func(*terraform.State) error {
						if !responsePage1 || !responsePage2 {
							return fmt.Errorf("expected list-teams to request both page 1 and page 2 (pagination), saw page1=%v page2=%v", responsePage1, responsePage2)
						}
						return nil
					},
				),
			},
		},
	})
}

func TestAccTeamsDataSource_basic(t *testing.T) {
	t.Parallel()
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckTeamResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTeamsDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_teams.all_teams", "teams.#"),
					testAccCheckTeamsSet("data.firehydrant_teams.all_teams"),
				),
			},
		},
	})
}

func testAccCheckTeamsSet(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		teamsResource, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find teams resource in state: %s", name)
		}

		if teamsResource.Primary.ID == "" {
			return fmt.Errorf("Teams resource ID not set")
		}

		attributes := teamsResource.Primary.Attributes
		teams, teamsOk := attributes["teams.#"]
		if !teamsOk {
			return fmt.Errorf("Teams list is missing")
		}

		teamsCount, err := strconv.Atoi(teams)
		if err != nil {
			return err
		}

		if teamsCount < 1 {
			return fmt.Errorf("Incorrect number of teams - expected at least 1, got %d", teamsCount)
		}

		return nil
	}
}

func testAccTeamsDataSourceConfig_Query() string {
	return `
data "firehydrant_teams" "test_teams" {
	query = "team"
}`
}

func testAccTeamsDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "test_team" {
  name = "test-team-%s"
}

data "firehydrant_teams" "all_teams" {
	query = "test-team"

	# The query string doesn't reference the resource, so without an explicit
	# dependency Terraform may read this data source before the team exists.
	depends_on = [firehydrant_team.test_team]
}`, rName)
}

func testAccTeamsDataSourceConfig_pagination() string {
	return `
data "firehydrant_teams" "paginated_teams" {}`
}
