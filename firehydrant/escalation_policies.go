package firehydrant

import (
	"context"
	"fmt"
	"time"

	"github.com/dghubble/sling"
)

type EscalationPolicies interface {
	Get(ctx context.Context, teamID, id string) (*EscalationPolicyResponse, error)
	Create(ctx context.Context, teamID string, createReq CreateEscalationPolicyRequest) (*EscalationPolicyResponse, error)
	Update(ctx context.Context, teamID, id string, updateReq UpdateEscalationPolicyRequest) (*EscalationPolicyResponse, error)
	Delete(ctx context.Context, teamID, id string) error
}

type RESTEscalationPoliciesClient struct {
	client *APIClient
}

var _ EscalationPolicies = &RESTEscalationPoliciesClient{}

type EscalationPolicyResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
	Repetitions int    `json:"repetitions"`

	Steps []EscalationPolicyStepWithTarget `json:"steps"`

	HandoffStep *EscalationPolicyHandoffStep `json:"handoff_step"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EscalationPolicyStepTarget struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name,omitempty"`
}

type EscalationPolicyHandoffStep struct {
	Target EscalationPolicyTarget `json:"target"`
}

type CreateEscalationPolicyHandoffStep struct {
	Type string `json:"target_type"`
	ID   string `json:"target_id"`
}

type EscalationPolicyStepWithTarget struct {
	ID       string                   `json:"id"`
	Position int                      `json:"position"`
	Timeout  string                   `json:"timeout"`
	Targets  []EscalationPolicyTarget `json:"targets"`
}

type EscalationPolicyTarget struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type EscalationPolicyStep struct {
	ID       string                   `json:"id"`
	Position int                      `json:"position"`
	Timeout  string                   `json:"timeout"`
	Targets  []EscalationPolicyTarget `json:"targets"`
}

type CreateEscalationPolicyRequest struct {
	Name        string                             `json:"name"`
	Description string                             `json:"description"`
	Default     bool                               `json:"default"`
	Repetitions int                                `json:"repetitions"`
	Steps       []EscalationPolicyStep             `json:"steps"`
	HandoffStep *CreateEscalationPolicyHandoffStep `json:"handoff_step,omitempty"`
}

type UpdateEscalationPolicyRequest struct {
	Name        string                             `json:"name"`
	Description string                             `json:"description"`
	Default     bool                               `json:"default"`
	Repetitions int                                `json:"repetitions"`
	Steps       []EscalationPolicyStep             `json:"steps"`
	HandoffStep *CreateEscalationPolicyHandoffStep `json:"handoff_step"`
}

func (c *RESTEscalationPoliciesClient) restClient() *sling.Sling {
	return c.client.client()
}

func (c *RESTEscalationPoliciesClient) Create(ctx context.Context, teamID string, createReq CreateEscalationPolicyRequest) (*EscalationPolicyResponse, error) {
	escalationPolicyResponse := &EscalationPolicyResponse{}
	apiError := &APIError{}

	response, err := c.restClient().Post(fmt.Sprintf("teams/%s/escalation_policies", teamID)).BodyJSON(createReq).Receive(escalationPolicyResponse, apiError)
	if err != nil {
		return nil, err
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return escalationPolicyResponse, nil
}

func (c *RESTEscalationPoliciesClient) Get(ctx context.Context, teamID, id string) (*EscalationPolicyResponse, error) {
	escalationPolicyResponse := &EscalationPolicyResponse{}
	apiError := &APIError{}

	response, err := c.restClient().Get(fmt.Sprintf("teams/%s/escalation_policies/%s", teamID, id)).Receive(escalationPolicyResponse, apiError)
	if err != nil {
		return nil, err
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return escalationPolicyResponse, nil
}

func (c *RESTEscalationPoliciesClient) Update(ctx context.Context, teamID, id string, updateReq UpdateEscalationPolicyRequest) (*EscalationPolicyResponse, error) {
	escalationPolicyResponse := &EscalationPolicyResponse{}
	apiError := &APIError{}

	response, err := c.restClient().Patch(fmt.Sprintf("teams/%s/escalation_policies/%s", teamID, id)).BodyJSON(updateReq).Receive(escalationPolicyResponse, apiError)
	if err != nil {
		return nil, err
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return escalationPolicyResponse, nil
}

func (c *RESTEscalationPoliciesClient) Delete(ctx context.Context, teamID, id string) error {
	apiError := &APIError{}

	_, err := c.restClient().Delete(fmt.Sprintf("teams/%s/escalation_policies/%s", teamID, id)).Receive(nil, apiError)
	if err != nil {
		return err
	}

	return nil
}
