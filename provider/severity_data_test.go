package provider

import (
	"fmt"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

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
						"data.firehydrant_severity.test_severity", "slug", "TESTSEVERITYBASIC"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_severity.test_severity", "description", "test-description"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_severity.test_severity", "type", string(firehydrant.SeverityTypeUnexpectedDowntime)),
				),
			},
		},
	})
}

func TestAccSeverityDataSource_allAttributes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSeverityDataSourceConfig_allAttributes(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_severity.test_severity", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_severity.test_severity", "slug", "TESTSEVERITYALL"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_severity.test_severity", "description", "test-description"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_severity.test_severity", "type", string(firehydrant.SeverityTypeGameday)),
				),
			},
		},
	})
}

func testAccSeverityDataSourceConfig_basic() string {
	return fmt.Sprintln(`
resource "firehydrant_severity" "test_severity" {
  slug        = "TESTSEVERITYBASIC"
  description = "test-description"
}

data "firehydrant_severity" "test_severity" {
  slug = firehydrant_severity.test_severity.id
}`)
}

func testAccSeverityDataSourceConfig_allAttributes() string {
	return fmt.Sprintln(`
resource "firehydrant_severity" "test_severity" {
  slug        = "TESTSEVERITYALL"
  description = "test-description"
  type        = "gameday"
}

data "firehydrant_severity" "test_severity" {
  slug = firehydrant_severity.test_severity.id
}`)
}
