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

func TestAccTeams(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testTeamDoesNotExist("firehydrant_team.terraform-acceptance-test-team"),
		Steps: []resource.TestStep{
			{
				Config: testTeamConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testTeamExists("firehydrant_team.terraform-acceptance-test-team"),
					resource.TestCheckResourceAttr("firehydrant_team.terraform-acceptance-test-team", "name", fmt.Sprintf("test-team-%s", rName)),
				),
			},
			{
				Config: testTeamConfig(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testTeamExists("firehydrant_team.terraform-acceptance-test-team"),
					resource.TestCheckResourceAttr("firehydrant_team.terraform-acceptance-test-team", "name", fmt.Sprintf("test-team-%s", rNameUpdated)),
					resource.TestCheckResourceAttr("firehydrant_team.terraform-acceptance-test-team", "services.#", "0"),
				),
			},
			{
				Config: testTeamConfigWithService(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testTeamExists("firehydrant_team.terraform-acceptance-test-team"),
					resource.TestCheckResourceAttr("firehydrant_team.terraform-acceptance-test-team", "name", fmt.Sprintf("test-team-%s", rNameUpdated)),
					resource.TestCheckResourceAttr("firehydrant_team.terraform-acceptance-test-team", "services.#", "1"),
					resource.TestCheckResourceAttrSet("firehydrant_team.terraform-acceptance-test-team", "services.0.id"),
					resource.TestCheckResourceAttr("firehydrant_team.terraform-acceptance-test-team", "services.0.name", fmt.Sprintf("test-service-%s", rNameUpdated)),
				),
			},
		},
	})
}

const testTeamConfigTemplate = `
resource "firehydrant_team" "terraform-acceptance-test-team" {
	name = "test-team-%s"
}
`

func testTeamConfig(rName string) string {
	return fmt.Sprintf(testTeamConfigTemplate, rName)
}

const testTeamWithService = `
resource "firehydrant_service" "service" {
	name = "test-service-%s"
}

resource "firehydrant_team" "terraform-acceptance-test-team" {
	name = "test-team-%s"

	services {
		id = firehydrant_service.service.id
	}
}
`

func testTeamConfigWithService(rName string) string {
	return fmt.Sprintf(testTeamWithService, rName, rName)
}

func testTeamExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID was not set")
		}

		c, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		svc, err := c.GetTeam(context.TODO(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if expected, got := rs.Primary.Attributes["name"], svc.Name; expected != got {
			return fmt.Errorf("Expected name %s, got %s", expected, got)
		}

		if expected, got := rs.Primary.Attributes["description"], svc.Description; expected != got {
			return fmt.Errorf("Expected description %s, got %s", expected, got)
		}

		return nil
	}
}

func testTeamDoesNotExist(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[resourceName]

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID was not set")
		}

		_, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		// TODO: Archives dont hide teams from the details endpoint
		// svc, err := c.GetTeam(context.TODO(), rs.Primary.ID)
		// if svc != nil {
		// 	return fmt.Errorf("The team existed, when it should not")
		// }

		// if _, isNotFound := err.(firehydrant.NotFound); !isNotFound {
		// 	return err
		// }

		return nil
	}
}
