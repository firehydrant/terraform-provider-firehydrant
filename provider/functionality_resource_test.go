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

func TestAccFunctionalityResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckFunctionalityResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionalityResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFunctionalityResourceExistsWithAttributes_basic("firehydrant_functionality.test_functionality"),
					resource.TestCheckResourceAttrSet("firehydrant_functionality.test_functionality", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "name", fmt.Sprintf("test-functionality-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_functionality.test_functionality", "services.#", "0"),
				),
			},
		},
	})
}

func TestAccFunctionalityResource_update(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckFunctionalityResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionalityResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFunctionalityResourceExistsWithAttributes_basic("firehydrant_functionality.test_functionality"),
					resource.TestCheckResourceAttrSet("firehydrant_functionality.test_functionality", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "name", fmt.Sprintf("test-functionality-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_functionality.test_functionality", "services.#", "0"),
				),
			},
			{
				Config: testAccFunctionalityResourceConfig_update(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionalityResourceExistsWithAttributes_update("firehydrant_functionality.test_functionality"),
					resource.TestCheckResourceAttrSet("firehydrant_functionality.test_functionality", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "name", fmt.Sprintf("test-functionality-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
					resource.TestCheckResourceAttr("firehydrant_functionality.test_functionality", "services.#", "1"),
					resource.TestCheckResourceAttrSet("firehydrant_functionality.test_functionality", "services.0.id"),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "services.0.name", fmt.Sprintf("test-service-%s", rNameUpdated)),
				),
			},
			{
				Config: testAccFunctionalityResourceConfig_basic(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFunctionalityResourceExistsWithAttributes_basic("firehydrant_functionality.test_functionality"),
					resource.TestCheckResourceAttrSet("firehydrant_functionality.test_functionality", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "name", fmt.Sprintf("test-functionality-%s", rNameUpdated)),
					resource.TestCheckResourceAttr("firehydrant_functionality.test_functionality", "services.#", "0"),
				),
			},
		},
	})
}

func TestAccFunctionalityResource_basicServiceIDs(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckFunctionalityResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionalityResourceConfig_basicServiceIDs(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFunctionalityResourceExistsWithAttributes_basicServiceIDs("firehydrant_functionality.test_functionality"),
					resource.TestCheckResourceAttrSet("firehydrant_functionality.test_functionality", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "name", fmt.Sprintf("test-functionality-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_functionality.test_functionality", "service_ids.#", "1"),
				),
			},
		},
	})
}

func TestAccFunctionalityResource_updateServiceIDs(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckFunctionalityResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionalityResourceConfig_basicServiceIDs(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFunctionalityResourceExistsWithAttributes_basicServiceIDs("firehydrant_functionality.test_functionality"),
					resource.TestCheckResourceAttrSet("firehydrant_functionality.test_functionality", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "name", fmt.Sprintf("test-functionality-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_functionality.test_functionality", "service_ids.#", "1"),
				),
			},
			{
				Config: testAccFunctionalityResourceConfig_updateServiceIDs(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionalityResourceExistsWithAttributes_updateServiceIDs("firehydrant_functionality.test_functionality"),
					resource.TestCheckResourceAttrSet("firehydrant_functionality.test_functionality", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "name", fmt.Sprintf("test-functionality-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "service_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccFunctionalityResource_updateServicesToServiceIDs(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckFunctionalityResourceDestroy(),
		Steps: []resource.TestStep{
			// Start with a config that has services
			{
				Config: testAccFunctionalityResourceConfig_basicServices(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFunctionalityResourceExistsWithAttributes_basicServices("firehydrant_functionality.test_functionality"),
					resource.TestCheckResourceAttrSet("firehydrant_functionality.test_functionality", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "name", fmt.Sprintf("test-functionality-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_functionality.test_functionality", "services.#", "1"),
				),
			},
			// Update the config by changing services to service_ids
			{
				Config: testAccFunctionalityResourceConfig_basicServiceIDs(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFunctionalityResourceExistsWithAttributes_basicServiceIDs("firehydrant_functionality.test_functionality"),
					resource.TestCheckResourceAttrSet("firehydrant_functionality.test_functionality", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "name", fmt.Sprintf("test-functionality-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_functionality.test_functionality", "service_ids.#", "1"),
				),
			},
			// Update the config by adding a new service id to service_ids
			{
				Config: testAccFunctionalityResourceConfig_updateServiceIDs(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionalityResourceExistsWithAttributes_updateServiceIDs("firehydrant_functionality.test_functionality"),
					resource.TestCheckResourceAttrSet("firehydrant_functionality.test_functionality", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "name", fmt.Sprintf("test-functionality-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "service_ids.#", "2"),
				),
			},
			// Update the config by removing a service id from service_ids
			{
				Config: testAccFunctionalityResourceConfig_basicServiceIDs(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFunctionalityResourceExistsWithAttributes_basicServiceIDs("firehydrant_functionality.test_functionality"),
					resource.TestCheckResourceAttrSet("firehydrant_functionality.test_functionality", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "name", fmt.Sprintf("test-functionality-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "service_ids.#", "1"),
				),
			},
			// Update the config by removing service_ids
			{
				Config: testAccFunctionalityResourceConfig_basic(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFunctionalityResourceExistsWithAttributes_basic("firehydrant_functionality.test_functionality"),
					resource.TestCheckResourceAttrSet("firehydrant_functionality.test_functionality", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_functionality.test_functionality", "name", fmt.Sprintf("test-functionality-%s", rNameUpdated)),
					resource.TestCheckResourceAttr("firehydrant_functionality.test_functionality", "service_ids.#", "0"),
				),
			},
		},
	})
}

func TestAccFunctionalityResourceImport_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionalityResourceConfig_basic(rName),
			},

			{
				ResourceName:      "firehydrant_functionality.test_functionality",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFunctionalityResourceImport_services(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionalityResourceConfig_basicServices(rName),
			},

			{
				ResourceName:            "firehydrant_functionality.test_functionality",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"services", "service_ids"},
			},
		},
	})
}

func TestAccFunctionalityResourceImport_serviceIDs(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionalityResourceConfig_basicServiceIDs(rName),
			},

			{
				ResourceName:      "firehydrant_functionality.test_functionality",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckFunctionalityResourceExistsWithAttributes_basic(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		functionalityResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if functionalityResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		functionalityResponse, err := client.Functionalities().Get(context.TODO(), functionalityResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := functionalityResource.Primary.Attributes["name"], functionalityResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		if functionalityResponse.Description != "" {
			return fmt.Errorf("Unexpected description. Expected no description, got: %s", functionalityResponse.Description)
		}

		if len(functionalityResponse.Services) != 0 {
			return fmt.Errorf("Unexpected number of services. Expected no services, got: %v", len(functionalityResponse.Services))
		}

		return nil
	}
}

func testAccCheckFunctionalityResourceExistsWithAttributes_update(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		functionalityResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if functionalityResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		functionalityResponse, err := client.Functionalities().Get(context.TODO(), functionalityResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := functionalityResource.Primary.Attributes["name"], functionalityResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = functionalityResource.Primary.Attributes["description"], functionalityResponse.Description
		if expected != got {
			return fmt.Errorf("Unexpected description. Expected: %s, got: %s", expected, got)
		}

		// TODO: Check the service ids
		if len(functionalityResponse.Services) != 1 {
			return fmt.Errorf("Unexpected number of services. Expected: 1, got: %v", len(functionalityResponse.Services))
		}

		return nil
	}
}

func testAccCheckFunctionalityResourceExistsWithAttributes_basicServices(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		functionalityResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if functionalityResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		functionalityResponse, err := client.Functionalities().Get(context.TODO(), functionalityResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := functionalityResource.Primary.Attributes["name"], functionalityResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		if functionalityResponse.Description != "" {
			return fmt.Errorf("Unexpected description. Expected no description, got: %s", functionalityResponse.Description)
		}

		// TODO: Check the service ids
		if len(functionalityResponse.Services) != 1 {
			return fmt.Errorf("Unexpected number of services. Expected: 1, got: %v", len(functionalityResponse.Services))
		}

		return nil
	}
}

func testAccCheckFunctionalityResourceExistsWithAttributes_basicServiceIDs(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		functionalityResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if functionalityResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		functionalityResponse, err := client.Functionalities().Get(context.TODO(), functionalityResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := functionalityResource.Primary.Attributes["name"], functionalityResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		if functionalityResponse.Description != "" {
			return fmt.Errorf("Unexpected description. Expected no description, got: %s", functionalityResponse.Description)
		}

		// TODO: Check the service ids
		if len(functionalityResponse.Services) != 1 {
			return fmt.Errorf("Unexpected number of services. Expected: 1, got: %v", len(functionalityResponse.Services))
		}

		return nil
	}
}

func testAccCheckFunctionalityResourceExistsWithAttributes_updateServiceIDs(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		functionalityResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if functionalityResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		functionalityResponse, err := client.Functionalities().Get(context.TODO(), functionalityResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := functionalityResource.Primary.Attributes["name"], functionalityResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = functionalityResource.Primary.Attributes["description"], functionalityResponse.Description
		if expected != got {
			return fmt.Errorf("Unexpected description. Expected: %s, got: %s", expected, got)
		}

		// TODO: Check the service ids
		if len(functionalityResponse.Services) != 2 {
			return fmt.Errorf("Unexpected number of services. Expected: 1, got: %v", len(functionalityResponse.Services))
		}

		return nil
	}
}

func testAccCheckFunctionalityResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		for _, functionalityResource := range s.RootModule().Resources {
			if functionalityResource.Type != "firehydrant_functionality" {
				continue
			}

			if functionalityResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			_, err := client.Functionalities().Get(context.TODO(), functionalityResource.Primary.ID)
			if err == nil {
				return fmt.Errorf("Functionality %s still exists", functionalityResource.Primary.ID)
			}
		}

		return nil
	}
}

func testAccFunctionalityResourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_functionality" "test_functionality" {
  name = "test-functionality-%s"
}`, rName)
}

func testAccFunctionalityResourceConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service" {
  name = "test-service-%s"
}

resource "firehydrant_functionality" "test_functionality" {
  name        = "test-functionality-%s"
  description = "test-description-%s"

  services {
    id = firehydrant_service.test_service.id
  }
}`, rName, rName, rName)
}

func testAccFunctionalityResourceConfig_basicServiceIDs(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service1" {
  name = "test-service1-%s"
}

resource "firehydrant_functionality" "test_functionality" {
  name = "test-functionality-%s"

  service_ids = [firehydrant_service.test_service1.id]
}`, rName, rName)
}

func testAccFunctionalityResourceConfig_updateServiceIDs(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service1" {
  name = "test-service1-%s"
}

resource "firehydrant_service" "test_service2" {
  name = "test-service2-%s"
}

resource "firehydrant_functionality" "test_functionality" {
  name        = "test-functionality-%s"
  description = "test-description-%s"

  service_ids = [
    firehydrant_service.test_service1.id,
    firehydrant_service.test_service2.id
  ]
}`, rName, rName, rName, rName)
}

func testAccFunctionalityResourceConfig_basicServices(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_service" "test_service1" {
  name = "test-service1-%s"
}

resource "firehydrant_functionality" "test_functionality" {
  name = "test-functionality-%s"

  services {
    id = firehydrant_service.test_service1.id
  }
}`, rName, rName)
}
