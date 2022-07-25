package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTaskList() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantTaskList,
		UpdateContext: updateResourceFireHydrantTaskList,
		ReadContext:   readResourceFireHydrantTaskList,
		DeleteContext: deleteResourceFireHydrantTaskList,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_list_items": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required
						"summary": {
							Type:     schema.TypeString,
							Required: true,
						},

						// Optional
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			// Optional
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func readResourceFireHydrantTaskList(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the task list
	taskListID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read task list: %s", taskListID), map[string]interface{}{
		"id": taskListID,
	})
	taskListResponse, err := firehydrantAPIClient.TaskLists().Get(ctx, taskListID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			tflog.Debug(ctx, fmt.Sprintf("Task list %s no longer exists", taskListID), map[string]interface{}{
				"id": taskListID,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading task list %s: %v", taskListID, err)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"name":        taskListResponse.Name,
		"description": taskListResponse.Description,
	}

	taskListItems := make([]interface{}, len(taskListResponse.TaskListItems))
	for index, currentTaskListItem := range taskListResponse.TaskListItems {
		taskListItems[index] = map[string]interface{}{
			"description": currentTaskListItem.Description,
			"summary":     currentTaskListItem.Summary,
		}
	}
	attributes["task_list_items"] = taskListItems

	// Set the resource attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for task list %s: %v", key, taskListID, err)
		}
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantTaskList(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get attributes from config and construct the create request
	createRequest := firehydrant.CreateTaskListRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Process any optional attributes and add to the create request if necessary
	taskListItems := d.Get("task_list_items").([]interface{})
	for _, currentTaskListItem := range taskListItems {
		taskListItem := currentTaskListItem.(map[string]interface{})

		createRequest.TaskListItems = append(createRequest.TaskListItems, firehydrant.TaskListItem{
			Description: taskListItem["description"].(string),
			Summary:     taskListItem["summary"].(string),
		})
	}

	// Create the new task list
	tflog.Debug(ctx, fmt.Sprintf("Create task list: %s", createRequest.Name), map[string]interface{}{
		"name": createRequest.Name,
	})
	taskListResponse, err := firehydrantAPIClient.TaskLists().Create(ctx, createRequest)
	if err != nil {
		return diag.Errorf("Error creating task list %s: %v", createRequest.Name, err)
	}

	// Set the new task list's ID in state
	d.SetId(taskListResponse.ID)

	// Update state with the latest information from the API
	return readResourceFireHydrantTaskList(ctx, d, m)
}

func updateResourceFireHydrantTaskList(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the update request
	updateRequest := firehydrant.UpdateTaskListRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Process any optional attributes and add to the update request if necessary
	taskListItems := d.Get("task_list_items").([]interface{})
	for _, currentTaskListItem := range taskListItems {
		taskListItem := currentTaskListItem.(map[string]interface{})

		updateRequest.TaskListItems = append(updateRequest.TaskListItems, firehydrant.TaskListItem{
			Description: taskListItem["description"].(string),
			Summary:     taskListItem["summary"].(string),
		})
	}

	// Update the task list
	tflog.Debug(ctx, fmt.Sprintf("Update task list: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})
	_, err := firehydrantAPIClient.TaskLists().Update(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating task list %s: %v", d.Id(), err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantTaskList(ctx, d, m)
}

func deleteResourceFireHydrantTaskList(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the task list
	taskListID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Delete task list: %s", taskListID), map[string]interface{}{
		"id": taskListID,
	})
	err := firehydrantAPIClient.TaskLists().Delete(ctx, taskListID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			return nil
		}
		return diag.Errorf("Error deleting task list %s: %v", taskListID, err)
	}

	return diag.Diagnostics{}
}
