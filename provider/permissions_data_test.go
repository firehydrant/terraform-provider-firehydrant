package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccPermissionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccPermissionsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_permissions.all", "id"),
					resource.TestCheckResourceAttrSet("data.firehydrant_permissions.all", "permissions.#"),
					// Check that we have at least some permissions (validate against real staging data)
					resource.TestCheckResourceAttr("data.firehydrant_permissions.all", "permissions.0.slug", "create_alerts"),
					resource.TestCheckResourceAttr("data.firehydrant_permissions.all", "permissions.0.display_name", "Create Alerts"),
					resource.TestCheckResourceAttr("data.firehydrant_permissions.all", "permissions.0.category_slug", "alerts_oncall"),
					// Verify the structure of permissions
					resource.TestCheckResourceAttrSet("data.firehydrant_permissions.all", "permissions.0.description"),
					resource.TestCheckResourceAttrSet("data.firehydrant_permissions.all", "permissions.0.category_display_name"),
					// Check second permission
					resource.TestCheckResourceAttr("data.firehydrant_permissions.all", "permissions.1.slug", "respond_to_alerts"),
					resource.TestCheckResourceAttr("data.firehydrant_permissions.all", "permissions.1.display_name", "Respond To Alerts"),
				),
			},
		},
	})
}

func TestAccUserPermissionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccUserPermissionsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_user_permissions.mine", "id"),
					resource.TestCheckResourceAttrSet("data.firehydrant_user_permissions.mine", "permissions.#"),
				),
			},
		},
	})
}

func TestAccTeamPermissionsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTeamPermissionsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_team_permissions.team_perms", "id"),
					resource.TestCheckResourceAttrSet("data.firehydrant_team_permissions.team_perms", "permissions.#"),
				),
			},
		},
	})
}

func testAccPermissionsDataSourceConfig() string {
	return `
data "firehydrant_permissions" "all" {}
`
}

func testAccUserPermissionsDataSourceConfig() string {
	return `
data "firehydrant_user_permissions" "mine" {}
`
}

func testAccTeamPermissionsDataSourceConfig() string {
	return `
data "firehydrant_team_permissions" "team_perms" {}
`
}

/** Offline Tests with Mock Server **********************************************************************************/

func permissionsMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/permissions":
			w.Write([]byte(`{
				"data": [
					{
						"slug": "incidents.read",
						"display_name": "Read Incidents",
						"description": "Can view incident details",
						"category_display_name": "Incidents",
						"category_slug": "incidents",
						"parent_slug": "",
						"available": true,
						"dependency_slugs": []
					},
					{
						"slug": "incidents.write",
						"display_name": "Write Incidents",
						"description": "Can create and modify incidents",
						"category_display_name": "Incidents",
						"category_slug": "incidents",
						"parent_slug": "incidents.read",
						"available": true,
						"dependency_slugs": ["incidents.read"]
					},
					{
						"slug": "roles.admin",
						"display_name": "Admin Roles",
						"description": "Full access to role management",
						"category_display_name": "Roles",
						"category_slug": "roles",
						"parent_slug": "",
						"available": false,
						"dependency_slugs": ["roles.read", "roles.write"]
					}
				]
			}`))
		case "/permissions/user":
			w.Write([]byte(`{
				"data": ["incidents.read", "incidents.write"]
			}`))
		case "/permissions/team":
			w.Write([]byte(`{
				"data": [
					{
						"slug": "team.manage",
						"display_name": "Manage Team",
						"description": "Can manage team settings",
						"category_display_name": "Teams",
						"category_slug": "teams",
						"parent_slug": "",
						"available": true,
						"dependency_slugs": []
					}
				]
			}`))
		default:
			http.NotFound(w, req)
		}
	}))
}

func TestOfflinePermissionsDataSources(t *testing.T) {
	ts := permissionsMockServer()
	defer ts.Close()

	c, err := firehydrant.NewRestClient("test-token", firehydrant.WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Error initializing API client: %s", err.Error())
	}

	// Test all permissions data source
	t.Run("AllPermissions", func(t *testing.T) {
		r := schema.TestResourceDataRaw(t, dataSourcePermissions().Schema, map[string]interface{}{})

		d := dataFireHydrantPermissions(context.Background(), r, c)
		if d.HasError() {
			t.Fatalf("Error reading permissions: %v", d)
		}

		permissions := r.Get("permissions").([]interface{})
		if len(permissions) != 3 {
			t.Fatalf("Expected 3 permissions, got %d", len(permissions))
		}

		firstPerm := permissions[0].(map[string]interface{})
		if firstPerm["slug"].(string) != "incidents.read" {
			t.Fatalf("Expected first permission slug to be 'incidents.read', got %s", firstPerm["slug"])
		}
		if firstPerm["available"].(bool) != true {
			t.Fatalf("Expected first permission to be available")
		}
	})

	// Test user permissions data source
	t.Run("UserPermissions", func(t *testing.T) {
		r := schema.TestResourceDataRaw(t, dataSourceUserPermissions().Schema, map[string]interface{}{})

		d := dataFireHydrantUserPermissions(context.Background(), r, c)
		if d.HasError() {
			t.Fatalf("Error reading current user permissions: %v", d)
		}

		permissions := r.Get("permissions").(*schema.Set)
		if permissions.Len() != 2 {
			t.Fatalf("Expected 2 permissions, got %d", permissions.Len())
		}

		if !permissions.Contains("incidents.read") {
			t.Fatalf("Expected permissions to contain 'incidents.read'")
		}
	})

	// Test team permissions data source
	t.Run("TeamPermissions", func(t *testing.T) {
		r := schema.TestResourceDataRaw(t, dataSourceTeamPermissions().Schema, map[string]interface{}{})

		d := dataFireHydrantTeamPermissions(context.Background(), r, c)
		if d.HasError() {
			t.Fatalf("Error reading team permissions: %v", d)
		}

		permissions := r.Get("permissions").([]interface{})
		if len(permissions) != 1 {
			t.Fatalf("Expected 1 team permission, got %d", len(permissions))
		}

		teamPerm := permissions[0].(map[string]interface{})
		if teamPerm["slug"].(string) != "team.manage" {
			t.Fatalf("Expected team permission slug to be 'team.manage', got %s", teamPerm["slug"])
		}
	})
}
