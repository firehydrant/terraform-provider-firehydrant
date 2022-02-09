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
				Config: testAccServiceDataSourceConfig(rName),
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

func testAccServiceDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service" {
  name = "test-service-%s"
}

data "firehydrant_service" "test_service" {
  id = firehydrant_service.test_service.id
}`, rName)
}
