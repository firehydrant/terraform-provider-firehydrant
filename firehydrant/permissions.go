package firehydrant

import (
	"context"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

type Permissions interface {
	List(ctx context.Context) (*PermissionsListResponse, error)
	ListUser(ctx context.Context) (*UserPermissionsResponse, error)
	ListTeamPermissions(ctx context.Context) (*PermissionsListResponse, error)
}

type RESTPermissionsClient struct {
	client *APIClient
}

var _ Permissions = &RESTPermissionsClient{}

// PermissionsListResponse represents the API response for listing permissions
type PermissionsListResponse struct {
	Data []Permission `json:"data"`
}

// UserPermissionsResponse represents the current user's permissions (just slugs)
type UserPermissionsResponse struct {
	Data []string `json:"data"`
}

func (c *RESTPermissionsClient) restClient() *sling.Sling {
	return c.client.client()
}

func (c *RESTPermissionsClient) List(ctx context.Context) (*PermissionsListResponse, error) {
	permissionsResponse := &PermissionsListResponse{}
	apiError := &APIError{}

	response, err := c.restClient().Get("permissions").Receive(permissionsResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not list permissions")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return permissionsResponse, nil
}

func (c *RESTPermissionsClient) ListUser(ctx context.Context) (*UserPermissionsResponse, error) {
	permissionsResponse := &UserPermissionsResponse{}
	apiError := &APIError{}

	response, err := c.restClient().Get("permissions/user").Receive(permissionsResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not list current user permissions")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return permissionsResponse, nil
}

func (c *RESTPermissionsClient) ListTeamPermissions(ctx context.Context) (*PermissionsListResponse, error) {
	permissionsResponse := &PermissionsListResponse{}
	apiError := &APIError{}

	response, err := c.restClient().Get("permissions/team").Receive(permissionsResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not list team permissions")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return permissionsResponse, nil
}
