package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/suite"
)

/** Suite *************************************************************************************************************/
type testOnCallScheduleDataSuite struct {
	suite.Suite
}

func TestOnCallScheduleData(t *testing.T) {
	suite.Run(t, new(testOnCallScheduleDataSuite))
}

func (s *testOnCallScheduleDataSuite) terraform(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "test_on_call_schedule_data_team" {
	name = "test-team-%s"
}

resource "firehydrant_on_call_schedule" "test_on_call_schedule_data_1" {
  name        = "test-on-call-schedule-%s"
  description = "test-description"
  team_id     = firehydrant_team.test_on_call_schedule_data_team.id
  time_zone   = "America/Los_Angeles"
  
  member_ids = ["test-user-id"]
  slack_user_group_id = "test-slack-user-group-id"

  strategy {
	type         = "weekly"
	handoff_time = "10:00:00"
	handoff_day  = "thursday"
  }

  restrictions {
    start_day  = "monday"
    start_time = "14:00:00"
    end_day    = "friday"
    end_time   = "17:00:00"
  }
}

data "firehydrant_on_call_schedule" "test_on_call_schedule_data" {
	id = firehydrant_on_call_schedule.test_on_call_schedule_data_1.id
	team_id = firehydrant_team.test_on_call_schedule_data_team.id
}`, rName, rName)
}

func (s *testOnCallScheduleDataSuite) terraformWithCustomStrategy(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "test_on_call_schedule_data_team" {
	name = "test-team-%s"
}

resource "firehydrant_on_call_schedule" "test_on_call_schedule_data_1" {
  name        = "test-on-call-schedule-%s"
  description = "test-description"
  team_id     = firehydrant_team.test_on_call_schedule_data_team.id
  time_zone   = "America/Los_Angeles"
  start_time  = "2024-01-01T10:00:00-08:00"
  
  member_ids = ["test-user-id"]
  slack_user_group_id = "test-slack-user-group-id"

  strategy {
	type           = "custom"
	shift_duration = "PT8H"
  }
}

data "firehydrant_on_call_schedule" "test_on_call_schedule_data" {
	id = firehydrant_on_call_schedule.test_on_call_schedule_data_1.id
	team_id = firehydrant_team.test_on_call_schedule_data_team.id
}`, rName, rName)
}

func (s *testOnCallScheduleDataSuite) terraformWithoutRestrictions(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "test_on_call_schedule_data_team" {
	name = "test-team-%s"
}

resource "firehydrant_on_call_schedule" "test_on_call_schedule_data_1" {
  name        = "test-on-call-schedule-%s"
  description = "test-description"
  team_id     = firehydrant_team.test_on_call_schedule_data_team.id
  time_zone   = "America/Los_Angeles"
  
  member_ids = ["test-user-id"]
  slack_user_group_id = "test-slack-user-group-id"

  strategy {
	type         = "weekly"
	handoff_time = "10:00:00"
	handoff_day  = "thursday"
  }
}

data "firehydrant_on_call_schedule" "test_on_call_schedule_data" {
	id = firehydrant_on_call_schedule.test_on_call_schedule_data_1.id
	team_id = firehydrant_team.test_on_call_schedule_data_team.id
}`, rName, rName)
}

func (s *testOnCallScheduleDataSuite) testResource(steps ...resource.TestStep) {
	resource.Test(s.T(), resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(s.T()) },
		ProviderFactories: defaultProviderFactories(),
		Steps:             steps,
	})
}

/** Tests *************************************************************************************************************/

func TestAccOnCallScheduleDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccOnCallScheduleDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_on_call_schedule.test_on_call_schedule", "id"),
					resource.TestCheckResourceAttrSet("data.firehydrant_on_call_schedule.test_on_call_schedule", "team_id"),
					resource.TestCheckResourceAttrSet("data.firehydrant_on_call_schedule.test_on_call_schedule", "name"),
					resource.TestCheckResourceAttrSet("data.firehydrant_on_call_schedule.test_on_call_schedule", "description"),
					resource.TestCheckResourceAttrSet("data.firehydrant_on_call_schedule.test_on_call_schedule", "time_zone"),
					resource.TestCheckResourceAttrSet("data.firehydrant_on_call_schedule.test_on_call_schedule", "strategy.0.type"),
					resource.TestCheckResourceAttrSet("data.firehydrant_on_call_schedule.test_on_call_schedule", "strategy.0.handoff_time"),
					resource.TestCheckResourceAttrSet("data.firehydrant_on_call_schedule.test_on_call_schedule", "strategy.0.handoff_day"),
				),
			},
		},
	})
}

func testAccOnCallScheduleDataSourceConfig_basic() string {
	return `
resource "firehydrant_team" "test_team" {
	name = "test-team-acc"
}

resource "firehydrant_on_call_schedule" "test_schedule" {
	name        = "test-on-call-schedule-acc"
	description = "test-description"
	team_id     = firehydrant_team.test_team.id
	time_zone   = "America/New_York"
	
	member_ids = ["test-user-id"]
	slack_user_group_id = "test-slack-user-group-id"

	strategy {
		type         = "weekly"
		handoff_time = "09:00:00"
		handoff_day  = "monday"
	}

	restrictions {
		start_day  = "monday"
		start_time = "09:00:00"
		end_day    = "friday"
		end_time   = "17:00:00"
	}
}

data "firehydrant_on_call_schedule" "test_on_call_schedule" {
	id = firehydrant_on_call_schedule.test_schedule.id
	team_id = firehydrant_team.test_team.id
}`
}

func (s *testOnCallScheduleDataSuite) TestSuccess() {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	s.testResource(
		// Test with restrictions
		resource.TestStep{
			Config: s.terraform(rName),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "id"),
				resource.TestCheckResourceAttrSet("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "team_id"),
				resource.TestCheckResourceAttr("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "name", fmt.Sprintf("test-on-call-schedule-%s", rName)),
				resource.TestCheckResourceAttr("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "description", "test-description"),
				resource.TestCheckResourceAttr("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "time_zone", "America/Los_Angeles"),
				resource.TestCheckResourceAttr("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "slack_user_group_id", "test-slack-user-group-id"),
				resource.TestCheckResourceAttr("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "strategy.0.type", "weekly"),
				resource.TestCheckResourceAttr("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "strategy.0.handoff_time", "10:00:00"),
				resource.TestCheckResourceAttr("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "strategy.0.handoff_day", "thursday"),
				resource.TestCheckResourceAttr("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "restrictions.0.start_day", "monday"),
				resource.TestCheckResourceAttr("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "restrictions.0.start_time", "14:00:00"),
				resource.TestCheckResourceAttr("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "restrictions.0.end_day", "friday"),
				resource.TestCheckResourceAttr("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "restrictions.0.end_time", "17:00:00"),
			),
		},
		// Test custom strategy with shift_duration
		resource.TestStep{
			Config: s.terraformWithCustomStrategy(rName),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "strategy.0.type", "custom"),
				resource.TestCheckResourceAttr("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "strategy.0.shift_duration", "PT8H"),
			),
		},
		// Test without restrictions
		resource.TestStep{
			Config: s.terraformWithoutRestrictions(rName),
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr("data.firehydrant_on_call_schedule.test_on_call_schedule_data", "restrictions.#", "0"),
			),
		},
	)
}
