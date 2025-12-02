package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"regexp"
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

func TestAccRotationResource_basic(t *testing.T) {
	sharedTeamID := getSharedTeamID(t)
	sharedScheduleID := getSharedOnCallScheduleID(t)
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckRotationResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRotationConfig_basic(rName, sharedTeamID, sharedScheduleID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation", "id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "name", fmt.Sprintf("test-rotation-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.type", "weekly"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_time", "10:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_day", "thursday"),
				),
			},
			{
				Config: testAccRotationConfig_restrictions(rName, sharedTeamID, sharedScheduleID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation_with_restrictions", "id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_with_restrictions", "name", fmt.Sprintf("test-rotation-restrictions-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_with_restrictions", "strategy.0.type", "weekly"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_with_restrictions", "strategy.0.handoff_time", "10:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_with_restrictions", "strategy.0.handoff_day", "thursday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_with_restrictions", "restrictions.0.start_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_with_restrictions", "restrictions.0.start_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_with_restrictions", "restrictions.0.end_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_with_restrictions", "restrictions.0.end_time", "14:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_with_restrictions", "restrictions.1.start_day", "tuesday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_with_restrictions", "restrictions.1.start_time", "12:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_with_restrictions", "restrictions.1.end_day", "tuesday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_with_restrictions", "restrictions.1.end_time", "18:00:00"),
				),
			},
		},
	})
}

func testAccRotationConfig_basic(rName, sharedTeamID, sharedScheduleID string) string {
	return fmt.Sprintf(`
	resource "firehydrant_rotation" "test_rotation" {
	  team_id = "%s"
		schedule_id = "%s"
		name = "test-rotation-%s"
		description = "test-description-%s"
		time_zone = "America/New_York"

		enable_slack_channel_notifications = false
		prevent_shift_deletion = true
		color = "#3192ff"

		strategy {
			type         = "weekly"
			handoff_time = "10:00:00"
			handoff_day  = "thursday"
		}
	}
	`, sharedTeamID, sharedScheduleID, rName, rName)
}

func testAccRotationConfig_restrictions(rName, sharedTeamID, sharedScheduleID string) string {
	return fmt.Sprintf(`
	resource "firehydrant_rotation" "test_rotation_with_restrictions" {
	  team_id = "%s"
		schedule_id = "%s"
		name = "test-rotation-restrictions-%s"
		description = "test-description-%s"
		time_zone = "America/New_York"

		enable_slack_channel_notifications = false
		prevent_shift_deletion = true
		color = "#3192ff"

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
	`, sharedTeamID, sharedScheduleID, rName, rName)
}

func testAccCheckRotationResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_rotation" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			// Check if the rotation still exists
			_, err := client.Sdk.Signals.GetOnCallScheduleRotation(context.TODO(), stateResource.Primary.ID, stateResource.Primary.Attributes["team_id"], stateResource.Primary.Attributes["schedule_id"])
			if err == nil {
				return fmt.Errorf("Rotation %s still exists", stateResource.Primary.ID)
			}
			errStr := err.Error()
			if !strings.Contains(errStr, "404") && !strings.Contains(errStr, "record not found") {
				return fmt.Errorf("Error checking rotation %s: %v", stateResource.Primary.ID, err)
			}
		}

		return nil
	}
}

func offlineRotationMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Handle GET request for reading rotation
		if req.Method == "GET" {
			w.Write([]byte(`{
  "id": "rotation-id",
  "name": "A pleasant rotation",
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
  "time_zone": "America/New_York",
  "enable_slack_channel_notifications": false,
  "prevent_shift_deletion": false,
  "strategy": {
    "type": "weekly",
    "handoff_time": "10:00:00",
    "handoff_day": "thursday"
  },
  "restrictions": [],
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}`))
		} else if req.Method == "POST" {
			// Handle POST request for creating rotation
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{
  "id": "rotation-id",
  "name": "A pleasant rotation",
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
  "time_zone": "America/New_York",
  "enable_slack_channel_notifications": false,
  "prevent_shift_deletion": false,
  "strategy": {
    "type": "weekly",
    "handoff_time": "10:00:00",
    "handoff_day": "thursday"
  },
  "restrictions": [],
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z"
}`))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
}

func TestOfflineRotationReadMemberID(t *testing.T) {
	ts := offlineRotationMockServer()
	defer ts.Close()

	client := &firehydrant.APIClient{}
	client.Sdk = fhsdk.New(
		fhsdk.WithServerURL(ts.URL),
		fhsdk.WithSecurity(components.Security{
			APIKey: "test-token-very-authorized",
		}),
	)

	r := schema.TestResourceDataRaw(t, resourceRotation().Schema, map[string]interface{}{
		"team_id":     "team-1",
		"schedule_id": "schedule-1",
		"id":          "rotation-id",
		"name":        "test-rotation",
		"description": "test-description",
		"time_zone":   "America/New_York",
		"members": []interface{}{
			map[string]interface{}{
				"user_id": "member-1",
			},
		},
	})

	d := readResourceFireHydrantRotation(context.Background(), r, client)
	if d.HasError() {
		t.Fatalf("error reading rotation: %v", d)
	}

	members := r.Get("members").([]interface{})
	if len(members) != 1 {
		t.Fatalf("expected 1 member, got %d", len(members))
	}

	memberMap, ok := members[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected member to be a map, got %T: %v", members[0], members[0])
	}

	userID, ok := memberMap["user_id"].(string)
	if !ok {
		t.Fatalf("expected user_id to be a string, got %T: %v", memberMap["user_id"], memberMap["user_id"])
	}
	if userID != "member-1" {
		t.Fatalf("expected user_id to be member-1, got %s", userID)
	}
}

func TestOfflineRotationCreate(t *testing.T) {
	ts := offlineRotationMockServer()
	defer ts.Close()

	client := &firehydrant.APIClient{}
	client.Sdk = fhsdk.New(
		fhsdk.WithServerURL(ts.URL),
		fhsdk.WithSecurity(components.Security{
			APIKey: "test-token-very-authorized",
		}),
	)

	r := schema.TestResourceDataRaw(t, resourceRotation().Schema, map[string]interface{}{
		"team_id":     "team-1",
		"schedule_id": "schedule-1",
		"name":        "test-rotation",
		"description": "test-description",
		"time_zone":   "America/New_York",
		"members": []interface{}{
			map[string]interface{}{
				"user_id": "member-1",
			},
		},
	})

	d := createResourceFireHydrantRotation(context.Background(), r, client)
	if d.HasError() {
		t.Fatalf("error creating rotation: %v", d)
	}

	// Read the resource to populate members in state (as Terraform would do)
	d = readResourceFireHydrantRotation(context.Background(), r, client)
	if d.HasError() {
		t.Fatalf("error reading rotation: %v", d)
	}

	members := r.Get("members").([]interface{})
	if len(members) != 1 {
		t.Fatalf("expected 1 member, got %d", len(members))
	}

	memberMap, ok := members[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected member to be a map, got %T: %v", members[0], members[0])
	}

	userID, ok := memberMap["user_id"].(string)
	if !ok {
		t.Fatalf("expected user_id to be a string, got %T: %v", memberMap["user_id"], memberMap["user_id"])
	}
	if userID != "member-1" {
		t.Fatalf("expected user_id to be member-1, got %s", userID)
	}
}

func TestAccRotationResource_updateHandoffAndRestrictions(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckRotationResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Initial configuration
				Config: testAccRotationConfig_withHandoff(rName, "monday", "09:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation", "id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "name", fmt.Sprintf("test-rotation-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "time_zone", "America/New_York"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "slack_user_group_id", "test-group-1"),
					// Strategy settings
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.type", "weekly"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_day", "monday"),
					// Initial configuration has no restrictions
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.#", "0"),
				),
			},
			{
				// Update handoff day/time and add restrictions
				Config: testAccRotationConfig_withHandoffAndRestrictions(rName, "wednesday", "13:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation", "id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "name", fmt.Sprintf("test-rotation-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "time_zone", "America/New_York"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "slack_user_group_id", "test-group-1"),
					// Changed handoff settings
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.type", "weekly"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_time", "13:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_day", "wednesday"),
					// Added restrictions
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.#", "2"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.0.start_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.0.start_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.0.end_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.0.end_time", "17:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.1.start_day", "tuesday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.1.start_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.1.end_day", "tuesday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.1.end_time", "17:00:00"),
				),
			},
			{
				// Update just handoff time, keeping restrictions
				Config: testAccRotationConfig_withHandoffAndRestrictions(rName, "wednesday", "15:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation", "id"),
					// Changed handoff time only
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_time", "15:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_day", "wednesday"),
					// Restrictions still present
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.#", "2"),
				),
			},
		},
	})
}

func testAccRotationConfig_withHandoff(rName, handoffDay, handoffTime string) string {
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
			handoff_time = "01:00:00"
			handoff_day  = "thursday"
		}
	}

	resource "firehydrant_rotation" "test_rotation" {
		team_id = firehydrant_team.test_team.id
		schedule_id = firehydrant_on_call_schedule.test_schedule.id
		name = "test-rotation-%s"
		time_zone = "America/New_York"
		slack_user_group_id = "test-group-1"
		color = "#3192ff"

		strategy {
			type         = "weekly"
			handoff_time = "%s"
			handoff_day  = "%s"
		}
	}
	`, rName, rName, rName, handoffTime, handoffDay)
}

func testAccRotationConfig_withHandoffAndRestrictions(rName, handoffDay, handoffTime string) string {
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
			handoff_time = "01:00:00"
			handoff_day  = "thursday"
		}
	}

	resource "firehydrant_rotation" "test_rotation" {
		team_id = firehydrant_team.test_team.id
		schedule_id = firehydrant_on_call_schedule.test_schedule.id
		name = "test-rotation-%s"
		time_zone = "America/New_York"
		slack_user_group_id = "test-group-1"
		color = "#3192ff"

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
	`, rName, rName, rName, handoffTime, handoffDay)
}

func TestAccRotationResource_scheduleModifications(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckRotationResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Initial configuration with restrictions
				Config: testAccRotationConfig_withHandoffAndRestrictions(rName, "monday", "09:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation", "id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.type", "weekly"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.#", "2"),
				),
			},
			{
				// Change just handoff day, keeping time and restrictions
				Config: testAccRotationConfig_withHandoffAndRestrictions(rName, "friday", "09:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation", "id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_day", "friday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.#", "2"),
				),
			},
			{
				// Remove all restrictions but keep handoff settings
				Config: testAccRotationConfig_withHandoff(rName, "friday", "09:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation", "id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_day", "friday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.#", "0"),
				),
			},
			{
				// Add different restriction pattern
				Config: testAccRotationConfig_withBusinessHours(rName, "friday", "09:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation", "id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_day", "friday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.#", "1"),
					// Check business hours restriction
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.0.start_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.0.start_time", "09:00:00"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.0.end_day", "friday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "restrictions.0.end_time", "17:00:00"),
				),
			},
		},
	})
}

func testAccRotationConfig_withBusinessHours(rName, handoffDay, handoffTime string) string {
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
			handoff_time = "01:00:00"
			handoff_day  = "thursday"
		}
	}

	resource "firehydrant_rotation" "test_rotation" {
		team_id = firehydrant_team.test_team.id
		schedule_id = firehydrant_on_call_schedule.test_schedule.id
		name = "test-rotation-%s"
		time_zone = "America/New_York"
		slack_user_group_id = "test-group-1"
		color = "#3192ff"

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
	`, rName, rName, rName, handoffTime, handoffDay)
}

func TestAccRotationResource_effectiveAt(t *testing.T) {
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	futureTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339) // Tomorrow
	pastTime := time.Now().Add(-24 * time.Hour).Format(time.RFC3339)  // Yesterday

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckRotationResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Initial schedule setup
				Config: testAccRotationConfig_withHandoff(rName, "monday", "09:00:00"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation", "id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_day", "monday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_time", "09:00:00"),
				),
			},
			{
				// Update with future effective_at
				Config: testAccRotationConfig_withEffectiveAt(rName, "friday", "13:00:00", futureTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation", "id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_day", "friday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_time", "13:00:00"),
					// effective_at shouldn't be in state
					resource.TestCheckNoResourceAttr("firehydrant_rotation.test_rotation", "effective_at"),
				),
			},
			{
				// Update with past effective_at (should apply immediately)
				Config: testAccRotationConfig_withEffectiveAt(rName, "wednesday", "15:00:00", pastTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation", "id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_day", "wednesday"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation", "strategy.0.handoff_time", "15:00:00"),
					// effective_at shouldn't be in state
					resource.TestCheckNoResourceAttr("firehydrant_rotation.test_rotation", "effective_at"),
				),
			},
			{
				// Test invalid timestamp format
				Config:      testAccRotationConfig_withEffectiveAt(rName, "thursday", "10:00:00", "invalid-timestamp"),
				ExpectError: regexp.MustCompile("effective_at must be a valid RFC3339 timestamp"),
			},
			{
				// Verify plan is empty when effective_at changes but nothing else does
				Config:   testAccRotationConfig_withEffectiveAt(rName, "wednesday", "15:00:00", futureTime),
				PlanOnly: true,
			},
		},
	})
}

func testAccRotationConfig_withEffectiveAt(rName, handoffDay, handoffTime, effectiveAt string) string {
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
			handoff_time = "01:00:00"
			handoff_day  = "thursday"
		}
	}

	resource "firehydrant_rotation" "test_rotation" {
		team_id = firehydrant_team.test_team.id
		schedule_id = firehydrant_on_call_schedule.test_schedule.id
		name = "test-rotation-%s"
		time_zone = "America/New_York"
		slack_user_group_id = "test-group-1"
		color = "#3192ff"

		strategy {
			type         = "weekly"
			handoff_time = "%s"
			handoff_day  = "%s"
		}

		effective_at = "%s"
	}
	`, rName, rName, rName, handoffTime, handoffDay, effectiveAt)
}

func TestAccRotationResourceImport_basic(t *testing.T) {
	sharedTeamID := getSharedTeamID(t)
	sharedScheduleID := getSharedOnCallScheduleID(t)
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	resourceName := "firehydrant_rotation.test_rotation"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckRotationResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccRotationConfig_basic(rName, sharedTeamID, sharedScheduleID),
			},
			{
				ResourceName: resourceName,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("Not found: %s", resourceName)
					}
					return fmt.Sprintf("%s:%s:%s", rs.Primary.Attributes["team_id"], rs.Primary.Attributes["schedule_id"], rs.Primary.Attributes["id"]), nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRotationResource_members(t *testing.T) {
	sharedTeamID := getSharedTeamID(t)
	sharedScheduleID := getSharedOnCallScheduleID(t)
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)
	futureTime := time.Now().Add(24 * time.Hour).Format(time.RFC3339) // Tomorrow

	existingUser := os.Getenv("EXISTING_USER_EMAIL")

	// Get a second user - use the same user for simplicity, but in real scenarios would be different
	// The API allows the same user to be added multiple times
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckRotationResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Step 1: Create rotation with 2 members
				Config: testAccRotationConfig_withTwoMembers(rName, sharedTeamID, sharedScheduleID, existingUser),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation_members", "id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_members", "name", fmt.Sprintf("test-rotation-members-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_members", "members.#", "2"),
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation_members", "members.0.user_id"),
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation_members", "members.1.user_id"),
				),
			},
			{
				// Step 2: Remove one member (go from 2 to 1)
				Config: testAccRotationConfig_withMember(rName, sharedTeamID, sharedScheduleID, existingUser, futureTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation_members", "id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_members", "members.#", "1"),
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation_members", "members.0.user_id"),
				),
			},
			{
				// Step 3: Add member back (go from 1 to 2)
				Config: testAccRotationConfig_withTwoMembers(rName, sharedTeamID, sharedScheduleID, existingUser, futureTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation_members", "id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_members", "members.#", "2"),
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation_members", "members.0.user_id"),
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation_members", "members.1.user_id"),
				),
			},
		},
	})
}

func TestAccRotationResource_membersWithUnassignedSlot(t *testing.T) {
	sharedTeamID := getSharedTeamID(t)
	sharedScheduleID := getSharedOnCallScheduleID(t)
	rName := acctest.RandStringFromCharSet(20, acctest.CharSetAlphaNum)

	existingUser := os.Getenv("EXISTING_USER_EMAIL")
	if existingUser == "" {
		existingUser = "ops+terraform-ci@firehydrant.io"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckRotationResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Step 1: Create rotation with a member, unassigned slot, and another member
				Config: testAccRotationConfig_withUnassignedSlot(rName, sharedTeamID, sharedScheduleID, existingUser),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation_members", "id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_members", "name", fmt.Sprintf("test-rotation-members-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_members", "members.#", "3"),
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation_members", "members.0.user_id"),
					resource.TestCheckResourceAttr("firehydrant_rotation.test_rotation_members", "members.1.user_id", ""), // Unassigned slot
					resource.TestCheckResourceAttrSet("firehydrant_rotation.test_rotation_members", "members.2.user_id"),
				),
			},
		},
	})
}

func testAccRotationConfig_withTwoMembers(rName, sharedTeamID, sharedScheduleID, userEmail string, effectiveAt ...string) string {
	effectiveAtStr := ""
	if len(effectiveAt) > 0 && effectiveAt[0] != "" {
		effectiveAtStr = fmt.Sprintf("\n\t\teffective_at = \"%s\"", effectiveAt[0])
	}

	return fmt.Sprintf(`
	data "firehydrant_user" "test_user" {
		email = "%s"
	}

	resource "firehydrant_rotation" "test_rotation_members" {
	  team_id = "%s"
		schedule_id = "%s"
		name = "test-rotation-members-%s"
		description = "test-description-%s"
		time_zone = "America/New_York"

		enable_slack_channel_notifications = false
		prevent_shift_deletion = true
		color = "#3192ff"

		members {
			user_id = data.firehydrant_user.test_user.id
		}

		members {
			user_id = data.firehydrant_user.test_user.id
		}

		strategy {
			type         = "weekly"
			handoff_time = "10:00:00"
			handoff_day  = "thursday"
		}%s
	}
	`, userEmail, sharedTeamID, sharedScheduleID, rName, rName, effectiveAtStr)
}

func testAccRotationConfig_withUnassignedSlot(rName, sharedTeamID, sharedScheduleID, userEmail string) string {
	return fmt.Sprintf(`
	data "firehydrant_user" "test_user" {
		email = "%s"
	}

	resource "firehydrant_rotation" "test_rotation_members" {
	  team_id = "%s"
		schedule_id = "%s"
		name = "test-rotation-members-%s"
		description = "test-description-%s"
		time_zone = "America/New_York"

		enable_slack_channel_notifications = false
		prevent_shift_deletion = true
		color = "#3192ff"

		members {
			user_id = data.firehydrant_user.test_user.id
		}

		members {
			user_id = ""
		}

		members {
			user_id = data.firehydrant_user.test_user.id
		}

		strategy {
			type         = "weekly"
			handoff_time = "10:00:00"
			handoff_day  = "thursday"
		}
	}
	`, userEmail, sharedTeamID, sharedScheduleID, rName, rName)
}

func testAccRotationConfig_withMember(rName, sharedTeamID, sharedScheduleID, userEmail string, effectiveAt ...string) string {
	effectiveAtStr := ""
	if len(effectiveAt) > 0 && effectiveAt[0] != "" {
		effectiveAtStr = fmt.Sprintf("\n\t\teffective_at = \"%s\"", effectiveAt[0])
	}

	return fmt.Sprintf(`
	data "firehydrant_user" "test_user" {
		email = "%s"
	}

	resource "firehydrant_rotation" "test_rotation_members" {
	  team_id = "%s"
		schedule_id = "%s"
		name = "test-rotation-members-%s"
		description = "test-description-%s"
		time_zone = "America/New_York"

		enable_slack_channel_notifications = false
		prevent_shift_deletion = true
		color = "#3192ff"

		members {
			user_id = data.firehydrant_user.test_user.id
		}

		strategy {
			type         = "weekly"
			handoff_time = "10:00:00"
			handoff_day  = "thursday"
		}%s
	}
	`, userEmail, sharedTeamID, sharedScheduleID, rName, rName, effectiveAtStr)
}
