package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLifecyclePhaseDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLifecyclePhaseDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_lifecycle_phase.started", "id"),
					resource.TestCheckResourceAttr("data.firehydrant_lifecycle_phase.started", "name", "started"),
				),
			},
		},
	})
}

func TestAccLifecyclePhaseDataSource_invalid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccLifecyclePhaseDataSourceConfig_invalid(),
				ExpectError: regexp.MustCompile(`Lifecycle phase foo is invalid.  Valid lifecycle phases are`),
			},
		},
	})
}

func testAccLifecyclePhaseDataSourceConfig_basic() string {
	return `
data "firehydrant_lifecycle_phase" "started" {
  name = "started"
}`
}

func testAccLifecyclePhaseDataSourceConfig_invalid() string {
	return `
data "firehydrant_lifecycle_phase" "invalid" {
  name = "foo"
}`
}
