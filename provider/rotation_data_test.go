package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/suite"
)

/** Suite *************************************************************************************************************/
type testRotationDataSuite struct {
	suite.Suite
}

func TestRotationData(t *testing.T) {
	suite.Run(t, new(testRotationDataSuite))
}

func (s *testRotationDataSuite) terraform(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "test_rotation_data_team" {
	name = "test-team-%s"
}

resource "firehydrant_on_call_schedule" "test_on_call_schedule_data_1" {
  name        = "test-on-call-schedule-%s"
  description = "test-description"
  team_id     = firehydrant_team.test_rotation_data_team.id
  time_zone   = "America/Los_Angeles"
  
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

resource "firehydrant_rotation" "test_rotation_data_1" {
  name = "test-rotation-%s"
	description = "test-description"
	team_id = firehydrant_team.test_rotation_data_team.id
	schedule_id = firehydrant_on_call_schedule.test_on_call_schedule_data_1.id
	time_zone = "America/Los_Angeles"

	slack_user_group_id = "test-slack-user-group-id"
	enable_slack_channel_notifications = true
	prevent_shift_deletion = true

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

data "firehydrant_rotation" "test_rotation_data" {
	id = firehydrant_rotation.test_rotation_data_1.id
	team_id = firehydrant_team.test_rotation_data_team.id
	schedule_id = firehydrant_on_call_schedule.test_on_call_schedule_data_1.id
}`, rName, rName, rName)
}

func (s *testRotationDataSuite) testResource(steps ...resource.TestStep) {
	resource.Test(s.T(), resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(s.T()) },
		ProviderFactories: defaultProviderFactories(),
		Steps:             steps,
	})
}

/** Tests *************************************************************************************************************/

func (s *testRotationDataSuite) TestSuccess() {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	s.testResource(resource.TestStep{
		Config: s.terraform(rName),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttrSet("data.firehydrant_rotation.test_rotation_data", "id"),
			resource.TestCheckResourceAttrSet("data.firehydrant_rotation.test_rotation_data", "team_id"),
			resource.TestCheckResourceAttrSet("data.firehydrant_rotation.test_rotation_data", "schedule_id"),
			resource.TestCheckResourceAttr("data.firehydrant_rotation.test_rotation_data", "name", fmt.Sprintf("test-rotation-%s", rName)),
			resource.TestCheckResourceAttr("data.firehydrant_rotation.test_rotation_data", "description", "test-description"),
			resource.TestCheckResourceAttr("data.firehydrant_rotation.test_rotation_data", "time_zone", "America/Los_Angeles"),
			resource.TestCheckResourceAttr("data.firehydrant_rotation.test_rotation_data", "slack_user_group_id", "test-slack-user-group-id"),
		),
	})
}
