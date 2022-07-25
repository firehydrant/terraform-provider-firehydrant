package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSeverityDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSeverityDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_severity.test_severity", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_severity.test_severity", "slug", "TESTSEVERITY"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_severity.test_severity", "description", "test-description"),
				),
			},
		},
	})
}

func testAccSeverityDataSourceConfig_basic() string {
	return fmt.Sprintln(`
resource "firehydrant_severity" "test_severity" {
  slug        = "TESTSEVERITY"
  description = "test-description"
}

data "firehydrant_severity" "test_severity" {
  slug = firehydrant_severity.test_severity.id
}`)
}
