package firehydrant

import (
	"context"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// EnvironmentsClient is an interface for interacting with environments on FireHydrant
type EnvironmentsClient interface {
	Get(ctx context.Context, id string) (*EnvironmentResponse, error)
	Create(ctx context.Context, createReq CreateEnvironmentRequest) (*EnvironmentResponse, error)
	Update(ctx context.Context, id string, updateReq UpdateEnvironmentRequest) (*EnvironmentResponse, error)
	Delete(ctx context.Context, id string) error
}

// RESTEnvironmentsClient implements the EnvironmentsClient interface
type RESTEnvironmentsClient struct {
	client *APIClient
}

var _ EnvironmentsClient = &RESTEnvironmentsClient{}

func (c *RESTEnvironmentsClient) restClient() *sling.Sling {
	return c.client.client()
}

// EnvironmentResponse is the payload for a single environment
// URL: GET https://api.firehydrant.io/v1/environments/{id}
type EnvironmentResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Slug        string    `json:"slug"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Get retrieves an environment from FireHydrant
func (c *RESTEnvironmentsClient) Get(ctx context.Context, id string) (*EnvironmentResponse, error) {
	envResponse := &EnvironmentResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Get("environments/"+id).Receive(envResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get environment")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return envResponse, nil
}

// CreateEnvironmentRequest is the payload for creating a service
// URL: POST https://api.firehydrant.io/v1/services
type CreateEnvironmentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Create creates an environment in FireHydrant
func (c *RESTEnvironmentsClient) Create(ctx context.Context, req CreateEnvironmentRequest) (*EnvironmentResponse, error) {
	envResponse := &EnvironmentResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Post("environments").BodyJSON(&req).Receive(envResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create environment")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return envResponse, nil
}

// UpdateEnvironmentRequest is the payload for updating a environment
// URL: PATCH https://api.firehydrant.io/v1/environments/{id}
type UpdateEnvironmentRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description"`
}

// Update updates a environment in FireHydrant
func (c *RESTEnvironmentsClient) Update(ctx context.Context, id string, req UpdateEnvironmentRequest) (*EnvironmentResponse, error) {
	envResponse := &EnvironmentResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Patch("environments/"+id).BodyJSON(&req).Receive(envResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update environment")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return envResponse, nil
}

// Delete deletes a environment from FireHydrant
func (c *RESTEnvironmentsClient) Delete(ctx context.Context, id string) error {
	apiError := &APIError{}
	response, err := c.restClient().Delete("environments/"+id).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete environment")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
