package firehydrant

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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
			Next:  0,
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

	qry := &TeamQuery{Query: "test-team"}

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

func TestListTeamsPaginated(t *testing.T) {
	responses := []TeamsResponse{
		TeamsResponse{
			Teams: []TeamResponse{
				{
					ID: "team-1",
				},
			},
			Pagination: &Pagination{
				Count: 2,
				Page:  1,
				Items: 1,
				Pages: 2,
				Last:  2,
				Next:  2,
			},
		},
		TeamsResponse{
			Teams: []TeamResponse{
				{
					ID: "team-2",
				},
			},
			Pagination: &Pagination{
				Count: 2,
				Page:  2,
				Items: 1,
				Pages: 2,
				Last:  2,
				Next:  0, // Technically null in JSON, but marshalled to zero-value of int in Go.
			},
		},
	}

	requestCount := 0

	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if requestCount >= len(responses) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		q := req.URL.Query()
		pageStr := q.Get("page")
		if pageStr == "" {
			pageStr = "1"
		}
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		requestCount++

		// Decrement page to get the correct response in slice.
		// URL pages are 1-indexed, but slices are 0-indexed.
		response := responses[page-1]

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

	qry := &TeamQuery{}

	result, err := c.Teams().List(context.TODO(), qry)
	if err != nil {
		t.Fatalf("Received error hitting ping endpoint: %s", err.Error())
	}
	if requestCount != 2 {
		t.Errorf("Expected 2 requests, got %d", requestCount)
	}
	if total := len(result.Teams); total != 2 {
		t.Errorf("Expected 2 results, got %d", total)
	}
	if result.Teams[0].ID != "team-1" {
		t.Errorf("Expected team-1, got %s", result.Teams[0].ID)
	}
	if result.Teams[1].ID != "team-2" {
		t.Errorf("Expected team-2, got %s", result.Teams[1].ID)
	}
}
