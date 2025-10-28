package provider

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccEscalationPolicyResource_basic(t *testing.T) {
	sharedTeamID := getSharedTeamID(t)
	sharedScheduleID := getSharedOnCallScheduleID(t)
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckEscalationPolicyResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEscalationPolicyConfig_basic(rName, sharedTeamID, sharedScheduleID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_escalation_policy.test_escalation_policy", "id"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "name", fmt.Sprintf("test-escalation-policy-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "description", fmt.Sprintf("test-description-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "step_strategy", "static"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "step.0.timeout", "PT1M"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "step.0.targets.0.type", "OnCallSchedule"),
					resource.TestCheckResourceAttrSet("firehydrant_escalation_policy.test_escalation_policy", "step.0.targets.0.id"),
				),
			},
		},
	})
}

func TestAccEscalationPolicyResource_dynamicWithPriorityPolicies(t *testing.T) {
	sharedTeamID := getSharedTeamID(t)
	sharedScheduleID := getSharedOnCallScheduleID(t)
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckEscalationPolicyResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEscalationPolicyConfig_dynamicPriority(rName, sharedTeamID, sharedScheduleID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_escalation_policy.test_escalation_policy", "id"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "name", fmt.Sprintf("test-escalation-policy-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "description", fmt.Sprintf("test-description-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "step_strategy", "dynamic_by_priority"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "repetitions", "1"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "step.0.timeout", "PT1M"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "step.0.targets.0.type", "OnCallSchedule"),
					resource.TestCheckResourceAttrSet("firehydrant_escalation_policy.test_escalation_policy", "step.0.targets.0.id"),
					// Test notification priority policies
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.priority", "HIGH"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.repetitions", "2"),
				),
			},
			{
				Config: testAccEscalationPolicyConfig_dynamicPriorityUpdated(rName, sharedTeamID, sharedScheduleID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_escalation_policy.test_escalation_policy", "id"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "name", fmt.Sprintf("test-escalation-policy-updated-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "description", fmt.Sprintf("test-description-updated-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "step_strategy", "dynamic_by_priority"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "repetitions", "1"),
					// Test multiple priority levels
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.priority", "HIGH"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.repetitions", "3"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.1.priority", "MEDIUM"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.1.repetitions", "1"),
				),
			},
		},
	})
}

func TestAccEscalationPolicyResource_dynamicWithHandoffSteps(t *testing.T) {
	sharedTeamID := getSharedTeamID(t)
	sharedScheduleID := getSharedOnCallScheduleID(t)
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckEscalationPolicyResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEscalationPolicyConfig_dynamicWithHandoffSteps(rName, sharedTeamID, sharedScheduleID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_escalation_policy.test_escalation_policy", "id"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "name", fmt.Sprintf("test-escalation-policy-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "step_strategy", "dynamic_by_priority"),
					// Test notification priority policies with handoff steps
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.priority", "HIGH"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.repetitions", "2"),
					resource.TestCheckResourceAttr("firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.handoff_step.0.target_type", "Team"),
					resource.TestCheckResourceAttrSet("firehydrant_escalation_policy.test_escalation_policy", "notification_priority_policies.0.handoff_step.0.target_id"),
				),
			},
		},
	})
}

func testAccEscalationPolicyConfig_basic(rName, sharedTeamID, sharedScheduleID string) string {
	return fmt.Sprintf(`
	resource "firehydrant_escalation_policy" "test_escalation_policy" {
		team_id = "%s"
		name = "test-escalation-policy-%s"
		description = "test-description-%s"
		repetitions = 1
		step_strategy = "static"
		
		step {
			timeout     = "PT1M"

			targets {
				type = "OnCallSchedule"
				id   = "%s"
			}
		}

		handoff_step {
			target_type = "Team"
			target_id   = "%s"
		}
	}
	`, sharedTeamID, rName, rName, sharedScheduleID, sharedTeamID)
}

func testAccEscalationPolicyConfig_dynamicPriority(rName, sharedTeamID, sharedScheduleID string) string {
	return fmt.Sprintf(`
	resource "firehydrant_escalation_policy" "test_escalation_policy" {
		team_id = "%s"
		name = "test-escalation-policy-%s"
		description = "test-description-%s"
		repetitions = 1
		step_strategy = "dynamic_by_priority"

		step {
			timeout     = "PT1M"
			priorities  = ["HIGH"]

			targets {
				type = "OnCallSchedule"
				id   = "%s"
			}
		}

		step {
			timeout     = "PT2M"
			priorities  = ["LOW"]

			targets {
				type = "OnCallSchedule"
				id   = "%s"
			}
		}

		notification_priority_policies {
			priority = "HIGH"
			repetitions = 2
		}

		notification_priority_policies {
			priority = "LOW"
			repetitions = 1
		}
	}
	`, sharedTeamID, rName, rName, sharedScheduleID, sharedScheduleID)
}

func testAccEscalationPolicyConfig_dynamicPriorityUpdated(rName, sharedTeamID, sharedScheduleID string) string {
	return fmt.Sprintf(`
	resource "firehydrant_escalation_policy" "test_escalation_policy" {
		team_id = "%s"
		name = "test-escalation-policy-updated-%s"
		description = "test-description-updated-%s"
		repetitions = 1
		step_strategy = "dynamic_by_priority"

		step {
			timeout     = "PT1M"
			priorities  = ["HIGH", "MEDIUM"]

			targets {
				type = "OnCallSchedule"
				id   = "%s"
			}
		}

		notification_priority_policies {
			priority = "HIGH"
			repetitions = 3
		}

		notification_priority_policies {
			priority = "MEDIUM"
			repetitions = 1
		}
	}
	`, sharedTeamID, rName, rName, sharedScheduleID)
}

func testAccEscalationPolicyConfig_dynamicWithHandoffSteps(rName, sharedTeamID, sharedScheduleID string) string {
	return fmt.Sprintf(`
	resource "firehydrant_escalation_policy" "test_escalation_policy" {
		team_id = "%s"
		name = "test-escalation-policy-%s"
		description = "test-description-%s"
		repetitions = 1
		step_strategy = "dynamic_by_priority"

		step {
			timeout     = "PT1M"
			priorities  = ["HIGH"]

			targets {
				type = "OnCallSchedule"
				id   = "%s"
			}
		}

		notification_priority_policies {
			priority = "HIGH"
			repetitions = 2
			
			handoff_step {
				target_type = "Team"
				target_id   = "%s"
			}
		}
	}
	`, sharedTeamID, rName, rName, sharedScheduleID, sharedTeamID)
}

func testAccCheckEscalationPolicyResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_escalation_policy" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			// Check if the escalation policy still exists
			_, err := client.Sdk.Signals.GetTeamEscalationPolicy(context.TODO(), stateResource.Primary.Attributes["team_id"], stateResource.Primary.ID)
			if err == nil {
				return fmt.Errorf("Escalation policy %s still exists", stateResource.Primary.ID)
			}
			errStr := err.Error()
			if !strings.Contains(errStr, "404") && !strings.Contains(errStr, "record not found") {
				return fmt.Errorf("Error checking escalation policy %s: %v", stateResource.Primary.ID, err)
			}
		}

		return nil
	}
}
