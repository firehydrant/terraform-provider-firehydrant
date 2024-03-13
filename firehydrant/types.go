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
	AlertOnAdd            bool               `json:"alert_on_add,omitempty"`
	AutoAddRespondingTeam bool               `json:"auto_add_responding_team,omitempty"`
	Description           string             `json:"description"`
	Labels                map[string]string  `json:"labels,omitempty"`
	Links                 []ServiceLink      `json:"links,omitempty"`
	Name                  string             `json:"name"`
	Owner                 *ServiceTeam       `json:"owner,omitempty"`
	ServiceTier           int                `json:"service_tier,int,omitempty"`
	Teams                 []ServiceTeam      `json:"teams,omitempty"`
	ExternalResources     []ExternalResource `json:"external_resources,omitempty"`
}

// ExternalResource is a nested object to link services to things like PagerDuty services
type ExternalResource struct {
	RemoteID       string `json:"remote_id"`
	ConnectionType string `json:"connection_type,omitempty"`
}

// RunbookTeam represents a team when creating a runbook
type RunbookTeam struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	AlertOnAdd                       bool               `json:"alert_on_add"`
	AutoAddRespondingTeam            bool               `json:"auto_add_responding_team"`
	Description                      string             `json:"description"`
	Labels                           map[string]string  `json:"labels"`
	Links                            []ServiceLink      `json:"links"`
	Name                             string             `json:"name,omitempty"`
	Owner                            *ServiceTeam       `json:"owner"`
	RemoveOwner                      bool               `json:"remove_owner,omitempty"`
	RemoveRemainingTeams             bool               `json:"remove_remaining_teams,omitempty"`
	RemoveRemainingExternalResources bool               `json:"remove_remaining_external_resources,omitempty"`
	ServiceTier                      int                `json:"service_tier,int"`
	Teams                            []ServiceTeam      `json:"teams"`
	ExternalResources                []ExternalResource `json:"external_resources,omitempty"`
}

// ServiceResponse is the payload for retrieving a service
// URL: GET https://api.firehydrant.io/v1/services/{id}
type ServiceResponse struct {
	ID                    string             `json:"id"`
	AlertOnAdd            bool               `json:"alert_on_add"`
	AutoAddRespondingTeam bool               `json:"auto_add_responding_team"`
	Description           string             `json:"description"`
	Labels                map[string]string  `json:"labels"`
	Links                 []ServiceLink      `json:"links"`
	Name                  string             `json:"name"`
	Owner                 *ServiceTeam       `json:"owner"`
	ServiceTier           int                `json:"service_tier"`
	Slug                  string             `json:"slug"`
	Teams                 []ServiceTeam      `json:"teams"`
	ExternalResources     []ExternalResource `json:"external_resources"`

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

// UserResponse is the payload for a user
// URL: GET https://api.firehydrant.io/v1/users
type GetUserParams struct {
	Query string `url:"query,omitempty"`
}

type User struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	SlackLinked bool      `json:"slack_linked?"`
	SlackUserId string    `json:"slack_user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserResponse struct {
	Users []User `json:"data"`
}

// ScheduleResponse is the payload for a schedule
// URL: GET https://api.firehydrant.io/v1/schedules
type GetScheduleParams struct {
	Query string `url:"query,omitempty"`
}

type Schedule struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Integration string `json:"integration"`
	Discarded   bool   `json:"discarded"`
}

type ScheduleResponse struct {
	Schedules []Schedule `json:"data"`
}

// TeamResponse is the payload for a single environment
// URL: GET https://api.firehydrant.io/v1/teams/{id}
type TeamResponse struct {
	ID            string               `json:"id"`
	Name          string               `json:"name"`
	Description   string               `json:"description"`
	Slug          string               `json:"slug"`
	OwnedServices []ServiceResponse    `json:"owned_services"`
	Services      []ServiceResponse    `json:"services"`
	Memberships   []MembershipResponse `json:"memberships"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// TeamsResponse is the payload for retrieving a list of teams
type TeamsResponse struct {
	Teams      []TeamResponse `json:"data"`
	Pagination *Pagination    `json:"pagination,omitempty"`
}

// MembershipResponse represents the response coming back from teams
// for membership
type MembershipResponse struct {
	DefaultIncidentRole IncidentRoleResponse `json:"default_incident_role,omitempty"`
	Schedule            Schedule             `json:"schedule,omitempty"`
	User                User                 `json:"user,omitempty"`
}

// Membership represents a user_id or schedule_id along with a
// incident_role_id for a team membership resource
type Membership struct {
	IncidentRoleId string `json:"incident_role_id,omitempty"`
	ScheduleId     string `json:"schedule_id,omitempty"`
	UserId         string `json:"user_id,omitempty"`
}

// TeamQuery is the query used to search for teams
type TeamQuery struct {
	Query string `url:"query,omitempty"`
	Page  int    `url:"page,omitempty"`
}

// CreateTeamRequest is the payload for creating a service
// URL: POST https://api.firehydrant.io/v1/services
type CreateTeamRequest struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Slug        string       `json:"slug,omitempty"`
	Memberships []Membership `json:"memberships,omitempty"`
}

// TeamService represents a service when creating a functionality
type TeamService struct {
	ID string `json:"id"`
}

// UpdateTeamRequest is the payload for updating a environment
// URL: PATCH https://api.firehydrant.io/v1/environments/{id}
type UpdateTeamRequest struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Slug        string       `json:"slug,omitempty"`
	Memberships []Membership `json:"memberships,omitempty"`
}

type Pagination struct {
	Count int `json:"count"`
	Page  int `json:"page"`
	Items int `json:"items"`
	Pages int `json:"pages"`
	Last  int `json:"last"`
	Prev  int `json:"prev,omitempty"`
	Next  int `json:"next,omitempty"`
}

// FunctionalityResponse is the payload for a single environment
// URL: GET https://api.firehydrant.io/v1/functionalities/{id}
type FunctionalityResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Slug        string                 `json:"slug"`
	Services    []FunctionalityService `json:"services"`
	Labels      map[string]string      `json:"labels"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// FunctionalityService represents a service when creating a functionality
type FunctionalityService struct {
	ID string `json:"id"`
}

// CreateFunctionalityRequest is the payload for creating a service
// URL: POST https://api.firehydrant.io/v1/services
type CreateFunctionalityRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Services    []FunctionalityService `json:"services,omitempty"`
	Labels      map[string]string      `json:"labels"`
}

// UpdateFunctionalityRequest is the payload for updating a environment
// URL: PATCH https://api.firehydrant.io/v1/environments/{id}
type UpdateFunctionalityRequest struct {
	Name                    string                 `json:"name,omitempty"`
	Description             string                 `json:"description"`
	RemoveRemainingServices bool                   `json:"remove_remaining_services"`
	Labels                  map[string]string      `json:"labels"`
	Services                []FunctionalityService `json:"services"`
}

// SlackChannelResponse is the response for retrieving Slack channel information, including FireHydrant ID.
// URL: GET https://api.firehydrant.io/v1/integrations/slack/channels?slack_channel_id={id}
type SlackChannelResponse struct {
	ID             string `json:"id"`
	SlackChannelID string `json:"slack_channel_id"`
	Name           string `json:"name"`
}

type SlackChannelsResponse struct {
	Channels   []*SlackChannelResponse `json:"data"`
	Pagination *Pagination             `json:"pagination,omitempty"`
}
