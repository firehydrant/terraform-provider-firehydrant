package provider

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSeverityResource_basic(t *testing.T) {
	rSlug := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckSeverityResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSeverityResourceConfig_basic(rSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSeverityResourceExistsWithAttributes_basic("firehydrant_severity.test_severity"),
					resource.TestCheckResourceAttrSet("firehydrant_severity.test_severity", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_severity.test_severity", "slug", fmt.Sprintf("TESTSEVERITY%s", rSlug)),
					resource.TestCheckResourceAttr(
						"firehydrant_severity.test_severity", "type", string(firehydrant.SeverityTypeUnexpectedDowntime)),
				),
			},
		},
	})
}

func TestAccSeverityResource_update(t *testing.T) {
	rSlug := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	rSlugUpdated := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckSeverityResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSeverityResourceConfig_basic(rSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSeverityResourceExistsWithAttributes_basic("firehydrant_severity.test_severity"),
					resource.TestCheckResourceAttrSet("firehydrant_severity.test_severity", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_severity.test_severity", "slug", fmt.Sprintf("TESTSEVERITY%s", rSlug)),
					resource.TestCheckResourceAttr(
						"firehydrant_severity.test_severity", "type", string(firehydrant.SeverityTypeUnexpectedDowntime)),
				),
			},
			{
				Config: testAccSeverityResourceConfig_update(rSlugUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSeverityResourceExistsWithAttributes_update("firehydrant_severity.test_severity"),
					resource.TestCheckResourceAttrSet("firehydrant_severity.test_severity", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_severity.test_severity", "slug", fmt.Sprintf("TESTSEVERITY%s", rSlugUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_severity.test_severity", "description", fmt.Sprintf("test-description-%s", rSlugUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_severity.test_severity", "type", string(firehydrant.SeverityTypeMaintenance)),
				),
			},
			{
				Config: testAccSeverityResourceConfig_basic(rSlugUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSeverityResourceExistsWithAttributes_basic("firehydrant_severity.test_severity"),
					resource.TestCheckResourceAttrSet("firehydrant_severity.test_severity", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_severity.test_severity", "slug", fmt.Sprintf("TESTSEVERITY%s", rSlugUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_severity.test_severity", "type", string(firehydrant.SeverityTypeUnexpectedDowntime)),
				),
			},
		},
	})
}

func TestAccSeverityResource_validateSchemaAttributesSlug(t *testing.T) {
	rSlug := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccSeverityResourceConfig_slugTooLong(rSlug),
				ExpectError: regexp.MustCompile(`expected length of slug to be in the range \(0 - 23\)`),
			},
			{
				Config:      testAccSeverityResourceConfig_slugWithInvalidCharacters(rSlug),
				ExpectError: regexp.MustCompile(`invalid value for slug \(must only include letters and numbers\)`),
			},
		},
	})
}

func TestAccSeverityResourceImport_basic(t *testing.T) {
	rSlug := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSeverityResourceConfig_basic(rSlug),
			},
			{
				ResourceName:      "firehydrant_severity.test_severity",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSeverityResourceImport_allAttributes(t *testing.T) {
	rSlug := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccSeverityResourceConfig_update(rSlug),
			},
			{
				ResourceName:      "firehydrant_severity.test_severity",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckSeverityResourceExistsWithAttributes_basic(resourceSlug string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		severityResource, ok := s.RootModule().Resources[resourceSlug]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceSlug)
		}
		if severityResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		severityResponse, err := client.Severities().Get(context.TODO(), severityResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := severityResource.Primary.Attributes["slug"], severityResponse.Slug
		if expected != got {
			return fmt.Errorf("Unexpected slug. Expected: %s, got: %s", expected, got)
		}

		if severityResponse.Description != "" {
			return fmt.Errorf("Unexpected description. Expected no description, got: %s", severityResponse.Description)
		}

		if severityResponse.Type != string(firehydrant.SeverityTypeUnexpectedDowntime) {
			return fmt.Errorf("Unexpected type. Expected default type of %s, got: %s", string(firehydrant.SeverityTypeUnexpectedDowntime), severityResponse.Type)
		}
		expected, got = severityResource.Primary.Attributes["type"], severityResponse.Type
		if expected != got {
			return fmt.Errorf("Unexpected type. Expected: %s, got: %s", expected, got)
		}

		return nil
	}
}

func testAccCheckSeverityResourceExistsWithAttributes_update(resourceSlug string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		severityResource, ok := s.RootModule().Resources[resourceSlug]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceSlug)
		}
		if severityResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		severityResponse, err := client.Severities().Get(context.TODO(), severityResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := severityResource.Primary.Attributes["slug"], severityResponse.Slug
		if expected != got {
			return fmt.Errorf("Unexpected slug. Expected: %s, got: %s", expected, got)
		}

		expected, got = severityResource.Primary.Attributes["description"], severityResponse.Description
		if expected != got {
			return fmt.Errorf("Unexpected description. Expected: %s, got: %s", expected, got)
		}

		expected, got = severityResource.Primary.Attributes["type"], severityResponse.Type
		if expected != got {
			return fmt.Errorf("Unexpected type. Expected: %s, got: %s", expected, got)
		}

		return nil
	}
}

func testAccCheckSeverityResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_severity" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			_, err := client.Severities().Get(context.TODO(), stateResource.Primary.ID)
			if err == nil {
				return fmt.Errorf("Severity %s still exists", stateResource.Primary.ID)
			}
		}

		return nil
	}
}

func testAccSeverityResourceConfig_basic(rSlug string) string {
	return fmt.Sprintf(`
resource "firehydrant_severity" "test_severity" {
  slug = "TESTSEVERITY%s"
}`, rSlug)
}

func testAccSeverityResourceConfig_update(rSlug string) string {
	return fmt.Sprintf(`
resource "firehydrant_severity" "test_severity" {
  slug        = "TESTSEVERITY%s"
  description = "test-description-%s"
  type        = "maintenance"
}`, rSlug, rSlug)
}

func testAccSeverityResourceConfig_slugTooLong(rSlug string) string {
	return fmt.Sprintf(`
resource "firehydrant_severity" "test_severity" {
  slug = "THISSLUGISWAYTOOLONG%s"
}`, rSlug)
}

func testAccSeverityResourceConfig_slugWithInvalidCharacters(rSlug string) string {
	return fmt.Sprintf(`
resource "firehydrant_severity" "test_severity" {
  slug = "INVALID-SLUG%s"
}`, rSlug)
}
