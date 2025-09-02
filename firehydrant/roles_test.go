package firehydrant

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func expectedRoleResponse() *RoleResponse {
	return &RoleResponse{
		ID:             "role-id",
		Name:           "Test Role",
		Slug:           "test-role",
		Description:    "A test role for unit testing",
		OrganizationID: "org-123",
		BuiltIn:        false,
		ReadOnly:       false,
		Permissions: []Permission{
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
		},
		CreatedAt: "2025-01-01T00:00:00Z",
		UpdatedAt: "2025-01-01T00:00:00Z",
	}
}

func expectedRoleResponseJSON() string {
	return `{
		"id": "role-id",
		"name": "Test Role",
		"slug": "test-role",
		"description": "A test role for unit testing",
		"organization_id": "org-123",
		"built_in": false,
		"read_only": false,
		"permissions": [
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
			}
		],
		"created_at": "2025-01-01T00:00:00Z",
		"updated_at": "2025-01-01T00:00:00Z"
	}`
}

func expectedRolesListResponseJSON() string {
	return `{
		"data": [` + expectedRoleResponseJSON() + `],
		"pagination": {
			"count": 1,
			"page": 1,
			"items": 1,
			"pages": 1,
			"prev": 0,
			"next": 0,
			"last": 1
		}
	}`
}

func TestRoleGet(t *testing.T) {
	var requestPath string
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPath = req.URL.Path
		w.Write([]byte(expectedRoleResponseJSON()))
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	res, err := c.Roles().Get(context.TODO(), "role-id")
	if err != nil {
		t.Fatalf("error retrieving role: %s", err.Error())
	}

	if expected := "/roles/role-id"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}

	expectedResponse := expectedRoleResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestRoleList(t *testing.T) {
	var requestPath string
	var queryParams string
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPath = req.URL.Path
		queryParams = req.URL.RawQuery
		w.Write([]byte(expectedRolesListResponseJSON()))
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	query := RolesQuery{
		Query: "test",
		Page:  2,
	}

	res, err := c.Roles().List(context.TODO(), query)
	if err != nil {
		t.Fatalf("error listing roles: %s", err.Error())
	}

	if expected := "/roles"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}

	if !strings.Contains(queryParams, "query=test") || !strings.Contains(queryParams, "page=2") {
		t.Fatalf("query params mismatch: expected to contain 'query=test' and 'page=2', got: '%s'", queryParams)
	}

	if len(res.Data) != 1 {
		t.Fatalf("expected 1 role in response, got %d", len(res.Data))
	}

	expectedRole := expectedRoleResponse()
	if !reflect.DeepEqual(expectedRole, &res.Data[0]) {
		t.Fatalf("role mismatch: expected '%+v', got: '%+v'", expectedRole, &res.Data[0])
	}
}

func TestRoleCreate(t *testing.T) {
	var requestPath string
	var requestBody map[string]interface{}
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPath = req.URL.Path
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			t.Fatalf("error unmarshalling request body: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(expectedRoleResponseJSON()))
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	createReq := CreateRoleRequest{
		Name:        "Test Role",
		Slug:        "test-role",
		Description: "A test role for unit testing",
		Permissions: []string{"incidents.read", "incidents.write"},
	}

	res, err := c.Roles().Create(context.TODO(), createReq)
	if err != nil {
		t.Fatalf("error creating role: %s", err.Error())
	}

	if expected := "/roles"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}

	// Verify request body contents
	if requestBody["name"].(string) != "Test Role" {
		t.Fatalf("expected name 'Test Role', got %s", requestBody["name"])
	}

	permissions := requestBody["permissions"].([]interface{})
	if len(permissions) != 2 {
		t.Fatalf("expected 2 permissions, got %d", len(permissions))
	}

	expectedResponse := expectedRoleResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestRoleUpdate(t *testing.T) {
	var requestPath string
	var requestBody map[string]interface{}
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPath = req.URL.Path
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			t.Fatalf("error unmarshalling request body: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(expectedRoleResponseJSON()))
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	updateReq := UpdateRoleRequest{
		Name:        "Updated Role Name",
		Description: "Updated description",
		Permissions: []string{"incidents.read"},
	}

	res, err := c.Roles().Update(context.TODO(), "role-id", updateReq)
	if err != nil {
		t.Fatalf("error updating role: %s", err.Error())
	}

	if expected := "/roles/role-id"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}

	// Verify only changed fields are in request body
	if requestBody["name"].(string) != "Updated Role Name" {
		t.Fatalf("expected name 'Updated Role Name', got %s", requestBody["name"])
	}

	expectedResponse := expectedRoleResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestRoleDelete(t *testing.T) {
	var requestPath string
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPath = req.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	err = c.Roles().Delete(context.TODO(), "role-id")
	if err != nil {
		t.Fatalf("error deleting role: %s", err.Error())
	}

	if expected := "/roles/role-id"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}
}
