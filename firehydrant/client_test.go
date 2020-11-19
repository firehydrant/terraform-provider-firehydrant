package firehydrant

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/go-querystring/query"
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

	res, err := c.Ping(context.TODO())
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
	res, err := c.Services().Get(context.TODO(), testServiceID)
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
	res, err := c.Services().Get(context.TODO(), testServiceID)
	if err != nil {
		t.Fatalf("Received error hitting get service endpoint: %s", err.Error())
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

func TestGetEnvironment(t *testing.T) {
	var requestPathRcvd string

	expectedEnvironment := EnvironmentResponse{
		ID:          "test-id",
		Name:        "test environment",
		Description: "this environment causes people to forget to share their screen",
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPathRcvd = req.URL.Path

		if err := json.NewEncoder(w).Encode(expectedEnvironment); err != nil {
			panic(err)
		}
	})
	ts := httptest.NewServer(h)

	defer ts.Close()

	testToken := "testing-123"
	c, err := NewRestClient(testToken, WithBaseURL(ts.URL))

	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	res, err := c.GetEnvironment(context.TODO(), expectedEnvironment.ID)
	if err != nil {
		t.Fatalf("Received error hitting environment get endpoint: %s", err.Error())
	}

	if expected := "/environments/" + expectedEnvironment.ID; expected != requestPathRcvd {
		t.Fatalf("Expected %s, Got: %s for request path", expected, requestPathRcvd)
	}

	if !reflect.DeepEqual(&expectedEnvironment, res) {
		t.Fatalf("Expected %+v, Got: %+v for response", expectedEnvironment, res)
	}
}

func TestCreateEnvironment(t *testing.T) {
	var requestPathRcvd string

	req := CreateEnvironmentRequest{
		Name:        "test environment",
		Description: "this environment causes people to forget to share their screen",
	}

	resp := EnvironmentResponse{
		ID:          "test-id",
		Name:        req.Name,
		Description: req.Description,
	}

	var rcvdEnv EnvironmentResponse

	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPathRcvd = req.URL.Path

		if err := json.NewDecoder(req.Body).Decode(&rcvdEnv); err != nil {
			panic(err)
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
	})
	ts := httptest.NewServer(h)

	defer ts.Close()

	testToken := "testing-123"
	c, err := NewRestClient(testToken, WithBaseURL(ts.URL))

	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	res, err := c.CreateEnvironment(context.TODO(), req)
	if err != nil {
		t.Fatalf("Received error hitting environment create endpoint: %s", err.Error())
	}

	if expected := "/environments"; expected != requestPathRcvd {
		t.Fatalf("Expected %s, Got: %s for request path", expected, requestPathRcvd)
	}

	if !reflect.DeepEqual(&resp, res) {
		t.Fatalf("Expected %+v, Got: %+v for response", resp, res)
	}
}

func TestGetServices(t *testing.T) {
	var requestPathRcvd string
	response := ServicesResponse{
		Services: []ServiceResponse{
			{
				ID: "hello-world",
			},
		},
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestPathRcvd = req.URL.Path + "?" + req.URL.RawQuery

		if err := json.NewEncoder(w).Encode(&response); err != nil {
			panic(err)
		}
	})
	ts := httptest.NewServer(h)

	defer ts.Close()

	testToken := "testing-123"
	c, err := NewRestClient(testToken, WithBaseURL(ts.URL))

	if err != nil {
		t.Fatalf("Received error initializing API client: %s", err.Error())
		return
	}

	qry := &ServiceQuery{
		Query: "hello-world",
		LabelsSelector: LabelsSelector{
			"key1": "val1",
			"key2": "val2",
		},
	}

	vs, err := query.Values(qry)
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log(vs)

	_, err = c.Services().Get(context.TODO(), qry)
	if err != nil {
		t.Fatalf("Received error hitting ping endpoint: %s", err.Error())
	}

	if expected := "/services?labels=key1%3Dval1%2Ckey2%3Dval2&query=hello-world"; expected != requestPathRcvd {
		t.Fatalf("Expected %s, Got: %s for request path", expected, requestPathRcvd)
	}
}
