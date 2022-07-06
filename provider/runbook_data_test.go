package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRunbookDataSource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRunbookDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_runbook.test_runbook", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_runbook.test_runbook", "name", fmt.Sprintf("test-runbook-%s", rName)),
				),
			},
		},
	})
}

func TestAccRunbookDataSource_allAttributes(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRunbookDataSourceConfig_allAttributes(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_runbook.test_runbook", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_runbook.test_runbook", "name", fmt.Sprintf("test-runbook-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_runbook.test_runbook", "description", fmt.Sprintf("test-description-%s", rName)),
					resource.TestCheckResourceAttrSet("data.firehydrant_runbook.test_runbook", "owner_id"),
				),
			},
		},
	})
}

func testAccRunbookDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "firehydrant_runbook_action" "create_incident_channel" {
  slug             = "create_incident_channel"
  integration_slug = "slack"
  type             = "incident"
}

resource "firehydrant_runbook" "test_runbook" {
  name        = "test-runbook-%s"
  type        = "incident"

  steps {
    name      = "Create Incident Channel"
    action_id = data.firehydrant_runbook_action.create_incident_channel.id
    config = {
      channel_name_format = "-inc-{{ number }}"
    }
  }
}

data "firehydrant_runbook" "test_runbook" {
  id = firehydrant_runbook.test_runbook.id
}`, rName)
}

func testAccRunbookDataSourceConfig_allAttributes(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "test_team1" {
	name = "test-team1-%s"
}

data "firehydrant_runbook_action" "create_incident_channel" {
  slug             = "create_incident_channel"
  integration_slug = "slack"
  type             = "incident"
}

resource "firehydrant_runbook" "test_runbook" {
  name        = "test-runbook-%s"
  type        = "incident"
  description = "test-description-%s"
	owner_id = firehydrant_team.test_team1.id

  steps {
    name      = "Create Incident Channel"
    action_id = data.firehydrant_runbook_action.create_incident_channel.id
    config = {
      channel_name_format = "-inc-{{ number }}"
    }
  }
}

data "firehydrant_runbook" "test_runbook" {
  id = firehydrant_runbook.test_runbook.id
}`, rName, rName, rName)
}
