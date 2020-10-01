package firehydrant

import (
	"context"
	"fmt"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

const (
	// MajorVersion is the major version
	MajorVersion = 0
	// MinorVersion is the minor version
	MinorVersion = 1
	// PatchVersion is the patch version
	PatchVersion = 0

	// UserAgentPrefix is the prefix of the User-Agent header that all terraform REST calls perform
	UserAgentPrefix = "firehydrant-terraform-provider"
)

// Version is the semver of this provider
var Version = fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, PatchVersion)

// APIClient is the client that accesses all of the api.firehydrant.io resources
type APIClient struct {
	baseURL string
	token   string
}

const (
	// DefaultBaseURL is the URL that is used to make requests to the FireHydrant API
	DefaultBaseURL = "https://api.firehydrant.io/v1/"
)

var _ Client = &APIClient{}

// Client is the client that makes requests to FireHydrant
type Client interface {
	Ping(ctx context.Context) (*PingResponse, error)

	// Services
	GetService(ctx context.Context, id string) (*ServiceResponse, error)
	CreateService(ctx context.Context, req CreateServiceRequest) (*ServiceResponse, error)
	UpdateService(ctx context.Context, serviceID string, req UpdateServiceRequest) (*ServiceResponse, error)
	DeleteService(ctx context.Context, serviceID string) error

	// Environments
	GetEnvironment(ctx context.Context, id string) (*EnvironmentResponse, error)
	CreateEnvironment(ctx context.Context, req CreateEnvironmentRequest) (*EnvironmentResponse, error)
	UpdateEnvironment(ctx context.Context, id string, req UpdateEnvironmentRequest) (*EnvironmentResponse, error)
	DeleteEnvironment(ctx context.Context, id string) error

	// Functionalities
	GetFunctionality(ctx context.Context, id string) (*FunctionalityResponse, error)
	CreateFunctionality(ctx context.Context, req CreateFunctionalityRequest) (*FunctionalityResponse, error)
	UpdateFunctionality(ctx context.Context, id string, req UpdateFunctionalityRequest) (*FunctionalityResponse, error)
	DeleteFunctionality(ctx context.Context, id string) error
}

// OptFunc is a function that sets a setting on a client
type OptFunc func(c *APIClient) error

// WithBaseURL modifies the base URL for all requests
func WithBaseURL(baseURL string) OptFunc {
	return func(c *APIClient) error {
		c.baseURL = baseURL
		return nil
	}
}

// NewRestClient initializes a new API client for FireHydrant
func NewRestClient(token string, opts ...OptFunc) (*APIClient, error) {
	c := &APIClient{
		baseURL: DefaultBaseURL,
		token:   token,
	}

	for _, f := range opts {
		if err := f(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *APIClient) client() *sling.Sling {
	return sling.New().Base(c.baseURL).
		Set("User-Agent", fmt.Sprintf("%s (%s)", UserAgentPrefix, Version)).
		Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
}

// Ping hits and verifies the HTTP of FireHydrant
// TODO: Check failure case
func (c *APIClient) Ping(ctx context.Context) (*PingResponse, error) {
	res := &PingResponse{}

	if _, err := c.client().Get("ping").Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not ping")
	}

	return res, nil
}

// GetService retrieves a service from the FireHydrant API
// TODO: Check failure case
func (c *APIClient) GetService(ctx context.Context, id string) (*ServiceResponse, error) {
	res := &ServiceResponse{}

	if _, err := c.client().Get("services/"+id).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not get service")
	}

	return res, nil
}

// CreateService creates a brand spankin new service in FireHydrant
// TODO: Check failure case
func (c *APIClient) CreateService(ctx context.Context, createReq CreateServiceRequest) (*ServiceResponse, error) {
	res := &ServiceResponse{}

	if _, err := c.client().Post("services").BodyJSON(&createReq).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not create service")
	}

	return res, nil
}

// UpdateService updates a old spankin service in FireHydrant
// TODO: Check failure case
func (c *APIClient) UpdateService(ctx context.Context, serviceID string, updateReq UpdateServiceRequest) (*ServiceResponse, error) {
	res := &ServiceResponse{}

	if _, err := c.client().Patch("services/"+serviceID).BodyJSON(&updateReq).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not update service")
	}

	return res, nil
}

// DeleteService updates a old spankin service in FireHydrant
// TODO: Check failure case
func (c *APIClient) DeleteService(ctx context.Context, serviceID string) error {
	if _, err := c.client().Delete("services/"+serviceID).Receive(nil, nil); err != nil {
		return errors.Wrap(err, "could not delete service")
	}

	return nil
}

// GetEnvironment retrieves an environment from the FireHydrant API
func (c *APIClient) GetEnvironment(ctx context.Context, id string) (*EnvironmentResponse, error) {
	var env EnvironmentResponse

	resp, err := c.client().Get("environments/"+id).Receive(&env, nil)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code was not a 200, got %d", resp.StatusCode)
	}

	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve environment")
	}

	return &env, nil
}

// CreateEnvironment creates an environment
func (c *APIClient) CreateEnvironment(ctx context.Context, req CreateEnvironmentRequest) (*EnvironmentResponse, error) {
	res := &EnvironmentResponse{}

	if _, err := c.client().Post("environments").BodyJSON(&req).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not create environment")
	}

	return res, nil
}

// UpdateEnvironment updates a environment in FireHydrant
func (c *APIClient) UpdateEnvironment(ctx context.Context, id string, req UpdateEnvironmentRequest) (*EnvironmentResponse, error) {
	res := &EnvironmentResponse{}

	if _, err := c.client().Patch("environments/"+id).BodyJSON(&req).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not update environment")
	}

	return res, nil
}

// DeleteEnvironment deletes a environment record from FireHydrant
func (c *APIClient) DeleteEnvironment(ctx context.Context, id string) error {
	if _, err := c.client().Delete("environments/"+id).Receive(nil, nil); err != nil {
		return errors.Wrap(err, "could not delete service")
	}

	return nil
}

// GetFunctionality retrieves an functionality from the FireHydrant API
func (c *APIClient) GetFunctionality(ctx context.Context, id string) (*FunctionalityResponse, error) {
	var fun FunctionalityResponse

	resp, err := c.client().Get("functionalities/"+id).Receive(&fun, nil)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code was not a 200, got %d", resp.StatusCode)
	}

	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve functionality")
	}

	return &fun, nil
}

// CreateFunctionality creates an functionality
func (c *APIClient) CreateFunctionality(ctx context.Context, req CreateFunctionalityRequest) (*FunctionalityResponse, error) {
	res := &FunctionalityResponse{}

	if _, err := c.client().Post("functionalities").BodyJSON(&req).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not create functionality")
	}

	return res, nil
}

// UpdateFunctionality updates a functionality in FireHydrant
func (c *APIClient) UpdateFunctionality(ctx context.Context, id string, req UpdateFunctionalityRequest) (*FunctionalityResponse, error) {
	res := &FunctionalityResponse{}

	if _, err := c.client().Patch("functionalities/"+id).BodyJSON(&req).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not update functionality")
	}

	return res, nil
}

// DeleteFunctionality deletes a functionality record from FireHydrant
func (c *APIClient) DeleteFunctionality(ctx context.Context, id string) error {
	if _, err := c.client().Delete("functionalities/"+id).Receive(nil, nil); err != nil {
		return errors.Wrap(err, "could not delete service")
	}

	return nil
}
