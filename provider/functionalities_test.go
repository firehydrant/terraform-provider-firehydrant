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

func TestAccFunctionalities(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testFunctionalityDoesNotExist("firehydrant_functionality.terraform-acceptance-test-functionality"),
		Steps: []resource.TestStep{
			{
				Config: testFunctionalityConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testFunctionalityExists("firehydrant_functionality.terraform-acceptance-test-functionality"),
					resource.TestCheckResourceAttr("firehydrant_functionality.terraform-acceptance-test-functionality", "name", rName),
				),
			},
			{
				Config: testFunctionalityConfig(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testFunctionalityExists("firehydrant_functionality.terraform-acceptance-test-functionality"),
					resource.TestCheckResourceAttr("firehydrant_functionality.terraform-acceptance-test-functionality", "name", rNameUpdated),
					resource.TestCheckResourceAttr("firehydrant_functionality.terraform-acceptance-test-functionality", "services.#", "0"),
				),
			},
			{
				Config: testFunctionalityConfigWithService(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testFunctionalityExists("firehydrant_functionality.terraform-acceptance-test-functionality"),
					resource.TestCheckResourceAttr("firehydrant_functionality.terraform-acceptance-test-functionality", "name", rNameUpdated),
					resource.TestCheckResourceAttr("firehydrant_functionality.terraform-acceptance-test-functionality", "services.#", "1"),
					resource.TestCheckResourceAttrSet("firehydrant_functionality.terraform-acceptance-test-functionality", "services.0.id"),
					resource.TestCheckResourceAttr("firehydrant_functionality.terraform-acceptance-test-functionality", "services.0.name", "test service from terraform"),
				),
			},
		},
	})
}

const testFunctionalityConfigTemplate = `
resource "firehydrant_functionality" "terraform-acceptance-test-functionality" {
	name = "%s"
}
`

func testFunctionalityConfig(rName string) string {
	return fmt.Sprintf(testFunctionalityConfigTemplate, rName)
}

const testFunctionalityWithService = `
resource "firehydrant_service" "service" {
	name = "test service from terraform"
}

resource "firehydrant_functionality" "terraform-acceptance-test-functionality" {
	name = "%s"

	services {
		id = firehydrant_service.service.id
	}
}
`

func testFunctionalityConfigWithService(rName string) string {
	return fmt.Sprintf(testFunctionalityWithService, rName)
}

func testFunctionalityExists(resourceName string) resource.TestCheckFunc {
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

		svc, err := c.GetFunctionality(context.TODO(), rs.Primary.ID)
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

func testFunctionalityDoesNotExist(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[resourceName]

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID was not set")
		}

		c, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		svc, err := c.GetFunctionality(context.TODO(), rs.Primary.ID)
		if svc != nil {
			return fmt.Errorf("The functionality existed, when it should not")
		}

		if _, isNotFound := err.(firehydrant.NotFound); !isNotFound {
			return err
		}

		return nil
	}
}
