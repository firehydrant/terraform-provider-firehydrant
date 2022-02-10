package provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSeverities(t *testing.T) {
	rName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testSeverityDoesNotExist("firehydrant_severity.terraform-acceptance-test-severity"),
		Steps: []resource.TestStep{
			{
				Config: testSeverityConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testSeverityExists("firehydrant_severity.terraform-acceptance-test-severity"),
					resource.TestCheckResourceAttr("firehydrant_severity.terraform-acceptance-test-severity", "slug", fmt.Sprintf("TESTSEVERITY%s", rName)),
				),
			},
			// TODO(bobbytables): Updating severities in Terraform is currently problematic because FireHydrant uses
			// slugs as the IDs and those are stored in Terraform state as the resource ID. Since updates can change a slug but terraform updates _wont_
			// update the resource with the new slug as the ID, it's technically not possible to perform a slug update in Terraform against FireHydrant.
			// {
			// 	Config: testSeverityConfig(rNameUpdated),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testSeverityExists("firehydrant_severity.terraform-acceptance-test-severity"),
			// 		resource.TestCheckResourceAttr("firehydrant_severity.terraform-acceptance-test-severity", "slug", strings.ToUpper(rNameUpdated)),
			// 	),
			// },
		},
	})
}

const testSeverityConfigTemplate = `
resource "firehydrant_severity" "terraform-acceptance-test-severity" {
	slug = "TESTSEVERITY%s"
}
`

func testSeverityConfig(rName string) string {
	return fmt.Sprintf(testSeverityConfigTemplate, strings.ToUpper(rName))
}

func testSeverityExists(resourceName string) resource.TestCheckFunc {
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

		svc, err := c.GetSeverity(context.TODO(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if expected, got := rs.Primary.Attributes["slug"], svc.Slug; expected != got {
			return fmt.Errorf("Expected slug %s, got %s", expected, got)
		}

		if expected, got := rs.Primary.Attributes["description"], svc.Description; expected != got {
			return fmt.Errorf("Expected description %s, got %s", expected, got)
		}

		return nil
	}
}

func testSeverityDoesNotExist(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[resourceName]

		if rs.Primary.ID == "" {
			return fmt.Errorf("ID was not set")
		}

		_, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		// TODO: Archives dont hide severitys from the details endpoint
		// svc, err := c.GetSeverity(context.TODO(), rs.Primary.ID)
		// if svc != nil {
		// 	return fmt.Errorf("The severity existed, when it should not")
		// }

		// if _, isNotFound := err.(firehydrant.NotFound); !isNotFound {
		// 	return err
		// }

		return nil
	}
}
