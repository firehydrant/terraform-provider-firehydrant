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

func TestAccTeamResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckTeamResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTeamResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTeamResourceExistsWithAttributes_basic("firehydrant_team.test_team"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "name", fmt.Sprintf("test-team-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_team.test_team", "services.#", "0"),
				),
			},
		},
	})
}

func TestAccTeamResource_update(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckTeamResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTeamResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTeamResourceExistsWithAttributes_basic("firehydrant_team.test_team"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "name", fmt.Sprintf("test-team-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_team.test_team", "services.#", "0"),
				),
			},
			{
				Config: testAccTeamResourceConfig_update(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamResourceExistsWithAttributes_update("firehydrant_team.test_team"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "name", fmt.Sprintf("test-team-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "services.#", "1"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "services.0.id"),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "services.0.name", fmt.Sprintf("test-service-%s", rNameUpdated)),
				),
			},
			{
				Config: testAccTeamResourceConfig_basic(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTeamResourceExistsWithAttributes_basic("firehydrant_team.test_team"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "name", fmt.Sprintf("test-team-%s", rNameUpdated)),
					resource.TestCheckResourceAttr("firehydrant_team.test_team", "services.#", "0"),
				),
			},
		},
	})
}

func TestAccTeamResource_basicServiceIDs(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckTeamResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTeamResourceConfig_basicServiceIDs(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTeamResourceExistsWithAttributes_basicServiceIDs("firehydrant_team.test_team"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "name", fmt.Sprintf("test-team-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_team.test_team", "service_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccTeamResource_updateServiceIDs(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckTeamResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTeamResourceConfig_basicServiceIDs(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTeamResourceExistsWithAttributes_basicServiceIDs("firehydrant_team.test_team"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "name", fmt.Sprintf("test-team-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_team.test_team", "service_ids.#", "1"),
				),
			},
			{
				Config: testAccTeamResourceConfig_updateServiceIDs(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamResourceExistsWithAttributes_updateServiceIDs("firehydrant_team.test_team"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "name", fmt.Sprintf("test-team-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "service_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccTeamResource_updateServicesToServiceIDs(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckTeamResourceDestroy(),
		Steps: []resource.TestStep{
			// Start with a config that has services
			{
				Config: testAccTeamResourceConfig_basicServices(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTeamResourceExistsWithAttributes_basicServices("firehydrant_team.test_team"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "name", fmt.Sprintf("test-team-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_team.test_team", "services.#", "1"),
				),
			},
			// Update the config by changing services to service_ids
			{
				Config: testAccTeamResourceConfig_basicServiceIDs(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTeamResourceExistsWithAttributes_basicServiceIDs("firehydrant_team.test_team"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "name", fmt.Sprintf("test-team-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_team.test_team", "service_ids.#", "1"),
				),
			},
			// Update the config by adding a new service id to service_ids
			{
				Config: testAccTeamResourceConfig_updateServiceIDs(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamResourceExistsWithAttributes_updateServiceIDs("firehydrant_team.test_team"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "name", fmt.Sprintf("test-team-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "service_ids.#", "2"),
				),
			},
			// Update the config by removing a service id from service_ids
			{
				Config: testAccTeamResourceConfig_basicServiceIDs(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTeamResourceExistsWithAttributes_basicServiceIDs("firehydrant_team.test_team"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "name", fmt.Sprintf("test-team-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "service_ids.#", "1"),
				),
			},
			// Update the config by removing service_ids
			{
				Config: testAccTeamResourceConfig_basic(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTeamResourceExistsWithAttributes_basic("firehydrant_team.test_team"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "name", fmt.Sprintf("test-team-%s", rNameUpdated)),
					resource.TestCheckResourceAttr("firehydrant_team.test_team", "services.#", "0"),
				),
			},
		},
	})
}

func TestAccTeamResourceImport_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTeamResourceConfig_basic(rName),
			},

			{
				ResourceName:      "firehydrant_team.test_team",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTeamResourceImport_services(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTeamResourceConfig_basicServices(rName),
			},

			{
				ResourceName:            "firehydrant_team.test_team",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"services", "service_ids"},
			},
		},
	})
}

func TestAccTeamResourceImport_serviceIDs(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTeamResourceConfig_basicServiceIDs(rName),
			},

			{
				ResourceName:      "firehydrant_team.test_team",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTeamResourceExistsWithAttributes_basic(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		teamResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if teamResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		teamResponse, err := client.GetTeam(context.TODO(), teamResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := teamResource.Primary.Attributes["name"], teamResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		if teamResponse.Description != "" {
			return fmt.Errorf("Unexpected description. Expected no description, got: %s", teamResponse.Description)
		}

		if len(teamResponse.Services) != 0 {
			return fmt.Errorf("Unexpected number of services. Expected no services, got: %v", len(teamResponse.Services))
		}

		return nil
	}
}

func testAccCheckTeamResourceExistsWithAttributes_update(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		teamResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if teamResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		teamResponse, err := client.GetTeam(context.TODO(), teamResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := teamResource.Primary.Attributes["name"], teamResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = teamResource.Primary.Attributes["description"], teamResponse.Description
		if expected != got {
			return fmt.Errorf("Unexpected description. Expected: %s, got: %s", expected, got)
		}

		// TODO: Check the service ids
		if len(teamResponse.Services) != 1 {
			return fmt.Errorf("Unexpected number of services. Expected: 1, got: %v", len(teamResponse.Services))
		}

		return nil
	}
}

func testAccCheckTeamResourceExistsWithAttributes_basicServices(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		teamResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if teamResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		teamResponse, err := client.GetTeam(context.TODO(), teamResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := teamResource.Primary.Attributes["name"], teamResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		if teamResponse.Description != "" {
			return fmt.Errorf("Unexpected description. Expected no description, got: %s", teamResponse.Description)
		}

		// TODO: Check the service ids
		if len(teamResponse.Services) != 1 {
			return fmt.Errorf("Unexpected number of services. Expected: 1, got: %v", len(teamResponse.Services))
		}

		return nil
	}
}

func testAccCheckTeamResourceExistsWithAttributes_basicServiceIDs(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		teamResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if teamResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		teamResponse, err := client.GetTeam(context.TODO(), teamResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := teamResource.Primary.Attributes["name"], teamResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		if teamResponse.Description != "" {
			return fmt.Errorf("Unexpected description. Expected no description, got: %s", teamResponse.Description)
		}

		// TODO: Check the service ids
		if len(teamResponse.Services) != 1 {
			return fmt.Errorf("Unexpected number of services. Expected: 1, got: %v", len(teamResponse.Services))
		}

		return nil
	}
}

func testAccCheckTeamResourceExistsWithAttributes_updateServiceIDs(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		teamResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if teamResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		teamResponse, err := client.GetTeam(context.TODO(), teamResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := teamResource.Primary.Attributes["name"], teamResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = teamResource.Primary.Attributes["description"], teamResponse.Description
		if expected != got {
			return fmt.Errorf("Unexpected description. Expected: %s, got: %s", expected, got)
		}

		// TODO: Check the service ids
		if len(teamResponse.Services) != 2 {
			return fmt.Errorf("Unexpected number of services. Expected: 1, got: %v", len(teamResponse.Services))
		}

		return nil
	}
}

func testAccCheckTeamResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// TODO: add this back once the bug in the API is fixed
		//client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		//if err != nil {
		//	return err
		//}
		//
		//for _, teamResource := range s.RootModule().Resources {
		//	if teamResource.Type != "firehydrant_team" {
		//		continue
		//	}
		//
		//	if teamResource.Primary.ID == "" {
		//		return fmt.Errorf("No instance ID is set")
		//	}
		//
		//	_, err := client.GetTeam(context.TODO(), teamResource.Primary.ID)
		//	if err == nil {
		//		return fmt.Errorf("Team %s still exists", teamResource.Primary.ID)
		//	}
		//}

		return nil
	}
}

func testAccTeamResourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "test_team" {
  name = "test-team-%s"
}`, rName)
}

func testAccTeamResourceConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service" {
  name = "test-service-%s"
}

resource "firehydrant_team" "test_team" {
  name        = "test-team-%s"
  description = "test-description-%s"

  services {
    id = firehydrant_service.test_service.id
  }
}`, rName, rName, rName)
}

func testAccTeamResourceConfig_basicServiceIDs(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service1" {
  name = "test-service1-%s"
}

resource "firehydrant_team" "test_team" {
  name = "test-team-%s"

  service_ids = [firehydrant_service.test_service1.id]
}`, rName, rName)
}

func testAccTeamResourceConfig_updateServiceIDs(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service1" {
  name = "test-service1-%s"
}

resource "firehydrant_service" "test_service2" {
  name = "test-service2-%s"
}

resource "firehydrant_team" "test_team" {
  name        = "test-team-%s"
  description = "test-description-%s"

  service_ids = [
    firehydrant_service.test_service1.id,
    firehydrant_service.test_service2.id
  ]
}`, rName, rName, rName, rName)
}

func testAccTeamResourceConfig_basicServices(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service1" {
  name = "test-service1-%s"
}

resource "firehydrant_team" "test_team" {
  name = "test-team-%s"

  services {
    id = firehydrant_service.test_service1.id
  }
}`, rName, rName)
}
