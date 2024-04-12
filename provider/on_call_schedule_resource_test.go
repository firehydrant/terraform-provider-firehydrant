package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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

func testAccOnCallScheduleConfig_customStrategy(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_team" "team_team" {
  name = "test-team-%s"
}

resource "firehydrant_on_call_schedule" "test_on_call_schedule_custom_strategy" {
  name        = "test-on-call-schedule-custom-strategy-%s"
  description = "test-description-%s"
  team_id     = firehydrant_team.team_team.id
  time_zone   = "America/Los_Angeles"
  start_time  = "2024-04-11T11:56:29-07:00"

  strategy {
    type           = "custom"
    shift_duration = "PT93600S"
  }

  restrictions {
    start_day  = "monday"
    start_time = "14:00:00"
    end_day    = "friday"
    end_time   = "17:00:00"
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

			// Normally we'd check if err == nil here, because we'd expect a 404 if we try to get a resource
			// that has been deleted. However, the incident role API will still return deleted/archived incident
			// roles instead of returning 404. So, to check for incident roles that are deleted, we have to check
			// for incident roles that have a DiscardedAt timestamp
			_, err := client.OnCallSchedules().Get(context.TODO(), stateResource.Primary.Attributes["team_id"], stateResource.Primary.ID)
			if err != nil && !errors.Is(err, firehydrant.ErrorNotFound) {
				return fmt.Errorf("On-call schedule %s still exists", stateResource.Primary.ID)
			}
		}

		return nil
	}
}

func offlineOnCallScheduleMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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
	}))
}

func TestOfflineOnCallScheduleReadMemberID(t *testing.T) {
	ts := offlineOnCallScheduleMockServer()
	defer ts.Close()

	c, err := firehydrant.NewRestClient("test-token-very-authorized", firehydrant.WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	r := schema.TestResourceDataRaw(t, resourceOnCallSchedule().Schema, map[string]interface{}{
		"team_id":     "team-1",
		"name":        "test-on-call-schedule",
		"description": "test-description",
		"time_zone":   "America/New_York",
		"member_ids":  []interface{}{"member-1"},
	})

	d := readResourceFireHydrantOnCallSchedule(context.Background(), r, c)
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

	c, err := firehydrant.NewRestClient("test-token-very-authorized", firehydrant.WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	r := schema.TestResourceDataRaw(t, resourceOnCallSchedule().Schema, map[string]interface{}{
		"team_id":     "team-1",
		"name":        "test-on-call-schedule",
		"description": "test-description",
		"time_zone":   "America/New_York",
		"member_ids":  []interface{}{"member-1"},
	})

	d := createResourceFireHydrantOnCallSchedule(context.Background(), r, c)
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

	c, err := firehydrant.NewRestClient("test-token-very-authorized", firehydrant.WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	r := schema.TestResourceDataRaw(t, resourceOnCallSchedule().Schema, map[string]interface{}{
		"team_id":     "team-1",
		"name":        "test-on-call-schedule",
		"description": "test-description",
		"time_zone":   "America/New_York",
		"members":     []interface{}{"member-1"},
	})

	d := createResourceFireHydrantOnCallSchedule(context.Background(), r, c)
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
