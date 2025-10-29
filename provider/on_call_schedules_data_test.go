package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/suite"
)

/** Suite *************************************************************************************************************/
type testOnCallSchedulesDataSuite struct {
	suite.Suite
}

func TestOnCallSchedulesData(t *testing.T) {
	suite.Run(t, new(testOnCallSchedulesDataSuite))
}

func (s *testOnCallSchedulesDataSuite) terraform(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "team" {
	name = "test-team-%s"
}

resource "firehydrant_on_call_schedule" "schedule_1" {
  name        = "test-on-call-schedule-%s"
  description = "test-description"
  team_id     = firehydrant_team.team.id
  time_zone   = "America/Los_Angeles"
  color       = "#ff0000" 
  
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

data "firehydrant_on_call_schedules" "schedules" {
	team_id = firehydrant_team.team.id
	query = "test-on-call-schedule-%s"
	depends_on = [firehydrant_on_call_schedule.schedule_1]
}`, rName, rName, rName)
}

func (s *testOnCallSchedulesDataSuite) testResource(steps ...resource.TestStep) {
	resource.Test(s.T(), resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(s.T()) },
		ProviderFactories: sharedProviderFactories(),
		Steps:             steps,
	})
}

/** Tests *************************************************************************************************************/

func (s *testOnCallSchedulesDataSuite) TestMatches() {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	s.testResource(resource.TestStep{
		Config: s.terraform(rName),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttrSet("firehydrant_team.team", "id"),
			resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.schedule_1", "id"),

			resource.TestCheckResourceAttrSet("data.firehydrant_on_call_schedules.schedules", "on_call_schedules.#"),
			resource.TestCheckResourceAttrSet("data.firehydrant_on_call_schedules.schedules", "on_call_schedules.0.id"),
			resource.TestCheckResourceAttrSet("data.firehydrant_on_call_schedules.schedules", "on_call_schedules.0.team_id"),
			resource.TestCheckResourceAttr("data.firehydrant_on_call_schedules.schedules", "on_call_schedules.0.name", fmt.Sprintf("test-on-call-schedule-%s", rName)),
			resource.TestCheckResourceAttr("data.firehydrant_on_call_schedules.schedules", "on_call_schedules.0.description", "test-description"),
			resource.TestCheckResourceAttr("data.firehydrant_on_call_schedules.schedules", "on_call_schedules.0.time_zone", "America/Los_Angeles"),
			resource.TestCheckResourceAttr("data.firehydrant_on_call_schedules.schedules", "on_call_schedules.0.slack_user_group_id", "test-slack-user-group-id"),
		),
	})
}
