package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccRoleResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckRoleResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_role.test_role", "id"),
					resource.TestCheckResourceAttr("firehydrant_role.test_role", "name", fmt.Sprintf("test-role-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_role.test_role", "description", "Test role for Terraform"),
					resource.TestCheckResourceAttrSet("firehydrant_role.test_role", "slug"),

					resource.TestCheckResourceAttr("firehydrant_role.test_role", "built_in", "false"),
					resource.TestCheckResourceAttr("firehydrant_role.test_role", "permissions.#", "1"),
				),
			},
			// Test update
			{
				Config: testAccRoleConfig_updated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_role.test_role", "id"),
					resource.TestCheckResourceAttr("firehydrant_role.test_role", "name", fmt.Sprintf("updated-role-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_role.test_role", "description", "Updated test role"),
				),
			},
		},
	})
}

func TestAccRoleResource_withPermissions(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckRoleResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRoleConfig_withPermissions(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_role.test_role", "id"),
					resource.TestCheckResourceAttr("firehydrant_role.test_role", "name", fmt.Sprintf("test-role-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_role.test_role", "permissions.#", "20"),
				),
			},
			// Update permissions
			{
				Config: testAccRoleConfig_withUpdatedPermissions(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_role.test_role", "id"),
					resource.TestCheckResourceAttr("firehydrant_role.test_role", "permissions.#", "19"),
				),
			},
		},
	})
}

func testAccRoleConfig_basic(rName string) string {
	return fmt.Sprintf(`
	resource "firehydrant_role" "test_role" {
		name        = "test-role-%s"
		description = "Test role for Terraform"
		permissions = [
			"read_users"
		]
	}
	`, rName)
}

func testAccRoleConfig_updated(rName string) string {
	return fmt.Sprintf(`
	resource "firehydrant_role" "test_role" {
		name        = "updated-role-%s"
		description = "Updated test role"
		permissions = [
			"read_users"
		]
	}
	`, rName)
}

func testAccRoleConfig_withPermissions(rName string) string {
	return fmt.Sprintf(`
	resource "firehydrant_role" "test_role" {
		name        = "test-role-%s"
		description = "Test role with permissions"
		permissions = [
			"read_alerts",
			"create_alerts",
			"read_escalation_policies",
			"read_on_call_schedules",
			"read_teams",
			"read_users",
			"read_incident_settings",
			"read_integrations",
			"read_incidents",
			"read_webhooks",
			"read_runbooks",
			"read_status_templates",
			"read_audiences",
			"read_change_events",
			"read_organization_settings",
			"read_service_catalog",
			"read_analytics",
			"read_alert_rules",
			"read_call_routes",
			"read_support_hours"
		]
	}
	`, rName)
}

func testAccRoleConfig_withUpdatedPermissions(rName string) string {
	return fmt.Sprintf(`
	resource "firehydrant_role" "test_role" {
		name        = "test-role-%s"
		description = "Test role with updated permissions"
		permissions = [
			"read_alerts",
			"read_alert_rules",
			"read_call_routes",
			"read_escalation_policies",
			"read_on_call_schedules",
			"read_incident_settings",
			"read_incidents",
			"read_integrations",
			"read_service_catalog",
			"read_support_hours",
			"read_status_templates",
			"read_audiences",
			"read_change_events",
			"read_organization_settings",
			"read_runbooks",
			"read_webhooks",
			"read_analytics",
			"read_teams",
			"read_users"
		]
	}
	`, rName)
}

// Offline test with mock server for faster unit testing
func offlineRoleMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`{
			"id": "role-123",
			"name": "Test Role",
			"slug": "test-role",
			"description": "A test role",

			"built_in": false,
			"read_only": false,
			"permissions": [
				{
					"slug": "read_incidents",
					"display_name": "Read Incidents",
					"description": "Can view incidents",
					"available": true
				}
			],
			"created_at": "2025-01-01T00:00:00Z",
			"updated_at": "2025-01-01T00:00:00Z"
		}`))
	}))
}

func TestOfflineRoleRead(t *testing.T) {
	ts := offlineRoleMockServer()
	defer ts.Close()

	c, err := firehydrant.NewRestClient("test-token", firehydrant.WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Error initializing API client: %s", err.Error())
	}

	r := schema.TestResourceDataRaw(t, resourceRole().Schema, map[string]interface{}{
		"name":        "Test Role",
		"description": "A test role",
	})

	d := readResourceFireHydrantRole(context.Background(), r, c)
	if d.HasError() {
		t.Fatalf("Error reading role: %v", d)
	}

	// Verify the role data was set correctly
	if r.Get("name").(string) != "Test Role" {
		t.Fatalf("Expected name 'Test Role', got %s", r.Get("name").(string))
	}
}
