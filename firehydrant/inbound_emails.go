package firehydrant

import (
	"context"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// InboundEmailClient is an interface for interacting with inbound emails on FireHydrant
type InboundEmailsClient interface {
	Create(ctx context.Context, createReq CreateInboundEmailRequest) (*InboundEmailResponse, error)
	Get(ctx context.Context, id string) (*InboundEmailResponse, error)
	Update(ctx context.Context, id string, updateReq UpdateInboundEmailRequest) (*InboundEmailResponse, error)
	Delete(ctx context.Context, id string) error
}

// RESTInboundEmailClient implements the InboundEmailClient interface
type RESTInboundEmailClient struct {
	client *APIClient
}

var _ InboundEmailsClient = &RESTInboundEmailClient{}

func (c *RESTInboundEmailClient) restClient() *sling.Sling {
	return c.client.client()
}

// CreateInboundEmailRequest is the payload for creating an inbound email
type CreateInboundEmailRequest struct {
	Name                 string   `json:"name"`
	Slug                 string   `json:"slug"`
	Description          string   `json:"description"`
	StatusCEL            string   `json:"status_cel"`
	LevelCEL             string   `json:"level_cel"`
	AllowedSenders       []string `json:"allowed_senders"`
	Target               Target   `json:"target"`
	Rules                []string `json:"rules"`
	RuleMatchingStrategy string   `json:"rule_matching_strategy"`
}

// Target represents the target for the inbound email
type Target struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// InboundEmailResponse is the response for an inbound email
type InboundEmailResponse struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	Slug                 string   `json:"slug"`
	Description          string   `json:"description"`
	StatusCEL            string   `json:"status_cel"`
	LevelCEL             string   `json:"level_cel"`
	AllowedSenders       []string `json:"allowed_senders"`
	Target               Target   `json:"target"`
	Rules                []string `json:"rules"`
	RuleMatchingStrategy string   `json:"rule_matching_strategy"`
}

// UpdateInboundEmailRequest is the payload for updating an inbound email
type UpdateInboundEmailRequest struct {
	Name                 string   `json:"name"`
	Slug                 string   `json:"slug"`
	Description          string   `json:"description"`
	StatusCEL            string   `json:"status_cel"`
	LevelCEL             string   `json:"level_cel"`
	AllowedSenders       []string `json:"allowed_senders"`
	Target               Target   `json:"target"`
	Rules                []string `json:"rules"`
	RuleMatchingStrategy string   `json:"rule_matching_strategy"`
}

// Create creates a new inbound email in FireHydrant
func (c *RESTInboundEmailClient) Create(ctx context.Context, createReq CreateInboundEmailRequest) (*InboundEmailResponse, error) {
	inboundEmailResponse := &InboundEmailResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Post("signals/email_targets").BodyJSON(&createReq).Receive(inboundEmailResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create inbound email")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return inboundEmailResponse, nil
}

// Get retrieves an inbound email from FireHydrant
func (c *RESTInboundEmailClient) Get(ctx context.Context, id string) (*InboundEmailResponse, error) {
	inboundEmailResponse := &InboundEmailResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Get("signals/email_targets/"+id).Receive(inboundEmailResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get inbound email")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return inboundEmailResponse, nil
}

// Update updates an inbound email in FireHydrant
func (c *RESTInboundEmailClient) Update(ctx context.Context, id string, updateReq UpdateInboundEmailRequest) (*InboundEmailResponse, error) {
	inboundEmailResponse := &InboundEmailResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Patch("signals/email_targets/"+id).BodyJSON(&updateReq).Receive(inboundEmailResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update inbound email")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return inboundEmailResponse, nil
}

// Delete deletes an inbound email from FireHydrant
func (c *RESTInboundEmailClient) Delete(ctx context.Context, id string) error {
	apiError := &APIError{}
	response, err := c.restClient().Delete("signals/email_targets/"+id).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete inbound email")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}

// Add this method to the APIClient struct in client.go
func (c *APIClient) InboundEmails() InboundEmailsClient {
	return &RESTInboundEmailClient{client: c}
}
