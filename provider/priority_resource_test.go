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

func TestAccPriorityResource_basic(t *testing.T) {
	rSlug := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckPriorityResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPriorityResourceConfig_basic(rSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPriorityResourceExistsWithAttributes_basic("firehydrant_priority.test_priority"),
					resource.TestCheckResourceAttrSet("firehydrant_priority.test_priority", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_priority.test_priority", "slug", fmt.Sprintf("test-priority-%s", rSlug)),
				),
			},
		},
	})
}

func TestAccPriorityResource_update(t *testing.T) {
	rSlug := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rSlugUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckPriorityResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPriorityResourceConfig_basic(rSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPriorityResourceExistsWithAttributes_basic("firehydrant_priority.test_priority"),
					resource.TestCheckResourceAttrSet("firehydrant_priority.test_priority", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_priority.test_priority", "slug", fmt.Sprintf("test-priority-%s", rSlug)),
				),
			},
			{
				Config: testAccPriorityResourceConfig_update(rSlugUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPriorityResourceExistsWithAttributes_update("firehydrant_priority.test_priority"),
					resource.TestCheckResourceAttrSet("firehydrant_priority.test_priority", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_priority.test_priority", "slug", fmt.Sprintf("test-priority-%s", rSlugUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_priority.test_priority", "description", fmt.Sprintf("test-description-%s", rSlugUpdated)),
				),
			},
			{
				Config: testAccPriorityResourceConfig_basic(rSlugUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckPriorityResourceExistsWithAttributes_basic("firehydrant_priority.test_priority"),
					resource.TestCheckResourceAttrSet("firehydrant_priority.test_priority", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_priority.test_priority", "slug", fmt.Sprintf("test-priority-%s", rSlugUpdated)),
				),
			},
		},
	})
}

func TestAccPriorityResourceImport_basic(t *testing.T) {
	rSlug := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPriorityResourceConfig_basic(rSlug),
			},
			{
				ResourceName:      "firehydrant_priority.test_priority",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPriorityResourceExistsWithAttributes_basic(resourceSlug string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		priorityResource, ok := s.RootModule().Resources[resourceSlug]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceSlug)
		}
		if priorityResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		priorityResponse, err := client.GetPriority(context.TODO(), priorityResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := priorityResource.Primary.Attributes["slug"], priorityResponse.Slug
		if expected != got {
			return fmt.Errorf("Unexpected slug. Expected: %s, got: %s", expected, got)
		}

		if priorityResponse.Description != "" {
			return fmt.Errorf("Unexpected description. Expected no description, got: %s", priorityResponse.Description)
		}

		return nil
	}
}

func testAccCheckPriorityResourceExistsWithAttributes_update(resourceSlug string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		priorityResource, ok := s.RootModule().Resources[resourceSlug]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceSlug)
		}
		if priorityResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		priorityResponse, err := client.GetPriority(context.TODO(), priorityResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := priorityResource.Primary.Attributes["slug"], priorityResponse.Slug
		if expected != got {
			return fmt.Errorf("Unexpected slug. Expected: %s, got: %s", expected, got)
		}

		expected, got = priorityResource.Primary.Attributes["description"], priorityResponse.Description
		if expected != got {
			return fmt.Errorf("Unexpected description. Expected: %s, got: %s", expected, got)
		}

		return nil
	}
}

func testAccCheckPriorityResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		for _, priorityResource := range s.RootModule().Resources {
			if priorityResource.Type != "firehydrant_priority" {
				continue
			}

			if priorityResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			_, err := client.GetPriority(context.TODO(), priorityResource.Primary.ID)
			if err == nil {
				return fmt.Errorf("Priority %s still exists", priorityResource.Primary.ID)
			}
		}

		return nil
	}
}

func testAccPriorityResourceConfig_basic(rSlug string) string {
	return fmt.Sprintf(`
resource "firehydrant_priority" "test_priority" {
  slug = "test-priority-%s"
}`, rSlug)
}

func testAccPriorityResourceConfig_update(rSlug string) string {
	return fmt.Sprintf(`
resource "firehydrant_priority" "test_priority" {
  slug        = "test-priority-%s"
  description = "test-description-%s"
}`, rSlug, rSlug)
}
