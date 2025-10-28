package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEnvironmentDataSource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_environment.test_environment", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_environment.test_environment", "name", fmt.Sprintf("test-environment-%s", rName)),
				),
			},
		},
	})
}

func TestAccEnvironmentDataSource_allAttributes(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentDataSourceConfig_allAttributes(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_environment.test_environment", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_environment.test_environment", "name", fmt.Sprintf("test-environment-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_environment.test_environment", "description", fmt.Sprintf("test-description-%s", rName)),
				),
			},
		},
	})
}

func testAccEnvironmentDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_environment" "test_environment" {
  name    = "test-environment-%s"
}

data "firehydrant_environment" "test_environment" {
  environment_id = firehydrant_environment.test_environment.id
}`, rName)
}

func testAccEnvironmentDataSourceConfig_allAttributes(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_environment" "test_environment" {
  name        = "test-environment-%s"
  description = "test-description-%s"
}

data "firehydrant_environment" "test_environment" {
  environment_id = firehydrant_environment.test_environment.id
}`, rName, rName)
}
