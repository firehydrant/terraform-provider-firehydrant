package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPriorityDataSource_basic(t *testing.T) {
	rSlug := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPriorityDataSourceConfig_basic(rSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_priority.test_priority", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_priority.test_priority", "slug", fmt.Sprintf("TESTPRIORITY%s", rSlug)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_priority.test_priority", "description", "test-description"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_priority.test_priority", "default", "false"),
				),
			},
		},
	})
}

func testAccPriorityDataSourceConfig_basic(rSlug string) string {
	return fmt.Sprintf(`
resource "firehydrant_priority" "test_priority" {
  slug        = "TESTPRIORITY%s"
  description = "test-description"
  default     = false
}

data "firehydrant_priority" "test_priority" {
  slug = firehydrant_priority.test_priority.id
}`, rSlug)
}
