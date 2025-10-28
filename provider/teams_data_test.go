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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/ping") && !strings.HasPrefix(r.URL.Path, "/v1/ping") && !strings.HasPrefix(r.URL.Path, "/v1/teams") {
			t.Errorf("Expected to request '/ping', '/v1/ping', or '/v1/teams', got: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"id":"123","name":"Test team 1", "description": "a description", "slug": "test-team-1"},{"id":"234","name":"Test team 2", "description": "a description", "slug": "test-team-2"}],"pagination":{"count":2,"page":1,"items":20,"pages":1,"last":2,"prev":null,"next":null}}`))
	}))

	defer server.Close()

	orig := os.Getenv("FIREHYDRANT_BASE_URL")
	os.Setenv("FIREHYDRANT_BASE_URL", server.URL)
	t.Cleanup(func() { os.Setenv("FIREHYDRANT_BASE_URL", orig) })

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTeamsDataSourceConfig_Query(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_teams.test_teams", "teams.#"),
					resource.TestCheckResourceAttr("data.firehydrant_teams.test_teams", "teams.0.id", "123"),
					resource.TestCheckResourceAttr("data.firehydrant_teams.test_teams", "teams.0.name", "Test team 1"),
					resource.TestCheckResourceAttr("data.firehydrant_teams.test_teams", "teams.1.id", "234"),
					resource.TestCheckResourceAttr("data.firehydrant_teams.test_teams", "teams.1.name", "Test team 2"),
				),
			},
		},
	})
}

func TestAccTeamsDataSource_basic(t *testing.T) {
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
	query = "Test team"
}`
}

func testAccTeamsDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "test_team" {
  name = "test-team-%s"
}

data "firehydrant_teams" "all_teams" {
}`, rName)
}
