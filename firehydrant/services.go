package firehydrant

import (
	"context"
	"fmt"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// ServicesClient is an interface for interacting with services in FireHydrant
type ServicesClient interface {
	Get(ctx context.Context, id string) (*ServiceResponse, error)
	List(ctx context.Context, req *ServiceQuery) (*ServicesResponse, error)
	Create(ctx context.Context, req CreateServiceRequest) (*ServiceResponse, error)
	Update(ctx context.Context, serviceID string, req UpdateServiceRequest) (*ServiceResponse, error)
	Delete(ctx context.Context, serviceID string) error
}

// RESTServicesClient implements the ServicesClient interface
type RESTServicesClient struct {
	client *APIClient
}

var _ ServicesClient = &RESTServicesClient{}

func (c *RESTServicesClient) restClient() *sling.Sling {
	return c.client.client()
}

// Get retrieves a service from the FireHydrant API
// TODO: Check failure case
func (c *RESTServicesClient) Get(ctx context.Context, id string) (*ServiceResponse, error) {
	res := &ServiceResponse{}
	resp, err := c.restClient().Get("services/"+id).Receive(res, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not get service")
	}

	if resp.StatusCode == 404 {
		return nil, NotFound(fmt.Sprintf("Could not find service with ID %s", id))
	}

	return res, nil
}

// List retrieves a list of services based on a service query
func (c *RESTServicesClient) List(ctx context.Context, req *ServiceQuery) (*ServicesResponse, error) {
	res := &ServicesResponse{}
	_, err := c.restClient().Get("services").QueryStruct(req).Receive(res, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not get service")
	}

	return res, nil
}

// Create creates a brand spankin new service in FireHydrant
// TODO: Check failure case
func (c *RESTServicesClient) Create(ctx context.Context, createReq CreateServiceRequest) (*ServiceResponse, error) {
	res := &ServiceResponse{}

	if _, err := c.restClient().Post("services").BodyJSON(&createReq).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not create service")
	}

	return res, nil
}

// UpdateService updates a old spankin service in FireHydrant
// TODO: Check failure case
func (c *RESTServicesClient) Update(ctx context.Context, serviceID string, updateReq UpdateServiceRequest) (*ServiceResponse, error) {
	res := &ServiceResponse{}

	if _, err := c.restClient().Patch("services/"+serviceID).BodyJSON(&updateReq).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not update service")
	}

	return res, nil
}

// DeleteService updates a old spankin service in FireHydrant
// TODO: Check failure case
func (c *RESTServicesClient) Delete(ctx context.Context, serviceID string) error {
	if _, err := c.restClient().Delete("services/"+serviceID).Receive(nil, nil); err != nil {
		return errors.Wrap(err, "could not delete service")
	}

	return nil
}
