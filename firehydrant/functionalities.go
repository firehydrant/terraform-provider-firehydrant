package firehydrant

import (
	"context"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// FunctionalitiesClient is an interface for interacting with task lists on FireHydrant
type FunctionalitiesClient interface {
	Get(ctx context.Context, id string) (*FunctionalityResponse, error)
	Create(ctx context.Context, createReq CreateFunctionalityRequest) (*FunctionalityResponse, error)
	Update(ctx context.Context, id string, updateReq UpdateFunctionalityRequest) (*FunctionalityResponse, error)
	Delete(ctx context.Context, id string) error
}

// RESTFunctionalitiesClient implements the FunctionalitiesClient interface
type RESTFunctionalitiesClient struct {
	client *APIClient
}

var _ FunctionalitiesClient = &RESTFunctionalitiesClient{}

func (c *RESTFunctionalitiesClient) restClient() *sling.Sling {
	return c.client.client()
}

// FunctionalityResponse is the payload for a single environment
// URL: GET https://api.firehydrant.io/v1/functionalities/{id}
type FunctionalityResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Slug        string                 `json:"slug"`
	Services    []FunctionalityService `json:"services"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// FunctionalityService represents a service when creating a functionality
type FunctionalityService struct {
	ID string `json:"id"`
}

// Get retrieves a functionality from FireHydrant
func (c *RESTFunctionalitiesClient) Get(ctx context.Context, id string) (*FunctionalityResponse, error) {
	funcResponse := &FunctionalityResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Get("functionalities/"+id).Receive(funcResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get functionality")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return funcResponse, nil
}

// CreateFunctionalityRequest is the payload for creating a service
// URL: POST https://api.firehydrant.io/v1/services
type CreateFunctionalityRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Services    []FunctionalityService `json:"services,omitempty"`
}

// Create creates a functionality in FireHydrant
func (c *RESTFunctionalitiesClient) Create(ctx context.Context, req CreateFunctionalityRequest) (*FunctionalityResponse, error) {
	funcResponse := &FunctionalityResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Post("functionalities").BodyJSON(&req).Receive(funcResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create functionality")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return funcResponse, nil
}

// UpdateFunctionalityRequest is the payload for updating a environment
// URL: PATCH https://api.firehydrant.io/v1/environments/{id}
type UpdateFunctionalityRequest struct {
	Name                    string                 `json:"name,omitempty"`
	Description             string                 `json:"description"`
	RemoveRemainingServices bool                   `json:"remove_remaining_services"`
	Services                []FunctionalityService `json:"services"`
}

// Update updates a functionality in FireHydrant
func (c *RESTFunctionalitiesClient) Update(ctx context.Context, id string, req UpdateFunctionalityRequest) (*FunctionalityResponse, error) {
	funcResponse := &FunctionalityResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Patch("functionalities/"+id).BodyJSON(&req).Receive(funcResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update functionality")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return funcResponse, nil
}

// Delete deletes a functionality from FireHydrant
func (c *RESTFunctionalitiesClient) Delete(ctx context.Context, id string) error {
	apiError := &APIError{}
	response, err := c.restClient().Delete("functionalities/"+id).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete functionality")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
