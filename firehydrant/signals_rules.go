package firehydrant

import (
	"context"
	"fmt"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

type SignalsRules interface {
	Get(ctx context.Context, teamID, id string) (*SignalsRuleResponse, error)
	Create(ctx context.Context, teamID string, createReq CreateSignalsRuleRequest) (*SignalsRuleResponse, error)
	Update(ctx context.Context, teamID, id string, updateReq UpdateSignalsRuleRequest) (*SignalsRuleResponse, error)
	Delete(ctx context.Context, teamID, id string) error
}

type RESTSignalsRulesClient struct {
	client *APIClient
}

var _ SignalsRules = &RESTSignalsRulesClient{}

type SignalRuleTarget struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type SignalRuleIncidentType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SignalsRuleResponse struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Expression   string                 `json:"expression"`
	Target       SignalRuleTarget       `json:"target"`
	IncidentType SignalRuleIncidentType `json:"incident_type"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateSignalsRuleRequest struct {
	Name           string `json:"name"`
	Expression     string `json:"expression"`
	TargetType     string `json:"target_type"`
	TargetID       string `json:"target_id"`
	IncidentTypeID string `json:"incident_type_id,omitempty"`
}

type UpdateSignalsRuleRequest struct {
	Name           string `json:"name"`
	Expression     string `json:"expression"`
	TargetType     string `json:"target_type"`
	TargetID       string `json:"target_id"`
	IncidentTypeID string `json:"incident_type_id,omitempty"`
}

func (c *RESTSignalsRulesClient) restClient() *sling.Sling {
	return c.client.client()
}

func (c *RESTSignalsRulesClient) Create(ctx context.Context, teamID string, createReq CreateSignalsRuleRequest) (*SignalsRuleResponse, error) {
	signalRuleResponse := &SignalsRuleResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Post(fmt.Sprintf("teams/%s/signal_rules", teamID)).BodyJSON(&createReq).Receive(signalRuleResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create signal rule")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return signalRuleResponse, nil
}

func (c *RESTSignalsRulesClient) Get(ctx context.Context, teamID, id string) (*SignalsRuleResponse, error) {
	signalRuleResponse := &SignalsRuleResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Get(fmt.Sprintf("teams/%s/signal_rules/%s", teamID, id)).Receive(signalRuleResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get signal rule")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return signalRuleResponse, nil
}

func (c *RESTSignalsRulesClient) Update(ctx context.Context, teamID, id string, updateReq UpdateSignalsRuleRequest) (*SignalsRuleResponse, error) {
	signalRuleResponse := &SignalsRuleResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Patch(fmt.Sprintf("teams/%s/signal_rules/%s", teamID, id)).BodyJSON(&updateReq).Receive(signalRuleResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update signal rule")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return signalRuleResponse, nil
}

func (c *RESTSignalsRulesClient) Delete(ctx context.Context, teamID, id string) error {
	apiError := &APIError{}
	response, err := c.restClient().Delete(fmt.Sprintf("teams/%s/signal_rules/%s", teamID, id)).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete signal rule")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
