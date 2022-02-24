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

func TestAccServiceResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
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
						"firehydrant_service.test_service", "service_tier", "5"),
					resource.TestCheckResourceAttr("firehydrant_team.test_service", "team_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccServiceResource_update(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
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
						"firehydrant_service.test_service", "service_tier", "5"),
					resource.TestCheckResourceAttr("firehydrant_team.test_service", "team_ids.#", "0"),
				),
			},
			{
				Config: testAccServiceResourceConfig_update(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceResourceExistsWithAttributes_update("firehydrant_service.test_service"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "alert_on_add", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "owner_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "service_tier", "1"),
					resource.TestCheckResourceAttr("firehydrant_team.test_service", "team_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccServiceResource_updateOwnerIDAndTeamIDs(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckServiceResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfig_update(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceResourceExistsWithAttributes_update("firehydrant_service.test_service"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "alert_on_add", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "description", fmt.Sprintf("test-description-%s", rName)),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "owner_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "service_tier", "1"),
					resource.TestCheckResourceAttr("firehydrant_team.test_service", "team_ids.#", "2"),
				),
			},
			{
				Config: testAccServiceResourceConfig_updateChangeOwnerIDAndTeamIDs(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceResourceExistsWithAttributes_update("firehydrant_service.test_service"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "alert_on_add", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "description", fmt.Sprintf("test-description-%s", rName)),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "owner_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "service_tier", "1"),
					resource.TestCheckResourceAttr("firehydrant_team.test_service", "team_ids.#", "2"),
				),
			},
			{
				Config: testAccServiceResourceConfig_updateRemoveOwnerIDAndTeamIDs(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceResourceExistsWithAttributes_updateRemoveOwnerIDAndTeamIDs("firehydrant_service.test_service"),
					resource.TestCheckResourceAttrSet("firehydrant_service.test_service", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "name", fmt.Sprintf("test-service-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "alert_on_add", "true"),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "description", fmt.Sprintf("test-description-%s", rName)),
					// Make sure owner_id is not set
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "owner_id", ""),
					resource.TestCheckResourceAttr(
						"firehydrant_service.test_service", "service_tier", "1"),
					resource.TestCheckResourceAttr("firehydrant_team.test_service", "team_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccServiceResourceImport_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
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
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceResourceConfig_update(rName),
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

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		serviceResponse, err := client.Services().Get(context.TODO(), serviceResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := serviceResource.Primary.Attributes["name"], serviceResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = serviceResource.Primary.Attributes["alert_on_add"], fmt.Sprintf("%t", serviceResponse.AlertOnAdd)
		if expected != got {
			return fmt.Errorf("Unexpected alert_on_add. Expected: %s, got: %s", expected, got)
		}

		if serviceResponse.Description != "" {
			return fmt.Errorf("Unexpected description. Expected no description, got: %s", serviceResponse.Description)
		}

		//if !reflect.DeepEqual(serviceResponse.Labels, []string{}) {
		//	return fmt.Errorf("Bad labels: %v", serviceResponse.Labels)
		//}

		if serviceResponse.Owner != nil {
			return fmt.Errorf("Unexpected owner. Expected no owner ID, got: %s", serviceResponse.Owner.ID)
		}

		expected, got = serviceResource.Primary.Attributes["service_tier"], fmt.Sprintf("%d", serviceResponse.ServiceTier)
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

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		serviceResponse, err := client.Services().Get(context.TODO(), serviceResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := serviceResource.Primary.Attributes["name"], serviceResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = serviceResource.Primary.Attributes["alert_on_add"], fmt.Sprintf("%t", serviceResponse.AlertOnAdd)
		if expected != got {
			return fmt.Errorf("Unexpected alert_on_add. Expected: %s, got: %s", expected, got)
		}

		expected, got = serviceResource.Primary.Attributes["description"], serviceResponse.Description
		if expected != got {
			return fmt.Errorf("Unexpected description. Expected: %s, got: %s", expected, got)
		}

		//if !reflect.DeepEqual(serviceResponse.Labels, []string{}) {
		//	return fmt.Errorf("Bad labels: %v", serviceResponse.Labels)
		//}

		if serviceResponse.Owner == nil {
			return fmt.Errorf("Unexpected owner. Expected owner to be set.")
		}
		expected, got = serviceResource.Primary.Attributes["owner_id"], serviceResponse.Owner.ID
		if expected != got {
			return fmt.Errorf("Unexpected owner ID. Expected:%s, got: %s", expected, got)
		}

		expected, got = serviceResource.Primary.Attributes["service_tier"], fmt.Sprintf("%d", serviceResponse.ServiceTier)
		if expected != got {
			return fmt.Errorf("Unexpected service_tier. Expected: %s, got: %s", expected, got)
		}

		// TODO: Check the team ids
		if len(serviceResponse.Teams) != 2 {
			return fmt.Errorf("Unexpected number of services. Expected: 2, got: %v", len(serviceResponse.Teams))
		}

		return nil
	}
}

func testAccCheckServiceResourceExistsWithAttributes_updateRemoveOwnerIDAndTeamIDs(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		serviceResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if serviceResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		serviceResponse, err := client.Services().Get(context.TODO(), serviceResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := serviceResource.Primary.Attributes["name"], serviceResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = serviceResource.Primary.Attributes["alert_on_add"], fmt.Sprintf("%t", serviceResponse.AlertOnAdd)
		if expected != got {
			return fmt.Errorf("Unexpected alert_on_add. Expected: %s, got: %s", expected, got)
		}

		expected, got = serviceResource.Primary.Attributes["description"], serviceResponse.Description
		if expected != got {
			return fmt.Errorf("Unexpected description. Expected: %s, got: %s", expected, got)
		}

		//if !reflect.DeepEqual(serviceResponse.Labels, []string{}) {
		//	return fmt.Errorf("Bad labels: %v", serviceResponse.Labels)
		//}

		if serviceResponse.Owner != nil {
			return fmt.Errorf("Unexpected owner. Expected owner to not be set, got: %s.", serviceResponse.Owner.ID)
		}

		expected, got = serviceResource.Primary.Attributes["service_tier"], fmt.Sprintf("%d", serviceResponse.ServiceTier)
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

func testAccCheckServiceResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
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

			_, err := client.Services().Get(context.TODO(), stateResource.Primary.ID)
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

func testAccServiceResourceConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "test_team1" {
  name = "test-team1-%s"
}

resource "firehydrant_team" "test_team2" {
  name = "test-team2-%s"
}

resource "firehydrant_team" "test_team3" {
  name = "test-team3-%s"
}

resource "firehydrant_service" "test_service" {
  name         = "test-service-%s"
  alert_on_add = true
  description  = "test-description-%s"
  owner_id     = firehydrant_team.test_team1.id
  service_tier = "1"
  team_ids = [
    firehydrant_team.test_team2.id,
    firehydrant_team.test_team3.id
  ]
}`, rName, rName, rName, rName, rName)
}

func testAccServiceResourceConfig_updateChangeOwnerIDAndTeamIDs(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "test_team1" {
  name = "test-team1-%s"
}

resource "firehydrant_team" "test_team2" {
  name = "test-team2-%s"
}

resource "firehydrant_team" "test_team3" {
  name = "test-team3-%s"
}

resource "firehydrant_service" "test_service" {
  name         = "test-service-%s"
  alert_on_add = true
  description  = "test-description-%s"
  owner_id     = firehydrant_team.test_team2.id
  service_tier = "1"
  team_ids = [
    firehydrant_team.test_team1.id,
    firehydrant_team.test_team3.id
  ]
}`, rName, rName, rName, rName, rName)
}

func testAccServiceResourceConfig_updateRemoveOwnerIDAndTeamIDs(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service" {
  name         = "test-service-%s"
  alert_on_add = true
  description  = "test-description-%s"
  service_tier = "1"
}`, rName, rName)
}
