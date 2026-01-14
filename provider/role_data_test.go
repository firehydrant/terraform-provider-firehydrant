package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRoleDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccBuiltInRoleDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Test data source lookup by slug
					resource.TestCheckResourceAttrSet("data.firehydrant_role.member", "id"),
					resource.TestCheckResourceAttr("data.firehydrant_role.member", "slug", "member"),
					resource.TestCheckResourceAttrSet("data.firehydrant_role.member", "name"),
					resource.TestCheckResourceAttrSet("data.firehydrant_role.member", "description"),

					resource.TestCheckResourceAttrSet("data.firehydrant_role.member", "permissions.#"),
				),
			},
		},
	})
}

func TestAccRoleDataSource_builtIn(t *testing.T) {
	// Test looking up a built-in role that should always exist
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccBuiltInRoleDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_role.member", "id"),
					resource.TestCheckResourceAttr("data.firehydrant_role.member", "slug", "member"),
					resource.TestCheckResourceAttrSet("data.firehydrant_role.member", "permissions.#"),
				),
			},
		},
	})
}

func testAccBuiltInRoleDataSourceConfig() string {
	return `
data "firehydrant_role" "member" {
	slug = "member"
}
`
}
