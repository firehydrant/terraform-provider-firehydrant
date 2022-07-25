package firehydrant

import (
	"context"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// PrioritiesClient is an interface for interacting with priorities on FireHydrant
type PrioritiesClient interface {
	Get(ctx context.Context, slug string) (*PriorityResponse, error)
	Create(ctx context.Context, createReq CreatePriorityRequest) (*PriorityResponse, error)
	Update(ctx context.Context, slug string, updateReq UpdatePriorityRequest) (*PriorityResponse, error)
	Delete(ctx context.Context, slug string) error
}

// RESTPrioritiesClient implements the PrioritiesClient interface
type RESTPrioritiesClient struct {
	client *APIClient
}

var _ PrioritiesClient = &RESTPrioritiesClient{}

func (c *RESTPrioritiesClient) restClient() *sling.Sling {
	return c.client.client()
}

// PriorityResponse is the payload for a single priority
// URL: GET https://api.firehydrant.io/v1/priorities/{id}
type PriorityResponse struct {
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}

// Get retrieves a priority from FireHydrant
func (c *RESTPrioritiesClient) Get(ctx context.Context, slug string) (*PriorityResponse, error) {
	priorityResponse := &PriorityResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Get("priorities/"+slug).Receive(priorityResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get priority")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return priorityResponse, nil
}

// CreatePriorityRequest is the payload for creating a priority
// URL: POST https://api.firehydrant.io/v1/priorities
type CreatePriorityRequest struct {
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}

// Create creates a priority
func (c *RESTPrioritiesClient) Create(ctx context.Context, createReq CreatePriorityRequest) (*PriorityResponse, error) {
	priorityResponse := &PriorityResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Post("priorities").BodyJSON(&createReq).Receive(priorityResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create priority")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return priorityResponse, nil
}

// UpdatePriorityRequest is the payload for updating a priority
// URL: PATCH https://api.firehydrant.io/v1/priorities/{id}
type UpdatePriorityRequest struct {
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}

// Update updates a priority in FireHydrant
func (c *RESTPrioritiesClient) Update(ctx context.Context, slug string, updateReq UpdatePriorityRequest) (*PriorityResponse, error) {
	priorityResponse := &PriorityResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Patch("priorities/"+slug).BodyJSON(&updateReq).Receive(priorityResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update priority")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return priorityResponse, nil
}

// Delete deletes a priority from FireHydrant
func (c *RESTPrioritiesClient) Delete(ctx context.Context, slug string) error {
	apiError := &APIError{}
	response, err := c.restClient().Delete("priorities/"+slug).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete priority")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
