package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTeamsDataSource_QueryMatch(t *testing.T) {
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

func testAccCheckTeamsContainsSharedTeams(name string, expectedTeamIDs []string) resource.TestCheckFunc {
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

		if teamsCount < len(expectedTeamIDs) {
			return fmt.Errorf("Incorrect number of teams - expected at least %d, got %d", len(expectedTeamIDs), teamsCount)
		}

		// Check that we can find our expected team IDs in the results
		foundTeams := make(map[string]bool)
		for i := 0; i < teamsCount; i++ {
			teamID := attributes[fmt.Sprintf("teams.%d.id", i)]
			foundTeams[teamID] = true
		}

		for _, expectedID := range expectedTeamIDs {
			if !foundTeams[expectedID] {
				return fmt.Errorf("Expected team ID %s not found in results", expectedID)
			}
		}

		return nil
	}
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
}`, rName)
}
