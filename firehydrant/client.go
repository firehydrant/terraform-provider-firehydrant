package firehydrant

import (
	"context"
	"fmt"
	"net/http"
	"os"

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

	Environments() EnvironmentsClient
	Functionalities() FunctionalitiesClient
	IncidentRoles() IncidentRolesClient
	Priorities() PrioritiesClient
	Runbooks() RunbooksClient
	RunbookActions() RunbookActionsClient
	ServiceDependencies() ServiceDependenciesClient
	Services() ServicesClient
	Severities() SeveritiesClient
	TaskLists() TaskListsClient

	// Teams
	GetTeam(ctx context.Context, id string) (*TeamResponse, error)
	CreateTeam(ctx context.Context, req CreateTeamRequest) (*TeamResponse, error)
	UpdateTeam(ctx context.Context, id string, req UpdateTeamRequest) (*TeamResponse, error)
	DeleteTeam(ctx context.Context, id string) error

	// Users
	GetUsers(ctx context.Context, params GetUserParams) (*UserResponse, error)

	// Schedules
	GetSchedules(ctx context.Context, params GetScheduleParams) (*ScheduleResponse, error)
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
	firehydrantBaseURL := os.Getenv("FIREHYDRANT_BASE_URL")
	if firehydrantBaseURL == "" {
		firehydrantBaseURL = DefaultBaseURL
	}

	c := &APIClient{
		baseURL: firehydrantBaseURL,
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

// Environments returns a EnvironmentsClient interface for interacting with environments in FireHydrant
func (c *APIClient) Environments() EnvironmentsClient {
	return &RESTEnvironmentsClient{client: c}
}

// Functionalities returns a FunctionalitiesClient interface for interacting with functionalities in FireHydrant
func (c *APIClient) Functionalities() FunctionalitiesClient {
	return &RESTFunctionalitiesClient{client: c}
}

// IncidentRoles returns a IncidentRolesClient interface for interacting with incident roles in FireHydrant
func (c *APIClient) IncidentRoles() IncidentRolesClient {
	return &RESTIncidentRolesClient{client: c}
}

// Priorities returns a PrioritiesClient interface for interacting with priorities in FireHydrant
func (c *APIClient) Priorities() PrioritiesClient {
	return &RESTPrioritiesClient{client: c}
}

// Runbooks returns a RunbooksClient interface for interacting with runbooks in FireHydrant
func (c *APIClient) Runbooks() RunbooksClient {
	return &RESTRunbooksClient{client: c}
}

// RunbookActions returns a RunbookActionsClient interface for interacting with runbook actions in FireHydrant
func (c *APIClient) RunbookActions() RunbookActionsClient {
	return &RESTRunbookActionsClient{client: c}
}

// ServiceDependencies returns a ServiceDependenciesClient interface for interacting with service dependencies in FireHydrant
func (c *APIClient) ServiceDependencies() ServiceDependenciesClient {
	return &RESTServiceDependenciesClient{client: c}
}

// Services returns a ServicesClient interface for interacting with services in FireHydrant
func (c *APIClient) Services() ServicesClient {
	return &RESTServicesClient{client: c}
}

// Severities returns a SeveritiesClient interface for interacting with severities in FireHydrant
func (c *APIClient) Severities() SeveritiesClient {
	return &RESTSeveritiesClient{client: c}
}

// TaskLists returns a TaskListsClient interface for interacting with task lists in FireHydrant
func (c *APIClient) TaskLists() TaskListsClient {
	return &RESTTaskListsClient{client: c}
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

// GetUsers gets matching users in FireHydrant
func (c *APIClient) GetUsers(ctx context.Context, params GetUserParams) (*UserResponse, error) {
	userResponse := &UserResponse{}
	apiError := &APIError{}
	response, err := c.client().Get("users").QueryStruct(params).Receive(userResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get users")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return userResponse, nil
}

// GetSchedules gets matching schedules in FireHydrant
func (c *APIClient) GetSchedules(ctx context.Context, params GetScheduleParams) (*ScheduleResponse, error) {
	scheduleResponse := &ScheduleResponse{}
	apiError := &APIError{}
	response, err := c.client().Get("schedules").QueryStruct(params).Receive(scheduleResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get schedules")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return scheduleResponse, nil
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
