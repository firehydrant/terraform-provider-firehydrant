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

func TestAccIncidentRoleResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckIncidentRoleResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIncidentRoleResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIncidentRoleResourceExistsWithAttributes_basic("firehydrant_incident_role.test_incident_role"),
					resource.TestCheckResourceAttrSet("firehydrant_incident_role.test_incident_role", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_role.test_incident_role", "name", fmt.Sprintf("test-incident-role-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_role.test_incident_role", "summary", fmt.Sprintf("test-summary-%s", rName)),
				),
			},
		},
	})
}

func TestAccIncidentRoleResource_update(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckIncidentRoleResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIncidentRoleResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIncidentRoleResourceExistsWithAttributes_basic("firehydrant_incident_role.test_incident_role"),
					resource.TestCheckResourceAttrSet("firehydrant_incident_role.test_incident_role", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_role.test_incident_role", "name", fmt.Sprintf("test-incident-role-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_role.test_incident_role", "summary", fmt.Sprintf("test-summary-%s", rName)),
				),
			},
			{
				Config: testAccIncidentRoleResourceConfig_update(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIncidentRoleResourceExistsWithAttributes_update("firehydrant_incident_role.test_incident_role"),
					resource.TestCheckResourceAttrSet("firehydrant_incident_role.test_incident_role", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_role.test_incident_role", "name", fmt.Sprintf("test-incident-role-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_role.test_incident_role", "summary", fmt.Sprintf("test-summary-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_role.test_incident_role", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
				),
			},
			{
				Config: testAccIncidentRoleResourceConfig_basic(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIncidentRoleResourceExistsWithAttributes_basic("firehydrant_incident_role.test_incident_role"),
					resource.TestCheckResourceAttrSet("firehydrant_incident_role.test_incident_role", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_role.test_incident_role", "name", fmt.Sprintf("test-incident-role-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_incident_role.test_incident_role", "summary", fmt.Sprintf("test-summary-%s", rNameUpdated)),
				),
			},
		},
	})
}

func TestAccIncidentRoleResourceImport_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccIncidentRoleResourceConfig_basic(rName),
			},
			{
				ResourceName:      "firehydrant_incident_role.test_incident_role",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIncidentRoleResourceImport_allAttributes(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccIncidentRoleResourceConfig_update(rName),
			},
			{
				ResourceName:      "firehydrant_incident_role.test_incident_role",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIncidentRoleResourceExistsWithAttributes_basic(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		incidentRoleResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if incidentRoleResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		incidentRole, err := client.Sdk.IncidentSettings.GetIncidentRole(context.TODO(), incidentRoleResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := incidentRoleResource.Primary.Attributes["name"], *incidentRole.GetName()
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = incidentRoleResource.Primary.Attributes["summary"], *incidentRole.GetSummary()
		if expected != got {
			return fmt.Errorf("Unexpected summary. Expected: %s, got: %s", expected, got)
		}

		if incidentRole.GetDescription() != nil && *incidentRole.GetDescription() != "" {
			return fmt.Errorf("Unexpected description. Expected no description, got: %s", *incidentRole.GetDescription())
		}

		return nil
	}
}

func testAccCheckIncidentRoleResourceExistsWithAttributes_update(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		incidentRoleResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if incidentRoleResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		incidentRole, err := client.Sdk.IncidentSettings.GetIncidentRole(context.TODO(), incidentRoleResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := incidentRoleResource.Primary.Attributes["name"], *incidentRole.GetName()
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = incidentRoleResource.Primary.Attributes["summary"], *incidentRole.GetSummary()
		if expected != got {
			return fmt.Errorf("Unexpected summary. Expected: %s, got: %s", expected, got)
		}

		expected, got = incidentRoleResource.Primary.Attributes["description"], *incidentRole.GetDescription()
		if expected != got {
			return fmt.Errorf("Unexpected description. Expected: %s, got: %s", expected, got)
		}

		return nil
	}
}

func testAccCheckIncidentRoleResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_incident_role" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			// Normally we'd check if err == nil here, because we'd expect a 404 if we try to get a resource
			// that has been deleted. However, the incident role API will still return deleted/archived incident
			// roles instead of returning 404. So, to check for incident roles that are deleted, we have to check
			// for incident roles that have a DiscardedAt timestamp
			incidentRole, _ := client.Sdk.IncidentSettings.GetIncidentRole(context.TODO(), stateResource.Primary.ID)
			if incidentRole.GetDiscardedAt() == nil || incidentRole.GetDiscardedAt().IsZero() {
				return fmt.Errorf("Incident role %s still exists", stateResource.Primary.ID)
			}
		}

		return nil
	}
}

func testAccIncidentRoleResourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_incident_role" "test_incident_role" {
  name    = "test-incident-role-%s"
  summary = "test-summary-%s"
}`, rName, rName)
}

func testAccIncidentRoleResourceConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_incident_role" "test_incident_role" {
  name        = "test-incident-role-%s"
  summary     = "test-summary-%s"
  description = "test-description-%s"
}`, rName, rName, rName)
}
