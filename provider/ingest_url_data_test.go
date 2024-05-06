package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func offlineIngestURLMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`{ "url":"https://signals.firehydrant.com/v1/process/some-long-jwt" }`))
	}))
}
func offlineTransposerMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(`{
			"data":[
				{"name": "Valid Transposer", "slug": "valid-transposer", "example_payload": "", "expression": "", "expected": "", 
					"website": "", "description": "", "tags": [""], "ingest_url": "https://signals.firehydrant.com/v1/transpose/valid-transposer/some-long-jwt"}
			]
		}`))
	}))
}

func TestOfflineIngestURL_UserID(t *testing.T) {
	tis := offlineIngestURLMockServer()
	defer tis.Close()

	c, err := firehydrant.NewRestClient("test-token-very-authorized", firehydrant.WithBaseURL(tis.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	r := schema.TestResourceDataRaw(t, dataSourceIngestURL().Schema, map[string]interface{}{
		"user_id": "00000000-0000-4000-8000-000000000000",
	})

	d := dataFireHydrantIngestURL(context.Background(), r, c)
	if d.HasError() {
		t.Fatalf("error reading ingest URL: %v", d)
	}

	url := r.Get("url")
	if url != nil {
		if url.(string) != "https://signals.firehydrant.com/v1/process/some-long-jwt" {
			t.Fatalf("expected URL to be https://signals.firehydrant.com/v1/process/some-long-jwt, got %s", url.(string))
		}
	} else {
		t.Fatal("attribute url not present")
	}
	id := r.Id()
	if id == "" {
		t.Fatal("ID cannot be empty")
	}
}

func TestOfflineIngestURL_ValidTransposer(t *testing.T) {
	tts := offlineTransposerMockServer()
	defer tts.Close()

	c, err := firehydrant.NewRestClient("test-token-very-authorized", firehydrant.WithBaseURL(tts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	r := schema.TestResourceDataRaw(t, dataSourceIngestURL().Schema, map[string]interface{}{
		"user_id":    "00000000-0000-4000-8000-000000000000",
		"transposer": "valid-transposer",
	})

	d := dataFireHydrantIngestURL(context.Background(), r, c)
	if d.HasError() {
		t.Fatalf("error reading ingest URL: %v", d)
	}

	url := r.Get("url")
	if url != nil {
		if url.(string) != "https://signals.firehydrant.com/v1/transpose/valid-transposer/some-long-jwt" {
			t.Fatalf("expected URL to be https://signals.firehydrant.com/v1/transpose/valid-transposer/some-long-jwt, got %s", url.(string))
		}
	} else {
		t.Fatal("attribute url not present")
	}
	id := r.Id()
	if id == "" {
		t.Fatal("ID cannot be empty")
	}
}

func TestOfflineIngestURL_InvalidAttributes(t *testing.T) {
	tts := offlineTransposerMockServer()
	defer tts.Close()

	c, err := firehydrant.NewRestClient("test-token-very-authorized", firehydrant.WithBaseURL(tts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	r := schema.TestResourceDataRaw(t, dataSourceIngestURL().Schema, map[string]interface{}{
		"on_call_schedule_id": "00000000-0000-4000-8000-000000000000",
		"transposer":          "valid-transposer",
	})

	d := dataFireHydrantIngestURL(context.Background(), r, c)
	if !d.HasError() {
		t.Fatalf("didnt fail on reading ingest URL: %v", d)
	}
}

func TestOfflineIngestURL_InvalidTransposer(t *testing.T) {
	tts := offlineTransposerMockServer()
	defer tts.Close()

	c, err := firehydrant.NewRestClient("test-token-very-authorized", firehydrant.WithBaseURL(tts.URL))
	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}
	r := schema.TestResourceDataRaw(t, dataSourceIngestURL().Schema, map[string]interface{}{
		"team_id":    "00000000-0000-4000-8000-000000000000",
		"transposer": "invalid-transposer",
	})

	d := dataFireHydrantIngestURL(context.Background(), r, c)
	if !d.HasError() {
		t.Fatalf("didnt fail on reading ingest URL: %v", d)
	}
}
