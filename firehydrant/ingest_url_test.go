package firehydrant

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func expectedIngestURLResponse() *IngestURLResponse {
	return &IngestURLResponse{
		URL: "https://signals.firehydrant.com/v1/process/some-long-jwt",
	}
}

func expectedIngestURLResponseJSON() string {
	return `{
	"url":"https://signals.firehydrant.com/v1/process/some-long-jwt" 
}`
}

func ingestURLMockServer(path, userQuery, teamQuery, escalationPolicyQuery, scheduleQuery *string) *httptest.Server {
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		*path = req.URL.Path
		*userQuery = req.URL.Query().Get("user_id")
		*teamQuery = req.URL.Query().Get("team_id")
		*escalationPolicyQuery = req.URL.Query().Get("escalation_policy_id")
		*scheduleQuery = req.URL.Query().Get("on_call_schedule_id")

		if *userQuery == "" && *teamQuery == "" {
			w.Write([]byte(expectedIngestURLResponseJSON()))
		} else if *userQuery == "00000000-0000-4000-8000-000000000000" {
			w.Write([]byte(expectedIngestURLResponseJSON()))
		} else if *teamQuery == "00000000-0000-4000-8000-000000000000" {
			if *escalationPolicyQuery == "00000000-0000-4000-8000-000000000000" {
				w.Write([]byte(expectedIngestURLResponseJSON()))
			} else if *scheduleQuery == "00000000-0000-4000-8000-000000000000" {
				w.Write([]byte(expectedIngestURLResponseJSON()))
			} else if *escalationPolicyQuery == "" && *scheduleQuery == "" {
				w.Write([]byte(expectedIngestURLResponseJSON()))
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	ts := httptest.NewServer(h)
	return ts
}

func TestIngestURLGet_Default(t *testing.T) {
	var requestPath, userQuery, teamQuery, escalationPolicyQuery, scheduleQuery string
	ts := ingestURLMockServer(&requestPath, &userQuery, &teamQuery, &escalationPolicyQuery, &scheduleQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := IngestURLParams{
		UserID:             "",
		TeamID:             "",
		EscalationPolicyID: "",
		OnCallScheduleID:   "",
	}
	res, err := c.IngestURL().Get(context.Background(), params)
	if err != nil {
		t.Fatalf("error retrieving ingest URL: %s", err.Error())
	}

	if expected := "/signals/ingest_url"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}
	if expected := ""; expected != userQuery {
		t.Fatalf("request user query params mismatch: expected '%s', got: '%s'", expected, userQuery)
	}
	if expected := ""; expected != teamQuery {
		t.Fatalf("request team query params mismatch: expected '%s', got: '%s'", expected, teamQuery)
	}

	expectedResponse := expectedIngestURLResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestIngestURLGet_UserID(t *testing.T) {
	var requestPath, userQuery, teamQuery, escalationPolicyQuery, scheduleQuery string
	ts := ingestURLMockServer(&requestPath, &userQuery, &teamQuery, &escalationPolicyQuery, &scheduleQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := IngestURLParams{
		UserID:             "00000000-0000-4000-8000-000000000000",
		TeamID:             "",
		EscalationPolicyID: "",
		OnCallScheduleID:   "",
	}
	res, err := c.IngestURL().Get(context.Background(), params)
	if err != nil {
		t.Fatalf("error retrieving ingest URL: %s", err.Error())
	}

	if expected := "/signals/ingest_url"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}
	if expected := "00000000-0000-4000-8000-000000000000"; expected != userQuery {
		t.Fatalf("request user query params mismatch: expected '%s', got: '%s'", expected, userQuery)
	}

	expectedResponse := expectedIngestURLResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestIngestURLGet_TeamID(t *testing.T) {
	var requestPath, userQuery, teamQuery, escalationPolicyQuery, scheduleQuery string
	ts := ingestURLMockServer(&requestPath, &userQuery, &teamQuery, &escalationPolicyQuery, &scheduleQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := IngestURLParams{
		UserID:             "",
		TeamID:             "00000000-0000-4000-8000-000000000000",
		EscalationPolicyID: "",
		OnCallScheduleID:   "",
	}
	res, err := c.IngestURL().Get(context.Background(), params)
	if err != nil {
		t.Fatalf("error retrieving ingest URL: %s", err.Error())
	}

	if expected := "/signals/ingest_url"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}
	if expected := "00000000-0000-4000-8000-000000000000"; expected != teamQuery {
		t.Fatalf("request team query params mismatch: expected '%s', got: '%s'", expected, teamQuery)
	}

	expectedResponse := expectedIngestURLResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestIngestURLGet_EscalationPolicyID(t *testing.T) {
	var requestPath, userQuery, teamQuery, escalationPolicyQuery, scheduleQuery string
	ts := ingestURLMockServer(&requestPath, &userQuery, &teamQuery, &escalationPolicyQuery, &scheduleQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := IngestURLParams{
		UserID:             "",
		TeamID:             "00000000-0000-4000-8000-000000000000",
		EscalationPolicyID: "00000000-0000-4000-8000-000000000000",
		OnCallScheduleID:   "",
	}
	res, err := c.IngestURL().Get(context.Background(), params)
	if err != nil {
		t.Fatalf("error retrieving ingest URL: %s", err.Error())
	}

	if expected := "/signals/ingest_url"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}
	if expected := "00000000-0000-4000-8000-000000000000"; expected != teamQuery {
		t.Fatalf("request team query params mismatch: expected '%s', got: '%s'", expected, teamQuery)
	}
	if expected := "00000000-0000-4000-8000-000000000000"; expected != escalationPolicyQuery {
		t.Fatalf("request escalation policy query params mismatch: expected '%s', got: '%s'", expected, escalationPolicyQuery)
	}

	expectedResponse := expectedIngestURLResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestIngestURLGet_ScheduleID(t *testing.T) {
	var requestPath, userQuery, teamQuery, escalationPolicyQuery, scheduleQuery string
	ts := ingestURLMockServer(&requestPath, &userQuery, &teamQuery, &escalationPolicyQuery, &scheduleQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := IngestURLParams{
		UserID:             "",
		TeamID:             "00000000-0000-4000-8000-000000000000",
		EscalationPolicyID: "",
		OnCallScheduleID:   "00000000-0000-4000-8000-000000000000",
	}
	res, err := c.IngestURL().Get(context.Background(), params)
	if err != nil {
		t.Fatalf("error retrieving ingest URL: %s", err.Error())
	}

	if expected := "/signals/ingest_url"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}
	if expected := "00000000-0000-4000-8000-000000000000"; expected != teamQuery {
		t.Fatalf("request team query params mismatch: expected '%s', got: '%s'", expected, teamQuery)
	}
	if expected := "00000000-0000-4000-8000-000000000000"; expected != scheduleQuery {
		t.Fatalf("request schedule query params mismatch: expected '%s', got: '%s'", expected, scheduleQuery)
	}

	expectedResponse := expectedIngestURLResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestIngestURLGet_NotFound(t *testing.T) {
	var requestPath, userQuery, teamQuery, escalationPolicyQuery, scheduleQuery string
	ts := ingestURLMockServer(&requestPath, &userQuery, &teamQuery, &escalationPolicyQuery, &scheduleQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := IngestURLParams{
		UserID:             "",
		TeamID:             "11111111-0000-4000-8000-000000000000",
		EscalationPolicyID: "",
		OnCallScheduleID:   "",
	}
	_, err = c.IngestURL().Get(context.Background(), params)
	if err == nil {
		t.Fatalf("expected ErrorNotFound in retrieving ingest URL, got nil")
	}
	if !errors.Is(err, ErrorNotFound) {
		t.Fatalf("expected ErrorNotFound in retrieving ingest URL, got: %s", err)
	}
}
