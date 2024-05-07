package firehydrant

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func expectedSignalsRuleResponse() *SignalsRuleResponse {
	timestamp, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00.000Z")
	return &SignalsRuleResponse{
		ID:         "00000000-0000-8000-4000-000000000000",
		Name:       "Test Rule",
		Expression: "signal.summary.contains(\"foo\")",
		Target: SignalRuleTarget{
			ID:   "00000000-0000-4000-8000-000000000000",
			Type: "User",
		},
		IncidentType: SignalRuleIncidentType{
			ID:   "00000000-0000-4000-8000-000000000000",
			Name: "Test incident type",
		},
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}
}

func expectedSignalsRuleJSON() string {
	return `{
		"id": "00000000-0000-8000-4000-000000000000",
		"name": "Test Rule",
		"expression": "signal.summary.contains(\"foo\")",
		"target": {
			"id": "00000000-0000-4000-8000-000000000000",
			"type": "User"
		},
		"created_at": "2024-01-01T12:00:00.000Z",
		"updated_at": "2024-01-01T12:00:00.000Z",
		"incident_type": {
			"id": "00000000-0000-4000-8000-000000000000",
			"name": "Test incident type"
		}
	}`
}

func signalsRulesMockServer(path *string) *httptest.Server {
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		*path = req.URL.Path

		if strings.Contains(*path, "teams/00000000-0000-4000-8000-000000000000/signal_rules/00000000-0000-8000-4000-000000000000") {
			w.Write([]byte(expectedSignalsRuleJSON()))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	ts := httptest.NewServer(h)
	return ts
}

func TestSignalsRulesGet(t *testing.T) {
	var requestPath string
	team_id := "00000000-0000-4000-8000-000000000000"
	rule_id := "00000000-0000-8000-4000-000000000000"
	ts := signalsRulesMockServer(&requestPath)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	res, err := c.SignalsRules().Get(context.Background(), team_id, rule_id)
	if err != nil {
		t.Fatalf("error retrieving signal rule: %s", err.Error())
	}

	if expected := fmt.Sprintf("/teams/%s/signal_rules/%s", team_id, rule_id); expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}

	expectedResponse := expectedSignalsRuleResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestSignalsRulesGet_NotFound(t *testing.T) {
	var requestPath string
	team_id := "invalid_team_id"
	rule_id := "00000000-0000-8000-4000-000000000000"
	ts := signalsRulesMockServer(&requestPath)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	_, err = c.SignalsRules().Get(context.Background(), team_id, rule_id)
	if err == nil {
		t.Fatalf("expected ErrorNotFound in retrieving slack channel, got nil")
	}
}
