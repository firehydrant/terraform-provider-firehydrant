package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestUserDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testUserDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_user.test_user", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_user.test_user", "email", "test-user@firehydrant.io"),
				),
			},
		},
	})
}

func testUserDataSourceConfig_basic() string {
	return fmt.Sprintln(`
data "firehydrant_user" "test_user" {
  email = "test-user@firehydrant.io"
}`)
}
