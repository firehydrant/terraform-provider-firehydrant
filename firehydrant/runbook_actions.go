package firehydrant

import (
	"context"
	"fmt"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// RunbooksClient is an interface for interacting with runbooks on FireHydrant
type RunbookActionsClient interface {
	Get(ctx context.Context, runbookType string, integrationSlug string, actionSlug string) (*RunbookAction, error)
}

// RESTRunbooksClient implements the RunbooksClient interface
type RESTRunbookActionsClient struct {
	client *APIClient
}

var _ RunbookActionsClient = &RESTRunbookActionsClient{}

func (c *RESTRunbookActionsClient) restClient() *sling.Sling {
	return c.client.client()
}

// RunbookActionsResponse is the payload for retrieving runbook actions
// URL: GET https://api.firehydrant.io/v1/runbooks/actions
type RunbookActionsResponse struct {
	Actions []RunbookAction `json:"data"`
}

type RunbookAction struct {
	ID          string       `json:"id"`
	Integration *Integration `json:"integration"`
	Name        string       `json:"name"`
	Slug        string       `json:"slug"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Integration struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
}

type RunbookActionsQuery struct {
	Type  string `url:"type,omitempty"`
	Items uint   `url:"per_page,omitempty"`
}

// Get returns a runbook action from the FireHydrant API
func (c *RESTRunbookActionsClient) Get(ctx context.Context, runbookType string, integrationSlug string, actionSlug string) (*RunbookAction, error) {
	runbookActionResponse := &RunbookActionsResponse{}
	apiError := &APIError{}
	query := RunbookActionsQuery{Type: runbookType, Items: 100}
	response, err := c.restClient().Get("runbooks/actions").QueryStruct(query).Receive(runbookActionResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get runbook")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	for _, action := range runbookActionResponse.Actions {
		if action.Slug == actionSlug && action.Integration.Slug == integrationSlug {
			return &action, nil
		}
	}

	return nil, fmt.Errorf("could not find runbook action")
}
