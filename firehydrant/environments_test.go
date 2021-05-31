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

	_, err = c.Services().List(context.TODO(), qry)
	if err != nil {
		t.Fatalf("Received error hitting ping endpoint: %s", err.Error())
	}

	if expected := "/services?labels=key1%3Dval1%2Ckey2%3Dval2&query=hello-world"; expected != requestPathRcvd {
		t.Fatalf("Expected %s, Got: %s for request path", expected, requestPathRcvd)
	}
}
