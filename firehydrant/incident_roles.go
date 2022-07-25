package firehydrant

import (
	"context"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// IncidentRolesClient is an interface for interacting with incident roles on FireHydrant
type IncidentRolesClient interface {
	Get(ctx context.Context, id string) (*IncidentRoleResponse, error)
	Create(ctx context.Context, createReq CreateIncidentRoleRequest) (*IncidentRoleResponse, error)
	Update(ctx context.Context, id string, updateReq UpdateIncidentRoleRequest) (*IncidentRoleResponse, error)
	Delete(ctx context.Context, id string) error
}

// RESTIncidentRolesClient implements the IncidentRolesClient interface
type RESTIncidentRolesClient struct {
	client *APIClient
}

func (c *RESTIncidentRolesClient) restClient() *sling.Sling {
	return c.client.client()
}

var _ IncidentRolesClient = &RESTIncidentRolesClient{}

// IncidentRoleResponse is the payload for retrieving an incident role
// URL: GET https://api.firehydrant.io/v1/incident_roles/{id}
type IncidentRoleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Summary     string `json:"summary"`

	CreatedAt   time.Time `json:"created_at"`
	DiscardedAt time.Time `json:"discarded_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Get returns an incident role from the FireHydrant API
func (c *RESTIncidentRolesClient) Get(ctx context.Context, id string) (*IncidentRoleResponse, error) {
	incidentRoleResponse := &IncidentRoleResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Get("incident_roles/"+id).Receive(incidentRoleResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get incident role")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return incidentRoleResponse, nil
}

// CreateIncidentRoleRequest is the payload for creating an incident role
// URL: POST https://api.firehydrant.io/v1/incident_roles
type CreateIncidentRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

// Create creates a brand spankin new incident role in FireHydrant
func (c *RESTIncidentRolesClient) Create(ctx context.Context, createReq CreateIncidentRoleRequest) (*IncidentRoleResponse, error) {
	incidentRoleResponse := &IncidentRoleResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Post("incident_roles").BodyJSON(&createReq).Receive(incidentRoleResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create incident role")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return incidentRoleResponse, nil
}

// UpdateIncidentRoleRequest is the payload for updating an incident role
// URL: PATCH https://api.firehydrant.io/v1/incident_roles/{id}
type UpdateIncidentRoleRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

// Update updates an incident role in FireHydrant
func (c *RESTIncidentRolesClient) Update(ctx context.Context, id string, updateReq UpdateIncidentRoleRequest) (*IncidentRoleResponse, error) {
	incidentRoleResponse := &IncidentRoleResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Patch("incident_roles/"+id).BodyJSON(updateReq).Receive(incidentRoleResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update incident role")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return incidentRoleResponse, nil
}

func (c *RESTIncidentRolesClient) Delete(ctx context.Context, id string) error {
	apiError := &APIError{}
	response, err := c.restClient().Delete("incident_roles/"+id).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete incident role")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
