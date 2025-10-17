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

func dataSourceOnCallSchedule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantOnCallSchedule,
		Schema: map[string]*schema.Schema{
			// Required
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slack_user_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"strategy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handoff_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handoff_day": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shift_duration": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"restrictions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_day": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_day": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataFireHydrantOnCallSchedule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Get the team
	id := d.Get("id").(string)
	teamID := d.Get("team_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read on-call schedule %s for team %s", id, teamID), map[string]interface{}{
		"id":      id,
		"team_id": teamID,
	})
	schedule, err := client.Sdk.Signals.GetTeamOnCallSchedule(ctx, teamID, id, nil, nil)
	if err != nil {
		return diag.Errorf("Error reading on-call schedule %s for team %s: %v", id, teamID, err)
	}

	// Gather values from API response
	attributes := dataFireHydrantOnCallScheduleToAttributesMap(teamID, *schedule)

	// Set the data source attributes to the values we got from the API
	for key, value := range attributes {
		if err = d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for on-call schedule %s: %v", key, id, err)
		}
	}

	// Set the schedule's ID in state
	d.SetId(*schedule.GetID())

	return diag.Diagnostics{}
}

func dataFireHydrantOnCallScheduleToAttributesMap(teamID string, schedule components.SignalsAPIOnCallScheduleEntity) map[string]interface{} {
	attributes := map[string]interface{}{
		"id":                  *schedule.GetID(),
		"team_id":             teamID,
		"name":                *schedule.GetName(),
		"description":         *schedule.GetDescription(),
		"time_zone":           *schedule.GetTimeZone(),
		"slack_user_group_id": *schedule.GetSlackUserGroupID(),
	}

	// Add strategy if available
	if strategy := schedule.GetStrategy(); strategy != nil {
		attributes["strategy"] = strategyToMapSDK(*strategy)
	}

	// Add restrictions if available
	attributes["restrictions"] = restrictionsToDataSDK(schedule.GetRestrictions())

	return attributes
}
