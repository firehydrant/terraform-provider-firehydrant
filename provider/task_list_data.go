package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Singular services data source
func dataSourceTaskList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantTaskList,
		Schema: map[string]*schema.Schema{
			// Required
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_list_items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Computed
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"summary": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func dataFireHydrantTaskList(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the task list
	taskListID := d.Get("id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read task list: %s", taskListID), map[string]interface{}{
		"id": taskListID,
	})
	taskListResponse, err := firehydrantAPIClient.TaskLists().Get(ctx, taskListID)
	if err != nil {
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

	// Set the task list's ID in state
	d.SetId(taskListResponse.ID)

	return diag.Diagnostics{}
}
