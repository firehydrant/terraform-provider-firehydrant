package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func defaultProviderFactories() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"firehydrant": func() (*schema.Provider, error) {
			return Provider(), nil
		},
	}
}

func TestAccService(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testServiceDoesNotExist("firehydrant_service.terraform-acceptance-test-service"),
		Steps: []resource.TestStep{
			{
				Config: testServiceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testServiceExists("firehydrant_service.terraform-acceptance-test-service"),
					resource.TestCheckResourceAttr("firehydrant_service.terraform-acceptance-test-service", "name", fmt.Sprintf("test-service-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_service.terraform-acceptance-test-service", "description", fmt.Sprintf("%s description", rName)),
				),
			},
			{
				Config: testServiceConfig(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testServiceExists("firehydrant_service.terraform-acceptance-test-service"),
					resource.TestCheckResourceAttr("firehydrant_service.terraform-acceptance-test-service", "name", fmt.Sprintf("test-service-%s", rNameUpdated)),
					resource.TestCheckResourceAttr("firehydrant_service.terraform-acceptance-test-service", "description", fmt.Sprintf("%s description", rNameUpdated)),
				),
			},
			{
				Config: testServiceDataSourceConfig(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.firehydrant_services.services", "services.0.name", fmt.Sprintf("test-service-%s", rNameUpdated)),
					resource.TestCheckResourceAttr("data.firehydrant_services.services", "services.0.service_tier", "5"),
				),
			},
		},
	})
}

func testFireHydrantIsSetup(t *testing.T) {
	if v := os.Getenv("FIREHYDRANT_API_KEY"); v == "" {
		t.Fatalf("Missing required environment variable: %s", "FIREHYDRANT_API_KEY")
	}
}

const testServiceConfigTemplate = `
resource "firehydrant_service" "terraform-acceptance-test-service" {
	name = "test-service-%s"
	description = "%s description"
	labels = {
		key1 = "value1"
	}
	service_tier = 5
}
`

func testServiceDataSourceConfig(rName string) string {
	return fmt.Sprintf(testServiceDataSourceConfigTemplate, rName, rName, rName)
}

const testServiceDataSourceConfigTemplate = `
resource "firehydrant_service" "terraform-acceptance-test-service" {
	name = "test-service-%s"
	description = "%s description"
	labels = {
		key1 = "value1"
	}
	service_tier = 5
}

data "firehydrant_services" "services" {
	query = "%s"
	labels = {
		key1 = "value1"
	}
}

output "services" {
	value = data.firehydrant_services.services
}
`

func testServiceConfig(rName string) string {
	return fmt.Sprintf(testServiceConfigTemplate, rName, rName)
}

func testServiceExists(resourceName string) resource.TestCheckFunc {
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

		svc, err := c.Services().Get(context.TODO(), rs.Primary.ID)
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

func testServiceDoesNotExist(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[resourceName]

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID was not set")
		}

		c, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		svc, err := c.Services().Get(context.TODO(), rs.Primary.ID)
		if svc != nil {
			return fmt.Errorf("The service existed, when it should not")
		}

		if _, isNotFound := err.(firehydrant.NotFound); !isNotFound {
			return err
		}

		return nil
	}
}
