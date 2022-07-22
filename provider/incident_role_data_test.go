package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIncidentRoleDataSource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccIncidentRoleDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_incident_role.test_incident_role", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_role.test_incident_role", "name", fmt.Sprintf("test-incident-role-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_role.test_incident_role", "summary", fmt.Sprintf("test-summary-%s", rName)),
				),
			},
		},
	})
}

func TestAccIncidentRoleDataSource_allAttributes(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccIncidentRoleDataSourceConfig_allAttributes(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_incident_role.test_incident_role", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_role.test_incident_role", "name", fmt.Sprintf("test-incident-role-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_role.test_incident_role", "description", fmt.Sprintf("test-description-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_role.test_incident_role", "summary", fmt.Sprintf("test-summary-%s", rName)),
				),
			},
		},
	})
}

func testAccIncidentRoleDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_incident_role" "test_incident_role" {
  name    = "test-incident-role-%s"
  summary = "test-summary-%s"
}

data "firehydrant_incident_role" "test_incident_role" {
  id = firehydrant_incident_role.test_incident_role.id
}`, rName, rName)
}

func testAccIncidentRoleDataSourceConfig_allAttributes(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_incident_role" "test_incident_role" {
  name        = "test-incident-role-%s"
  description = "test-description-%s"
  summary     = "test-summary-%s"
}

data "firehydrant_incident_role" "test_incident_role" {
  id = firehydrant_incident_role.test_incident_role.id
}`, rName, rName, rName)
}
