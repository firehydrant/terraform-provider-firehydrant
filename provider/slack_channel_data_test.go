package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func offlineSlackChannelsMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`{
  "data": [{"id":"00000000-0000-4000-8000-000000000000","name":"#team-rocket","slack_channel_id":"C01010101Z"}],
  "pagination": {"count":1,"page":1,"items":1,"pages":1,"last":1,"prev":null,"next":null}
}`))
	}))
}

func TestOfflineSlackChannelsReadMemberID(t *testing.T) {
	ts := offlineSlackChannelsMockServer()
	defer ts.Close()

	c, err := firehydrant.NewRestClient("test-token-very-authorized", firehydrant.WithBaseURL(ts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	r := schema.TestResourceDataRaw(t, dataSourceSlackChannel().Schema, map[string]interface{}{
		"slack_channel_id":   "C01010101Z",
		"slack_channel_name": "team-rocket",
	})

	d := dataFireHydrantSlackChannelRead(context.Background(), r, c)
	if d.HasError() {
		t.Fatalf("error reading on-call schedule: %v", d)
	}
	if id := r.Id(); id != "00000000-0000-4000-8000-000000000000" {
		t.Fatalf("expected ID to be 00000000-0000-4000-8000-000000000000, got %s", id)
	}
}
