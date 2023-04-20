package firehydrant

import (
	"context"
	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// TeamsClient is an interface for interacting with teams on FireHydrant
type TeamsClient interface {
	Get(ctx context.Context, slug string) (*TeamResponse, error)
	List(ctx context.Context, req *TeamQuery) (*TeamsResponse, error)
	Create(ctx context.Context, createReq CreateTeamRequest) (*TeamResponse, error)
	Update(ctx context.Context, slug string, updateReq UpdateTeamRequest) (*TeamResponse, error)
	Archive(ctx context.Context, slug string) error
}

// RESTTeamsClient implements the TeamsClient interface
type RESTTeamsClient struct {
	client *APIClient
}

var _ TeamsClient = &RESTTeamsClient{}

func (c *RESTTeamsClient) restClient() *sling.Sling {
	return c.client.client()
}

// Get retrieves a team from FireHydrant
func (c *RESTTeamsClient) Get(ctx context.Context, id string) (*TeamResponse, error) {
	teamResponse := &TeamResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Get("teams/"+id).Receive(teamResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get team")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return teamResponse, nil
}

// List retrieves a list of teams based on a team query
func (c *RESTTeamsClient) List(ctx context.Context, req *TeamQuery) (*TeamsResponse, error) {
	teamsResponse := &TeamsResponse{}
	apiError := &APIError{}
	curPage := 1

	for {
		req.Page = curPage
		response, err := c.restClient().Get("teams").QueryStruct(req).Receive(teamsResponse, apiError)
		if err != nil {
			return nil, errors.Wrap(err, "could not get teams")
		}

		err = checkResponseStatusCode(response, apiError)
		if err != nil {
			return nil, err
		}

		if teamsResponse.Pagination == nil || teamsResponse.Pagination.Last == curPage {
			break
		}
		curPage += 1
	}

	return teamsResponse, nil
}

// Create creates a team
func (c *RESTTeamsClient) Create(ctx context.Context, createReq CreateTeamRequest) (*TeamResponse, error) {
	teamResponse := &TeamResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Post("teams").BodyJSON(&createReq).Receive(teamResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create team")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return teamResponse, nil
}

// Update updates a team in FireHydrant
func (c *RESTTeamsClient) Update(ctx context.Context, slug string, updateReq UpdateTeamRequest) (*TeamResponse, error) {
	teamResponse := &TeamResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Patch("teams/"+slug).BodyJSON(&updateReq).Receive(teamResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update team")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return teamResponse, nil
}

// Archive archives a team in FireHydrant
func (c *RESTTeamsClient) Archive(ctx context.Context, slug string) error {
	apiError := &APIError{}
	response, err := c.restClient().Delete("teams/"+slug).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not archive team")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
