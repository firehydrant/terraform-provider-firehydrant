package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccServiceDependencyResource_basic(t *testing.T) {
	sharedServiceID1 := getSharedServiceID(t)
	sharedServiceID2 := getSharedServiceID2(t)
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckServiceDependencyResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDependencyResourceConfig_basic(rName, sharedServiceID1, sharedServiceID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceDependencyResourceExistsWithAttributes_basic("firehydrant_service_dependency.test_service_dependency"),
					resource.TestCheckResourceAttrSet("firehydrant_service_dependency.test_service_dependency", "id"),
					resource.TestCheckResourceAttrSet("firehydrant_service_dependency.test_service_dependency", "connected_service_id"),
					resource.TestCheckResourceAttrSet("firehydrant_service_dependency.test_service_dependency", "service_id"),
				),
			},
		},
	})
}

func TestAccServiceDependencyResource_update(t *testing.T) {
	sharedServiceID1 := getSharedServiceID(t)
	sharedServiceID2 := getSharedServiceID2(t)
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckServiceDependencyResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDependencyResourceConfig_basic(rName, sharedServiceID1, sharedServiceID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceDependencyResourceExistsWithAttributes_basic("firehydrant_service_dependency.test_service_dependency"),
					resource.TestCheckResourceAttrSet("firehydrant_service_dependency.test_service_dependency", "id"),
					resource.TestCheckResourceAttrSet("firehydrant_service_dependency.test_service_dependency", "connected_service_id"),
					resource.TestCheckResourceAttrSet("firehydrant_service_dependency.test_service_dependency", "service_id"),
				),
			},
			{
				Config: testAccServiceDependencyResourceConfig_update(rNameUpdated, sharedServiceID1, sharedServiceID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceDependencyResourceExistsWithAttributes_update("firehydrant_service_dependency.test_service_dependency"),
					resource.TestCheckResourceAttrSet("firehydrant_service_dependency.test_service_dependency", "id"),
					resource.TestCheckResourceAttrSet("firehydrant_service_dependency.test_service_dependency", "connected_service_id"),
					resource.TestCheckResourceAttrSet("firehydrant_service_dependency.test_service_dependency", "service_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_service_dependency.test_service_dependency", "notes", fmt.Sprintf("test-notes-%s", rNameUpdated)),
				),
			},
			{
				Config: testAccServiceDependencyResourceConfig_basic(rNameUpdated, sharedServiceID1, sharedServiceID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceDependencyResourceExistsWithAttributes_basic("firehydrant_service_dependency.test_service_dependency"),
					resource.TestCheckResourceAttrSet("firehydrant_service_dependency.test_service_dependency", "id"),
					resource.TestCheckResourceAttrSet("firehydrant_service_dependency.test_service_dependency", "connected_service_id"),
					resource.TestCheckResourceAttrSet("firehydrant_service_dependency.test_service_dependency", "service_id"),
				),
			},
		},
	})
}

func TestAccServiceDependencyResourceImport_basic(t *testing.T) {
	sharedServiceID1 := getSharedServiceID(t)
	sharedServiceID2 := getSharedServiceID2(t)
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDependencyResourceConfig_basic(rName, sharedServiceID1, sharedServiceID2),
			},
			{
				ResourceName:      "firehydrant_service_dependency.test_service_dependency",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceDependencyResourceImport_allAttributes(t *testing.T) {
	sharedServiceID1 := getSharedServiceID(t)
	sharedServiceID2 := getSharedServiceID2(t)
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceDependencyResourceConfig_update(rName, sharedServiceID1, sharedServiceID2),
			},
			{
				ResourceName:      "firehydrant_service_dependency.test_service_dependency",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckServiceDependencyResourceExistsWithAttributes_basic(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		serviceDependencyResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if serviceDependencyResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		serviceDependencyResponse, err := client.ServiceDependencies().Get(context.TODO(), serviceDependencyResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := serviceDependencyResource.Primary.Attributes["connected_service_id"], serviceDependencyResponse.ConnectedService.ID
		if expected != got {
			return fmt.Errorf("Unexpected connected_service_id. Expected: %s, got: %s", expected, got)
		}

		expected, got = serviceDependencyResource.Primary.Attributes["service_id"], serviceDependencyResponse.Service.ID
		if expected != got {
			return fmt.Errorf("Unexpected service_id. Expected: %s, got: %s", expected, got)
		}

		if serviceDependencyResponse.Notes != "" {
			return fmt.Errorf("Unexpected notes. Expected no notes, got: %s", serviceDependencyResponse.Notes)
		}

		return nil
	}
}

func testAccCheckServiceDependencyResourceExistsWithAttributes_update(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		serviceDependencyResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if serviceDependencyResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		serviceDependencyResponse, err := client.ServiceDependencies().Get(context.TODO(), serviceDependencyResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := serviceDependencyResource.Primary.Attributes["connected_service_id"], serviceDependencyResponse.ConnectedService.ID
		if expected != got {
			return fmt.Errorf("Unexpected connected_service_id. Expected: %s, got: %s", expected, got)
		}

		expected, got = serviceDependencyResource.Primary.Attributes["service_id"], serviceDependencyResponse.Service.ID
		if expected != got {
			return fmt.Errorf("Unexpected service_id. Expected: %s, got: %s", expected, got)
		}

		expected, got = serviceDependencyResource.Primary.Attributes["notes"], serviceDependencyResponse.Notes
		if expected != got {
			return fmt.Errorf("Unexpected notes. Expected: %s, got: %s", expected, got)
		}

		return nil
	}
}

func testAccCheckServiceDependencyResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_service_dependency" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			_, err := client.ServiceDependencies().Get(context.TODO(), stateResource.Primary.ID)
			if err == nil {
				return fmt.Errorf("Service dependency %s still exists", stateResource.Primary.ID)
			}
		}

		return nil
	}
}

func testAccServiceDependencyResourceConfig_basic(rName, sharedServiceID1, sharedServiceID2 string) string {
	return fmt.Sprintf(`
resource "firehydrant_service_dependency" "test_service_dependency" {
  service_id = "%s"
  connected_service_id = "%s"
}`, sharedServiceID1, sharedServiceID2)
}

func testAccServiceDependencyResourceConfig_update(rName, sharedServiceID1, sharedServiceID2 string) string {
	return fmt.Sprintf(`
resource "firehydrant_service_dependency" "test_service_dependency" {
  service_id = "%s"
  connected_service_id = "%s"
  notes = "test-notes-%s"
}`, sharedServiceID1, sharedServiceID2, rName)
}
