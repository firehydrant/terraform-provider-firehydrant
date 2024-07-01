package firehydrant

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
		AssertRequestJSONBody(t, CreateServiceRequest{Name: "fake-service"}),
		AssertRequestMethod(t, "POST"),
	)

	require.NoError(t, err)
	defer teardown()

	_, err = c.Services().Create(context.TODO(), CreateServiceRequest{Name: "fake-service"})
	require.NoError(t, err, "error creating a service")
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

	if expected := "/services?labels=key1%3Dval1%2Ckey2%3Dval2&page=1&query=hello-world"; expected != requestPathRcvd {
		t.Fatalf("Expected %s, Got: %s for request path", expected, requestPathRcvd)
	}
}
