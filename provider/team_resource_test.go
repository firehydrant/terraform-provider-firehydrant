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
				),
			},
			{
				Config: testAccTeamResourceConfig_basic(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTeamResourceExistsWithAttributes_basic("firehydrant_team.test_team"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "name", fmt.Sprintf("test-team-%s", rNameUpdated)),
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

		return nil
	}
}

func testAccCheckTeamResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		for _, teamResource := range s.RootModule().Resources {
			if teamResource.Type != "firehydrant_team" {
				continue
			}

			if teamResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			_, err := client.GetTeam(context.TODO(), teamResource.Primary.ID)
			if err == nil {
				return fmt.Errorf("Team %s still exists", teamResource.Primary.ID)
			}
		}

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
resource "firehydrant_team" "test_team" {
  name        = "test-team-%s"
  description = "test-description-%s"
}`, rName, rName)
}
