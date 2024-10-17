package provider

import (
	"fmt"
	"os"
	"regexp"
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
				Config: testAccFireHydrantSignalRuleConfigBasic("MEDIUM"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "notification_priority_override", "MEDIUM"),
				),
			},
			{
				Config: testAccFireHydrantSignalRuleConfigBasic("LOW"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "notification_priority_override", "LOW"),
				),
			},
			{
				Config: testAccFireHydrantSignalRuleConfigBasic("HIGH"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "notification_priority_override", "HIGH"),
				),
			},
		},
	})
}

func TestAccFireHydrantSignalRule_invalidPriority(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccFireHydrantSignalRuleConfigBasic("INVALID"),
				ExpectError: regexp.MustCompile(`expected notification_priority_override to be one of \[LOW MEDIUM HIGH\], got INVALID`),
			},
		},
	})
}

func testAccFireHydrantSignalRuleConfigBasic(priority string) string {
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
			notification_priority_override = "%s"
	}
	`, os.Getenv("EXISTING_USER_EMAIL"), priority)
}

func TestAccFireHydrantSignalRule_IncidentTypeIDMissing(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccFireHydrantSignalRuleConfigIncidentTypeIDMissing(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
				),
			},
		},
	})
}

func testAccFireHydrantSignalRuleConfigIncidentTypeIDMissing() string {
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
