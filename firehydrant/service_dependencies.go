package firehydrant

import (
	"context"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// ServiceDependenciesClient is an interface for interacting with service dependencies on FireHydrant
type ServiceDependenciesClient interface {
	Get(ctx context.Context, id string) (*ServiceDependencyResponse, error)
	Create(ctx context.Context, createReq CreateServiceDependencyRequest) (*ServiceDependencyResponse, error)
	Update(ctx context.Context, id string, updateReq UpdateServiceDependencyRequest) (*ServiceDependencyResponse, error)
	Delete(ctx context.Context, id string) error
}

// RESTServiceDependenciesClient implements the ServiceDependenciesClient interface
type RESTServiceDependenciesClient struct {
	client *APIClient
}

var _ ServiceDependenciesClient = &RESTServiceDependenciesClient{}

func (c *RESTServiceDependenciesClient) restClient() *sling.Sling {
	return c.client.client()
}

// ServiceDependencyResponse is the payload for retrieving a service dependency
// URL: GET https://api.firehydrant.io/v1/service_dependencies/{id}
type ServiceDependencyResponse struct {
	ID    string `json:"id"`
	Notes string `json:"notes"`

	ConnectedService *ServiceDependencyService `json:"connected_service"`
	Service          *ServiceDependencyService `json:"service"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ServiceDependencyService is a service in a service dependency
type ServiceDependencyService struct {
	ID string `json:"id"`
}

// Get returns a service dependency from the FireHydrant API
func (c *RESTServiceDependenciesClient) Get(ctx context.Context, id string) (*ServiceDependencyResponse, error) {
	serviceDependencyResponse := &ServiceDependencyResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Get("service_dependencies/"+id).Receive(serviceDependencyResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get service dependency")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return serviceDependencyResponse, nil
}

// CreateServiceDependencyRequest is the payload for creating a service dependency
// URL: POST https://api.firehydrant.io/v1/service_dependencies
type CreateServiceDependencyRequest struct {
	ConnectedServiceID string `json:"connected_service_id"`
	ServiceID          string `json:"service_id"`
	Notes              string `json:"notes,omitempty"`
}

// Create creates a brand spankin new service dependency in FireHydrant
func (c *RESTServiceDependenciesClient) Create(ctx context.Context, createReq CreateServiceDependencyRequest) (*ServiceDependencyResponse, error) {
	serviceDependencyResponse := &ServiceDependencyResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Post("service_dependencies").BodyJSON(&createReq).Receive(serviceDependencyResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create service dependency")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return serviceDependencyResponse, nil
}

// UpdateServiceDependencyRequest is the payload for updating a service dependency
// URL: PATCH https://api.firehydrant.io/v1/service_dependencies/{id}
type UpdateServiceDependencyRequest struct {
	Notes string `json:"notes"`
}

// Update updates a service dependency in FireHydrant
func (c *RESTServiceDependenciesClient) Update(ctx context.Context, id string, updateReq UpdateServiceDependencyRequest) (*ServiceDependencyResponse, error) {
	serviceDependencyResponse := &ServiceDependencyResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Patch("service_dependencies/"+id).BodyJSON(updateReq).Receive(serviceDependencyResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update service dependency")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return serviceDependencyResponse, nil
}

func (c *RESTServiceDependenciesClient) Delete(ctx context.Context, id string) error {
	apiError := &APIError{}
	response, err := c.restClient().Delete("service_dependencies/"+id).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete service dependency")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
