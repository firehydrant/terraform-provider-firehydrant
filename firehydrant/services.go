package firehydrant

import (
	"context"

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
func (c *RESTServicesClient) Get(ctx context.Context, id string) (*ServiceResponse, error) {
	serviceResponse := &ServiceResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Get("services/"+id).Receive(serviceResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get service")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return serviceResponse, nil
}

// List retrieves a list of services based on a service query
func (c *RESTServicesClient) List(ctx context.Context, req *ServiceQuery) (*ServicesResponse, error) {
	servicesResponse := &ServicesResponse{}
	apiError := &APIError{}
	curPage := 1

	for {
		req.Page = curPage
		var pageResponse ServicesResponse
		response, err := c.restClient().Get("services").QueryStruct(req).Receive(&pageResponse, apiError)
		if err != nil {
			return nil, errors.Wrap(err, "could not get services")
		}

		err = checkResponseStatusCode(response, apiError)
		if err != nil {
			return nil, err
		}

		servicesResponse.Services = append(servicesResponse.Services, pageResponse.Services...)

		if pageResponse.Pagination == nil || pageResponse.Pagination.Next == 0 {
			break
		}

		curPage = pageResponse.Pagination.Next
	}

	return servicesResponse, nil
}

// Create creates a brand spankin new service in FireHydrant
func (c *RESTServicesClient) Create(ctx context.Context, createReq CreateServiceRequest) (*ServiceResponse, error) {
	serviceResponse := &ServiceResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Post("services").BodyJSON(&createReq).Receive(serviceResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create service")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return serviceResponse, nil
}

// UpdateService updates a old spankin service in FireHydrant
func (c *RESTServicesClient) Update(ctx context.Context, serviceID string, updateReq UpdateServiceRequest) (*ServiceResponse, error) {
	serviceResponse := &ServiceResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Patch("services/"+serviceID).BodyJSON(&updateReq).Receive(serviceResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update service")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return serviceResponse, nil
}

// DeleteService updates a old spankin service in FireHydrant
func (c *RESTServicesClient) Delete(ctx context.Context, serviceID string) error {
	apiError := &APIError{}
	response, err := c.restClient().Delete("services/"+serviceID).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete service")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
