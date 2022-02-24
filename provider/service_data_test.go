package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccServiceDataSource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_service.test_service", "alert_on_add", "false"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_service.test_service", "service_tier", "5"),
				),
			},
		},
	})
}

func TestAccServiceDataSource_allAttributes(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDataSourceConfig_allAttributes(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_service.test_service", "alert_on_add", "true"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_service.test_service", "description", fmt.Sprintf("test-description-%s", rName)),
					resource.TestCheckResourceAttrSet("data.firehydrant_service.test_service", "owner_id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_service.test_service", "service_tier", "1"),
					resource.TestCheckResourceAttr("data.firehydrant_service.test_service", "team_ids.#", "2"),
				),
			},
		},
	})
}

func testAccServiceDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service" {
  name = "test-service-%s"
}

data "firehydrant_service" "test_service" {
  id = firehydrant_service.test_service.id
}`, rName)
}

func testAccServiceDataSourceConfig_allAttributes(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "test_team1" {
  name = "test-team1-%s"
}

resource "firehydrant_team" "test_team2" {
  name = "test-team2-%s"
}

resource "firehydrant_team" "test_team3" {
  name = "test-team3-%s"
}

resource "firehydrant_service" "test_service" {
  name         = "test-service-%s"
  alert_on_add = true
  description  = "test-description-%s"
  owner_id     = firehydrant_team.test_team1.id
  service_tier = "1"
  team_ids = [
    firehydrant_team.test_team2.id,
    firehydrant_team.test_team3.id
  ]
}

data "firehydrant_service" "test_service" {
  id = firehydrant_service.test_service.id
}`, rName, rName, rName, rName, rName)
}
