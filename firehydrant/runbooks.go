package firehydrant

import (
	"context"
	"fmt"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// CreateRunbookRequest is the payload for creating a service
// URL: POST https://api.firehydrant.io/v1/runbooks
type CreateRunbookRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// UpdateRunbookRequest is the payload for updating a service
// URL: PATCH https://api.firehydrant.io/v1/runbooks/{id}
type UpdateRunbookRequest struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// RunbookResponse is the payload for retrieving a service
// URL: GET https://api.firehydrant.io/v1/runbooks/{id}
type RunbookResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RunbooksClient is an interface for interacting with runbooks on FireHydrant
type RunbooksClient interface {
	Get(ctx context.Context, id string) (*RunbookResponse, error)
	Create(ctx context.Context, createReq CreateRunbookRequest) (*RunbookResponse, error)
	Update(ctx context.Context, id string, updateReq UpdateRunbookRequest) (*RunbookResponse, error)
	Delete(ctx context.Context, id string) error
}

// RESTRunbooksClient implements the RunbooksClient interface
type RESTRunbooksClient struct {
	client *APIClient
}

var _ RunbooksClient = &RESTRunbooksClient{}

func (c *RESTRunbooksClient) restClient() *sling.Sling {
	return c.client.client()
}

// Get returns a runbook from the FireHydrant API
func (c *RESTRunbooksClient) Get(ctx context.Context, id string) (*RunbookResponse, error) {
	res := &RunbookResponse{}
	resp, err := c.restClient().Get("runbooks/"+id).Receive(res, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not get runbook")
	}

	if resp.StatusCode == 404 {
		return nil, NotFound(fmt.Sprintf("Could not find runbook with ID %s", id))
	}

	return res, nil
}

// Create creates a brand spankin new runbook in FireHydrant
// TODO: Check failure case
func (c *RESTRunbooksClient) Create(ctx context.Context, createReq CreateRunbookRequest) (*RunbookResponse, error) {
	res := &RunbookResponse{}
	resp, err := c.restClient().Post("runbooks").BodyJSON(&createReq).Receive(res, nil)

	if err != nil {
		return nil, errors.Wrap(err, "could not create runbook")
	}

	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("error creating runbook: status %d", resp.StatusCode)
	}

	return res, nil
}

// Update updates a runbook in FireHydrant
func (c *RESTRunbooksClient) Update(ctx context.Context, id string, updateReq UpdateRunbookRequest) (*RunbookResponse, error) {
	res := &RunbookResponse{}
	resp, err := c.restClient().Put("runbooks/"+id).BodyJSON(updateReq).Receive(res, nil)

	if err != nil {
		return nil, errors.Wrap(err, "could not update runbook")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error creating runbook: status %d", resp.StatusCode)
	}

	return res, nil
}

func (c *RESTRunbooksClient) Delete(ctx context.Context, id string) error {
	_, err := c.restClient().Delete("runbooks/"+id).Receive(nil, nil)

	if err != nil {
		return errors.Wrap(err, "could not delete runbook")
	}

	return nil
}
