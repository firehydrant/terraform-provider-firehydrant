package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRunbookActionDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRunbookActionDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_runbook_action.test_runbook_action", "id"),
					resource.TestCheckResourceAttrSet("data.firehydrant_runbook_action.test_runbook_action", "name"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_runbook_action.test_runbook_action", "slug", "create_incident_channel"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_runbook_action.test_runbook_action", "integration_slug", "slack"),
				),
			},
		},
	})
}

func TestAccRunbookActionDataSource_multipleActionsForSlug(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccRunbookActionDataSourceConfig_multipleActionsForSlug(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_runbook_action.test_runbook_action", "id"),
					resource.TestCheckResourceAttrSet("data.firehydrant_runbook_action.test_runbook_action", "name"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_runbook_action.test_runbook_action", "slug", "create_incident_issue"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_runbook_action.test_runbook_action", "integration_slug", "shortcut"),
				),
			},
		},
	})
}

func TestAccRunbookActionDataSource_validateSchemaAttributesSlug(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccRunbookActionDataSourceConfig_slugInvalid(),
				ExpectError: regexp.MustCompile(`expected slug to be one of`),
			},
		},
	})
}

func TestAccRunbookActionDataSource_validateSchemaAttributesIntegrationSlug(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccRunbookActionDataSourceConfig_integrationSlugInvalid(),
				ExpectError: regexp.MustCompile(`expected integration_slug to be one of`),
			},
		},
	})
}

func testAccRunbookActionDataSourceConfig_basic() string {
	return fmt.Sprintf(`
data "firehydrant_runbook_action" "test_runbook_action" {
  integration_slug = "slack"
  slug             = "create_incident_channel"
  type             = "incident"
}`)
}

func testAccRunbookActionDataSourceConfig_multipleActionsForSlug() string {
	return fmt.Sprintf(`
data "firehydrant_runbook_action" "test_runbook_action" {
  integration_slug = "shortcut"
  slug             = "create_incident_issue"
  type             = "incident"
}`)
}

func testAccRunbookActionDataSourceConfig_slugInvalid() string {
	return fmt.Sprintf(`
data "firehydrant_runbook_action" "test_runbook_action" {
  integration_slug = "shortcut"
  slug             = "slug_invalid"
  type             = "incident"
}`)
}

func testAccRunbookActionDataSourceConfig_integrationSlugInvalid() string {
	return fmt.Sprintf(`
data "firehydrant_runbook_action" "test_runbook_action" {
  integration_slug = "integration_slug_invalid"
  slug             = "create_incident_issue"
  type             = "incident"
}`)
}
