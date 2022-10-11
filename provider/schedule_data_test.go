package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestScheduleDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testScheduleDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_schedule.test_schedule", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_schedule.test_schedule", "name", "My Rotation"),
				),
			},
		},
	})
}

func testScheduleDataSourceConfig_basic() string {
	return fmt.Sprintln(`
data "firehydrant_schedule" "test_schedule" {
  name = "My Rotation"
}`)
}
