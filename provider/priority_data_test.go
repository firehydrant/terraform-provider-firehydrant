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
					resource.TestCheckResourceAttrSet("data.firehydrant_priority.test_priority", "slug"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_priority.test_priority", "description", "test priority"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_priority.test_priority", "default", "false"),
				),
			},
		},
	})
}

func testAccPriorityDataSourceConfig_basic() string {
	return fmt.Sprintln(`
data "firehydrant_priority" "test_priority" {
  slug          = "test_priority"
  description   = "test priority"
  default       = false
}`)
}
