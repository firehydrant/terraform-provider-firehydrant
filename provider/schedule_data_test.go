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

func testScheduleDataSourceConfig_basic() string {
	return fmt.Sprintln(`
data "firehydrant_schedule" "test_schedule" {
  name = "My Rotation"
}`)
}

func TestScheduleDataSource_OneMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/v1/ping" && r.URL.Path != "/schedules" {
			t.Errorf("Expected to request '/ping' or '/schedules', got: %s", r.URL.Path)
		}

		if r.URL.Path == "/schedules" && r.URL.Query().Get("query") != "My Rotation" {
			t.Errorf("Expected query param 'query' to be 'My Rotation', got: %s", r.URL.Query().Get("query"))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":[{"id": "123", "name":"My Rotation", "integration" : "", "discarded" : false }]}`))
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
				Config: testScheduleDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_schedule.test_schedule", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_schedule.test_schedule", "name", "My Rotation"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_schedule.test_schedule", "id", "123"),
				),
			},
		},
	})
}

func TestScheduleDataSource_MultipleMatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/v1/ping" && r.URL.Path != "/schedules" {
			t.Errorf("Expected to request '/ping' or '/schedules', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/schedules" && r.URL.Query().Get("query") != "My Rotation" {
			t.Errorf("Expected query param 'query' to be 'My Rotation', got: %s", r.URL.Query().Get("query"))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":[{"id": "123", "name":"My Rotation", "integration" : "", "discarded" : false }, {"id": "123", "name":"My Rotation", "integration" : "", "discarded" : false }]}`))
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
				Config:      testScheduleDataSourceConfig_basic(),
				ExpectError: regexp.MustCompile(`Found multiple matching schedules for 'My Rotation'`),
			},
		},
	})
}

func TestScheduleDataSource_NoMatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/v1/ping" && r.URL.Path != "/schedules" {
			t.Errorf("Expected to request '/ping' or '/schedules', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/schedules" && r.URL.Query().Get("query") != "My Rotation" {
			t.Errorf("Expected query param 'query' to be 'My Rotation', got: %s", r.URL.Query().Get("query"))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
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
				Config:      testScheduleDataSourceConfig_basic(),
				ExpectError: regexp.MustCompile(`Did not find schedule matching 'My Rotation'`),
			},
		},
	})
}
