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
	PatchVersion = 4

	// UserAgentPrefix is the prefix of the User-Agent header that all terraform REST calls perform
	UserAgentPrefix = "firehydrant-terraform-provider"
)

type NotFound string

func (nf NotFound) Error() string {
	return string(nf)
}

// Version is the semver of this provider
var Version = fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, PatchVersion)

// APIClient is the client that accesses all of the api.firehydrant.io resources
type APIClient struct {
	baseURL         string
	token           string
	userAgentSuffix string
}

const (
	// DefaultBaseURL is the URL that is used to make requests to the FireHydrant API
	DefaultBaseURL = "https://api.firehydrant.io/v1/"
)

var _ Client = &APIClient{}

// Client is the client that makes requests to FireHydrant
type Client interface {
	Ping(ctx context.Context) (*PingResponse, error)

	Services() ServicesClient
	Runbooks() RunbooksClient
	RunbookActions() RunbookActionsClient

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

	// Teams
	GetTeam(ctx context.Context, id string) (*TeamResponse, error)
	CreateTeam(ctx context.Context, req CreateTeamRequest) (*TeamResponse, error)
	UpdateTeam(ctx context.Context, id string, req UpdateTeamRequest) (*TeamResponse, error)
	DeleteTeam(ctx context.Context, id string) error

	// Severities
	GetSeverity(ctx context.Context, slug string) (*SeverityResponse, error)
	CreateSeverity(ctx context.Context, req CreateSeverityRequest) (*SeverityResponse, error)
	UpdateSeverity(ctx context.Context, slug string, req UpdateSeverityRequest) (*SeverityResponse, error)
	DeleteSeverity(ctx context.Context, slug string) error

	// ServiceDependencies
	GetServiceDependency(ctx context.Context, id string) (*ServiceDependencyResponse, error)
	CreateServiceDependency(ctx context.Context, req CreateServiceDependencyRequest) (*ServiceDependencyResponse, error
	UpdateServiceDependency(ctx context.Context, id string, req UpdateServiceDependencyRequest) (*ServiceDependencyResponse, error)
	DeleteServiceDependency(ctx context.Context, id string) error
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

func WithUserAgentSuffix(suffix string) OptFunc {
	return func(c *APIClient) error {
		c.userAgentSuffix = suffix
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
		Set("User-Agent", fmt.Sprintf("%s (%s)/%s", UserAgentPrefix, Version, c.userAgentSuffix)).
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

// Services returns a ServicesClient interface for interacting with services in FireHydrant
func (c *APIClient) Services() ServicesClient {
	return &RESTServicesClient{client: c}
}

// Runbooks returns a RunbooksClient interface for interacting with runbooks in FireHydrant
func (c *APIClient) Runbooks() RunbooksClient {
	return &RESTRunbooksClient{client: c}
}

func (c *APIClient) RunbookActions() RunbookActionsClient {
	return &RESTRunbookActionsClient{client: c}
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

	if resp.StatusCode == 404 {
		return nil, NotFound(fmt.Sprintf("Could not find functionality with ID %s", id))
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

// GetTeam retrieves an team from the FireHydrant API
func (c *APIClient) GetTeam(ctx context.Context, id string) (*TeamResponse, error) {
	var fun TeamResponse

	resp, err := c.client().Get("teams/"+id).Receive(&fun, nil)

	if resp.StatusCode == 404 {
		return nil, NotFound(fmt.Sprintf("Could not find team with ID %s", id))
	}

	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve team")
	}

	return &fun, nil
}

// CreateTeam creates an team
func (c *APIClient) CreateTeam(ctx context.Context, req CreateTeamRequest) (*TeamResponse, error) {
	res := &TeamResponse{}

	if _, err := c.client().Post("teams").BodyJSON(&req).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not create team")
	}

	return res, nil
}

// UpdateTeam updates a team in FireHydrant
func (c *APIClient) UpdateTeam(ctx context.Context, id string, req UpdateTeamRequest) (*TeamResponse, error) {
	res := &TeamResponse{}

	if _, err := c.client().Patch("teams/"+id).BodyJSON(&req).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not update team")
	}

	return res, nil
}

// DeleteTeam deletes a team record from FireHydrant
func (c *APIClient) DeleteTeam(ctx context.Context, id string) error {
	if _, err := c.client().Delete("teams/"+id).Receive(nil, nil); err != nil {
		return errors.Wrap(err, "could not delete service")
	}

	return nil
}

// GetSeverity retrieves an severity from the FireHydrant API
func (c *APIClient) GetSeverity(ctx context.Context, slug string) (*SeverityResponse, error) {
	var fun SeverityResponse

	resp, err := c.client().Get("severities/"+slug).Receive(&fun, nil)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, NotFound(fmt.Sprintf("Could not find severity with ID %s", slug))
	}

	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve severity")
	}

	return &fun, nil
}

// CreateSeverity creates an severity
func (c *APIClient) CreateSeverity(ctx context.Context, req CreateSeverityRequest) (*SeverityResponse, error) {
	res := &SeverityResponse{}

	resp, err := c.client().Post("severities").BodyJSON(&req).Receive(res, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create severity")
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("Could not create severity %s", req.Slug)
	}

	return res, nil
}

// UpdateSeverity updates a severity in FireHydrant
func (c *APIClient) UpdateSeverity(ctx context.Context, slug string, req UpdateSeverityRequest) (*SeverityResponse, error) {
	res := &SeverityResponse{}

	if _, err := c.client().Patch("severities/"+slug).BodyJSON(&req).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not update severity")
	}

	return res, nil
}

// DeleteSeverity deletes a severity record from FireHydrant
func (c *APIClient) DeleteSeverity(ctx context.Context, slug string) error {
	if _, err := c.client().Delete("severities/"+slug).Receive(nil, nil); err != nil {
		return errors.Wrap(err, "could not delete service")
	}

	return nil
}

// GetSeverity retrieves an severity from the FireHydrant API
func (c *APIClient) GetServiceDependency(ctx context.Context, id string) (*ServiceDependencyResponse, error) {
	var fun SeverityResponse

	resp, err := c.client().Get("service_dependencies/"+id).Receive(&fun, nil)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, NotFound(fmt.Sprintf("Could not find ServiceDependency with ID %s", id))
	}

	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve ServiceDependency")
	}

	return &fun, nil
}

// CreateSeverity creates an severity
func (c *APIClient) CreateServiceDependency(ctx context.Context, req CreateServiceDependencyRequest) (*ServiceDependencyResponse, error) {
	res := &ServiceDependencyResponse{}

	resp, err := c.client().Post("service_dependencies").BodyJSON(&req).Receive(res, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create ServiceDependency")
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("Could not create ServiceDependency %s", req.ID)
	}

	return res, nil
}

// UpdateSeverity updates a severity in FireHydrant
func (c *APIClient) UpdateServiceDependency(ctx context.Context, id string, req UpdateServiceDependencyRequest) (*ServiceDependencyResponse, error) {
	res := &ServiceDependencyResponse{}

	if _, err := c.client().Patch("service_dependencies/"+id).BodyJSON(&req).Receive(res, nil); err != nil {
		return nil, errors.Wrap(err, "could not update ServiceDependency")
	}

	return res, nil
}

// DeleteSeverity deletes a severity record from FireHydrant
func (c *APIClient) DeleteServiceDependency(ctx context.Context, id string) error {
	if _, err := c.client().Delete("service_dependencies/"+id).Receive(nil, nil); err != nil {
		return errors.Wrap(err, "could not delete ServiceDependency")
	}

	return nil
}
