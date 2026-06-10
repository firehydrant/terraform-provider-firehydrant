package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIncidentTypeDataSource_basic(t *testing.T) {
	t.Parallel()
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckIncidentTypeResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIncidentTypeDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_incident_type.test_incident_type", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "name", fmt.Sprintf("test-incident-type-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "description", fmt.Sprintf("test-description-%s", rName)),
				),
			},
		},
	})
}

func TestAccIncidentTypeDataSource_allAttributes(t *testing.T) {
	t.Parallel()
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckIncidentTypeResourceDestroy(),
			testAccCheckTeamResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccIncidentTypeDataSourceConfig_allAttributes(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_incident_type.test_incident_type", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "name", fmt.Sprintf("test-incident-type-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "description", fmt.Sprintf("test-description-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "template.0.description", "test-template-description"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "template.0.customer_impact_summary", "test-summary"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "template.0.severity_slug", "SEV1"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "template.0.priority_slug", "TESTPRIORITY"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "template.0.private_incident", "false"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "template.0.tags.0", "foo"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "template.0.tags.1", "bar"),
					resource.TestCheckResourceAttrSet(
						"data.firehydrant_incident_type.test_incident_type", "template.0.team_ids.0"),
					resource.TestCheckResourceAttrSet(
						"data.firehydrant_incident_type.test_incident_type", "template.0.team_ids.1"),
					resource.TestCheckResourceAttrSet(
						"data.firehydrant_incident_type.test_incident_type", "template.0.runbook_ids.0"),
					resource.TestCheckResourceAttrSet(
						"data.firehydrant_incident_type.test_incident_type", "template.0.runbook_ids.1"),
				),
			},
		},
	})
}

func testAccIncidentTypeDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_incident_type" "test_incident_type" {
  name        = "test-incident-type-%s"
  description = "test-description-%s"

	template {}
}

data "firehydrant_incident_type" "test_incident_type" {
  id = firehydrant_incident_type.test_incident_type.id
}`, rName, rName)
}

// The runbooks and services these templates reference are created in-config so
// the tests don't depend on org-specific fixture IDs. severity_slug/priority_slug
// (SEV1/TESTPRIORITY) and the impact condition_id are still pre-existing platform
// fixtures — see TESTS.md.

func testAccIncidentTypeDataSourceConfig_allAttributes(rName string) string {
	return fmt.Sprintf(`

resource "firehydrant_team" "test_team_1" {
  name = "test-team-1-%s"
}

resource "firehydrant_team" "test_team_2" {
  name = "test-team-2-%s"
}




data "firehydrant_runbook_action" "create_incident_channel" {
  slug             = "create_incident_channel"
  integration_slug = "slack"
}

resource "firehydrant_runbook" "test_runbook_1" {
  name = "test-runbook-1-%s"

  steps {
    name      = "Create Incident Channel"
    action_id = data.firehydrant_runbook_action.create_incident_channel.id

    config = jsonencode({
      channel_name_format = "-inc-{{ number }}"
    })
  }
}

resource "firehydrant_runbook" "test_runbook_2" {
  name = "test-runbook-2-%s"

  steps {
    name      = "Create Incident Channel"
    action_id = data.firehydrant_runbook_action.create_incident_channel.id

    config = jsonencode({
      channel_name_format = "-inc2-{{ number }}"
    })
  }
}

resource "firehydrant_service" "test_service_1" {
  name = "test-service-1-%s"
}

resource "firehydrant_service" "test_service_2" {
  name = "test-service-2-%s"
}

resource "firehydrant_incident_type" "test_incident_type" {
  name        = "test-incident-type-%s"
  description = "test-description-%s"
	template {
	  description = "test-template-description"
		customer_impact_summary = "test-summary"
		severity_slug = "SEV1"
		priority_slug = "TESTPRIORITY"
		private_incident = false

		tags = [ "foo", "bar" ]
		runbook_ids = [ firehydrant_runbook.test_runbook_1.id, firehydrant_runbook.test_runbook_2.id ]
		team_ids = [ firehydrant_team.test_team_1.id, firehydrant_team.test_team_2.id ]

		impacts {
			impact_id    = firehydrant_service.test_service_1.id
			condition_id = "99762c0c-1ee0-44a0-a3a7-d1316dd902ca"
		}

		impacts {
			impact_id    = firehydrant_service.test_service_2.id
			condition_id = "99762c0c-1ee0-44a0-a3a7-d1316dd902ca"
		}
	}
}

data "firehydrant_incident_type" "test_incident_type" {
  id = firehydrant_incident_type.test_incident_type.id
}`, rName, rName, rName, rName, rName, rName, rName, rName)
}
