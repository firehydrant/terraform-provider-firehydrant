package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFireHydrantStatusUpdateTemplate_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccFireHydrantStatusUpdateTemplateConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantStatusUpdateTemplateExists("firehydrant_status_update_template.test"),
				),
			},
		},
	})
}

func testAccFireHydrantStatusUpdateTemplateConfigBasic() string {
	return `
	resource "firehydrant_status_update_template" "test" {
		name = "test-signal-rule"
		body = "This is the template body"
	}
	`
}

func testAccCheckFireHydrantStatusUpdateTemplateExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource ID is set")
		}

		return nil
	}
}
