package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/firehydrant-go-sdk/models/sdkerrors"
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
	client := m.(*firehydrant.APIClient)

	id := d.Get("id").(string)
	teamID := d.Get("team_id").(string)
	scheduleID := d.Get("schedule_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read rotation %s for schedule %s for team %s", id, scheduleID, teamID), map[string]interface{}{
		"id":          id,
		"team_id":     teamID,
		"schedule_id": scheduleID,
	})
	rotation, err := client.Sdk.Signals.GetOnCallScheduleRotation(ctx, id, teamID, scheduleID)
	if err != nil {
		if sdkErr, ok := err.(*sdkerrors.SDKError); ok && sdkErr.StatusCode == 404 {
			return diag.Errorf("Rotation %s not found", id)
		}
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
	d.SetId(*rotation.GetID())

	return diag.Diagnostics{}
}

func dataFireHydrantRotationToAttributesMap(teamID string, scheduleID string, rotation *components.SignalsAPIOnCallRotationEntity) map[string]interface{} {
	attributes := map[string]interface{}{
		"id":          *rotation.GetID(),
		"team_id":     teamID,
		"schedule_id": scheduleID,
		"name":        *rotation.GetName(),
		"time_zone":   *rotation.GetTimeZone(),
	}

	// Handle optional description field
	if description := rotation.GetDescription(); description != nil {
		attributes["description"] = *description
	}

	// Handle optional slack_user_group_id field
	if slackUserGroupID := rotation.GetSlackUserGroupID(); slackUserGroupID != nil && *slackUserGroupID != "" {
		attributes["slack_user_group_id"] = *slackUserGroupID
	}

	return attributes
}
