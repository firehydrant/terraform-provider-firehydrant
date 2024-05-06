package firehydrant

import (
	"net/http"
	"net/http/httptest"
	"time"
)

func expectedSignalsRuleResponse() *SignalsRuleResponse {
	return &SignalsRuleResponse{
		ID:         "00000000-0000-4000-8000-000000000000",
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func expectedSignalsRuleJSON() string {
	return `{
		"id": "00000000-0000-4000-8000-000000000000",
		"name": "Test Rule",
		"expression": "signal.summary.contains(\"foo\")",
		"target": {
			"id": "00000000-0000-4000-8000-000000000000",
			"type": "User",
		},
		"created_at": "2024-01-01T12:00:00.000Z",
		"updated_at": "2024-01-01T12:00:00.000Z",
		"incident_type": {
			"id": "00000000-0000-4000-8000-000000000000",
			"name": "Test incident type"
		}
	}`
}

func transposerMockServer(path *string) *httptest.Server {
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		*path = req.URL.Path

		if *path == "teams/00000000-0000-4000-8000-000000000000/signal_rules/00000000-0000-4000-8000-000000000000" ||
			*path == "teams/00000000-0000-4000-8000-000000000000/signal_rules" {
			w.Write([]byte(expectedSignalsRuleJSON()))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	ts := httptest.NewServer(h)
	return ts
}
