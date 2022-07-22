package firehydrant

import (
	"context"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// RunbookType represents the type of the runbook.
type RunbookType string

// List of valid runbook types
const (
	RunbookTypeDefault RunbookType = "incident"
)

const RunbookAttachmentRuleDefaultJSON = `
{
  "logic": {
    "manually": [
      {
        "var": "when_invoked"
      }
    ]
  },
  "user_data": {}
}`

// CreateRunbookRequest is the payload for creating a service
// URL: POST https://api.firehydrant.io/v1/runbooks
type CreateRunbookRequest struct {
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Description string       `json:"description"`
	Owner       *RunbookTeam `json:"owner,omitempty"`

	AttachmentRule map[string]interface{} `json:"attachment_rule,omitempty"`

	Severities []RunbookRelation `json:"severities"`

	Steps []RunbookStep `json:"steps,omitempty"`
}

// RunbookRelation associates a runbook to a type in FireHydrant (such as a severity)
type RunbookRelation struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Owner       *RunbookTeam  `json:"owner,omitempty"`
	Steps       []RunbookStep `json:"steps"`
	Type        string        `json:"type"`
}

// RunbookStep is a step inside of a runbook that can automate something (like creating a incident slack channel)
type RunbookStep struct {
	Name            string                 `json:"name"`
	ActionID        string                 `json:"action_id"`
	StepID          string                 `json:"step_id,omitempty"`
	Config          map[string]interface{} `json:"config,omitempty"`
	Automatic       bool                   `json:"automatic,omitempty"`
	Repeats         bool                   `json:"repeats,omitempty"`
	RepeatsDuration string                 `json:"repeats_duration,omitempty"`
	Rule            map[string]interface{} `json:"rule,omitempty"`
}

// UpdateRunbookRequest is the payload for updating a service
// URL: PATCH https://api.firehydrant.io/v1/runbooks/{id}
type UpdateRunbookRequest struct {
	Name           string                 `json:"name,omitempty"`
	Description    string                 `json:"description"`
	Owner          *RunbookTeam           `json:"owner,omitempty"`
	Steps          []RunbookStep          `json:"steps,omitempty"`
	AttachmentRule map[string]interface{} `json:"attachment_rule,omitempty"`
}

// RunbookResponse is the payload for retrieving a service
// URL: GET https://api.firehydrant.io/v1/runbooks/{id}
type RunbookResponse struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Owner       *RunbookTeam  `json:"owner"`
	Steps       []RunbookStep `json:"steps"`

	AttachmentRule map[string]interface{} `json:"attachment_rule"`

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
	runbookResponse := &RunbookResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Get("runbooks/"+id).Receive(runbookResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get runbook")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return runbookResponse, nil
}

// Create creates a brand spankin new runbook in FireHydrant
func (c *RESTRunbooksClient) Create(ctx context.Context, createReq CreateRunbookRequest) (*RunbookResponse, error) {
	// Set the default type of the runbook
	createReq.Type = string(RunbookTypeDefault)

	runbookResponse := &RunbookResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Post("runbooks").BodyJSON(&createReq).Receive(runbookResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create runbook")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return runbookResponse, nil
}

// Update updates a runbook in FireHydrant
func (c *RESTRunbooksClient) Update(ctx context.Context, id string, updateReq UpdateRunbookRequest) (*RunbookResponse, error) {
	runbookResponse := &RunbookResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Put("runbooks/"+id).BodyJSON(updateReq).Receive(runbookResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update runbook")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return runbookResponse, nil
}

func (c *RESTRunbooksClient) Delete(ctx context.Context, id string) error {
	apiError := &APIError{}
	response, err := c.restClient().Delete("runbooks/"+id).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete runbook")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
