package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFunctionalityDataSource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionalityDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_functionality.test_functionality", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_functionality.test_functionality", "name", fmt.Sprintf("test-functionality-%s", rName)),
				),
			},
		},
	})
}

func TestAccFunctionalityDataSource_allAttributes(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionalityDataSourceConfig_allAttributes(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_functionality.test_functionality", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_functionality.test_functionality", "name", fmt.Sprintf("test-functionality-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_functionality.test_functionality", "description", fmt.Sprintf("test-description-%s", rName)),
					resource.TestCheckResourceAttr("data.firehydrant_functionality.test_functionality", "service_ids.#", "1"),
				),
			},
		},
	})
}

func testAccFunctionalityDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_functionality" "test_functionality" {
  name = "test-functionality-%s"
}

data "firehydrant_functionality" "test_functionality" {
  functionality_id = firehydrant_functionality.test_functionality.id
}`, rName)
}

func testAccFunctionalityDataSourceConfig_allAttributes(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service" {
  name = "test-service-%s"
}

resource "firehydrant_functionality" "test_functionality" {
  name = "test-functionality-%s"
  description = "test-description-%s"
  service_ids = [firehydrant_service.test_service.id]
  labels = {
    test1 = "test-label1-foo",
  }
}

data "firehydrant_functionality" "test_functionality" {
  functionality_id = firehydrant_functionality.test_functionality.id
}`, rName, rName, rName)
}
