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

func TestAccEnvironmentResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckEnvironmentResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEnvironmentResourceExistsWithAttributes_basic("firehydrant_environment.test_environment"),
					resource.TestCheckResourceAttrSet("firehydrant_environment.test_environment", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_environment.test_environment", "name", fmt.Sprintf("test-environment-%s", rName)),
				),
			},
		},
	})
}

func TestAccEnvironmentResource_update(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckEnvironmentResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEnvironmentResourceExistsWithAttributes_basic("firehydrant_environment.test_environment"),
					resource.TestCheckResourceAttrSet("firehydrant_environment.test_environment", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_environment.test_environment", "name", fmt.Sprintf("test-environment-%s", rName)),
				),
			},
			{
				Config: testAccEnvironmentResourceConfig_update(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEnvironmentResourceExistsWithAttributes_update("firehydrant_environment.test_environment"),
					resource.TestCheckResourceAttrSet("firehydrant_environment.test_environment", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_environment.test_environment", "name", fmt.Sprintf("test-environment-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_environment.test_environment", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
				),
			},
			{
				Config: testAccEnvironmentResourceConfig_basic(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEnvironmentResourceExistsWithAttributes_basic("firehydrant_environment.test_environment"),
					resource.TestCheckResourceAttrSet("firehydrant_environment.test_environment", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_environment.test_environment", "name", fmt.Sprintf("test-environment-%s", rNameUpdated)),
				),
			},
		},
	})
}

func TestAccEnvironmentResourceImport_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentResourceConfig_basic(rName),
			},
			{
				ResourceName:      "firehydrant_environment.test_environment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEnvironmentResourceImport_allAttributes(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentResourceConfig_update(rName),
			},
			{
				ResourceName:      "firehydrant_environment.test_environment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckEnvironmentResourceExistsWithAttributes_basic(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		environmentResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if environmentResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		environmentResponse, err := client.Environments().Get(context.TODO(), environmentResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := environmentResource.Primary.Attributes["name"], environmentResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		if environmentResponse.Description != "" {
			return fmt.Errorf("Unexpected description. Expected no description, got: %s", environmentResponse.Description)
		}

		return nil
	}
}

func testAccCheckEnvironmentResourceExistsWithAttributes_update(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		environmentResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if environmentResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		environmentResponse, err := client.Environments().Get(context.TODO(), environmentResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := environmentResource.Primary.Attributes["name"], environmentResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = environmentResource.Primary.Attributes["description"], environmentResponse.Description
		if expected != got {
			return fmt.Errorf("Unexpected description. Expected: %s, got: %s", expected, got)
		}

		return nil
	}
}

func testAccCheckEnvironmentResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_environment" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			_, err := client.Environments().Get(context.TODO(), stateResource.Primary.ID)
			if err == nil {
				return fmt.Errorf("Environment %s still exists", stateResource.Primary.ID)
			}
		}

		return nil
	}
}

func testAccEnvironmentResourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_environment" "test_environment" {
  name    = "test-environment-%s"
}`, rName)
}

func testAccEnvironmentResourceConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_environment" "test_environment" {
  name        = "test-environment-%s"
  description = "test-description-%s"
}`, rName, rName)
}
