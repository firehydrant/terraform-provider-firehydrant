package firehydrant

import (
	"context"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// TaskListsClient is an interface for interacting with task lists on FireHydrant
type TaskListsClient interface {
	Get(ctx context.Context, id string) (*TaskListResponse, error)
	Create(ctx context.Context, createReq CreateTaskListRequest) (*TaskListResponse, error)
	Update(ctx context.Context, id string, updateReq UpdateTaskListRequest) (*TaskListResponse, error)
	Delete(ctx context.Context, id string) error
}

// RESTTaskListsClient implements the TaskListsClient interface
type RESTTaskListsClient struct {
	client *APIClient
}

var _ TaskListsClient = &RESTTaskListsClient{}

func (c *RESTTaskListsClient) restClient() *sling.Sling {
	return c.client.client()
}

// TaskListResponse is the payload for retrieving a service
// URL: GET https://api.firehydrant.io/v1/task_lists/{id}
type TaskListResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	TaskListItems []TaskListItem `json:"task_list_items"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TaskListItem is an item in a task list
type TaskListItem struct {
	Description string `json:"description"`
	Summary     string `json:"summary"`
}

// Get returns a task list from the FireHydrant API
func (c *RESTTaskListsClient) Get(ctx context.Context, id string) (*TaskListResponse, error) {
	taskListResponse := &TaskListResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Get("task_lists/"+id).Receive(taskListResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get task list")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return taskListResponse, nil
}

// CreateTaskListRequest is the payload for creating a task list
// URL: POST https://api.firehydrant.io/v1/task_list
type CreateTaskListRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	TaskListItems []TaskListItem `json:"task_list_items"`
}

// Create creates a brand spankin new task list in FireHydrant
func (c *RESTTaskListsClient) Create(ctx context.Context, createReq CreateTaskListRequest) (*TaskListResponse, error) {
	taskListResponse := &TaskListResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Post("task_lists").BodyJSON(&createReq).Receive(taskListResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create task list")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return taskListResponse, nil
}

// UpdateTaskListRequest is the payload for updating a task list
// URL: PATCH https://api.firehydrant.io/v1/task_lists/{id}
type UpdateTaskListRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description"`

	TaskListItems []TaskListItem `json:"task_list_items"`
}

// Update updates a task list in FireHydrant
func (c *RESTTaskListsClient) Update(ctx context.Context, id string, updateReq UpdateTaskListRequest) (*TaskListResponse, error) {
	taskListResponse := &TaskListResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Patch("task_lists/"+id).BodyJSON(updateReq).Receive(taskListResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update task list")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return taskListResponse, nil
}

func (c *RESTTaskListsClient) Delete(ctx context.Context, id string) error {
	apiError := &APIError{}
	response, err := c.restClient().Delete("task_lists/"+id).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete task list")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
