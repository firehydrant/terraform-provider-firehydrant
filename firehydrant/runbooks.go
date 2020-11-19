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

	Severities []RunbookRelation `json:"severities"`

	Steps []RunbookStep `json:"steps,omitempty"`
}

// RunbookRelation associates a runbook to a type in FireHydrant (such as a severity)
type RunbookRelation struct {
	ID string `json:"id"`
}

// RunbookStep is a step inside of a runbook that can automate something (like creating a incident slack channel)
type RunbookStep struct {
	Name            string            `json:"name"`
	ActionID        string            `json:"action_id"`
	StepID          string            `json:"step_id,omitempty"`
	Config          map[string]string `json:"config,omitempty"`
	Automatic       bool              `json:"automatic,omitempty"`
	Repeats         bool              `json:"repeats,omitempty"`
	RepeatsDuration string            `json:"repeats_duration,omitempty"`
	DelayDuration   string            `json:"delay_duration,omitempty"`
}

// UpdateRunbookRequest is the payload for updating a service
// URL: PATCH https://api.firehydrant.io/v1/runbooks/{id}
type UpdateRunbookRequest struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Steps       []RunbookStep     `json:"steps,omitempty"`
	Severities  []RunbookRelation `json:"severities"`
}

// RunbookResponse is the payload for retrieving a service
// URL: GET https://api.firehydrant.io/v1/runbooks/{id}
type RunbookResponse struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Description string        `json:"description"`
	Steps       []RunbookStep `json:"steps"`

	Severities []RunbookRelation `json:"severities"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

	res, err = c.Update(ctx, res.ID, UpdateRunbookRequest{
		Steps:      createReq.Steps,
		Severities: createReq.Severities,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not update created runbook")
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
