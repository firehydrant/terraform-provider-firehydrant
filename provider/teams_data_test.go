package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTeamsDataSource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
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

		if teamsCount <= 1 {
			return fmt.Errorf("Incorrect number of teams - expected at least 1, got %d", teamsCount)
		}

		return nil
	}
}

func testAccTeamsDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "test_team" {
  name = "test-team-%s"
}

data "firehydrant_teams" "all_teams" {
}`, rName)
}
