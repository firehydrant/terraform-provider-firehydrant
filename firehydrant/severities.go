package firehydrant

import (
	"context"
	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// SeveritiesClient is an interface for interacting with severities on FireHydrant
type SeveritiesClient interface {
	Get(ctx context.Context, slug string) (*SeverityResponse, error)
	Create(ctx context.Context, createReq CreateSeverityRequest) (*SeverityResponse, error)
	Update(ctx context.Context, slug string, updateReq UpdateSeverityRequest) (*SeverityResponse, error)
	Delete(ctx context.Context, slug string) error
}

// RESTSeveritiesClient implements the SeveritiesClient interface
type RESTSeveritiesClient struct {
	client *APIClient
}

var _ SeveritiesClient = &RESTSeveritiesClient{}

func (c *RESTSeveritiesClient) restClient() *sling.Sling {
	return c.client.client()
}

// SeverityResponse is the payload for a single environment
// URL: GET https://api.firehydrant.io/v1/severities/{id}
type SeverityResponse struct {
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// Get retrieves a severity from FireHydrant
func (c *RESTSeveritiesClient) Get(ctx context.Context, slug string) (*SeverityResponse, error) {
	sevResponse := &SeverityResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Get("severities/"+slug).Receive(sevResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get severity")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return sevResponse, nil
}

// CreateSeverityRequest is the payload for creating a severity
// URL: POST https://api.firehydrant.io/v1/severities
type CreateSeverityRequest struct {
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// Create creates a severity
func (c *RESTSeveritiesClient) Create(ctx context.Context, createReq CreateSeverityRequest) (*SeverityResponse, error) {
	sevResponse := &SeverityResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Post("severities").BodyJSON(&createReq).Receive(sevResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create severity")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return sevResponse, nil
}

// UpdateSeverityRequest is the payload for updating a severity
// URL: PATCH https://api.firehydrant.io/v1/severities/{id}
type UpdateSeverityRequest struct {
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

// Update updates a severity in FireHydrant
func (c *RESTSeveritiesClient) Update(ctx context.Context, slug string, updateReq UpdateSeverityRequest) (*SeverityResponse, error) {
	sevResponse := &SeverityResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Patch("severities/"+slug).BodyJSON(&updateReq).Receive(sevResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update severity")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return sevResponse, nil
}

// Delete deletes a severity from FireHydrant
func (c *RESTSeveritiesClient) Delete(ctx context.Context, slug string) error {
	apiError := &APIError{}
	response, err := c.restClient().Delete("severities/"+slug).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete severity")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
