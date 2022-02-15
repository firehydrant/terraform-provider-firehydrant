package provider

import (
	"fmt"
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
