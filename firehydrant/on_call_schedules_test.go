package firehydrant

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func expectedOnCallScheduleResponse() *OnCallScheduleResponse {
	return &OnCallScheduleResponse{
		ID:          "schedule-id",
		Name:        "A pleasant on-call schedule",
		Description: "Managed by Terraform. Contact @platform-eng for changes.",
		TimeZone:    "America/New_York",
		Strategy: OnCallScheduleStrategy{
			Type:        "weekly",
			HandoffTime: "09:00:00",
			HandoffDay:  "monday",
		},
		Members: []OnCallScheduleMember{
			{
				ID:   "77779528-690b-4161-84ca-312e932c626e",
				Name: "Frederick Graff",
			},
		},
		Restrictions: []OnCallScheduleRestriction{
			{
				StartDay:  "monday",
				StartTime: "09:00:00",
				EndDay:    "friday",
				EndTime:   "17:00:00",
			},
		},
	}
}

func expectedOnCallScheduleResponseJSON() string {
	return `{
  "id": "schedule-id",
  "name": "A pleasant on-call schedule",
  "description": "Managed by Terraform. Contact @platform-eng for changes.",
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

func TestOnCallSchedulesGet(t *testing.T) {
	var requestPath string
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPath = req.URL.Path
		w.Write([]byte(expectedOnCallScheduleResponseJSON()))
	})

	ts := httptest.NewServer(h)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	res, err := c.OnCallSchedules().Get(context.TODO(), "team-id", "schedule-id")
	if err != nil {
		t.Fatalf("error retrieving on-call schedule: %s", err.Error())
	}

	if expected := "/teams/team-id/on_call_schedules/schedule-id"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}

	expectedResponse := expectedOnCallScheduleResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}
