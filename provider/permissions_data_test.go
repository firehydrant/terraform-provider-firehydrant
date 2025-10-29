package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPermissionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_permissions.all", "id"),
					resource.TestCheckResourceAttrSet("data.firehydrant_permissions.all", "permissions.#"),
					// Check that we get permissions and they have the expected structure
					resource.TestCheckResourceAttrSet("data.firehydrant_permissions.all", "permissions.0.slug"),
					resource.TestCheckResourceAttrSet("data.firehydrant_permissions.all", "permissions.0.display_name"),
					resource.TestCheckResourceAttrSet("data.firehydrant_permissions.all", "permissions.0.description"),
					resource.TestCheckResourceAttrSet("data.firehydrant_permissions.all", "permissions.0.category_slug"),
					resource.TestCheckResourceAttrSet("data.firehydrant_permissions.all", "permissions.0.category_display_name"),
					resource.TestCheckResourceAttrSet("data.firehydrant_permissions.all", "permissions.0.available"),
				),
			},
		},
	})
}

func testAccPermissionsDataSourceConfig() string {
	return `
data "firehydrant_permissions" "all" {}
`
}
