package firehydrant

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

// Actor represents an actor doing things in the FireHydrant API
type Actor struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Type  string `json:"type"`
}

// PingResponse is the response the ping endpoint gives from FireHydrant
// URL: GET https://api.firehydrant.io/v1/ping
type PingResponse struct {
	Actor Actor `json:"actor"`
}

// CreateServiceRequest is the payload for creating a service
// URL: POST https://api.firehydrant.io/v1/services
type CreateServiceRequest struct {
	AlertOnAdd  bool              `json:"alert_on_add,omitempty"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels,omitempty"`
	Links       []ServiceLink     `json:"links,omitempty"`
	Name        string            `json:"name"`
	Owner       *ServiceTeam      `json:"owner,omitempty"`
	ServiceTier int               `json:"service_tier,int,omitempty"`
	Teams       []ServiceTeam     `json:"teams,omitempty"`
}

// ServiceTeam represents a team when creating a service
type ServiceTeam struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ServiceLink represents a link when creating/updating a service
type ServiceLink struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	HrefURL string `json:"href_url"`
}

// UpdateServiceRequest is the payload for updating a service
// URL: PATCH https://api.firehydrant.io/v1/services/{id}
type UpdateServiceRequest struct {
	AlertOnAdd           bool              `json:"alert_on_add"`
	Description          string            `json:"description"`
	Labels               map[string]string `json:"labels"`
	Links                []ServiceLink     `json:"links"`
	Name                 string            `json:"name,omitempty"`
	Owner                *ServiceTeam      `json:"owner"`
	RemoveOwner          bool              `json:"remove_owner,omitempty"`
	RemoveRemainingTeams bool              `json:"remove_remaining_teams,omitempty"`
	ServiceTier          int               `json:"service_tier,int"`
	Teams                []ServiceTeam     `json:"teams"`
}

// ServiceResponse is the payload for retrieving a service
// URL: GET https://api.firehydrant.io/v1/services/{id}
type ServiceResponse struct {
	ID          string            `json:"id"`
	AlertOnAdd  bool              `json:"alert_on_add"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	Links       []ServiceLink     `json:"links"`
	Name        string            `json:"name"`
	Owner       *ServiceTeam      `json:"owner"`
	ServiceTier int               `json:"service_tier"`
	Slug        string            `json:"slug"`
	Teams       []ServiceTeam     `json:"teams"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ServiceQuery is the query used to search for services
type ServiceQuery struct {
	Query          string         `url:"query,omitempty"`
	ServiceTier    int            `url:"int,service_tier,omitempty"`
	LabelsSelector LabelsSelector `url:"labels,omitempty"`
}

type LabelsSelector map[string]string

// EncodeValues implements Encoder
// https://github.com/google/go-querystring/blob/v1.0.0/query/encode.go#L39
func (sq LabelsSelector) EncodeValues(key string, v *url.Values) error {
	var labels []string

	keys, i := make([]string, len(sq)), 0
	for k := range sq {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, k := range keys {
		labels = append(labels, fmt.Sprintf("%s=%s", k, sq[k]))
	}

	v.Set(key, strings.Join(labels, ","))

	return nil
}

var _ query.Encoder = LabelsSelector{}

// ServicesResponse is the payload for retrieving a list of services
type ServicesResponse struct {
	Services []ServiceResponse `json:"data"`
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

// CreateEnvironmentRequest is the payload for creating a service
// URL: POST https://api.firehydrant.io/v1/services
type CreateEnvironmentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateEnvironmentRequest is the payload for updating a environment
// URL: PATCH https://api.firehydrant.io/v1/environments/{id}
type UpdateEnvironmentRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// FunctionalityResponse is the payload for a single environment
// URL: GET https://api.firehydrant.io/v1/functionalities/{id}
type FunctionalityResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Slug        string            `json:"slug"`
	Services    []ServiceResponse `json:"services"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// CreateFunctionalityRequest is the payload for creating a service
// URL: POST https://api.firehydrant.io/v1/services
type CreateFunctionalityRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Services    []FunctionalityService `json:"services,omitempty"`
}

// FunctionalityService represents a service when creating a functionality
type FunctionalityService struct {
	ID string `json:"id"`
}

// UpdateFunctionalityRequest is the payload for updating a environment
// URL: PATCH https://api.firehydrant.io/v1/environments/{id}
type UpdateFunctionalityRequest struct {
	Name                    string                 `json:"name,omitempty"`
	Description             string                 `json:"description"`
	RemoveRemainingServices bool                   `json:"remove_remaining_services"`
	Services                []FunctionalityService `json:"services"`
}

// TeamResponse is the payload for a single environment
// URL: GET https://api.firehydrant.io/v1/teams/{id}
type TeamResponse struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Slug          string            `json:"slug"`
	OwnedServices []ServiceResponse `json:"owned_services"`
	Services      []ServiceResponse `json:"services"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// CreateTeamRequest is the payload for creating a service
// URL: POST https://api.firehydrant.io/v1/services
type CreateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TeamService represents a service when creating a functionality
type TeamService struct {
	ID string `json:"id"`
}

// UpdateTeamRequest is the payload for updating a environment
// URL: PATCH https://api.firehydrant.io/v1/environments/{id}
type UpdateTeamRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// SeverityResponse is the payload for a single environment
// URL: GET https://api.firehydrant.io/v1/severities/{id}
type SeverityResponse struct {
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// CreateSeverityRequest is the payload for creating a service
// URL: POST https://api.firehydrant.io/v1/severities
type CreateSeverityRequest struct {
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// UpdateSeverityRequest is the payload for updating a environment
// URL: PATCH https://api.firehydrant.io/v1/severities/{id}
type UpdateSeverityRequest struct {
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

// PriorityResponse is the payload for a single priority
// URL: GET https://api.firehydrant.io/v1/priorities/{id}
type PriorityResponse struct {
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}

// CreatePriorityRequest is the payload for creating a priority
// URL: POST https://api.firehydrant.io/v1/priorities
type CreatePriorityRequest struct {
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}

// UpdatePriorityRequest is the payload for updating a priority
// URL: PATCH https://api.firehydrant.io/v1/priorities/{id}
type UpdatePriorityRequest struct {
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
}
