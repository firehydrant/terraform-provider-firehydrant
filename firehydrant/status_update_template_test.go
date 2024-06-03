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

func expectedStatusUpdateTemplateResponse() *StatusUpdateTemplateResponse {
	timestamp, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00.000Z")
	return &StatusUpdateTemplateResponse{
		ID:        "00000000-0000-8000-4000-000000000000",
		Name:      "New Template",
		Body:      "The template body",
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}
}

func expectedStatusUpdateTemplateJSON() string {
	return `{
		"id": "00000000-0000-8000-4000-000000000000",
		"name": "New Template",
		"body": "The template body",
		"created_at": "2024-01-01T12:00:00.000Z",
		"updated_at": "2024-01-01T12:00:00.000Z"
	}`
}

func statusUpdateTemplateMockServer(path *string) *httptest.Server {
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		*path = req.URL.Path

		if strings.Contains(*path, "status_update_templates/00000000-0000-8000-4000-000000000000") {
			w.Write([]byte(expectedStatusUpdateTemplateJSON()))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	ts := httptest.NewServer(h)
	return ts
}

func TestStatusUpdateTemplateGet(t *testing.T) {
	var requestPath string
	id := "00000000-0000-8000-4000-000000000000"
	ts := statusUpdateTemplateMockServer(&requestPath)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	res, err := c.StatusUpdateTemplates().Get(context.Background(), id)
	if err != nil {
		t.Fatalf("error retrieving status update template: %s", err.Error())
	}

	if expected := fmt.Sprintf("/status_update_templates/%s", id); expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}

	expectedResponse := expectedStatusUpdateTemplateResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestStatusUpdateTemplateGet_NotFound(t *testing.T) {
	var requestPath string
	id := "00000000-0000-4000-4000-000000000000"
	ts := statusUpdateTemplateMockServer(&requestPath)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	_, err = c.StatusUpdateTemplates().Get(context.Background(), id)
	if err == nil {
		t.Fatalf("expected ErrorNotFound in retrieving status update template, got nil")
	}
}
