package firehydrant

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// we're only using this to for ingest URLs at this point.  If we need to get other transposer information, we should add
// additional tests and test data to confirm that is being handled properly

func expectedTransposersResponse() *TransposersResponse {
	t := Transposer{
		Name:           "Valid Transposer",
		Slug:           "valid-transposer",
		ExamplePayload: "",
		Expression:     "",
		Expected:       "",
		Website:        "",
		Description:    "",
		Tags:           []string{""},
		IngestURL:      "https://signals.firehydrant.com/v1/transpose/valid-transposer/some-long-jwt",
	}
	return &TransposersResponse{
		Transposers: []Transposer{t},
	}
}

func expectedTransposerResponseJSON() string {
	return `{
	"data":[
		{"name": "Valid Transposer", "slug": "valid-transposer", "example_payload": "", "expression": "", "expected": "", 
			"website": "", "description": "", "tags": [""], "ingest_url": "https://signals.firehydrant.com/v1/transpose/valid-transposer/some-long-jwt"}
	]
}`
}

func transposerMockServer(path, userQuery, teamQuery, escalationPolicyQuery, scheduleQuery *string) *httptest.Server {
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		*path = req.URL.Path
		*userQuery = req.URL.Query().Get("user_id")
		*teamQuery = req.URL.Query().Get("team_id")
		*escalationPolicyQuery = req.URL.Query().Get("escalation_policy_id")
		*scheduleQuery = req.URL.Query().Get("on_call_schedule_id")

		if *userQuery == "" && *teamQuery == "" {
			w.Write([]byte(expectedTransposerResponseJSON()))
		} else if *userQuery == "00000000-0000-4000-8000-000000000000" {
			w.Write([]byte(expectedTransposerResponseJSON()))
		} else if *teamQuery == "00000000-0000-4000-8000-000000000000" {
			if *escalationPolicyQuery == "00000000-0000-4000-8000-000000000000" {
				w.Write([]byte(expectedTransposerResponseJSON()))
			} else if *scheduleQuery == "00000000-0000-4000-8000-000000000000" {
				w.Write([]byte(expectedTransposerResponseJSON()))
			} else if *escalationPolicyQuery == "" && *scheduleQuery == "" {
				w.Write([]byte(expectedTransposerResponseJSON()))
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	ts := httptest.NewServer(h)
	return ts
}

func TestTransposerGet_Default(t *testing.T) {
	var requestPath, userQuery, teamQuery, escalationPolicyQuery, scheduleQuery string
	ts := transposerMockServer(&requestPath, &userQuery, &teamQuery, &escalationPolicyQuery, &scheduleQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := TransposersParams{
		UserID:             "",
		TeamID:             "",
		EscalationPolicyID: "",
		OnCallScheduleID:   "",
	}
	res, err := c.Transposers().Get(context.Background(), params)
	if err != nil {
		t.Fatalf("error retrieving transposers: %s", err.Error())
	}

	if expected := "/signals/transposers"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}
	if expected := ""; expected != userQuery {
		t.Fatalf("request user query params mismatch: expected '%s', got: '%s'", expected, userQuery)
	}
	if expected := ""; expected != teamQuery {
		t.Fatalf("request team query params mismatch: expected '%s', got: '%s'", expected, teamQuery)
	}

	expectedResponse := expectedTransposersResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestTransposersGet_UserID(t *testing.T) {
	var requestPath, userQuery, teamQuery, escalationPolicyQuery, scheduleQuery string
	ts := transposerMockServer(&requestPath, &userQuery, &teamQuery, &escalationPolicyQuery, &scheduleQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := TransposersParams{
		UserID:             "00000000-0000-4000-8000-000000000000",
		TeamID:             "",
		EscalationPolicyID: "",
		OnCallScheduleID:   "",
	}
	res, err := c.Transposers().Get(context.Background(), params)
	if err != nil {
		t.Fatalf("error retrieving transposers: %s", err.Error())
	}

	if expected := "/signals/transposers"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}
	if expected := "00000000-0000-4000-8000-000000000000"; expected != userQuery {
		t.Fatalf("request user query params mismatch: expected '%s', got: '%s'", expected, userQuery)
	}

	expectedResponse := expectedTransposersResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestTransposersGet_TeamID(t *testing.T) {
	var requestPath, userQuery, teamQuery, escalationPolicyQuery, scheduleQuery string
	ts := transposerMockServer(&requestPath, &userQuery, &teamQuery, &escalationPolicyQuery, &scheduleQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := TransposersParams{
		UserID:             "",
		TeamID:             "00000000-0000-4000-8000-000000000000",
		EscalationPolicyID: "",
		OnCallScheduleID:   "",
	}
	res, err := c.Transposers().Get(context.Background(), params)
	if err != nil {
		t.Fatalf("error retrieving transposersL: %s", err.Error())
	}

	if expected := "/signals/transposers"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}
	if expected := "00000000-0000-4000-8000-000000000000"; expected != teamQuery {
		t.Fatalf("request team query params mismatch: expected '%s', got: '%s'", expected, teamQuery)
	}

	expectedResponse := expectedTransposersResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestTransposersGet_EscalationPolicyID(t *testing.T) {
	var requestPath, userQuery, teamQuery, escalationPolicyQuery, scheduleQuery string
	ts := transposerMockServer(&requestPath, &userQuery, &teamQuery, &escalationPolicyQuery, &scheduleQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := TransposersParams{
		UserID:             "",
		TeamID:             "00000000-0000-4000-8000-000000000000",
		EscalationPolicyID: "00000000-0000-4000-8000-000000000000",
		OnCallScheduleID:   "",
	}
	res, err := c.Transposers().Get(context.Background(), params)
	if err != nil {
		t.Fatalf("error retrieving transposers: %s", err.Error())
	}

	if expected := "/signals/transposers"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}
	if expected := "00000000-0000-4000-8000-000000000000"; expected != teamQuery {
		t.Fatalf("request team query params mismatch: expected '%s', got: '%s'", expected, teamQuery)
	}
	if expected := "00000000-0000-4000-8000-000000000000"; expected != escalationPolicyQuery {
		t.Fatalf("request escalation policy query params mismatch: expected '%s', got: '%s'", expected, escalationPolicyQuery)
	}

	expectedResponse := expectedTransposersResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestTransposersGet_ScheduleID(t *testing.T) {
	var requestPath, userQuery, teamQuery, escalationPolicyQuery, scheduleQuery string
	ts := transposerMockServer(&requestPath, &userQuery, &teamQuery, &escalationPolicyQuery, &scheduleQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := TransposersParams{
		UserID:             "",
		TeamID:             "00000000-0000-4000-8000-000000000000",
		EscalationPolicyID: "",
		OnCallScheduleID:   "00000000-0000-4000-8000-000000000000",
	}
	res, err := c.Transposers().Get(context.Background(), params)
	if err != nil {
		t.Fatalf("error retrieving transposers: %s", err.Error())
	}

	if expected := "/signals/transposers"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}
	if expected := "00000000-0000-4000-8000-000000000000"; expected != teamQuery {
		t.Fatalf("request team query params mismatch: expected '%s', got: '%s'", expected, teamQuery)
	}
	if expected := "00000000-0000-4000-8000-000000000000"; expected != scheduleQuery {
		t.Fatalf("request schedule query params mismatch: expected '%s', got: '%s'", expected, scheduleQuery)
	}

	expectedResponse := expectedTransposersResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestTransposersGet_NotFound(t *testing.T) {
	var requestPath, userQuery, teamQuery, escalationPolicyQuery, scheduleQuery string
	ts := transposerMockServer(&requestPath, &userQuery, &teamQuery, &escalationPolicyQuery, &scheduleQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := TransposersParams{
		UserID:             "",
		TeamID:             "11111111-0000-4000-8000-000000000000",
		EscalationPolicyID: "",
		OnCallScheduleID:   "",
	}
	_, err = c.Transposers().Get(context.Background(), params)
	if err == nil {
		t.Fatalf("expected ErrorNotFound in retrieving transposers, got nil")
	}
	if !errors.Is(err, ErrorNotFound) {
		t.Fatalf("expected ErrorNotFound in retrieving transposers, got: %s", err)
	}
}
