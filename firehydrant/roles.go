package firehydrant

import (
	"context"
	"fmt"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

type Roles interface {
	Get(ctx context.Context, id string) (*RoleResponse, error)
	List(ctx context.Context, params RolesQuery) (*RolesListResponse, error)
	Create(ctx context.Context, createReq CreateRoleRequest) (*RoleResponse, error)
	Update(ctx context.Context, id string, updateReq UpdateRoleRequest) (*RoleResponse, error)
	Delete(ctx context.Context, id string) error
}

type RESTRolesClient struct {
	client *APIClient
}

var _ Roles = &RESTRolesClient{}

// RoleResponse represents the API response for a single role
type RoleResponse struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	Slug           string       `json:"slug"`
	Description    string       `json:"description"`
	OrganizationID string       `json:"organization_id"`
	BuiltIn        bool         `json:"built_in"`
	ReadOnly       bool         `json:"read_only"`
	Permissions    []Permission `json:"permissions"`
	CreatedAt      string       `json:"created_at"`
	UpdatedAt      string       `json:"updated_at"`
}

// RolesListResponse represents the API response for listing roles
type RolesListResponse struct {
	Data       []RoleResponse `json:"data"`
	Pagination struct {
		Count int `json:"count"`
		Page  int `json:"page"`
		Items int `json:"items"`
		Pages int `json:"pages"`
		Prev  int `json:"prev"`
		Next  int `json:"next"`
		Last  int `json:"last"`
	} `json:"pagination"`
}

// RolesQuery represents query parameters for listing roles
type RolesQuery struct {
	Query string `url:"query,omitempty"`
	Page  int    `url:"page,omitempty"`
}

// CreateRoleRequest represents the request body for creating a role
type CreateRoleRequest struct {
	Name        string   `json:"name"`
	Slug        string   `json:"slug,omitempty"`
	Description string   `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// UpdateRoleRequest represents the request body for updating a role
type UpdateRoleRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

func (c *RESTRolesClient) restClient() *sling.Sling {
	return c.client.client()
}

func (c *RESTRolesClient) Get(ctx context.Context, id string) (*RoleResponse, error) {
	roleResponse := &RoleResponse{}
	apiError := &APIError{}

	response, err := c.restClient().Get(fmt.Sprintf("roles/%s", id)).Receive(roleResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get role")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return roleResponse, nil
}

func (c *RESTRolesClient) List(ctx context.Context, params RolesQuery) (*RolesListResponse, error) {
	rolesResponse := &RolesListResponse{}
	apiError := &APIError{}

	response, err := c.restClient().Get("roles").QueryStruct(params).Receive(rolesResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not list roles")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return rolesResponse, nil
}

func (c *RESTRolesClient) Create(ctx context.Context, createReq CreateRoleRequest) (*RoleResponse, error) {
	roleResponse := &RoleResponse{}
	apiError := &APIError{}

	response, err := c.restClient().Post("roles").BodyJSON(createReq).Receive(roleResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create role")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return roleResponse, nil
}

func (c *RESTRolesClient) Update(ctx context.Context, id string, updateReq UpdateRoleRequest) (*RoleResponse, error) {
	roleResponse := &RoleResponse{}
	apiError := &APIError{}

	response, err := c.restClient().Patch(fmt.Sprintf("roles/%s", id)).BodyJSON(updateReq).Receive(roleResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update role")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return roleResponse, nil
}

func (c *RESTRolesClient) Delete(ctx context.Context, id string) error {
	apiError := &APIError{}

	response, err := c.restClient().Delete(fmt.Sprintf("roles/%s", id)).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete role")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
