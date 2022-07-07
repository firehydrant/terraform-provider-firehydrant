package firehydrant

import (
	"context"
	"fmt"
	"net/http"

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

// checkResponseStatusCode checks to see if the response's status
// code corresponds to an error or not. An error is returned for
// all status codes 300 and above
func checkResponseStatusCode(response *http.Response, apiError *APIError) error {
	switch code := response.StatusCode; {
	case code >= 200 && code <= 299:
		return nil
	case code == 404:
		return ErrorNotFound
	case code == 401:
		return fmt.Errorf("%s\n%s", ErrorUnauthorized, apiError)
	default:
		return fmt.Errorf("%d request failed with error\n%s", code, apiError)
	}
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

	// Priorities
	GetPriority(ctx context.Context, slug string) (*PriorityResponse, error)
	CreatePriority(ctx context.Context, req CreatePriorityRequest) (*PriorityResponse, error)
	UpdatePriority(ctx context.Context, slug string, req UpdatePriorityRequest) (*PriorityResponse, error)
	DeletePriority(ctx context.Context, slug string) error
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
func (c *APIClient) Ping(ctx context.Context) (*PingResponse, error) {
	pingResponse := &PingResponse{}
	apiError := &APIError{}
	response, err := c.client().Get("ping").Receive(pingResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not ping")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return pingResponse, nil
}

// Services returns a ServicesClient interface for interacting with services in FireHydrant
func (c *APIClient) Services() ServicesClient {
	return &RESTServicesClient{client: c}
}

// Runbooks returns a RunbooksClient interface for interacting with runbooks in FireHydrant
func (c *APIClient) Runbooks() RunbooksClient {
	return &RESTRunbooksClient{client: c}
}

// RunbookActions returns a RunbookActionsClient interface for interacting with runbook actions in FireHydrant
func (c *APIClient) RunbookActions() RunbookActionsClient {
	return &RESTRunbookActionsClient{client: c}
}

// GetEnvironment retrieves an environment from FireHydrant
func (c *APIClient) GetEnvironment(ctx context.Context, id string) (*EnvironmentResponse, error) {
	envResponse := &EnvironmentResponse{}
	apiError := &APIError{}
	response, err := c.client().Get("environments/"+id).Receive(envResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get environment")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return envResponse, nil
}

// CreateEnvironment creates an environment in FireHydrant
func (c *APIClient) CreateEnvironment(ctx context.Context, req CreateEnvironmentRequest) (*EnvironmentResponse, error) {
	envResponse := &EnvironmentResponse{}
	apiError := &APIError{}
	response, err := c.client().Post("environments").BodyJSON(&req).Receive(envResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create environment")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return envResponse, nil
}

// UpdateEnvironment updates a environment in FireHydrant
func (c *APIClient) UpdateEnvironment(ctx context.Context, id string, req UpdateEnvironmentRequest) (*EnvironmentResponse, error) {
	envResponse := &EnvironmentResponse{}
	apiError := &APIError{}
	response, err := c.client().Patch("environments/"+id).BodyJSON(&req).Receive(envResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update environment")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return envResponse, nil
}

// DeleteEnvironment deletes a environment from FireHydrant
func (c *APIClient) DeleteEnvironment(ctx context.Context, id string) error {
	apiError := &APIError{}
	response, err := c.client().Delete("environments/"+id).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete environment")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}

// GetFunctionality retrieves a functionality from FireHydrant
func (c *APIClient) GetFunctionality(ctx context.Context, id string) (*FunctionalityResponse, error) {
	funcResponse := &FunctionalityResponse{}
	apiError := &APIError{}
	response, err := c.client().Get("functionalities/"+id).Receive(funcResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get functionality")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return funcResponse, nil
}

// CreateFunctionality creates a functionality in FireHydrant
func (c *APIClient) CreateFunctionality(ctx context.Context, req CreateFunctionalityRequest) (*FunctionalityResponse, error) {
	funcResponse := &FunctionalityResponse{}
	apiError := &APIError{}
	response, err := c.client().Post("functionalities").BodyJSON(&req).Receive(funcResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create functionality")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return funcResponse, nil
}

// UpdateFunctionality updates a functionality in FireHydrant
func (c *APIClient) UpdateFunctionality(ctx context.Context, id string, req UpdateFunctionalityRequest) (*FunctionalityResponse, error) {
	funcResponse := &FunctionalityResponse{}
	apiError := &APIError{}
	response, err := c.client().Patch("functionalities/"+id).BodyJSON(&req).Receive(funcResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update functionality")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return funcResponse, nil
}

// DeleteFunctionality deletes a functionality from FireHydrant
func (c *APIClient) DeleteFunctionality(ctx context.Context, id string) error {
	apiError := &APIError{}
	response, err := c.client().Delete("functionalities/"+id).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete functionality")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}

// GetTeam retrieves a team from FireHydrant
func (c *APIClient) GetTeam(ctx context.Context, id string) (*TeamResponse, error) {
	teamResponse := &TeamResponse{}
	apiError := &APIError{}
	response, err := c.client().Get("teams/"+id).Receive(teamResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get team")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return teamResponse, nil
}

// CreateTeam creates a team in FireHydrant
func (c *APIClient) CreateTeam(ctx context.Context, req CreateTeamRequest) (*TeamResponse, error) {
	teamResponse := &TeamResponse{}
	apiError := &APIError{}
	response, err := c.client().Post("teams").BodyJSON(&req).Receive(teamResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create team")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return teamResponse, nil
}

// UpdateTeam updates a team in FireHydrant
func (c *APIClient) UpdateTeam(ctx context.Context, id string, req UpdateTeamRequest) (*TeamResponse, error) {
	teamResponse := &TeamResponse{}
	apiError := &APIError{}
	response, err := c.client().Patch("teams/"+id).BodyJSON(&req).Receive(teamResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update team")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return teamResponse, nil
}

// DeleteTeam deletes a team from FireHydrant
func (c *APIClient) DeleteTeam(ctx context.Context, id string) error {
	apiError := &APIError{}
	response, err := c.client().Delete("teams/"+id).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete team")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}

// GetSeverity retrieves a severity from FireHydrant
func (c *APIClient) GetSeverity(ctx context.Context, slug string) (*SeverityResponse, error) {
	sevResponse := &SeverityResponse{}
	apiError := &APIError{}
	response, err := c.client().Get("severities/"+slug).Receive(sevResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get severity")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return sevResponse, nil
}

// CreateSeverity creates a severity
func (c *APIClient) CreateSeverity(ctx context.Context, req CreateSeverityRequest) (*SeverityResponse, error) {
	sevResponse := &SeverityResponse{}
	apiError := &APIError{}
	response, err := c.client().Post("severities").BodyJSON(&req).Receive(sevResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create severity")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return sevResponse, nil
}

// UpdateSeverity updates a severity in FireHydrant
func (c *APIClient) UpdateSeverity(ctx context.Context, slug string, req UpdateSeverityRequest) (*SeverityResponse, error) {
	sevResponse := &SeverityResponse{}
	apiError := &APIError{}
	response, err := c.client().Patch("severities/"+slug).BodyJSON(&req).Receive(sevResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update severity")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return sevResponse, nil
}

// DeleteSeverity deletes a severity from FireHydrant
func (c *APIClient) DeleteSeverity(ctx context.Context, slug string) error {
	apiError := &APIError{}
	response, err := c.client().Delete("severities/"+slug).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete severity")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}

// GetPriority retrieves a priority from FireHydrant
func (c *APIClient) GetPriority(ctx context.Context, slug string) (*PriorityResponse, error) {
	priorityResponse := &PriorityResponse{}
	apiError := &APIError{}
	response, err := c.client().Get("priorities/"+slug).Receive(priorityResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get priority")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return priorityResponse, nil
}

// CreatePriority creates a priority
func (c *APIClient) CreatePriority(ctx context.Context, req CreatePriorityRequest) (*PriorityResponse, error) {
	priorityResponse := &PriorityResponse{}
	apiError := &APIError{}
	response, err := c.client().Post("priorities").BodyJSON(&req).Receive(priorityResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create priority")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return priorityResponse, nil
}

// UpdatePriority updates a priority in FireHydrant
func (c *APIClient) UpdatePriority(ctx context.Context, slug string, req UpdatePriorityRequest) (*PriorityResponse, error) {
	priorityResponse := &PriorityResponse{}
	apiError := &APIError{}
	response, err := c.client().Patch("priorities/"+slug).BodyJSON(&req).Receive(priorityResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update priority")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return priorityResponse, nil
}

// DeletePriority deletes a priority from FireHydrant
func (c *APIClient) DeletePriority(ctx context.Context, slug string) error {
	apiError := &APIError{}
	response, err := c.client().Delete("priorities/"+slug).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete priority")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
