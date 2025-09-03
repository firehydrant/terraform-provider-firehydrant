package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPermissionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
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

// Test that demonstrates the workflow: get permissions, then create a role with them
func TestAccPermissionsAndRoleWorkflow(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckRoleResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionsAndRoleWorkflowConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify we can create a role with permissions
					resource.TestCheckResourceAttrSet("firehydrant_role.test_role", "id"),
					resource.TestCheckResourceAttr("firehydrant_role.test_role", "name", "test-permissions-workflow"),
					resource.TestCheckResourceAttr("firehydrant_role.test_role", "permissions.#", "1"),
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

func testAccPermissionsAndRoleWorkflowConfig() string {
	return `
data "firehydrant_permissions" "all" {}

resource "firehydrant_role" "test_role" {
  name        = "test-permissions-workflow"
  description = "Test role created via permissions workflow"
  permissions = ["read_users"]
}
`
}
