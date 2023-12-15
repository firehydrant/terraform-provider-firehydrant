package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOnCallScheduleResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckOnCallScheduleResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOnCallScheduleConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_on_call_schedule", "id"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule", "name", fmt.Sprintf("test-on-call-schedule-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule", "strategy.0.type", "weekly"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule", "strategy.0.handoff_time", "10:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule", "strategy.0.handoff_day", "thursday"),
				),
			},
			{
				Config: testAccOnCallScheduleConfig_restrictions(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "id"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "name", fmt.Sprintf("test-on-call-schedule-restrictions-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "strategy.0.type", "weekly"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "strategy.0.handoff_time", "10:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "strategy.0.handoff_day", "thursday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.0.start_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.0.start_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.0.end_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.0.end_time", "14:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.1.start_day", "tuesday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.1.start_time", "12:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.1.end_day", "tuesday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.1.end_time", "18:00:00"),
				),
			},
		},
	})
}

func testAccOnCallScheduleConfig_basic(rName string) string {
	return fmt.Sprintf(`
	resource "firehydrant_team" "team_team" {
		name = "test-team-%s"
	}

	resource "firehydrant_on_call_schedule" "test_on_call_schedule" {
		team_id = firehydrant_team.team_team.id
		name = "test-on-call-schedule-%s"
		description = "test-description-%s"
		time_zone = "America/New_York"

		strategy {
			type         = "weekly"
			handoff_time = "10:00:00"
			handoff_day  = "thursday"
		}
	}
	`, rName, rName, rName)
}

func testAccOnCallScheduleConfig_restrictions(rName string) string {
	return fmt.Sprintf(`
	resource "firehydrant_team" "team_team" {
		name = "test-team-%s"
	}

	resource "firehydrant_on_call_schedule" "test_on_call_schedule_with_restrictions" {
		team_id = firehydrant_team.team_team.id
		name = "test-on-call-schedule-restrictions-%s"
		description = "test-description-%s"
		time_zone = "America/New_York"

		strategy {
			type         = "weekly"
			handoff_time = "10:00:00"
			handoff_day  = "thursday"
		}

		restrictions {
			start_day = "monday"
			start_time = "09:00:00"
			end_day = "monday"
			end_time = "14:00:00"
		}

		restrictions {
			start_day = "tuesday"
			start_time = "12:00:00"
			end_day = "tuesday"
			end_time = "18:00:00"
		}
	}
	`, rName, rName, rName)
}

func testAccCheckOnCallScheduleResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_on_call_schedule" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			// Normally we'd check if err == nil here, because we'd expect a 404 if we try to get a resource
			// that has been deleted. However, the incident role API will still return deleted/archived incident
			// roles instead of returning 404. So, to check for incident roles that are deleted, we have to check
			// for incident roles that have a DiscardedAt timestamp
			_, err := client.OnCallSchedules().Get(context.TODO(), stateResource.Primary.Attributes["team_id"], stateResource.Primary.ID)
			if err != nil && !errors.Is(err, firehydrant.ErrorNotFound) {
				return fmt.Errorf("On-call schedule %s still exists", stateResource.Primary.ID)
			}
		}

		return nil
	}
}
