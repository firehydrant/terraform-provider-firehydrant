package firehydrant

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func expectedRotationResponse() *RotationResponse {
	return &RotationResponse{
		ID:                              "rotation-id",
		Name:                            "A pleasant rotation",
		Description:                     "Managed by Terraform. Contact @platform-eng for changes.",
		TimeZone:                        "America/New_York",
		Color:                           "#ff0000",
		SlackUserGroupID:                "slack-group-1",
		PreventShiftDeletion:            true,
		CoverageGapNotificationInterval: "an-appropriate-interval",
		Strategy: RotationStrategy{
			Type:        "weekly",
			HandoffTime: "09:00:00",
			HandoffDay:  "monday",
		},
		Members: []RotationMember{
			{
				UserID: func() *string { s := "77779528-690b-4161-84ca-312e932c626e"; return &s }(),
				Name:   func() *string { s := "Frederick Graff"; return &s }(),
			},
		},
		Restrictions: []RotationRestriction{
			{
				StartDay:  "monday",
				StartTime: "09:00:00",
				EndDay:    "friday",
				EndTime:   "17:00:00",
			},
		},
	}
}

func expectedRotationResponseJSON() string {
	return `{
  "id": "rotation-id",
  "name": "A pleasant rotation",
  "description": "Managed by Terraform. Contact @platform-eng for changes.",
  "color": "#ff0000",
  "slack_user_group_id": "slack-group-1",
	"prevent_shift_deletion": true,
	"coverage_gap_notification_interval": "an-appropriate-interval",
  "members": [
    {
      "id": "77779528-690b-4161-84ca-312e932c626e",
      "name": "Frederick Graff"
    }
  ],
  "team": {
    "id": "44498724-9ccf-4e9a-b18f-5458ffad979a",
    "name": "Philadelphia"
  },
 "strategy": {
    "handoff_day": "monday",
    "handoff_time": "09:00:00",
    "type": "weekly"
  },
  "time_zone": "America/New_York",
  "restrictions": [
    {
      "end_day": "friday",
      "end_time": "17:00:00",
      "start_day": "monday",
      "start_time": "09:00:00"
    }
  ]
}`
}

func TestRotationGet(t *testing.T) {
	var requestPath string
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPath = req.URL.Path
		w.Write([]byte(expectedRotationResponseJSON()))
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	res, err := c.Rotations().Get(context.TODO(), "team-id", "schedule-id", "rotation-id")
	if err != nil {
		t.Fatalf("error retrieving rotation: %s", err.Error())
	}

	if expected := "/teams/team-id/on_call_schedules/schedule-id/rotations/rotation-id"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}

	expectedResponse := expectedRotationResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestRotationUpdateHasEffectiveAt(t *testing.T) {
	var requestPath string
	var requestBody map[string]any
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPath = req.URL.Path
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			t.Fatalf("error unmarshalling request body: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(expectedRotationResponseJSON()))
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	updateReq := UpdateRotationRequest{
		Name:             "A pleasant on-call schedule",
		Description:      "Managed by Terraform. Contact @platform-eng for changes.",
		SlackUserGroupID: "slack-group-1",
		Members:          []RotationMember{{UserID: func() *string { s := "77779528-690b-4161-84ca-312e932c626e"; return &s }()}},
	}
	if _, err := c.Rotations().Update(context.TODO(), "team-id", "schedule-id", "rotation-id", updateReq); err != nil {
		t.Fatalf("error retrieving on-call schedule: %s", err.Error())
	}

	if expected := "/teams/team-id/on_call_schedules/schedule-id/rotations/rotation-id"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}
	effectiveAt := requestBody["effective_at"].(string)
	if effectiveAt == "" {
		t.Fatalf("expected effective_at to be set")
	}
	e, err := time.Parse(time.RFC3339, effectiveAt)
	if err != nil {
		t.Fatalf("error parsing effective_at: %s", err.Error())
	}
	if dur := time.Since(e); dur > time.Minute {
		t.Fatalf("expected effective_at to be now-ish, found %s ago", dur)
	}
}
