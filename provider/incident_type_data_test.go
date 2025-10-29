package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIncidentTypeDataSource_basic(t *testing.T) {
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

// TODO: add helpers for some static attributes (runbook_ids, severities, priorities, etc) so we aren't making a new
// everything for testing every resource.  In the meantime, I'm hardcoding some shit from the acceptance test instance.

func testAccIncidentTypeDataSourceConfig_allAttributes(rName string) string {
	return fmt.Sprintf(`

resource "firehydrant_team" "test_team_1" {
  name = "test-team-1-%s"
}

resource "firehydrant_team" "test_team_2" {
  name = "test-team-2-%s"
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
		runbook_ids = [ "88f9f172-cc07-477e-9a80-b1ae7669ec3d", "39de1363-4ae3-4aa3-913b-d63312c76afd" ]
		team_ids = [ firehydrant_team.test_team_1.id, firehydrant_team.test_team_2.id ]

		impacts {
          impact_id = "8c6731c8-a49a-415e-91c9-61378d526c58"
            condition_id = "99762c0c-1ee0-44a0-a3a7-d1316dd902ca"
        }
        
        impacts {
          impact_id = "500d9e2e-ea7c-4834-a81f-e336de24dbb1"
            condition_id = "99762c0c-1ee0-44a0-a3a7-d1316dd902ca"
    }
	}
}

data "firehydrant_incident_type" "test_incident_type" {
  id = firehydrant_incident_type.test_incident_type.id
}`, rName, rName, rName, rName)
}
