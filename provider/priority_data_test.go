package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPriorityDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPriorityDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_priority.test_priority", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_priority.test_priority", "slug", "TESTPRIORITY"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_priority.test_priority", "description", "test-description"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_priority.test_priority", "default", "true"),
				),
			},
		},
	})
}

func testAccPriorityDataSourceConfig_basic() string {
	return fmt.Sprintln(`
resource "firehydrant_priority" "test_priority" {
  slug        = "TESTPRIORITY"
  description = "test-description"
  default     = true
}

data "firehydrant_priority" "test_priority" {
  slug = firehydrant_priority.test_priority.id
}`)
}
