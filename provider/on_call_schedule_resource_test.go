package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	fhsdk "github.com/firehydrant/firehydrant-go-sdk"
	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOnCallScheduleResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckOnCallScheduleResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOnCallScheduleConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_on_call_schedule", "id"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule", "name", fmt.Sprintf("test-on-call-schedule-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule", "description", fmt.Sprintf("test-description-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule", "time_zone", "America/New_York"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule", "strategy.0.type", "weekly"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule", "strategy.0.handoff_time", "10:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule", "strategy.0.handoff_day", "thursday"),
				),
			},
			{
				Config: testAccOnCallScheduleConfig_restrictions(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "id"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "name", fmt.Sprintf("test-on-call-schedule-restrictions-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "strategy.0.type", "weekly"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "strategy.0.handoff_time", "10:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "strategy.0.handoff_day", "thursday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.0.start_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.0.start_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.0.end_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.0.end_time", "14:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.1.start_day", "tuesday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.1.start_time", "12:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.1.end_day", "tuesday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_on_call_schedule_with_restrictions", "restrictions.1.end_time", "18:00:00"),
				),
			},
			{
				Config: testAccOnCallScheduleConfig_customStrategy(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_custom_schedule", "id"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_custom_schedule", "name", fmt.Sprintf("test-custom-schedule-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_custom_schedule", "strategy.0.type", "custom"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_custom_schedule", "strategy.0.shift_duration", "PT8H"),
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_custom_schedule", "start_time"),
				),
			},
		},
	})
}

func testAccOnCallScheduleConfig_basic(rName string) string {
	return fmt.Sprintf(`
	resource "firehydrant_team" "team_team" {
		name = "test-team-%s"
	}

	resource "firehydrant_on_call_schedule" "test_on_call_schedule" {
		team_id = firehydrant_team.team_team.id
		name = "test-on-call-schedule-%s"
		description = "test-description-%s"
		time_zone = "America/New_York"

		strategy {
			type         = "weekly"
			handoff_time = "10:00:00"
			handoff_day  = "thursday"
		}
	}
	`, rName, rName, rName)
}

func testAccOnCallScheduleConfig_customStrategy(rName string) string {
	return fmt.Sprintf(`
	resource "firehydrant_team" "team_team" {
		name = "test-team-%s"
	}

	resource "firehydrant_on_call_schedule" "test_custom_schedule" {
		team_id = firehydrant_team.team_team.id
		name = "test-custom-schedule-%s"
		description = "test-description-%s"
		time_zone = "America/New_York"
		start_time = "%s"

		strategy {
			type           = "custom"
			shift_duration = "PT8H"
		}
	}
	`, rName, rName, rName, time.Now().Add(24*time.Hour).Format(time.RFC3339))
}

func testAccOnCallScheduleConfig_restrictions(rName string) string {
	return fmt.Sprintf(`
	resource "firehydrant_team" "team_team" {
		name = "test-team-%s"
	}

	resource "firehydrant_on_call_schedule" "test_on_call_schedule_with_restrictions" {
		team_id = firehydrant_team.team_team.id
		name = "test-on-call-schedule-restrictions-%s"
		description = "test-description-%s"
		time_zone = "America/New_York"
		slack_user_group_id = "test-group-1"

		strategy {
			type         = "weekly"
			handoff_time = "10:00:00"
			handoff_day  = "thursday"
		}

		restrictions {
			start_day = "monday"
			start_time = "09:00:00"
			end_day = "monday"
			end_time = "14:00:00"
		}

		restrictions {
			start_day = "tuesday"
			start_time = "12:00:00"
			end_day = "tuesday"
			end_time = "18:00:00"
		}
	}
	`, rName, rName, rName)
}

func testAccCheckOnCallScheduleResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_on_call_schedule" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			// Check if the on-call schedule still exists
			_, err := client.Sdk.Signals.GetTeamOnCallSchedule(context.TODO(), stateResource.Primary.Attributes["team_id"], stateResource.Primary.ID, nil, nil)
			if err == nil {
				return fmt.Errorf("On-call schedule %s still exists", stateResource.Primary.ID)
			}
			errStr := err.Error()
			if !strings.Contains(errStr, "404") && !strings.Contains(errStr, "record not found") {
				return fmt.Errorf("Error checking on-call schedule %s: %v", stateResource.Primary.ID, err)
			}
		}

		return nil
	}
}

func offlineOnCallScheduleMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Handle different endpoints
		if req.Method == "GET" {
			// GET request for reading
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
  "id": "schedule-id",
  "name": "A pleasant on-call schedule",
  "description": "Managed by Terraform. Contact @platform-eng for changes.",
  "members": [
    {
      "id": "member-1",
      "name": "Frederick Graff"
    }
  ],
  "team": {
    "id": "team-1",
    "name": "Philadelphia"
  },
  "time_zone": "America/New_York"
}`))
		} else if req.Method == "POST" {
			// POST request for creating
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{
  "id": "schedule-id",
  "name": "A pleasant on-call schedule",
  "description": "Managed by Terraform. Contact @platform-eng for changes.",
  "members": [
    {
      "id": "member-1",
      "name": "Frederick Graff"
    }
  ],
  "team": {
    "id": "team-1",
    "name": "Philadelphia"
  },
  "time_zone": "America/New_York"
}`))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func TestOfflineOnCallScheduleReadMemberID(t *testing.T) {
	ts := offlineOnCallScheduleMockServer()
	defer ts.Close()

	client := &firehydrant.APIClient{}
	client.Sdk = fhsdk.New(
		fhsdk.WithServerURL(ts.URL),
		fhsdk.WithSecurity(components.Security{
			APIKey: "test-token-very-authorized",
		}),
	)

	r := schema.TestResourceDataRaw(t, resourceOnCallSchedule().Schema, map[string]interface{}{
		"team_id":     "team-1",
		"name":        "test-on-call-schedule",
		"description": "test-description",
		"time_zone":   "America/New_York",
		"member_ids":  []interface{}{"member-1"},
	})

	d := readResourceFireHydrantOnCallSchedule(context.Background(), r, client)
	if d.HasError() {
		t.Fatalf("error reading on-call schedule: %v", d)
	}

	memberIDs := r.Get("member_ids").([]interface{})
	if len(memberIDs) != 1 {
		t.Fatalf("expected 1 member ID, got %d", len(memberIDs))
	}

	memberID, ok := memberIDs[0].(string)
	if !ok {
		t.Fatalf("expected member ID to be a string, got %T: %v", memberIDs[0], memberIDs[0])
	}
	if memberID != "member-1" {
		t.Fatalf("expected member ID to be member-1, got %s", memberIDs[0].(string))
	}
}

func TestOfflineOnCallScheduleCreate(t *testing.T) {
	ts := offlineOnCallScheduleMockServer()
	defer ts.Close()

	client := &firehydrant.APIClient{}
	client.Sdk = fhsdk.New(
		fhsdk.WithServerURL(ts.URL),
		fhsdk.WithSecurity(components.Security{
			APIKey: "test-token-very-authorized",
		}),
	)

	r := schema.TestResourceDataRaw(t, resourceOnCallSchedule().Schema, map[string]interface{}{
		"team_id":     "team-1",
		"name":        "test-on-call-schedule",
		"description": "test-description",
		"time_zone":   "America/New_York",
		"member_ids":  []interface{}{"member-1"},
	})

	d := createResourceFireHydrantOnCallSchedule(context.Background(), r, client)
	if d.HasError() {
		t.Fatalf("error creating on-call schedule: %v", d)
	}

	memberIDs := r.Get("member_ids").([]interface{})
	if len(memberIDs) != 1 {
		t.Fatalf("expected 1 member ID, got %d", len(memberIDs))
	}

	memberID, ok := memberIDs[0].(string)
	if !ok {
		t.Fatalf("expected member ID to be a string, got %T: %v", memberIDs[0], memberIDs[0])
	}
	if memberID != "member-1" {
		t.Fatalf("expected member ID to be member-1, got %s", memberIDs[0].(string))
	}
}

// Deprecated, but ensure it still works until we officially remove support.
func TestOfflineOnCallScheduleCreateDeprecated(t *testing.T) {
	ts := offlineOnCallScheduleMockServer()
	defer ts.Close()

	client := &firehydrant.APIClient{}
	client.Sdk = fhsdk.New(
		fhsdk.WithServerURL(ts.URL),
		fhsdk.WithSecurity(components.Security{
			APIKey: "test-token-very-authorized",
		}),
	)

	r := schema.TestResourceDataRaw(t, resourceOnCallSchedule().Schema, map[string]interface{}{
		"team_id":     "team-1",
		"name":        "test-on-call-schedule",
		"description": "test-description",
		"time_zone":   "America/New_York",
		"members":     []interface{}{"member-1"},
	})

	d := createResourceFireHydrantOnCallSchedule(context.Background(), r, client)
	if d.HasError() {
		t.Fatalf("error creating on-call schedule: %v", d)
	}

	memberIDs := r.Get("member_ids").([]interface{})
	if len(memberIDs) != 1 {
		t.Fatalf("expected 1 member ID, got %d", len(memberIDs))
	}

	memberID, ok := memberIDs[0].(string)
	if !ok {
		t.Fatalf("expected member ID to be a string, got %T: %v", memberIDs[0], memberIDs[0])
	}
	if memberID != "member-1" {
		t.Fatalf("expected member ID to be member-1, got %s", memberIDs[0].(string))
	}
}

func TestAccOnCallScheduleResource_updateHandoffAndRestrictions(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckOnCallScheduleResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Initial configuration
				Config: testAccOnCallScheduleConfig_withHandoff(rName, "monday", "09:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_schedule", "id"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "name", fmt.Sprintf("test-schedule-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "time_zone", "America/New_York"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "slack_user_group_id", "test-group-1"),
					// Strategy settings
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.type", "weekly"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_day", "monday"),
					// Initial configuration has no restrictions
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.#", "0"),
				),
			},
			{
				// Update handoff day/time and add restrictions
				Config: testAccOnCallScheduleConfig_withHandoffAndRestrictions(rName, "wednesday", "13:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_schedule", "id"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "name", fmt.Sprintf("test-schedule-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "time_zone", "America/New_York"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "slack_user_group_id", "test-group-1"),
					// Changed handoff settings
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.type", "weekly"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_time", "13:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_day", "wednesday"),
					// Added restrictions
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.#", "2"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.0.start_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.0.start_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.0.end_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.0.end_time", "17:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.1.start_day", "tuesday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.1.start_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.1.end_day", "tuesday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.1.end_time", "17:00:00"),
				),
			},
			{
				// Update just handoff time, keeping restrictions
				Config: testAccOnCallScheduleConfig_withHandoffAndRestrictions(rName, "wednesday", "15:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_schedule", "id"),
					// Changed handoff time only
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_time", "15:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_day", "wednesday"),
					// Restrictions still present
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.#", "2"),
				),
			},
		},
	})
}

func testAccOnCallScheduleConfig_withHandoff(rName, handoffDay, handoffTime string) string {
	return fmt.Sprintf(`
	resource "firehydrant_team" "test_team" {
		name = "test-team-%s"
	}

	resource "firehydrant_on_call_schedule" "test_schedule" {
		team_id = firehydrant_team.test_team.id
		name = "test-schedule-%s"
		time_zone = "America/New_York"
		slack_user_group_id = "test-group-1"

		strategy {
			type         = "weekly"
			handoff_time = "%s"
			handoff_day  = "%s"
		}
	}
	`, rName, rName, handoffTime, handoffDay)
}

func testAccOnCallScheduleConfig_withHandoffAndRestrictions(rName, handoffDay, handoffTime string) string {
	return fmt.Sprintf(`
	resource "firehydrant_team" "test_team" {
		name = "test-team-%s"
	}

	resource "firehydrant_on_call_schedule" "test_schedule" {
		team_id = firehydrant_team.test_team.id
		name = "test-schedule-%s"
		time_zone = "America/New_York"
		slack_user_group_id = "test-group-1"

		strategy {
			type         = "weekly"
			handoff_time = "%s"
			handoff_day  = "%s"
		}

		restrictions {
			start_day = "monday"
			start_time = "09:00:00"
			end_day = "monday"
			end_time = "17:00:00"
		}

		restrictions {
			start_day = "tuesday"
			start_time = "09:00:00"
			end_day = "tuesday"
			end_time = "17:00:00"
		}
	}
	`, rName, rName, handoffTime, handoffDay)
}

func TestAccOnCallScheduleResource_scheduleModifications(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckOnCallScheduleResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Initial configuration with restrictions
				Config: testAccOnCallScheduleConfig_withHandoffAndRestrictions(rName, "monday", "09:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_schedule", "id"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.type", "weekly"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.#", "2"),
				),
			},
			{
				// Change just handoff day, keeping time and restrictions
				Config: testAccOnCallScheduleConfig_withHandoffAndRestrictions(rName, "friday", "09:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_schedule", "id"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_day", "friday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.#", "2"),
				),
			},
			{
				// Remove all restrictions but keep handoff settings
				Config: testAccOnCallScheduleConfig_withHandoff(rName, "friday", "09:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_schedule", "id"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_day", "friday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.#", "0"),
				),
			},
			{
				// Add different restriction pattern
				Config: testAccOnCallScheduleConfig_withBusinessHours(rName, "friday", "09:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_schedule", "id"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_day", "friday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.#", "1"),
					// Check business hours restriction
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.0.start_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.0.start_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.0.end_day", "friday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "restrictions.0.end_time", "17:00:00"),
				),
			},
		},
	})
}

func testAccOnCallScheduleConfig_withBusinessHours(rName, handoffDay, handoffTime string) string {
	return fmt.Sprintf(`
	resource "firehydrant_team" "test_team" {
		name = "test-team-%s"
	}

	resource "firehydrant_on_call_schedule" "test_schedule" {
		team_id = firehydrant_team.test_team.id
		name = "test-schedule-%s"
		time_zone = "America/New_York"
		slack_user_group_id = "test-group-1"

		strategy {
			type         = "weekly"
			handoff_time = "%s"
			handoff_day  = "%s"
		}

		restrictions {
			start_day = "monday"
			start_time = "09:00:00"
			end_day = "friday"
			end_time = "17:00:00"
		}
	}
	`, rName, rName, handoffTime, handoffDay)
}

func TestAccOnCallScheduleResource_effectiveAt(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	futureTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339) // Tomorrow
	pastTime := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)  // Yesterday

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckOnCallScheduleResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Initial schedule setup
				Config: testAccOnCallScheduleConfig_withHandoff(rName, "monday", "09:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_schedule", "id"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_time", "09:00:00"),
				),
			},
			{
				// Update with future effective_at
				Config: testAccOnCallScheduleConfig_withEffectiveAt(rName, "friday", "13:00:00", futureTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_schedule", "id"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_day", "friday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_time", "13:00:00"),
					// effective_at shouldn't be in state
					resource.TestCheckNoResourceAttr("firehydrant_on_call_schedule.test_schedule", "effective_at"),
				),
			},
			{
				// Update with past effective_at (should apply immediately)
				Config: testAccOnCallScheduleConfig_withEffectiveAt(rName, "wednesday", "15:00:00", pastTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_on_call_schedule.test_schedule", "id"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_day", "wednesday"),
					resource.TestCheckResourceAttr("firehydrant_on_call_schedule.test_schedule", "strategy.0.handoff_time", "15:00:00"),
					// effective_at shouldn't be in state
					resource.TestCheckNoResourceAttr("firehydrant_on_call_schedule.test_schedule", "effective_at"),
				),
			},
			{
				// Test invalid timestamp format
				Config:      testAccOnCallScheduleConfig_withEffectiveAt(rName, "thursday", "10:00:00", "invalid-timestamp"),
				ExpectError: regexp.MustCompile("effective_at must be a valid RFC3339 timestamp"),
			},
			{
				// Verify plan is empty when effective_at changes but nothing else does
				Config:   testAccOnCallScheduleConfig_withEffectiveAt(rName, "wednesday", "15:00:00", futureTime),
				PlanOnly: true,
			},
		},
	})
}

func testAccOnCallScheduleConfig_withEffectiveAt(rName, handoffDay, handoffTime, effectiveAt string) string {
	return fmt.Sprintf(`
	resource "firehydrant_team" "test_team" {
		name = "test-team-%s"
	}

	resource "firehydrant_on_call_schedule" "test_schedule" {
		team_id = firehydrant_team.test_team.id
		name = "test-schedule-%s"
		time_zone = "America/New_York"
		slack_user_group_id = "test-group-1"

		strategy {
			type         = "weekly"
			handoff_time = "%s"
			handoff_day  = "%s"
		}

		effective_at = "%s"
	}
	`, rName, rName, handoffTime, handoffDay, effectiveAt)
}

func TestAccOnCallScheduleResourceImport_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resourceName := "firehydrant_on_call_schedule.test_on_call_schedule"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckOnCallScheduleResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccOnCallScheduleConfig_basic(rName),
			},
			{
				ResourceName: resourceName,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("Not found: %s", resourceName)
					}
					return fmt.Sprintf("%s:%s", rs.Primary.Attributes["team_id"], rs.Primary.Attributes["id"]), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
