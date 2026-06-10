package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIncidentTypeResource_basic(t *testing.T) {
	t.Parallel()
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckIncidentTypeResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIncidentTypeResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIncidentTypeResourceExistsWithAttributes_basic("firehydrant_incident_type.test_incident_type"),
					resource.TestCheckResourceAttrSet("firehydrant_incident_type.test_incident_type", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "name", fmt.Sprintf("test-incident-type-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "description", fmt.Sprintf("test-description-%s", rName)),
				),
			},
		},
	})
}

func TestAccIncidentTypeResource_update(t *testing.T) {
	t.Parallel()
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckIncidentTypeResourceDestroy(),
			testAccCheckTeamResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccIncidentTypeResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIncidentTypeResourceExistsWithAttributes_basic("firehydrant_incident_type.test_incident_type"),
					resource.TestCheckResourceAttrSet("firehydrant_incident_type.test_incident_type", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "name", fmt.Sprintf("test-incident-type-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "description", fmt.Sprintf("test-description-%s", rName)),
				),
			},
			{
				Config: testAccIncidentTypeResourceConfig_basic(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIncidentTypeResourceExistsWithAttributes_basic("firehydrant_incident_type.test_incident_type"),
					resource.TestCheckResourceAttrSet("firehydrant_incident_type.test_incident_type", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "name", fmt.Sprintf("test-incident-type-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
				),
			},
			{
				Config: testAccIncidentTypeResourceConfig_update(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIncidentTypeResourceExistsWithAttributes_update("firehydrant_incident_type.test_incident_type"),
					resource.TestCheckResourceAttrSet("firehydrant_incident_type.test_incident_type", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "name", fmt.Sprintf("test-incident-type-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "template.0.description", "test-template-description"),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "template.0.customer_impact_summary", "test-summary"),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "template.0.severity_slug", "SEV1"),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "template.0.priority_slug", "TESTPRIORITY"),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "template.0.private_incident", "false"),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "template.0.tags.0", "foo"),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "template.0.tags.1", "bar"),
					resource.TestCheckResourceAttrSet(
						"firehydrant_incident_type.test_incident_type", "template.0.team_ids.0"),
					resource.TestCheckResourceAttrSet(
						"firehydrant_incident_type.test_incident_type", "template.0.team_ids.1"),
					resource.TestCheckResourceAttrSet(
						"firehydrant_incident_type.test_incident_type", "template.0.runbook_ids.0"),
					resource.TestCheckResourceAttrSet(
						"firehydrant_incident_type.test_incident_type", "template.0.runbook_ids.1"),
				),
			},
		},
	})
}

func TestAccIncidentTypeResourceImport_basic(t *testing.T) {
	t.Parallel()
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccIncidentTypeResourceConfig_basic(rName),
			},
			{
				ResourceName:      "firehydrant_incident_type.test_incident_type",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIncidentTypeResourceImport_allAttributes(t *testing.T) {
	t.Parallel()
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccIncidentTypeResourceConfig_update(rName),
			},
			{
				ResourceName:      "firehydrant_incident_type.test_incident_type",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIncidentTypeResourceExistsWithAttributes_basic(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		incidentTypeResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if incidentTypeResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		incidentTypeResponse, err := client.Sdk.IncidentSettings.GetIncidentType(context.TODO(), incidentTypeResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := incidentTypeResource.Primary.Attributes["name"], incidentTypeResponse.Name
		if expected != *got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, *got)
		}

		expected, got = incidentTypeResource.Primary.Attributes["description"], incidentTypeResponse.Description
		if expected != *got {
			return fmt.Errorf("Unexpected summary. Expected: %s, got: %s", expected, *got)
		}

		return nil
	}
}

func testAccCheckIncidentTypeResourceExistsWithAttributes_update(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		incidentTypeResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if incidentTypeResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		incidentTypeResponse, err := client.Sdk.IncidentSettings.GetIncidentType(context.TODO(), incidentTypeResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := incidentTypeResource.Primary.Attributes["name"], incidentTypeResponse.Name
		if expected != *got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, *got)
		}

		expected, got = incidentTypeResource.Primary.Attributes["description"], incidentTypeResponse.Description
		if expected != *got {
			return fmt.Errorf("Unexpected description. Expected: %s, got: %s", expected, *got)
		}

		return nil
	}
}

func testAccCheckIncidentTypeResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_incident_type" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			_, err := client.Sdk.IncidentSettings.GetIncidentType(context.TODO(), stateResource.Primary.ID)
			if err == nil {
				return fmt.Errorf("Incident Type %s still exists", stateResource.Primary.ID)
			}
		}

		return nil
	}
}

func testAccIncidentTypeResourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_incident_type" "test_incident_type" {
  name        = "test-incident-type-%s"
  description = "test-description-%s"

	template {}
}`, rName, rName)
}

func testAccIncidentTypeResourceConfig_update(rName string) string {
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
}`, rName, rName, rName, rName, rName, rName, rName, rName)
}
