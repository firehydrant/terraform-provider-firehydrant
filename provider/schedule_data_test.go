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
  name = "My Schedule"
}`)
}

func testScheduleDataSourceConfig_exactMatch() string {
	return fmt.Sprintln(`
data "firehydrant_schedule" "test_schedule" {
  name = "My Schedule"
}`)
}

func TestScheduleDataSource_OneMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/v1/ping" && r.URL.Path != "/v1/schedules" {
			t.Errorf("Expected to request '/ping' or '/v1/schedules', got: %s", r.URL.Path)
		}

		if r.URL.Path == "/v1/schedules" && r.URL.Query().Get("query") != "My Schedule" {
			t.Errorf("Expected query param 'query' to be 'My Schedule', got: %s", r.URL.Query().Get("query"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"id": "123", "name":"My Schedule", "integration" : "", "discarded" : false }]}`))
	}))

	defer server.Close()

	orig := os.Getenv("FIREHYDRANT_BASE_URL")
	os.Setenv("FIREHYDRANT_BASE_URL", server.URL)
	t.Cleanup(func() { os.Setenv("FIREHYDRANT_BASE_URL", orig) })

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: mockProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testScheduleDataSourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_schedule.test_schedule", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_schedule.test_schedule", "name", "My Schedule"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_schedule.test_schedule", "id", "123"),
				),
			},
		},
	})
}

func TestScheduleDataSource_MultipleMatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/v1/ping" && r.URL.Path != "/v1/schedules" {
			t.Errorf("Expected to request '/ping' or '/v1/schedules', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/v1/schedules" && r.URL.Query().Get("query") != "My Schedule" {
			t.Errorf("Expected query param 'query' to be 'My Schedule', got: %s", r.URL.Query().Get("query"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"id": "456", "name":"My Schedule 2", "integration" : "", "discarded" : false }, {"id": "123", "name":"My Schedule", "integration" : "", "discarded" : false }, {"id": "789", "name":"My Schedule 3", "integration" : "", "discarded" : false }]}`))
	}))

	defer server.Close()

	orig := os.Getenv("FIREHYDRANT_BASE_URL")
	os.Setenv("FIREHYDRANT_BASE_URL", server.URL)
	t.Cleanup(func() { os.Setenv("FIREHYDRANT_BASE_URL", orig) })

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: mockProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testScheduleDataSourceConfig_exactMatch(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_schedule.test_schedule", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_schedule.test_schedule", "name", "My Schedule"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_schedule.test_schedule", "id", "123"),
				),
			},
		},
	})
}

func TestScheduleDataSource_MultipleMatchesNoExactMatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/v1/ping" && r.URL.Path != "/v1/schedules" {
			t.Errorf("Expected to request '/ping' or '/v1/schedules', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/v1/schedules" && r.URL.Query().Get("query") != "My Schedule" {
			t.Errorf("Expected query param 'query' to be 'My Schedule', got: %s", r.URL.Query().Get("query"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[{"id": "456", "name":"My Schedule 2", "integration" : "", "discarded" : false }, {"id": "789", "name":"My Schedule 3", "integration" : "", "discarded" : false }]}`))
	}))

	defer server.Close()

	orig := os.Getenv("FIREHYDRANT_BASE_URL")
	os.Setenv("FIREHYDRANT_BASE_URL", server.URL)
	t.Cleanup(func() { os.Setenv("FIREHYDRANT_BASE_URL", orig) })

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: mockProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testScheduleDataSourceConfig_exactMatch(),
				ExpectError: regexp.MustCompile(`Did not find schedule matching 'My Schedule'`),
			},
		},
	})
}

func TestScheduleDataSource_NoMatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/ping" && r.URL.Path != "/v1/ping" && r.URL.Path != "/v1/schedules" {
			t.Errorf("Expected to request '/ping' or '/v1/schedules', got: %s", r.URL.Path)
		}
		if r.URL.Path == "/v1/schedules" && r.URL.Query().Get("query") != "My Schedule" {
			t.Errorf("Expected query param 'query' to be 'My Schedule', got: %s", r.URL.Query().Get("query"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":[]}`))
	}))

	defer server.Close()

	orig := os.Getenv("FIREHYDRANT_BASE_URL")
	os.Setenv("FIREHYDRANT_BASE_URL", server.URL)
	t.Cleanup(func() { os.Setenv("FIREHYDRANT_BASE_URL", orig) })

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: mockProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testScheduleDataSourceConfig_basic(),
				ExpectError: regexp.MustCompile(`Did not find schedule matching 'My Schedule'`),
			},
		},
	})
}
