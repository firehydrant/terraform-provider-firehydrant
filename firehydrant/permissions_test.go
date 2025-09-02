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
			Slug:                "incidents.read",
			DisplayName:         "Read Incidents",
			Description:         "Can view incident details",
			CategoryDisplayName: "Incidents",
			CategorySlug:        "incidents",
			ParentSlug:          "",
			Available:           true,
			DependencySlugs:     []string{},
		},
		{
			Slug:                "incidents.write",
			DisplayName:         "Write Incidents",
			Description:         "Can create and modify incidents",
			CategoryDisplayName: "Incidents",
			CategorySlug:        "incidents",
			ParentSlug:          "incidents.read",
			Available:           true,
			DependencySlugs:     []string{"incidents.read"},
		},
		{
			Slug:                "roles.admin",
			DisplayName:         "Admin Roles",
			Description:         "Full access to role management",
			CategoryDisplayName: "Roles & Permissions",
			CategorySlug:        "roles",
			ParentSlug:          "",
			Available:           false,
			DependencySlugs:     []string{"roles.read", "roles.write"},
		},
	}
}

func expectedPermissionsListResponseJSON() string {
	return `{
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
				"category_display_name": "Roles & Permissions",
				"category_slug": "roles",
				"parent_slug": "",
				"available": false,
				"dependency_slugs": ["roles.read", "roles.write"]
			}
		]
	}`
}

func expectedCurrentUserPermissionsJSON() string {
	return `{
		"data": ["incidents.read", "incidents.write"]
	}`
}

func expectedTeamPermissionsResponseJSON() string {
	return `{
		"data": [
			{
				"slug": "team.manage",
				"display_name": "Manage Team",
				"description": "Can manage team settings and members",
				"category_display_name": "Teams",
				"category_slug": "teams",
				"parent_slug": "",
				"available": true,
				"dependency_slugs": []
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

func TestPermissionsListCurrentUser(t *testing.T) {
	var requestPath string
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPath = req.URL.Path
		w.Write([]byte(expectedCurrentUserPermissionsJSON()))
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	res, err := c.Permissions().ListUser(context.TODO())
	if err != nil {
		t.Fatalf("error listing current user permissions: %s", err.Error())
	}

	if expected := "/permissions/user"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}

	expectedSlugs := []string{"incidents.read", "incidents.write"}
	if !reflect.DeepEqual(expectedSlugs, res.Data) {
		t.Fatalf("permissions mismatch: expected '%+v', got: '%+v'", expectedSlugs, res.Data)
	}
}

func TestPermissionsListTeam(t *testing.T) {
	var requestPath string
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPath = req.URL.Path
		w.Write([]byte(expectedTeamPermissionsResponseJSON()))
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	res, err := c.Permissions().ListTeamPermissions(context.TODO())
	if err != nil {
		t.Fatalf("error listing team permissions: %s", err.Error())
	}

	if expected := "/permissions/team"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}

	if len(res.Data) != 1 {
		t.Fatalf("expected 1 team permission, got %d", len(res.Data))
	}

	teamPerm := res.Data[0]
	if teamPerm.Slug != "team.manage" {
		t.Fatalf("expected team permission slug 'team.manage', got %s", teamPerm.Slug)
	}
	if teamPerm.CategorySlug != "teams" {
		t.Fatalf("expected team permission category 'teams', got %s", teamPerm.CategorySlug)
	}
}
