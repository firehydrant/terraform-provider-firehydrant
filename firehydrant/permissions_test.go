package firehydrant

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func expectedPermissionsList() []Permission {
	return []Permission{
		{
			Slug:                "read_incidents",
			DisplayName:         "Read Incidents",
			Description:         "Can view incident details",
			CategoryDisplayName: "Incidents",
			CategorySlug:        "incidents",
			ParentSlug:          "",
			Available:           true,
			DependencySlugs:     []string{},
		},
		{
			Slug:                "manage_incidents",
			DisplayName:         "Write Incidents",
			Description:         "Can create and modify incidents",
			CategoryDisplayName: "Incidents",
			CategorySlug:        "incidents",
			ParentSlug:          "read_incidents",
			Available:           true,
			DependencySlugs:     []string{"read_incidents"},
		},
		{
			Slug:                "manage_roles",
			DisplayName:         "Manage Roles",
			Description:         "Full access to role management",
			CategoryDisplayName: "Roles & Permissions",
			CategorySlug:        "roles",
			ParentSlug:          "",
			Available:           false,
			DependencySlugs:     []string{"read_roles", "manage_roles"},
		},
	}
}

func expectedPermissionsListResponseJSON() string {
	return `{
		"data": [
			{
				"slug": "read_incidents",
				"display_name": "Read Incidents",
				"description": "Can view incident details",
				"category_display_name": "Incidents",
				"category_slug": "incidents",
				"parent_slug": "",
				"available": true,
				"dependency_slugs": []
			},
			{
				"slug": "manage_incidents",
				"display_name": "Write Incidents", 
				"description": "Can create and modify incidents",
				"category_display_name": "Incidents",
				"category_slug": "incidents",
				"parent_slug": "read_incidents",
				"available": true,
				"dependency_slugs": ["read_incidents"]
			},
			{
				"slug": "manage_roles",
				"display_name": "Manage Roles",
				"description": "Full access to role management",
				"category_display_name": "Roles & Permissions",
				"category_slug": "roles",
				"parent_slug": "",
				"available": false,
				"dependency_slugs": ["read_roles", "manage_roles"]
			}
		]
	}`
}

func TestPermissionsList(t *testing.T) {
	var requestPath string
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPath = req.URL.Path
		w.Write([]byte(expectedPermissionsListResponseJSON()))
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	res, err := c.Permissions().List(context.TODO())
	if err != nil {
		t.Fatalf("error listing permissions: %s", err.Error())
	}

	if expected := "/permissions"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}

	expectedPermissions := expectedPermissionsList()
	if !reflect.DeepEqual(expectedPermissions, res.Data) {
		t.Fatalf("permissions mismatch: expected '%+v', got: '%+v'", expectedPermissions, res.Data)
	}
}
