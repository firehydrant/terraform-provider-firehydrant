package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRunbookResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckRunbookResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRunbookResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRunbookResourceExistsWithAttributes_basic("firehydrant_runbook.test_runbook"),
					resource.TestCheckResourceAttrSet("firehydrant_runbook.test_runbook", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "name", fmt.Sprintf("test-runbook-%s", rName)),
					resource.TestCheckResourceAttrSet("firehydrant_runbook.test_runbook", "attachment_rule"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.#", "1"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.0.name", "Create Incident Channel"),
					resource.TestCheckResourceAttrSet(
						"firehydrant_runbook.test_runbook", "steps.0.action_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.0.automatic", "false"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.0.repeats", "false"),
				),
			},
		},
	})
}

func TestAccRunbookResource_update(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckRunbookResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRunbookResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRunbookResourceExistsWithAttributes_basic("firehydrant_runbook.test_runbook"),
					resource.TestCheckResourceAttrSet("firehydrant_runbook.test_runbook", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "name", fmt.Sprintf("test-runbook-%s", rName)),
					resource.TestCheckResourceAttrSet("firehydrant_runbook.test_runbook", "attachment_rule"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.#", "1"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.0.name", "Create Incident Channel"),
					resource.TestCheckResourceAttrSet(
						"firehydrant_runbook.test_runbook", "steps.0.action_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.0.automatic", "false"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.0.repeats", "false"),
				),
			},
			{
				Config: testAccRunbookResourceConfig_update(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRunbookResourceExistsWithAttributes_update("firehydrant_runbook.test_runbook"),
					resource.TestCheckResourceAttrSet("firehydrant_runbook.test_runbook", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "name", fmt.Sprintf("test-runbook-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
					resource.TestCheckResourceAttrSet("firehydrant_runbook.test_runbook", "owner_id"),
					resource.TestCheckResourceAttrSet("firehydrant_runbook.test_runbook", "attachment_rule"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.#", "2"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.0.name", "Notify Channel"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.0.automatic", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.0.repeats", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.0.repeats_duration", "PT15M"),
					resource.TestCheckResourceAttrSet(
						"firehydrant_runbook.test_runbook", "steps.0.action_id"),
					resource.TestCheckResourceAttrSet("firehydrant_runbook.test_runbook", "steps.0.rule"),
				),
			},
			{
				Config: testAccRunbookResourceConfig_basic(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRunbookResourceExistsWithAttributes_basic("firehydrant_runbook.test_runbook"),
					resource.TestCheckResourceAttrSet("firehydrant_runbook.test_runbook", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "name", fmt.Sprintf("test-runbook-%s", rNameUpdated)),
					resource.TestCheckResourceAttrSet("firehydrant_runbook.test_runbook", "attachment_rule"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.#", "1"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.0.name", "Create Incident Channel"),
					resource.TestCheckResourceAttrSet(
						"firehydrant_runbook.test_runbook", "steps.0.action_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.0.automatic", "false"),
					resource.TestCheckResourceAttr(
						"firehydrant_runbook.test_runbook", "steps.0.repeats", "false"),
				),
			},
		},
	})
}

func TestAccRunbookResource_validateSchemaAttributesStepsConfig(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccRunbookResourceConfig_stepsConfigInvalidJSON(rName),
				ExpectError: regexp.MustCompile(`"config" contains an invalid JSON`),
			},
		},
	})
}

func TestAccRunbookResource_validateSchemaAttributesAttachmentRule(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccRunbookResourceConfig_attachmentRuleInvalidJSON(rName),
				ExpectError: regexp.MustCompile(`"attachment_rule" contains an invalid JSON`),
			},
		},
	})
}

func TestAccRunbookResource_validateSchemaAttributesStepsRule(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccRunbookResourceConfig_stepsRuleInvalidJSON(rName),
				ExpectError: regexp.MustCompile(`"rule" contains an invalid JSON`),
			},
		},
	})
}

func TestAccRunbookResource_validateSchemaAttributesStepsRepeatsDuration(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccRunbookResourceConfig_stepsRequiredRepeatsDurationNotSet(rName),
				ExpectError: regexp.MustCompile("Error: step repeats requires step repeats_duration to be set"),
			},
		},
	})
}

func TestAccRunbookResourceImport_validateSchemaAttributesStepsRepeats(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccRunbookResourceConfig_stepsRequiredRepeatsNotSet(rName),
				ExpectError: regexp.MustCompile("Error: step repeats_duration requires step repeats to be set to true"),
			},
		},
	})
}

func TestAccRunbookResourceImport_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRunbookResourceConfig_basic(rName),
			},
			{
				ResourceName:      "firehydrant_runbook.test_runbook",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRunbookResourceImport_allAttributes(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRunbookResourceConfig_update(rName),
			},
			{
				ResourceName:      "firehydrant_runbook.test_runbook",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckRunbookResourceExistsWithAttributes_basic(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		runbookResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if runbookResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		runbookResponse, err := client.Runbooks().Get(context.TODO(), runbookResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := runbookResource.Primary.Attributes["name"], runbookResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		if runbookResponse.Description != "" {
			return fmt.Errorf("Unexpected description. Expected no description, got: %s", runbookResponse.Description)
		}

		if runbookResponse.Owner != nil {
			return fmt.Errorf("Unexpected owner. Expected no owner ID, got: %s", runbookResponse.Owner.ID)
		}

		if runbookResponse.AttachmentRule == nil {
			return fmt.Errorf("Unexpected attachment_rule. Expected attachment_rule to be set.")
		}
		var attachmentRule []byte
		if len(runbookResponse.AttachmentRule) > 0 {
			attachmentRule, err = json.Marshal(runbookResponse.AttachmentRule)
			if err != nil {
				return fmt.Errorf("Unexpected error converting attachment_rule to JSON: %v", err)
			}
		}
		normalizedAttachmentRuleJSON, _ := structure.NormalizeJsonString(firehydrant.RunbookAttachmentRuleDefaultJSON)
		if err != nil {
			return fmt.Errorf("Unexpected error normalizing runbook default attachment_rule JSON: %v", err)
		}
		if string(attachmentRule) != normalizedAttachmentRuleJSON {
			return fmt.Errorf("Unexpected attachment_rule. Expected attachment_rule to be set to the default value %s, got: %s", firehydrant.RunbookAttachmentRuleDefaultJSON, string(attachmentRule))
		}
		if runbookResource.Primary.Attributes["attachment_rule"] != string(attachmentRule) {
			return fmt.Errorf("Unexpected attachment_rule. Expected %s, got: %s", runbookResponse.AttachmentRule, runbookResource.Primary.Attributes["attachment_rule"])
		}

		if len(runbookResponse.Steps) != 1 {
			return fmt.Errorf("Unexpected number of steps. Expected 1 step, got: %v", len(runbookResponse.Steps))
		}

		for index, step := range runbookResponse.Steps {
			key := fmt.Sprintf("steps.%d", index)
			if runbookResource.Primary.Attributes[key+".name"] != step.Name {
				return fmt.Errorf("Unexpected runbook step name. Expected %s, got: %s", step.Name, runbookResource.Primary.Attributes[key+".name"])
			}

			if runbookResource.Primary.Attributes[key+".action_id"] != step.ActionID {
				return fmt.Errorf("Unexpected runbook step action_id. Expected %s, got: %s", step.Name, runbookResource.Primary.Attributes[key+".action_id"])
			}

			var config []byte
			if len(step.Config) > 0 {
				config, err = json.Marshal(step.Config)
				if err != nil {
					return fmt.Errorf("Unexpected error converting runbook step config to JSON: %v", err)
				}
			}
			if runbookResource.Primary.Attributes[key+".config"] != string(config) {
				return fmt.Errorf("Unexpected runbook step config. Expected %s, got: %s", step.Config, runbookResource.Primary.Attributes[key+".config"])
			}

			if runbookResource.Primary.Attributes[key+".automatic"] != fmt.Sprintf("%t", step.Automatic) {
				return fmt.Errorf("Unexpected runbook step automatic. Expected %t, got: %s", step.Automatic, runbookResource.Primary.Attributes[key+".automatic"])
			}

			if runbookResource.Primary.Attributes[key+".repeats"] != fmt.Sprintf("%t", step.Repeats) {
				return fmt.Errorf("Unexpected runbook step repeats. Expected %t, got: %s", step.Repeats, runbookResource.Primary.Attributes[key+".repeats"])
			}

			if runbookResource.Primary.Attributes[key+".repeats_duration"] != step.RepeatsDuration {
				return fmt.Errorf("Unexpected runbook step repeats_duration. Expected %s, got: %s", step.RepeatsDuration, runbookResource.Primary.Attributes[key+".repeats_duration"])
			}

			var rule []byte
			if len(step.Rule) > 0 {
				rule, err = json.Marshal(step.Rule)
				if err != nil {
					return fmt.Errorf("Unexpected error converting runbook step rule to JSON: %v", err)
				}
			}
			if runbookResource.Primary.Attributes[key+".rule"] != string(rule) {
				return fmt.Errorf("Unexpected runbook step rule. Expected %s, got: %s", step.Rule, runbookResource.Primary.Attributes[key+".rule"])
			}
		}

		return nil
	}
}

func testAccCheckRunbookResourceExistsWithAttributes_update(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		runbookResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if runbookResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		runbookResponse, err := client.Runbooks().Get(context.TODO(), runbookResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := runbookResource.Primary.Attributes["name"], runbookResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = runbookResource.Primary.Attributes["description"], runbookResponse.Description
		if expected != got {
			return fmt.Errorf("Unexpected description. Expected: %s, got: %s", expected, got)
		}

		if runbookResponse.AttachmentRule == nil {
			return fmt.Errorf("Unexpected attachment_rule. Expected attachment_rule to be set.")
		}
		var attachmentRule []byte
		if len(runbookResponse.AttachmentRule) > 0 {
			attachmentRule, err = json.Marshal(runbookResponse.AttachmentRule)
			if err != nil {
				return fmt.Errorf("Unexpected error converting attachment_rule to JSON: %v", err)
			}
		}
		if runbookResource.Primary.Attributes["attachment_rule"] != string(attachmentRule) {
			return fmt.Errorf("Unexpected attachment_rule. Expected %s, got: %s", runbookResponse.AttachmentRule, runbookResource.Primary.Attributes["attachment_rule"])
		}

		if runbookResponse.Owner == nil {
			return fmt.Errorf("Unexpected owner. Expected owner to be set.")
		}
		expected, got = runbookResource.Primary.Attributes["owner_id"], runbookResponse.Owner.ID
		if expected != got {
			return fmt.Errorf("Unexpected owner ID. Expected:%s, got: %s", expected, got)
		}

		if len(runbookResponse.Steps) != 2 {
			return fmt.Errorf("Unexpected number of steps. Expected 2 steps, got: %v", len(runbookResponse.Steps))
		}

		for index, step := range runbookResponse.Steps {
			key := fmt.Sprintf("steps.%d", index)
			if runbookResource.Primary.Attributes[key+".name"] != step.Name {
				return fmt.Errorf("Unexpected runbook step name. Expected %s, got: %s", step.Name, runbookResource.Primary.Attributes[key+".name"])
			}

			if runbookResource.Primary.Attributes[key+".action_id"] != step.ActionID {
				return fmt.Errorf("Unexpected runbook step action_id. Expected %s, got: %s", step.Name, runbookResource.Primary.Attributes[key+".action_id"])
			}

			var config []byte
			if len(step.Config) > 0 {
				config, err = json.Marshal(step.Config)
				if err != nil {
					return fmt.Errorf("Unexpected error converting runbook step config to JSON: %v", err)
				}
			}
			if runbookResource.Primary.Attributes[key+".config"] != string(config) {
				return fmt.Errorf("Unexpected runbook step config. Expected %s, got: %s", step.Config, runbookResource.Primary.Attributes[key+".config"])
			}

			if runbookResource.Primary.Attributes[key+".automatic"] != fmt.Sprintf("%t", step.Automatic) {
				return fmt.Errorf("Unexpected runbook step automatic. Expected %t, got: %s", step.Automatic, runbookResource.Primary.Attributes[key+".automatic"])
			}

			if runbookResource.Primary.Attributes[key+".repeats"] != fmt.Sprintf("%t", step.Repeats) {
				return fmt.Errorf("Unexpected runbook step repeats. Expected %t, got: %s", step.Repeats, runbookResource.Primary.Attributes[key+".repeats"])
			}

			if runbookResource.Primary.Attributes[key+".repeats_duration"] != step.RepeatsDuration {
				return fmt.Errorf("Unexpected runbook step repeats_duration. Expected %s, got: %s", step.RepeatsDuration, runbookResource.Primary.Attributes[key+".repeats_duration"])
			}

			var rule []byte
			if len(step.Rule) > 0 {
				rule, err = json.Marshal(step.Rule)
				if err != nil {
					return fmt.Errorf("Unexpected error converting runbook step rule to JSON: %v", err)
				}
			}
			if runbookResource.Primary.Attributes[key+".rule"] != string(rule) {
				return fmt.Errorf("Unexpected runbook step rule. Expected %s, got: %s", step.Rule, runbookResource.Primary.Attributes[key+".rule"])
			}
		}

		return nil
	}
}

func testAccCheckRunbookResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_runbook" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			_, err := client.Runbooks().Get(context.TODO(), stateResource.Primary.ID)
			if err == nil {
				return fmt.Errorf("Runbook %s still exists", stateResource.Primary.ID)
			}
		}

		return nil
	}
}

func testAccRunbookResourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "firehydrant_runbook_action" "create_incident_channel" {
  slug             = "create_incident_channel"
  integration_slug = "slack"
}

resource "firehydrant_runbook" "test_runbook" {
  name = "test-runbook-%s"

  steps {
    name      = "Create Incident Channel"
    action_id = data.firehydrant_runbook_action.create_incident_channel.id

    config = jsonencode({
      channel_name_format = "-inc-{{ number }}"
    })
  }
}`, rName)
}

func testAccRunbookResourceConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "test_team1" {
  name = "test-team1-%s"
}

data "firehydrant_runbook_action" "notify_channel" {
  slug             = "notify_channel"
  integration_slug = "slack"
}

data "firehydrant_runbook_action" "archive_channel" {
  slug             = "archive_incident_channel"
  integration_slug = "slack"
}

resource "firehydrant_runbook" "test_runbook" {
  name        = "test-runbook-%s"
  description = "test-description-%s"
  owner_id    = firehydrant_team.test_team1.id
  attachment_rule = jsonencode({
    logic = {
      eq = [
        {
          var = "incident_current_milestone"
        },
        {
          var = "usr.1"
        }
      ]
    }
    user_data = {
      "1" = {
        type  = "Milestone"
        value = "started"
        label = "Started"
      }
    }
  })

  steps {
    name             = "Notify Channel"
    action_id        = data.firehydrant_runbook_action.notify_channel.id
    automatic        = true
    repeats          = true
    repeats_duration = "PT15M"

    config = jsonencode({
      channels = "#incidents"
    })
    rule = jsonencode({
      logic = {
        eq = [
          {
            var = "incident_current_milestone",
          },
          {
            var = "usr.1"
          }
        ]
      },
      user_data = {
        "1" = {
          type  = "Milestone",
          value = "resolved",
          label = "Resolved"
        }
      }
    })
  }

  steps {
    name      = "Archive Channel"
    action_id = data.firehydrant_runbook_action.archive_channel.id
  }
}
`, rName, rName, rName)
}

func testAccRunbookResourceConfig_stepsRequiredRepeatsDurationNotSet(rName string) string {
	return fmt.Sprintf(`
data "firehydrant_runbook_action" "create_incident_channel" {
  slug             = "create_incident_channel"
  integration_slug = "slack"
}

resource "firehydrant_runbook" "test_runbook" {
  name = "test-runbook-%s"
  attachment_rule = jsonencode({
    logic = {
      eq = [
        {
          var = "incident_current_milestone"
        },
        {
          var = "usr.1"
        }
      ]
    }
    user_data = {
      "1" = {
        type  = "Milestone"
        value = "started"
        label = "Started"
      }
    }
  })

  steps {
    name      = "Create Incident Channel"
    repeats   = true
    action_id = data.firehydrant_runbook_action.create_incident_channel.id

    config = jsonencode({
      channel_name_format = "-inc-{{ number }}"
    })
    rule = jsonencode({
      logic = {
        eq = [
          {
            var = "incident_current_milestone",
          },
          {
            var = "usr.1"
          }
        ]
      },
      user_data = {
        "1" = {
          type  = "Milestone",
          value = "resolved",
          label = "Resolved"
        }
      }
    })
  }
}`, rName)
}

func testAccRunbookResourceConfig_stepsRequiredRepeatsNotSet(rName string) string {
	return fmt.Sprintf(`
data "firehydrant_runbook_action" "create_incident_channel" {
  slug             = "create_incident_channel"
  integration_slug = "slack"
}

resource "firehydrant_runbook" "test_runbook" {
  name = "test-runbook-%s"
  attachment_rule = jsonencode({
    logic = {
      eq = [
        {
          var = "incident_current_milestone"
        },
        {
          var = "usr.1"
        }
      ]
    }
    user_data = {
      "1" = {
        type  = "Milestone"
        value = "started"
        label = "Started"
      }
    }
  })

  steps {
    name             = "Create Incident Channel"
    repeats_duration = "PT15M"
    action_id        = data.firehydrant_runbook_action.create_incident_channel.id

    config = jsonencode({
      channel_name_format = "-inc-{{ number }}"
    })
    rule = jsonencode({
      logic = {
        eq = [
          {
            var = "incident_current_milestone",
          },
          {
            var = "usr.1"
          }
        ]
      },
      user_data = {
        "1" = {
          type  = "Milestone",
          value = "resolved",
          label = "Resolved"
        }
      }
    })
  }
}`, rName)
}

func testAccRunbookResourceConfig_stepsConfigInvalidJSON(rName string) string {
	return fmt.Sprintf(`
data "firehydrant_runbook_action" "create_incident_channel" {
  slug             = "create_incident_channel"
  integration_slug = "slack"
}

resource "firehydrant_runbook" "test_runbook" {
  name = "test-runbook-%s"

  steps {
    name      = "Create Incident Channel"
    action_id = data.firehydrant_runbook_action.create_incident_channel.id

    config = "{invalid_json = {{}}"
  }
}`, rName)
}

func testAccRunbookResourceConfig_attachmentRuleInvalidJSON(rName string) string {
	return fmt.Sprintf(`
data "firehydrant_runbook_action" "create_incident_channel" {
  slug             = "create_incident_channel"
  integration_slug = "slack"
}

resource "firehydrant_runbook" "test_runbook" {
  name            = "test-runbook-%s"
  attachment_rule = "{invalid_json = {{}}"

  steps {
    name      = "Create Incident Channel"
    action_id = data.firehydrant_runbook_action.create_incident_channel.id
  }
}`, rName)
}

func testAccRunbookResourceConfig_stepsRuleInvalidJSON(rName string) string {
	return fmt.Sprintf(`
data "firehydrant_runbook_action" "create_incident_channel" {
  slug             = "create_incident_channel"
  integration_slug = "slack"
}

resource "firehydrant_runbook" "test_runbook" {
  name = "test-runbook-%s"

  steps {
    name      = "Create Incident Channel"
    action_id = data.firehydrant_runbook_action.create_incident_channel.id
		rule = "{invalid_json = {{}}"
  }
}`, rName)
}
