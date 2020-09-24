package firehydrant

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	pingResponseJSON    = `{"response":"pong","actor":{"id":"2af3339f-9d81-434b-a208-427d6d85c124","name":"Bobby Tables","email":"bobby+dalmatians@firehydrant.io","type":"firehydrant_user"}}`
	serviceResponseJSON = `{"id": "da4bd45b-2b68-4c05-8564-d08dc7725291", "name": "Chow Hall", "description": "", "slug": "chow-hall", "created_at": "2019-07-30T13:02:22.243Z", "updated_at": "2019-12-09T23:59:18.094Z", "labels": {}}`
)

func TestClientInitialization(t *testing.T) {
	var requestPathRcvd, token string

	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPathRcvd = req.URL.Path
		token = req.Header.Get("Authorization")

		w.Write([]byte(pingResponseJSON))
	})
	ts := httptest.NewServer(h)

	defer ts.Close()

	testToken := "testing-123"
	c, err := NewRestClient(testToken, WithBaseURL(ts.URL))

	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	res, err := c.Ping()
	if err != nil {
		t.Fatalf("Received error hitting ping endpoint: %s", err.Error())
	}

	actorID := res.Actor.ID
	actorEmail := res.Actor.Email

	if expected := "/ping"; expected != requestPathRcvd {
		t.Fatalf("Expected %s, Got: %s for request path", expected, requestPathRcvd)
	}

	if expected := "Bearer " + testToken; expected != token {
		t.Fatalf("Expected %s, Got: %s for bearer token", expected, token)
	}

	if expected := "2af3339f-9d81-434b-a208-427d6d85c124"; expected != actorID {
		t.Fatalf("Expected %s, Got: %s for actor ID", expected, actorID)
	}

	if expected := "bobby+dalmatians@firehydrant.io"; expected != actorEmail {
		t.Fatalf("Expected %s, Got: %s for actor email", expected, actorEmail)
	}
}

func TestGetService(t *testing.T) {
	var requestPathRcvd string

	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPathRcvd = req.URL.Path

		w.Write([]byte(serviceResponseJSON))
	})
	ts := httptest.NewServer(h)

	defer ts.Close()

	testToken := "testing-123"
	c, err := NewRestClient(testToken, WithBaseURL(ts.URL))

	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	testServiceID := "test-service-id"
	res, err := c.GetService(testServiceID)
	if err != nil {
		t.Fatalf("Received error hitting ping endpoint: %s", err.Error())
	}

	serviceID := res.ID
	serviceName := res.Name

	if expected := "/services/" + testServiceID; expected != requestPathRcvd {
		t.Fatalf("Expected %s, Got: %s for request path", expected, requestPathRcvd)
	}

	if expected := "da4bd45b-2b68-4c05-8564-d08dc7725291"; expected != serviceID {
		t.Fatalf("Expected %s, Got: %s for service ID", expected, serviceID)
	}

	if expected := "Chow Hall"; expected != serviceName {
		t.Fatalf("Expected %s, Got: %s for service name", expected, serviceName)
	}
}

func TestCreateService(t *testing.T) {
	var requestPathRcvd string

	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPathRcvd = req.URL.Path

		w.Write([]byte(serviceResponseJSON))
	})
	ts := httptest.NewServer(h)

	defer ts.Close()

	testToken := "testing-123"
	c, err := NewRestClient(testToken, WithBaseURL(ts.URL))

	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	testServiceID := "test-service-id"
	res, err := c.GetService(testServiceID)
	if err != nil {
		t.Fatalf("Received error hitting ping endpoint: %s", err.Error())
	}

	serviceID := res.ID
	serviceName := res.Name

	if expected := "/services/" + testServiceID; expected != requestPathRcvd {
		t.Fatalf("Expected %s, Got: %s for request path", expected, requestPathRcvd)
	}

	if expected := "da4bd45b-2b68-4c05-8564-d08dc7725291"; expected != serviceID {
		t.Fatalf("Expected %s, Got: %s for service ID", expected, serviceID)
	}

	if expected := "Chow Hall"; expected != serviceName {
		t.Fatalf("Expected %s, Got: %s for service name", expected, serviceName)
	}
}
