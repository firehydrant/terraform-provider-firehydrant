package provider

import (
	"fmt"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSeverityDataSource_basic(t *testing.T) {
	t.Parallel()
	slug := "TESTSEV" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSeverityDataSourceConfig_basic(slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_severity.test_severity", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_severity.test_severity", "slug", slug),
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
	t.Parallel()
	slug := "TESTSEV" + acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSeverityDataSourceConfig_allAttributes(slug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_severity.test_severity", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_severity.test_severity", "slug", slug),
					resource.TestCheckResourceAttr(
						"data.firehydrant_severity.test_severity", "description", "test-description"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_severity.test_severity", "type", string(firehydrant.SeverityTypeGameday)),
				),
			},
		},
	})
}

func testAccSeverityDataSourceConfig_basic(slug string) string {
	return fmt.Sprintf(`
resource "firehydrant_severity" "test_severity" {
  slug        = "%s"
  description = "test-description"
}

data "firehydrant_severity" "test_severity" {
  slug = firehydrant_severity.test_severity.id
}`, slug)
}

func testAccSeverityDataSourceConfig_allAttributes(slug string) string {
	return fmt.Sprintf(`
resource "firehydrant_severity" "test_severity" {
  slug        = "%s"
  description = "test-description"
  type        = "gameday"
}

data "firehydrant_severity" "test_severity" {
  slug = firehydrant_severity.test_severity.id
}`, slug)
}
