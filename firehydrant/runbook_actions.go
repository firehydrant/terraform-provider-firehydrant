package firehydrant

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

type RunbookActionsResponse struct {
	Actions []RunbookAction `json:"data"`
}

// RunbookResponse is the payload for retrieving a service
// URL: GET https://api.firehydrant.io/v1/runbooks/{id}
type RunbookAction struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RunbookActionsQuery struct {
	Type string `url:"type,omitempty"`
	Depaginate string `url:"per_page,omitempty"`
}

// RunbooksClient is an interface for interacting with runbooks on FireHydrant
type RunbookActionsClient interface {
	Get(ctx context.Context, typ, integrationAndSlug string) (*RunbookAction, error)
}

// RESTRunbooksClient implements the RunbooksClient interface
type RESTRunbookActionsClient struct {
	client *APIClient
}

var _ RunbookActionsClient = &RESTRunbookActionsClient{}

func (c *RESTRunbookActionsClient) restClient() *sling.Sling {
	return c.client.client()
}

// Get returns a runbook action from the FireHydrant API
func (c *RESTRunbookActionsClient) Get(ctx context.Context, typ, integrationAndSlug string) (*RunbookAction, error) {
	res := &RunbookActionsResponse{}
	query := RunbookActionsQuery{Type: typ, Depaginate: "40"}
	_, err := c.restClient().Get("runbooks/actions").QueryStruct(query).Receive(res, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not get runbook")
	}

	split := strings.Split(integrationAndSlug, ".")
	_, slug := split[0], split[1]

	for _, v := range res.Actions {
		if v.Slug == slug {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("could not find runbook action")
}
