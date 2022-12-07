package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testUserDataSourceConfig_basic() string {
	return fmt.Sprintln(`
data "firehydrant_user" "test_user" {
  email = "test-user@firehydrant.io"
}`)
}

func TestUserDataSource_OneMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/users" {
			t.Errorf("Expected to request '/ping' or '/users', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/users" && r.URL.Query().Get("query") != "test-user@firehydrant.io" {
			t.Errorf("Expected query param 'query' to be 'test-user@firehydrant.io', got: %s", r.URL.Query().Get("query"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"id": "123", "email":"test-user@firehydrant.io"}]}`))
	}))

	defer server.Close()

	orig := os.Getenv("FIREHYDRANT_BASE_URL")
	os.Setenv("FIREHYDRANT_BASE_URL", server.URL)
	t.Cleanup(func() { os.Setenv("FIREHYDRANT_BASE_URL", orig) })

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testUserDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_user.test_user", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_user.test_user", "email", "test-user@firehydrant.io"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_user.test_user", "id", "123"),
				),
			},
		},
	})
}

func TestUserDataSource_MultipleMatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/users" {
			t.Errorf("Expected to request '/ping' or '/users', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/users" && r.URL.Query().Get("query") != "test-user@firehydrant.io" {
			t.Errorf("Expected query param 'query' to be 'test-user@firehydrant.io', got: %s", r.URL.Query().Get("query"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"id": "123", "email":"test-user@firehydrant.io"},{"id": "456", "email":"test-user@example.io"}]}`))
	}))

	defer server.Close()

	orig := os.Getenv("FIREHYDRANT_BASE_URL")
	os.Setenv("FIREHYDRANT_BASE_URL", server.URL)
	t.Cleanup(func() { os.Setenv("FIREHYDRANT_BASE_URL", orig) })

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testUserDataSourceConfig_basic(),
				ExpectError: regexp.MustCompile(`Found multiple matching users for 'test-user@firehydrant.io'`),
			},
		},
	})
}

func TestUserDataSource_NoMatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/users" {
			t.Errorf("Expected to request '/ping' or '/users', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/users" && r.URL.Query().Get("query") != "test-user@firehydrant.io" {
			t.Errorf("Expected query param 'query' to be 'test-user@firehydrant.io', got: %s", r.URL.Query().Get("query"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[]}`))
	}))

	defer server.Close()

	orig := os.Getenv("FIREHYDRANT_BASE_URL")
	os.Setenv("FIREHYDRANT_BASE_URL", server.URL)
	t.Cleanup(func() { os.Setenv("FIREHYDRANT_BASE_URL", orig) })

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testUserDataSourceConfig_basic(),
				ExpectError: regexp.MustCompile(`Did not find user matching 'test-user@firehydrant.io'`),
			},
		},
	})
}
