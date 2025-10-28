package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccServicesDataSource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServicesDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_services.all_services", "services.#"),
					testAccCheckServicesSet("data.firehydrant_services.all_services"),
				),
			},
		},
	})
}

func testAccCheckServicesSet(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		servicesResource, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Can't find services resource in state: %s", name)
		}

		if servicesResource.Primary.ID == "" {
			return fmt.Errorf("Services resource ID not set")
		}

		attributes := servicesResource.Primary.Attributes
		services, servicesOk := attributes["services.#"]
		if !servicesOk {
			return fmt.Errorf("Services list is missing")
		}

		servicesCount, err := strconv.Atoi(services)
		if err != nil {
			return err
		}

		if servicesCount <= 1 {
			return fmt.Errorf("Incorrect number of services - expected at least 1, got %d", servicesCount)
		}

		return nil
	}
}

func testAccServicesDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service" {
  name = "test-service-%s"
}

data "firehydrant_services" "all_services" {
}`, rName)
}
