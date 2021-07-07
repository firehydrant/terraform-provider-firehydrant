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
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ServiceTier int32             `json:"service_tier,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// UpdateServiceRequest is the payload for updating a service
// URL: PATCH https://api.firehydrant.io/v1/services/{id}
type UpdateServiceRequest struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	ServiceTier int32             `json:"service_tier,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// ServiceResponse is the payload for retrieving a service
// URL: GET https://api.firehydrant.io/v1/services/{id}
type ServiceResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ServiceTier string            `json:"service_tier"`
	Slug        string            `json:"slug"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Labels      map[string]string `json:"labels"`
}

// ServiceQuery is the query used to search for services
type ServiceQuery struct {
	Query          string         `url:"query,omitempty"`
	ServiceTier    int32          `url:"service_tier,omitempty"`
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
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Services    []FunctionalityService `json:"services,omitempty"`
}

// TeamResponse is the payload for a single environment
// URL: GET https://api.firehydrant.io/v1/teams/{id}
type TeamResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Slug        string            `json:"slug"`
	Services    []ServiceResponse `json:"services"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// CreateTeamRequest is the payload for creating a service
// URL: POST https://api.firehydrant.io/v1/services
type CreateTeamRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ServiceIDs  []string `json:"service_ids,omitempty"`
}

// TeamService represents a service when creating a functionality
type TeamService struct {
	ID string `json:"id"`
}

// UpdateTeamRequest is the payload for updating a environment
// URL: PATCH https://api.firehydrant.io/v1/environments/{id}
type UpdateTeamRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	ServiceIDs  []string `json:"service_ids,omitempty"`
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
