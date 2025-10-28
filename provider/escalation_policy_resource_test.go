package provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccEscalationPolicyResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckEscalationPolicyResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEscalationPolicyConfig_basic(rName),
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
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckEscalationPolicyResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEscalationPolicyConfig_dynamicPriority(rName),
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
				Config: testAccEscalationPolicyConfig_dynamicPriorityUpdated(rName),
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
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckEscalationPolicyResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEscalationPolicyConfig_dynamicWithHandoffSteps(rName),
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

func testAccEscalationPolicyConfig_basic(rName string) string {
	return fmt.Sprintf(`
	resource "firehydrant_team" "test-team" {
		name = "test-team-%s"
	}

	resource "firehydrant_on_call_schedule" "test_on_call_schedule" {
		team_id = firehydrant_team.test-team.id
		name = "test-on-call-schedule-restrictions-%s"
		time_zone = "America/New_York"
		slack_user_group_id = "test-group-1"

		strategy {
			type         = "weekly"
			handoff_time = "10:00:00"
			handoff_day  = "thursday"
		}
	}

	resource "firehydrant_escalation_policy" "test_escalation_policy" {
		team_id = firehydrant_team.test-team.id
		name = "test-escalation-policy-%s"
		description = "test-description-%s"
		repetitions = 1
		step_strategy = "static"
		
		step {
			timeout     = "PT1M"

			targets {
				type = "OnCallSchedule"
				id   = firehydrant_on_call_schedule.test_on_call_schedule.id
			}
		}

		handoff_step {
			target_type = "Team"
			target_id   = firehydrant_team.test-team.id
		}
	}
	`, rName, rName, rName, rName)
}

func testAccEscalationPolicyConfig_dynamicPriority(rName string) string {
	return fmt.Sprintf(`
	resource "firehydrant_team" "test-team" {
		name = "test-team-%s"
	}

	resource "firehydrant_on_call_schedule" "test_on_call_schedule" {
		team_id = firehydrant_team.test-team.id
		name = "test-on-call-schedule-restrictions-%s"
		time_zone = "America/New_York"
		slack_user_group_id = "test-group-1"

		strategy {
			type         = "weekly"
			handoff_time = "10:00:00"
			handoff_day  = "thursday"
		}
	}

	resource "firehydrant_escalation_policy" "test_escalation_policy" {
		team_id = firehydrant_team.test-team.id
		name = "test-escalation-policy-%s"
		description = "test-description-%s"
		repetitions = 1
		step_strategy = "dynamic_by_priority"

		step {
			timeout     = "PT1M"
			priorities  = ["HIGH"]

			targets {
				type = "OnCallSchedule"
				id   = firehydrant_on_call_schedule.test_on_call_schedule.id
			}
		}

		step {
			timeout     = "PT2M"
			priorities  = ["LOW"]

			targets {
				type = "OnCallSchedule"
				id   = firehydrant_on_call_schedule.test_on_call_schedule.id
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
	`, rName, rName, rName, rName)
}

func testAccEscalationPolicyConfig_dynamicPriorityUpdated(rName string) string {
	return fmt.Sprintf(`
	resource "firehydrant_team" "test-team" {
		name = "test-team-%s"
	}

	resource "firehydrant_on_call_schedule" "test_on_call_schedule" {
		team_id = firehydrant_team.test-team.id
		name = "test-on-call-schedule-restrictions-%s"
		time_zone = "America/New_York"
		slack_user_group_id = "test-group-1"

		strategy {
			type         = "weekly"
			handoff_time = "10:00:00"
			handoff_day  = "thursday"
		}
	}

	resource "firehydrant_escalation_policy" "test_escalation_policy" {
		team_id = firehydrant_team.test-team.id
		name = "test-escalation-policy-updated-%s"
		description = "test-description-updated-%s"
		repetitions = 1
		step_strategy = "dynamic_by_priority"

		step {
			timeout     = "PT1M"
			priorities  = ["HIGH", "MEDIUM"]

			targets {
				type = "OnCallSchedule"
				id   = firehydrant_on_call_schedule.test_on_call_schedule.id
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
	`, rName, rName, rName, rName)
}

func testAccEscalationPolicyConfig_dynamicWithHandoffSteps(rName string) string {
	return fmt.Sprintf(`
	resource "firehydrant_team" "test-team" {
		name = "test-team-%s"
	}

	resource "firehydrant_on_call_schedule" "test_on_call_schedule" {
		team_id = firehydrant_team.test-team.id
		name = "test-on-call-schedule-restrictions-%s"
		time_zone = "America/New_York"
		slack_user_group_id = "test-group-1"

		strategy {
			type         = "weekly"
			handoff_time = "10:00:00"
			handoff_day  = "thursday"
		}
	}

	resource "firehydrant_escalation_policy" "test_escalation_policy" {
		team_id = firehydrant_team.test-team.id
		name = "test-escalation-policy-%s"
		description = "test-description-%s"
		repetitions = 1
		step_strategy = "dynamic_by_priority"

		step {
			timeout     = "PT1M"
			priorities  = ["HIGH"]

			targets {
				type = "OnCallSchedule"
				id   = firehydrant_on_call_schedule.test_on_call_schedule.id
			}
		}

		notification_priority_policies {
			priority = "HIGH"
			repetitions = 2
			
			handoff_step {
				target_type = "Team"
				target_id   = firehydrant_team.test-team.id
			}
		}
	}
	`, rName, rName, rName, rName)
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
