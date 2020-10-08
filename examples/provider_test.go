package examples

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/firehydrant/terraform-provider-firehydrant/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func defaultProviderFactories() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"firehydrant": func() (*schema.Provider, error) {
			return provider.Provider(), nil
		},
	}
}

func TestAccService(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testServiceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testServiceExists("firehydrant_service.terraform-acceptance-test-service"),
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
  name = "%s"
}
`

func testServiceConfig(rName string) string {
	return fmt.Sprintf(testServiceConfigTemplate, rName)
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

		svc, err := c.GetService(context.TODO(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if expected, got := rs.Primary.Attributes["name"], svc.Name; expected != got {
			return fmt.Errorf("Expected name %s, got %s", expected, got)
		}

		return nil
	}
}
