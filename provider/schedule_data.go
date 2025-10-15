package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
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
	client := m.(*firehydrant.APIClient)
	var schedule *components.ScheduleEntity
	// Get the schedule name
	name := d.Get("name").(string)
	tflog.Debug(ctx, fmt.Sprintf("Fetch schedule: %s", name), map[string]interface{}{
		"id": name,
	})

	scheduleResponse, err := client.Sdk.Teams.ListSchedules(ctx, &name, nil, nil)
	if err != nil {
		return diag.Errorf("Error fetching schedule '%s': %v", name, err)
	}

	schedules := scheduleResponse.GetData()
	if len(schedules) == 0 {
		return diag.Errorf("Did not find schedule matching '%s'", name)
	}
	if len(schedules) > 1 {
		for _, s := range schedules {
			if s.GetName() != nil && *s.GetName() == name {
				// we do not allow multiple schedules with the same name so we can return the first match
				schedule = &s
				break
			}
		}
		if schedule == nil {
			return diag.Errorf("Did not find schedule matching '%s'", name)
		}
	} else {
		schedule = &schedules[0]
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"id": *schedule.GetID(),
	}

	// Set the data source attributes to the values we got from the API
	for key, val := range attributes {
		if err := d.Set(key, val); err != nil {
			return diag.Errorf("Error setting %s for schedule %s: %v", key, name, err)
		}
	}

	// Set the schedule's ID in state
	d.SetId(*schedule.GetID())

	return diag.Diagnostics{}
}
