package firehydrant

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func expectedSlackChannelResponse() *SlackChannelResponse {
	return &SlackChannelResponse{
		ID:             "00000000-0000-4000-8000-000000000000",
		Name:           "#team-rocket",
		SlackChannelID: "C01010101Z",
	}
}

func expectedSlackChannelsResponseJSON() string {
	return `{
	"data": [{"id":"00000000-0000-4000-8000-000000000000","name":"#team-rocket","slack_channel_id":"C01010101Z"}],
	"pagination": {"count":1,"page":1,"items":1,"pages":1,"last":1,"prev":null,"next":null}
}`
}

func slackChannelMockServer(path, query *string) *httptest.Server {
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		*path = req.URL.Path
		*query = req.URL.Query().Get("slack_channel_id")
		if *query == "C01010101Z" {
			w.Write([]byte(expectedSlackChannelsResponseJSON()))
		} else {
			*query = req.URL.Query().Get("name")
			if *query == "team-rocket" {
				w.Write([]byte(expectedSlackChannelsResponseJSON()))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}
	})

	ts := httptest.NewServer(h)
	return ts
}

func TestSlackChannelGet_ID(t *testing.T) {
	var requestPath, requestQuery string
	ts := slackChannelMockServer(&requestPath, &requestQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := SlackChannelParams{ID: "C01010101Z"}
	res, err := c.SlackChannels().Get(context.Background(), params)
	if err != nil {
		t.Fatalf("error retrieving slack channel: %s", err.Error())
	}

	if expected := "/integrations/slack/channels"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}
	if expected := "C01010101Z"; expected != requestQuery {
		t.Fatalf("request query params mismatch: expected '%s', got: '%s'", expected, requestQuery)
	}

	expectedResponse := expectedSlackChannelResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestSlackChannelGet_Name(t *testing.T) {
	var requestPath, requestQuery string
	ts := slackChannelMockServer(&requestPath, &requestQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := SlackChannelParams{Name: "#team-rocket"}
	res, err := c.SlackChannels().Get(context.Background(), params)
	if err != nil {
		t.Fatalf("error retrieving slack channel: %s", err.Error())
	}

	if expected := "/integrations/slack/channels"; expected != requestPath {
		t.Fatalf("request path mismatch: expected '%s', got: '%s'", expected, requestPath)
	}
	if expected := "team-rocket"; expected != requestQuery {
		t.Fatalf("request query params mismatch: expected '%s', got: '%s'", expected, requestQuery)
	}

	expectedResponse := expectedSlackChannelResponse()
	if !reflect.DeepEqual(expectedResponse, res) {
		t.Fatalf("response mismatch: expected '%+v', got: '%+v'", expectedResponse, res)
	}
}

func TestSlackChannelGetNotFound(t *testing.T) {
	var requestPath, requestQuery string
	ts := slackChannelMockServer(&requestPath, &requestQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := SlackChannelParams{ID: "C11111111"}
	_, err = c.SlackChannels().Get(context.Background(), params)
	if err == nil {
		t.Fatalf("expected ErrorNotFound in retrieving slack channel, got nil")
	}
	if !errors.Is(err, ErrorNotFound) {
		t.Fatalf("expected ErrorNotFound in retrieving slack channel, got: %s", err)
	}
}

func TestSlackChannelGetNotFound_noParams(t *testing.T) {
	var requestPath, requestQuery string
	ts := slackChannelMockServer(&requestPath, &requestQuery)
	defer ts.Close()

	c, err := NewRestClient("test-token-very-authorized", WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	params := SlackChannelParams{ID: "", Name: ""}
	_, err = c.SlackChannels().Get(context.Background(), params)
	if err == nil {
		t.Fatalf("expected ErrorNotFound in retrieving slack channel, got nil")
	}
}
