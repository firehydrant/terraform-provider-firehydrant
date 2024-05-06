package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFireHydrantSignalRule_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccFireHydrantSignalRuleConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
				),
			},
		},
	})
}

func testAccFireHydrantSignalRuleConfigBasic() string {
	return fmt.Sprintf(`
	data "firehydrant_user" "test_user" {
		email = "%s"
	}

	resource "firehydrant_team" "test" {
		name = "test-team"
	}

	resource "firehydrant_signal_rule" "test" {
		team_id = firehydrant_team.test.id
		name = "test-signal-rule"
		expression = "signal.summary == 'test-signal-summary'"
		target_type = "User"
		target_id = data.firehydrant_user.test_user.id
		incident_type_id = "00000000-0000-4000-8000-000000000000"
	}
	`, os.Getenv("EXISTING_USER_EMAIL"))
}

func testAccCheckFireHydrantSignalRuleExists(resourceName string) resource.TestCheckFunc {
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
