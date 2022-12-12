package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSchedule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantSchedule,
		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataFireHydrantSchedule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the user
	name := d.Get("name").(string)
	tflog.Debug(ctx, fmt.Sprintf("Fetch schedule: %s", name), map[string]interface{}{
		"id": name,
	})

	params := firehydrant.GetScheduleParams{Query: name}
	scheduleResponse, err := firehydrantAPIClient.GetSchedules(ctx, params)
	if err != nil {
		return diag.Errorf("Error fetching schedule '%s': %v", name, err)
	}

	if len(scheduleResponse.Schedules) == 0 {
		return diag.Errorf("Did not find schedule matching '%s'", name)
	}
	if len(scheduleResponse.Schedules) > 1 {
		return diag.Errorf("Found multiple matching schedules for '%s'", name)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"id": scheduleResponse.Schedules[0].ID,
	}

	// Set the data source attributes to the values we got from the API
	for key, val := range attributes {
		if err := d.Set(key, val); err != nil {
			return diag.Errorf("Error setting %s for user %s: %v", key, name, err)
		}
	}

	// Set the schedule's ID in state
	d.SetId(scheduleResponse.Schedules[0].ID)

	return diag.Diagnostics{}
}
