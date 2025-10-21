package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFireHydrantSignalRule_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckFireHydrantSignalRuleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFireHydrantSignalRuleConfigBasic(rName, "MEDIUM"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "notification_priority_override", "MEDIUM"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "create_incident_condition_when", "WHEN_UNSPECIFIED"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_name"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_is_pageable"),
				),
			},
			{
				Config: testAccFireHydrantSignalRuleConfigBasic(rName, "LOW"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "notification_priority_override", "LOW"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "create_incident_condition_when", "WHEN_UNSPECIFIED"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_name"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_is_pageable"),
				),
			},
			{
				Config: testAccFireHydrantSignalRuleConfigBasic(rName, "HIGH"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "notification_priority_override", "HIGH"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "create_incident_condition_when", "WHEN_UNSPECIFIED"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_name"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_is_pageable"),
				),
			},
		},
	})
}

func TestAccFireHydrantSignalRule_invalidPriority(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckFireHydrantSignalRuleDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccFireHydrantSignalRuleConfigBasic(rName, "INVALID"),
				ExpectError: regexp.MustCompile(`expected notification_priority_override to be one of \[LOW MEDIUM HIGH\], got INVALID`),
			},
		},
	})
}

func TestAccFireHydrantSignalRule_createIncidentConditionWhen(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckFireHydrantSignalRuleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFireHydrantSignalRuleConfigWithIncidentCondition(rName, "WHEN_ALWAYS", "PT30M"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "create_incident_condition_when", "WHEN_ALWAYS"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "deduplication_expiry", "PT30M"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "notification_priority_override", "MEDIUM"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_name"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_rule.test", "target_is_pageable"),
				),
			},
		},
	})
}

func testAccFireHydrantSignalRuleConfigBasic(rName, priority string) string {
	existingUser := os.Getenv("EXISTING_USER_EMAIL")
	if existingUser == "" {
		existingUser = "local@firehydrant.io"
	}

	return fmt.Sprintf(`
	data "firehydrant_user" "test_user" {
			email = "%s"
	}

	resource "firehydrant_team" "test" {
			name = "test-team-%s"
	}

	resource "firehydrant_signal_rule" "test" {
		team_id = firehydrant_team.test.id
		name = "test-signal-rule-%s"
		expression = "signal.summary == 'test-signal-summary-%s'"
		target_type = "User"
		target_id = data.firehydrant_user.test_user.id
		notification_priority_override = "%s"
		create_incident_condition_when = "WHEN_UNSPECIFIED"
	}
	`, existingUser, rName, rName, rName, priority)
}

func TestAccFireHydrantSignalRule_IncidentTypeIDMissing(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckFireHydrantSignalRuleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFireHydrantSignalRuleConfigIncidentTypeIDMissing(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "notification_priority_override", "LOW"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "create_incident_condition_when", "WHEN_UNSPECIFIED"),
				),
			},
		},
	})
}

func TestAccFireHydrantSignalRule_NotificationPriorityAddRemove(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckFireHydrantSignalRuleDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFireHydrantSignalRuleConfigWithPriority(rName, "HIGH"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "notification_priority_override", "HIGH"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "create_incident_condition_when", "WHEN_UNSPECIFIED"),
				),
			},
			{
				Config: testAccFireHydrantSignalRuleConfigWithPriority(rName, "MEDIUM"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "notification_priority_override", "MEDIUM"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "create_incident_condition_when", "WHEN_UNSPECIFIED"),
				),
			},
			{
				Config: testAccFireHydrantSignalRuleConfigWithPriority(rName, "HIGH"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFireHydrantSignalRuleExists("firehydrant_signal_rule.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "notification_priority_override", "HIGH"),
					resource.TestCheckResourceAttr("firehydrant_signal_rule.test", "create_incident_condition_when", "WHEN_UNSPECIFIED"),
				),
			},
		},
	})
}

// TODO: After the Go SDK omitempty limitation is fixed, add a test to verify that
// notification_priority_override can be properly cleared
func testAccFireHydrantSignalRuleConfigWithPriority(rName, priority string) string {
	existingUser := os.Getenv("EXISTING_USER_EMAIL")
	if existingUser == "" {
		existingUser = "local@firehydrant.io"
	}

	return fmt.Sprintf(`
	data "firehydrant_user" "test_user" {
		email = "%s"
	}

	resource "firehydrant_team" "test" {
		name = "test-team-%s"
	}

	resource "firehydrant_signal_rule" "test" {
		team_id = firehydrant_team.test.id
		name = "test-signal-rule-%s"
		expression = "signal.summary == 'test-signal-summary-%s'"
		target_type = "User"
		target_id = data.firehydrant_user.test_user.id
		notification_priority_override = "%s"
		create_incident_condition_when = "WHEN_UNSPECIFIED"
	}
	`, existingUser, rName, rName, rName, priority)
}

func testAccFireHydrantSignalRuleConfigWithoutPriority(rName string) string {
	existingUser := os.Getenv("EXISTING_USER_EMAIL")
	if existingUser == "" {
		existingUser = "local@firehydrant.io"
	}

	return fmt.Sprintf(`
	data "firehydrant_user" "test_user" {
		email = "%s"
	}

	resource "firehydrant_team" "test" {
		name = "test-team-%s"
	}

	resource "firehydrant_signal_rule" "test" {
		team_id = firehydrant_team.test.id
		name = "test-signal-rule-%s"
		expression = "signal.summary == 'test-signal-summary-%s'"
		target_type = "User"
		target_id = data.firehydrant_user.test_user.id
		create_incident_condition_when = "WHEN_UNSPECIFIED"
	}
	`, existingUser, rName, rName, rName)
}

func testAccFireHydrantSignalRuleConfigIncidentTypeIDMissing(rName string) string {
	existingUser := os.Getenv("EXISTING_USER_EMAIL")
	if existingUser == "" {
		existingUser = "local@firehydrant.io"
	}

	return fmt.Sprintf(`
	data "firehydrant_user" "test_user" {
		email = "%s"
	}

	resource "firehydrant_team" "test" {
		name = "test-team-%s"
	}

	resource "firehydrant_signal_rule" "test" {
		team_id = firehydrant_team.test.id
		name = "test-signal-rule-%s"
		expression = "signal.summary == 'test-signal-summary-%s'"
		target_type = "User"
		target_id = data.firehydrant_user.test_user.id
		notification_priority_override = "LOW"
		create_incident_condition_when = "WHEN_UNSPECIFIED"
	}
	`, existingUser, rName, rName, rName)
}

func testAccFireHydrantSignalRuleConfigWithIncidentCondition(rName, condition, expiry string) string {
	existingUser := os.Getenv("EXISTING_USER_EMAIL")
	if existingUser == "" {
		existingUser = "local@firehydrant.io"
	}

	return fmt.Sprintf(`
	data "firehydrant_user" "test_user" {
		email = "%s"
	}

	resource "firehydrant_team" "test" {
		name = "test-team-%s"
	}

	resource "firehydrant_signal_rule" "test" {
		team_id = firehydrant_team.test.id
		name = "test-signal-rule-%s"
		expression = "signal.summary == 'test-signal-summary-%s'"
		target_type = "User"
		target_id = data.firehydrant_user.test_user.id
		create_incident_condition_when = "%s"
		deduplication_expiry = "%s"
		notification_priority_override = "MEDIUM"
	}
	`, existingUser, rName, rName, rName, condition, expiry)
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

func testAccCheckFireHydrantSignalRuleDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "firehydrant_signal_rule" {
				continue
			}

			// Check if the signal rule still exists
			_, err := client.SignalsRules().Get(context.TODO(), rs.Primary.Attributes["team_id"], rs.Primary.ID)
			if err != nil && !errors.Is(err, firehydrant.ErrorNotFound) {
				return fmt.Errorf("Signal rule %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}
