package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccServiceResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckServiceResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceResourceExistsWithAttributes_basic("firehydrant_service.test_service"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "alert_on_add", "false"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "auto_add_responding_team", "false"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "service_tier", "5"),
					resource.TestCheckResourceAttr("firehydrant_service.test_service", "team_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccServiceResource_update(t *testing.T) {
	sharedTeamID := getSharedTeamID(t)
	sharedTeamID2 := getSharedTeamID2(t)
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckServiceResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceResourceExistsWithAttributes_basic("firehydrant_service.test_service"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "alert_on_add", "false"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "auto_add_responding_team", "false"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "service_tier", "5"),
					resource.TestCheckResourceAttr("firehydrant_service.test_service", "team_ids.#", "0"),
				),
			},
			{
				Config: testAccServiceResourceConfig_update(rNameUpdated, sharedTeamID, sharedTeamID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceResourceExistsWithAttributes_update("firehydrant_service.test_service"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "alert_on_add", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "auto_add_responding_team", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "labels.test1", fmt.Sprintf("test-label1-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "links.#", "2"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "owner_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "service_tier", "1"),
					resource.TestCheckResourceAttr("firehydrant_service.test_service", "team_ids.#", "2"),
				),
			},
			{
				Config: testAccServiceResourceConfig_basic(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceResourceExistsWithAttributes_basic("firehydrant_service.test_service"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "alert_on_add", "false"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "auto_add_responding_team", "false"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "service_tier", "5"),
					resource.TestCheckResourceAttr("firehydrant_service.test_service", "team_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccServiceResource_updateLabels(t *testing.T) {
	sharedTeamID := getSharedTeamID(t)
	sharedTeamID2 := getSharedTeamID2(t)
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckServiceResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfig_update(rName, sharedTeamID, sharedTeamID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceResourceExistsWithAttributes_update("firehydrant_service.test_service"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "alert_on_add", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "auto_add_responding_team", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "description", fmt.Sprintf("test-description-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "labels.test1", fmt.Sprintf("test-label1-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "links.#", "2"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "owner_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "service_tier", "1"),
					resource.TestCheckResourceAttr("firehydrant_service.test_service", "team_ids.#", "2"),
				),
			},
			{
				Config: testAccServiceResourceConfig_updateChangeLabels(rNameUpdated, sharedTeamID, sharedTeamID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceResourceExistsWithAttributes_update("firehydrant_service.test_service"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "alert_on_add", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "auto_add_responding_team", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "labels.test1", fmt.Sprintf("test-label1-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "labels.test2", fmt.Sprintf("test-label2-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "links.#", "2"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "owner_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "service_tier", "1"),
					resource.TestCheckResourceAttr("firehydrant_service.test_service", "team_ids.#", "2"),
				),
			},
			{
				Config: testAccServiceResourceConfig_basic(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceResourceExistsWithAttributes_basic("firehydrant_service.test_service"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "alert_on_add", "false"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "auto_add_responding_team", "false"),
					// Make sure the labels are not set
					resource.TestCheckNoResourceAttr(
						"firehydrant_service.test_service", "labels.test1"),
					resource.TestCheckNoResourceAttr(
						"firehydrant_service.test_service", "labels.test2"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "service_tier", "5"),
					resource.TestCheckResourceAttr("firehydrant_service.test_service", "team_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccServiceResource_updateOwnerIDAndTeamIDs(t *testing.T) {
	sharedTeamID := getSharedTeamID(t)
	sharedTeamID2 := getSharedTeamID2(t)
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckServiceResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfig_update(rName, sharedTeamID, sharedTeamID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceResourceExistsWithAttributes_update("firehydrant_service.test_service"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "alert_on_add", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "auto_add_responding_team", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "description", fmt.Sprintf("test-description-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "labels.test1", fmt.Sprintf("test-label1-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "links.#", "2"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "owner_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "service_tier", "1"),
					resource.TestCheckResourceAttr("firehydrant_service.test_service", "team_ids.#", "2"),
				),
			},
			{
				Config: testAccServiceResourceConfig_update(rNameUpdated, sharedTeamID, sharedTeamID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceResourceExistsWithAttributes_update("firehydrant_service.test_service"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "alert_on_add", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "auto_add_responding_team", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "labels.test1", fmt.Sprintf("test-label1-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "links.#", "2"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "owner_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "service_tier", "1"),
					resource.TestCheckResourceAttr("firehydrant_service.test_service", "team_ids.#", "2"),
				),
			},
			{
				Config: testAccServiceResourceConfig_updateChangeOwnerIDAndTeamIDs(rNameUpdated, sharedTeamID, sharedTeamID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceResourceExistsWithAttributes_update("firehydrant_service.test_service"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "alert_on_add", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "auto_add_responding_team", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "labels.test1", fmt.Sprintf("test-label1-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "links.#", "2"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "owner_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "service_tier", "1"),
					resource.TestCheckResourceAttr("firehydrant_service.test_service", "team_ids.#", "2"),
				),
			},
			{
				Config: testAccServiceResourceConfig_basic(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceResourceExistsWithAttributes_basic("firehydrant_service.test_service"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "alert_on_add", "false"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "auto_add_responding_team", "false"),
					// Make sure owner_id is not set
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "owner_id", ""),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "service_tier", "5"),
					resource.TestCheckResourceAttr("firehydrant_service.test_service", "team_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccServiceResourceImport_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfig_basic(rName),
			},
			{
				ResourceName:      "firehydrant_service.test_service",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceResourceImport_allAttributes(t *testing.T) {
	sharedTeamID := getSharedTeamID(t)
	sharedTeamID2 := getSharedTeamID2(t)
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfig_update(rName, sharedTeamID, sharedTeamID2),
			},
			{
				ResourceName:      "firehydrant_service.test_service",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckServiceResourceExistsWithAttributes_basic(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		serviceResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if serviceResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		serviceResponse, err := client.Sdk.CatalogEntries.GetService(context.TODO(), serviceResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := serviceResource.Primary.Attributes["name"], *serviceResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = serviceResource.Primary.Attributes["alert_on_add"], fmt.Sprintf("%t", *serviceResponse.AlertOnAdd)
		if expected != got {
			return fmt.Errorf("Unexpected alert_on_add. Expected: %s, got: %s", expected, got)
		}

		expected, got = serviceResource.Primary.Attributes["auto_add_responding_team"], fmt.Sprintf("%t", *serviceResponse.AutoAddRespondingTeam)
		if expected != got {
			return fmt.Errorf("Unexpected auto_add_responding_team. Expected: %s, got: %s", expected, got)
		}

		if serviceResponse.Description != nil && *serviceResponse.Description != "" {
			return fmt.Errorf("Unexpected description. Expected no description, got: %s", *serviceResponse.Description)
		}

		if serviceResponse.Labels != nil {
			return fmt.Errorf("Unexpected labels. Expected no labels")
		}

		if len(serviceResponse.Links) != 0 {
			return fmt.Errorf("Unexpected number of links. Expected no links, got: %v", len(serviceResponse.Links))
		}

		if serviceResponse.Owner != nil {
			return fmt.Errorf("Unexpected owner. Expected no owner ID, got: %s", *serviceResponse.Owner.ID)
		}

		expected, got = serviceResource.Primary.Attributes["service_tier"], fmt.Sprintf("%d", *serviceResponse.ServiceTier)
		if expected != got {
			return fmt.Errorf("Unexpected service_tier. Expected: %s, got: %s", expected, got)
		}

		// TODO: Check the team ids
		if len(serviceResponse.Teams) != 0 {
			return fmt.Errorf("Unexpected number of teams. Expected: 0, got: %v", len(serviceResponse.Teams))
		}

		return nil
	}
}

func testAccCheckServiceResourceExistsWithAttributes_update(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		serviceResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if serviceResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		serviceResponse, err := client.Sdk.CatalogEntries.GetService(context.TODO(), serviceResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := serviceResource.Primary.Attributes["name"], *serviceResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = serviceResource.Primary.Attributes["alert_on_add"], fmt.Sprintf("%t", *serviceResponse.AlertOnAdd)
		if expected != got {
			return fmt.Errorf("Unexpected alert_on_add. Expected: %s, got: %s", expected, got)
		}

		expected, got = serviceResource.Primary.Attributes["auto_add_responding_team"], fmt.Sprintf("%t", *serviceResponse.AutoAddRespondingTeam)
		if expected != got {
			return fmt.Errorf("Unexpected auto_add_responding_team. Expected: %s, got: %s", expected, got)
		}

		expected, got = serviceResource.Primary.Attributes["description"], *serviceResponse.Description
		if expected != got {
			return fmt.Errorf("Unexpected description. Expected: %s, got: %s", expected, got)
		}

		if serviceResponse.Labels == nil {
			return fmt.Errorf("Unexpected labels. Expected labels to be set")
		}

		// TODO: check link attributes
		if len(serviceResponse.Links) == 0 {
			return fmt.Errorf("Unexpected number of links. Expected at least 1 link, got: %v", len(serviceResponse.Links))
		}

		if serviceResponse.Owner == nil {
			return fmt.Errorf("Unexpected owner. Expected owner to be set.")
		}
		expected, got = serviceResource.Primary.Attributes["owner_id"], *serviceResponse.Owner.ID
		if expected != got {
			return fmt.Errorf("Unexpected owner ID. Expected:%s, got: %s", expected, got)
		}

		expected, got = serviceResource.Primary.Attributes["service_tier"], fmt.Sprintf("%d", *serviceResponse.ServiceTier)
		if expected != got {
			return fmt.Errorf("Unexpected service_tier. Expected: %s, got: %s", expected, got)
		}

		// TODO: Check the team ids
		if len(serviceResponse.Teams) != 2 {
			return fmt.Errorf("Unexpected number of teams. Expected: 2, got: %v", len(serviceResponse.Teams))
		}

		return nil
	}
}

func testAccCheckServiceResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_service" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			_, err := client.Sdk.CatalogEntries.GetService(context.TODO(), stateResource.Primary.ID)
			if err == nil {
				return fmt.Errorf("Service %s still exists", stateResource.Primary.ID)
			}
		}

		return nil
	}
}

func testAccServiceResourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service" {
  name = "test-service-%s"
}`, rName)
}

func testAccServiceResourceConfig_update(rName, sharedTeamID, sharedTeamID2 string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service" {
  name                     = "test-service-%s"
  alert_on_add             = true
  auto_add_responding_team = true
  description              = "test-description-%s"
  labels = {
    test1 = "test-label1-%s",
  }
  links {
    href_url = "https://example.com/test-link1-%s"
    name = "test-link1-%s"
  }
  links {
    href_url = "https://example.com/test-link2-%s"
    name = "test-link2-%s"
  }
  owner_id     = "%s"
  service_tier = "1"
  team_ids = [
    "%s",
    "%s"
  ]
}`, rName, rName, rName, rName, rName, rName, rName, sharedTeamID, sharedTeamID, sharedTeamID2)
}

func testAccServiceResourceConfig_updateChangeLabels(rName, sharedTeamID, sharedTeamID2 string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service" {
  name                     = "test-service-%s"
  alert_on_add             = true
  auto_add_responding_team = true
  description              = "test-description-%s"
  labels = {
    test1 = "test-label1-%s",
    test2 = "test-label2-%s"
  }
  links {
    href_url = "https://example.com/test-link1-%s"
    name = "test-link1-%s"
  }
  links {
    href_url = "https://example.com/test-link2-%s"
    name = "test-link2-%s"
  }
  owner_id     = "%s"
  service_tier = "1"
  team_ids = [
    "%s",
    "%s"
  ]
}`, rName, rName, rName, rName, rName, rName, rName, rName, sharedTeamID, sharedTeamID, sharedTeamID2)
}

func testAccServiceResourceConfig_updateChangeOwnerIDAndTeamIDs(rName, sharedTeamID, sharedTeamID2 string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service" {
  name                     = "test-service-%s"
  alert_on_add             = true
  auto_add_responding_team = true
  description              = "test-description-%s"
  labels = {
    test1 = "test-label1-%s"
  }
  links {
    href_url = "https://example.com/test-link1-%s"
    name = "test-link1-%s"
  }
  links {
    href_url = "https://example.com/test-link2-%s"
    name = "test-link2-%s"
  }
  owner_id     = "%s"
  service_tier = "1"
  team_ids = [
    "%s",
    "%s"
  ]
}`, rName, rName, rName, rName, rName, rName, rName, sharedTeamID, sharedTeamID, sharedTeamID2)
}
