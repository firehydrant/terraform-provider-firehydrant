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

func TestGetTeam(t *testing.T) {
	resp := &TeamResponse{}
	testTeamID := "test-team-id"
	c, teardown, err := setupClient("/teams/"+testTeamID, resp)
	require.NoError(t, err)
	defer teardown()

	res, err := c.Teams().Get(context.TODO(), testTeamID)
	require.NoError(t, err, "error retrieving a team")
	assert.Equal(t, resp.ID, res.ID, "returned team did not match")
	assert.Equal(t, resp.Name, res.Name, "returned team did not match")
}

func TestCreateTeam(t *testing.T) {
	resp := &TeamResponse{}
	c, teardown, err := setupClient("/teams", resp,
		AssertRequestJSONBody(t, CreateTeamRequest{Name: "fake-team", Description: "fake description"}),
		AssertRequestMethod(t, "POST"),
	)

	require.NoError(t, err)
	defer teardown()

	_, err = c.Teams().Create(context.TODO(), CreateTeamRequest{Name: "fake-team", Description: "fake description"})
	require.NoError(t, err, "error creating a team")
}

func TestGetTeams(t *testing.T) {
	var requestPathRcvd string
	response := TeamsResponse{
		Teams: []TeamResponse{
			{
				ID: "test-team",
			},
		},
		Pagination: &Pagination{
			Count: 1,
			Page:  1,
			Items: 1,
			Pages: 1,
			Last:  1,
			Next:  1,
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

	qry := &TeamQuery{
		Query: "test-team",
	}

	vs, err := query.Values(qry)
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log(vs)

	_, err = c.Teams().List(context.TODO(), qry)
	if err != nil {
		t.Fatalf("Received error hitting ping endpoint: %s", err.Error())
	}

	if expected := "/teams?page=1&query=test-team"; expected != requestPathRcvd {
		t.Fatalf("Expected %s, Got: %s for request path", expected, requestPathRcvd)
	}
}
