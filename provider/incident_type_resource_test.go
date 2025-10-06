package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIncidentTypeResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
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
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
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
						"firehydrant_incident_type.test_incident_type", "summary", fmt.Sprintf("test-summary-%s", rName)),
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
						"firehydrant_incident_type.test_incident_type", "summary", fmt.Sprintf("test-summary-%s", rNameUpdated)),
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
						"firehydrant_incident_type.test_incident_type", "summary", fmt.Sprintf("test-summary-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_type.test_incident_type", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
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
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "template.0.runbook_ids.0", "88f9f172-cc07-477e-9a80-b1ae7669ec3d"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "template.0.runbook_ids.1", "39de1363-4ae3-4aa3-913b-d63312c76afd"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "template.0.team_ids.0", "cd2b6d18-e616-4990-9065-ec0cab037680"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "template.0.team_ids.1", "2af88e28-0205-4d39-b304-133bd2e19358"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "template.0.impacts.0.impact_id", "500d9e2e-ea7c-4834-a81f-e336de24dbb1"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_incident_type.test_incident_type", "template.0.impacts.0.condition_id", "99762c0c-1ee0-44a0-a3a7-d1316dd902ca"),
				),
			},
		},
	})
}

func TestAccIncidentTypeResourceImport_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
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
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
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

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
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

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
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
		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
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
}`, rName, rName)
}

func testAccIncidentTypeResourceConfig_update(rName string) string {
	return fmt.Sprintf(`
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
		team_ids = [ "cd2b6d18-e616-4990-9065-ec0cab037680", "2af88e28-0205-4d39-b304-133bd2e19358" ]
		impacts = [
		  {
		    impact_id = "500d9e2e-ea7c-4834-a81f-e336de24dbb1"
				condition_id = "99762c0c-1ee0-44a0-a3a7-d1316dd902ca"
			},
		  {
			  condition_id = "99762c0c-1ee0-44a0-a3a7-d1316dd902ca"
				impact_id = "8c6731c8-a49a-415e-91c9-61378d526c58"
			},
		]
	}
}`, rName, rName)
}
