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
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_name"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_team_id"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_is_pageable"),
				),
			},
			{
				Config: testAccFireHydrantSignalRuleConfigBasic("LOW"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "notification_priority_override", "LOW"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_name"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_team_id"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_is_pageable"),
				),
			},
			{
				Config: testAccFireHydrantSignalRuleConfigBasic("HIGH"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "notification_priority_override", "HIGH"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_name"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_team_id"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_is_pageable"),
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

func TestAccFireHydrantSignalRule_createIncidentConditionWhen(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccFireHydrantSignalRuleConfigWithIncidentCondition("ALWAYS", "PT30M"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "create_incident_condition_when", "ALWAYS"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "deduplication_expiry", "PT30M"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_name"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_team_id"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_is_pageable"),
				),
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

func TestAccFireHydrantSignalRule_NotificationPriorityAddRemove(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				// First create with a priority set
				Config: testAccFireHydrantSignalRuleConfigWithPriority("HIGH"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "notification_priority_override", "HIGH"),
				),
			},
			{
				// Then update to remove the priority
				Config: testAccFireHydrantSignalRuleConfigWithoutPriority(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "notification_priority_override", ""),
				),
			},
		},
	})
}

func testAccFireHydrantSignalRuleConfigWithPriority(priority string) string {
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

func testAccFireHydrantSignalRuleConfigWithoutPriority() string {
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

func testAccFireHydrantSignalRuleConfigWithIncidentCondition(condition, expiry string) string {
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
		create_incident_condition_when = "%s"
		deduplication_expiry = "%s"
	}
	`, os.Getenv("EXISTING_USER_EMAIL"), condition, expiry)
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
