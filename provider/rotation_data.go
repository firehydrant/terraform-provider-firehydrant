package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRotation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantRotation,
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
			"schedule_id": {
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
		},
	}
}

func dataFireHydrantRotation(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	id := d.Get("id").(string)
	teamID := d.Get("team_id").(string)
	scheduleID := d.Get("schedule_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read rotation %s for schedule %s for team %s", id, scheduleID, teamID), map[string]interface{}{
		"id":          id,
		"team_id":     teamID,
		"schedule_id": scheduleID,
	})
	rotation, err := firehydrantAPIClient.Rotations().Get(ctx, teamID, scheduleID, id)
	if err != nil {
		return diag.Errorf("Error reading rotation %s for schedule %s for team %s: %v", id, scheduleID, teamID, err)
	}

	// Gather values from API response
	attributes := dataFireHydrantRotationToAttributesMap(teamID, scheduleID, rotation)

	// Set the data source attributes to the values we got from the API
	for key, value := range attributes {
		if err = d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for on-call schedule %s: %v", key, id, err)
		}
	}

	// Set the schedule's ID in state
	d.SetId(rotation.ID)

	return diag.Diagnostics{}
}

func dataFireHydrantRotationToAttributesMap(teamID string, scheduleID string, rotation *firehydrant.RotationResponse) map[string]interface{} {
	return map[string]interface{}{
		"id":                  rotation.ID,
		"team_id":             teamID,
		"schedule_id":         scheduleID,
		"name":                rotation.Name,
		"description":         rotation.Description,
		"time_zone":           rotation.TimeZone,
		"slack_user_group_id": rotation.SlackUserGroupID,
	}
}
