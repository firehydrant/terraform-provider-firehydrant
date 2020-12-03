package firehydrant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	pingResponseJSON    = `{"response":"pong","actor":{"id":"2af3339f-9d81-434b-a208-427d6d85c124","name":"Bobby Tables","email":"bobby+dalmatians@firehydrant.io","type":"firehydrant_user"}}`
	serviceResponseJSON = `{"id": "da4bd45b-2b68-4c05-8564-d08dc7725291", "name": "Chow Hall", "description": "", "slug": "chow-hall", "created_at": "2019-07-30T13:02:22.243Z", "updated_at": "2019-12-09T23:59:18.094Z", "labels": {}}`
)

type RequestTest func(req *http.Request)

func AssertRequestJSONBody(t *testing.T, src interface{}) RequestTest {
	return func(req *http.Request) {
		req.Body = ioutil.NopCloser(req.Body)

		buf := new(bytes.Buffer)
		require.NoError(t, json.NewEncoder(buf).Encode(src))

		// Read the body out so we can compare to what we received
		b, err := ioutil.ReadAll(req.Body)
		require.NoError(t, err)

		assert.Equal(t, buf.Bytes(), b)
	}
}

func AssertRequestMethod(t *testing.T, method string) RequestTest {
	return func(req *http.Request) {
		assert.Equal(t, method, req.Method)
	}
}

func setupClient(requestPath string, mockedResponse interface{}, requestTests ...RequestTest) (*APIClient, func(), error) {
	if err := faker.FakeData(mockedResponse); err != nil {
		return nil, nil, err
	}

	// We only handle the request path passed in the setup, this ensures that we serve
	// a 404 on any other request, failing the client in a more predictable and easier to
	// debug way
	mux := http.NewServeMux()
	mux.HandleFunc(requestPath, func(w http.ResponseWriter, req *http.Request) {
		if err := json.NewEncoder(w).Encode(mockedResponse); err != nil {
			panic(fmt.Errorf("could not encode JSON: %w", err))
		}

		for _, test := range requestTests {
			test(req)
		}
	})

	ts := httptest.NewServer(mux)

	c, err := NewRestClient("fake-token", WithBaseURL(ts.URL))
	if err != nil {
		return nil, nil, fmt.Errorf("could not generate rest client: %w", err)
	}

	teardown := func() {
		ts.Close()
	}

	return c, teardown, nil
}

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
	resp := &ServiceResponse{}
	testServiceID := "test-service-id"
	c, teardown, err := setupClient("/services/"+testServiceID, resp)
	require.NoError(t, err)
	defer teardown()

	res, err := c.Services().Get(context.TODO(), testServiceID)
	require.NoError(t, err, "error retrieving a service")
	assert.Equal(t, resp.ID, res.ID, "returned service did not match")
	assert.Equal(t, resp.Name, res.Name, "returned service did not match")
}

func TestCreateService(t *testing.T) {
	resp := &ServiceResponse{}
	c, teardown, err := setupClient("/services", resp,
		AssertRequestMethod(t, "POST"),
		AssertRequestJSONBody(t, CreateServiceRequest{Name: "fake-service"}),
	)

	require.NoError(t, err)
	defer teardown()

	_, err = c.Services().Create(context.TODO(), CreateServiceRequest{Name: "fake-service"})
	require.NoError(t, err, "error creating a service")
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

	_, err = c.Services().List(context.TODO(), qry)
	if err != nil {
		t.Fatalf("Received error hitting ping endpoint: %s", err.Error())
	}

	if expected := "/services?labels=key1%3Dval1%2Ckey2%3Dval2&query=hello-world"; expected != requestPathRcvd {
		t.Fatalf("Expected %s, Got: %s for request path", expected, requestPathRcvd)
	}
}
