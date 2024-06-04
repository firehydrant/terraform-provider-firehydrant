package firehydrant

import (
	"context"
	"fmt"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

type StatusUpdateTemplates interface {
	Get(ctx context.Context, id string) (*StatusUpdateTemplateResponse, error)
	Create(ctx context.Context, createReq CreateStatusUpdateTemplateRequest) (*StatusUpdateTemplateResponse, error)
	Update(ctx context.Context, id string, updateReq UpdateStatusUpdateTemplateRequest) (*StatusUpdateTemplateResponse, error)
	Delete(ctx context.Context, id string) error
}

type RESTStatusUpdateTemplateClient struct {
	client *APIClient
}

var _ StatusUpdateTemplates = &RESTStatusUpdateTemplateClient{}

type StatusUpdateTemplateResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Body string `json:"body"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateStatusUpdateTemplateRequest struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

type UpdateStatusUpdateTemplateRequest struct {
	Name string `json:"name"`
	Body string `json:"body"`
}

func (c *RESTStatusUpdateTemplateClient) restClient() *sling.Sling {
	return c.client.client()
}

func (c *RESTStatusUpdateTemplateClient) Create(ctx context.Context, createReq CreateStatusUpdateTemplateRequest) (*StatusUpdateTemplateResponse, error) {
	statusUpdateTemplateResponse := &StatusUpdateTemplateResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Post("status_update_templates").BodyJSON(&createReq).Receive(statusUpdateTemplateResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create status update template")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return statusUpdateTemplateResponse, nil
}

func (c *RESTStatusUpdateTemplateClient) Get(ctx context.Context, id string) (*StatusUpdateTemplateResponse, error) {
	statusUpdateTemplateResponse := &StatusUpdateTemplateResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Get(fmt.Sprintf("status_update_templates/%s", id)).Receive(statusUpdateTemplateResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get status update template")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return statusUpdateTemplateResponse, nil
}

func (c *RESTStatusUpdateTemplateClient) Update(ctx context.Context, id string, updateReq UpdateStatusUpdateTemplateRequest) (*StatusUpdateTemplateResponse, error) {
	statusUpdateTemplateResponse := &StatusUpdateTemplateResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Patch(fmt.Sprintf("status_update_templates/%s", id)).BodyJSON(&updateReq).Receive(statusUpdateTemplateResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update status update template")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return statusUpdateTemplateResponse, nil
}

func (c *RESTStatusUpdateTemplateClient) Delete(ctx context.Context, id string) error {
	apiError := &APIError{}
	response, err := c.restClient().Delete(fmt.Sprintf("status_update_templates/%s", id)).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete status update template")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
