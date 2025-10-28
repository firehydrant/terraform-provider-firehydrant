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

func TestAccRoleDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccBuiltInRoleDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Test data source lookup by slug
					resource.TestCheckResourceAttrSet("data.firehydrant_role.member", "id"),
					resource.TestCheckResourceAttr("data.firehydrant_role.member", "slug", "member"),
					resource.TestCheckResourceAttrSet("data.firehydrant_role.member", "name"),
					resource.TestCheckResourceAttrSet("data.firehydrant_role.member", "description"),

					resource.TestCheckResourceAttrSet("data.firehydrant_role.member", "permissions.#"),
				),
			},
		},
	})
}

func TestAccRoleDataSource_builtIn(t *testing.T) {
	// Test looking up a built-in role that should always exist
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccBuiltInRoleDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_role.member", "id"),
					resource.TestCheckResourceAttr("data.firehydrant_role.member", "slug", "member"),
					resource.TestCheckResourceAttrSet("data.firehydrant_role.member", "permissions.#"),
				),
			},
		},
	})
}

func testAccBuiltInRoleDataSourceConfig() string {
	return `
data "firehydrant_role" "member" {
	slug = "member"
}
`
}

/** Offline Unit Tests *******************************************************************************************/

func TestOfflineRoleDataSource(t *testing.T) {
	// Mock server that returns roles list for slug lookup
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/roles":
			// Mock response for List() call used in slug lookup
			w.Write([]byte(`{
				"data": [
					{
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
								"available": true,
								"dependency_slugs": []
							}
						],
						"created_at": "2025-01-01T00:00:00Z",
						"updated_at": "2025-01-01T00:00:00Z"
					}
				],
				"pagination": {
					"count": 1,
					"page": 1,
					"items": 1,
					"pages": 1
				}
			}`))
		case "/roles/role-123":
			// Mock response for Get() call used in ID lookup
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
						"available": true,
						"dependency_slugs": []
					}
				],
				"created_at": "2025-01-01T00:00:00Z",
				"updated_at": "2025-01-01T00:00:00Z"
			}`))
		default:
			http.NotFound(w, req)
		}
	}))
	defer ts.Close()

	c, err := firehydrant.NewRestClient("test-token", firehydrant.WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Error initializing API client: %s", err.Error())
	}

	// Test ID lookup
	t.Run("LookupByID", func(t *testing.T) {
		r := schema.TestResourceDataRaw(t, dataSourceRole().Schema, map[string]interface{}{
			"id": "role-123",
		})

		d := dataFireHydrantRole(context.Background(), r, c)
		if d.HasError() {
			t.Fatalf("Error reading role by ID: %v", d)
		}

		if r.Get("name").(string) != "Test Role" {
			t.Fatalf("Expected name 'Test Role', got %s", r.Get("name").(string))
		}
		if r.Get("slug").(string) != "test-role" {
			t.Fatalf("Expected slug 'test-role', got %s", r.Get("slug").(string))
		}
	})

	// Test slug lookup
	t.Run("LookupBySlug", func(t *testing.T) {
		r := schema.TestResourceDataRaw(t, dataSourceRole().Schema, map[string]interface{}{
			"slug": "test-role",
		})

		d := dataFireHydrantRole(context.Background(), r, c)
		if d.HasError() {
			t.Fatalf("Error reading role by slug: %v", d)
		}

		if r.Get("name").(string) != "Test Role" {
			t.Fatalf("Expected name 'Test Role', got %s", r.Get("name").(string))
		}
		if r.Get("id").(string) != "role-123" {
			t.Fatalf("Expected ID 'role-123', got %s", r.Get("id").(string))
		}
	})

	// Test missing slug
	t.Run("SlugNotFound", func(t *testing.T) {
		r := schema.TestResourceDataRaw(t, dataSourceRole().Schema, map[string]interface{}{
			"slug": "nonexistent-role",
		})

		d := dataFireHydrantRole(context.Background(), r, c)
		if !d.HasError() {
			t.Fatalf("Expected error for nonexistent slug, but got none")
		}
	})
}
