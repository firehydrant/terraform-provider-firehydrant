package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	fhsdk "github.com/firehydrant/firehydrant-go-sdk"
	"github.com/firehydrant/firehydrant-go-sdk/models/components"

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

func TestAccTeamResource_withMembership(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	existingUser := os.Getenv("EXISTING_USER_EMAIL")

	if existingUser == "" {
		existingUser = "local@firehydrant.io"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckTeamResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTeamResourceConfig_withMembership(rName, existingUser),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTeamResourceExistsWithAttributes_basic("firehydrant_team.test_team"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_team.test_team", "name", fmt.Sprintf("test-team-%s", rName)),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "memberships.0.user_id"),
					resource.TestCheckResourceAttrSet("firehydrant_team.test_team", "memberships.0.default_incident_role_id"),
				),
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

		client := fhsdk.New(fhsdk.WithSecurity(components.Security{APIKey: os.Getenv("FIREHYDRANT_API_KEY")}))

		teamResponse, err := client.Teams.GetTeam(context.TODO(), teamResource.Primary.ID, nil)
		if err != nil {
			return err
		}

		expected, got := teamResource.Primary.Attributes["name"], *teamResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		if *teamResponse.Description != "" {
			return fmt.Errorf("Unexpected description. Expected no description, got: %s", *teamResponse.Description)
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

		client := fhsdk.New(fhsdk.WithSecurity(components.Security{APIKey: os.Getenv("FIREHYDRANT_API_KEY")}))

		teamResponse, err := client.Teams.GetTeam(context.TODO(), teamResource.Primary.ID, nil)
		if err != nil {
			return err
		}

		expected, got := teamResource.Primary.Attributes["name"], *teamResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = teamResource.Primary.Attributes["description"], *teamResponse.Description
		if expected != got {
			return fmt.Errorf("Unexpected description. Expected: %s, got: %s", expected, got)
		}

		return nil
	}
}

func testAccCheckTeamResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := fhsdk.New(fhsdk.WithSecurity(components.Security{APIKey: os.Getenv("FIREHYDRANT_API_KEY")}))

		for _, teamResource := range s.RootModule().Resources {
			if teamResource.Type != "firehydrant_team" {
				continue
			}

			if teamResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			_, err := client.Teams.GetTeam(context.TODO(), teamResource.Primary.ID, nil)
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

func testAccTeamResourceConfig_withMembership(rName string, existingUser string) string {
	return fmt.Sprintf(`
data "firehydrant_user" "test_user" {
	email = "%s"
}

resource "firehydrant_incident_role" "test_incident_role" {
	name    = "test-incident-role-%s"
	summary = "test-summary-%s"
}

resource "firehydrant_team" "test_team" {
	name = "test-team-%s"

	memberships {
		user_id                  = data.firehydrant_user.test_user.id
		default_incident_role_id = resource.firehydrant_incident_role.test_incident_role.id
	}
}`, existingUser, rName, rName, rName)
}
