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

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantTeam,
		UpdateContext: updateResourceFireHydrantTeam,
		ReadContext:   readResourceFireHydrantTeam,
		DeleteContext: deleteResourceFireHydrantTeam,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Optional
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"memberships": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_incident_role_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"schedule_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"slug": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func readResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Get the team
	teamID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read team: %s", teamID), map[string]interface{}{
		"id": teamID,
	})
	teamResponse, err := client.Sdk.Teams.GetTeam(ctx, teamID, nil)
	if err != nil {
		if sdkErr, ok := err.(*sdkerrors.SDKError); ok && sdkErr.StatusCode == 404 {
			tflog.Debug(ctx, fmt.Sprintf("Team %s no longer exists", teamID), map[string]interface{}{
				"id": teamID,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading team %s: %v", teamID, err)
	}

	// Set values in state
	attributes := map[string]interface{}{
		"name":        *teamResponse.Name,
		"description": *teamResponse.Description,
		"slug":        *teamResponse.Slug,
	}

	// Process memberships
	memberships := make([]map[string]interface{}, 0)
	for _, currentMembership := range teamResponse.Memberships {
		membership := map[string]interface{}{}

		// Handle default incident role
		if currentMembership.DefaultIncidentRole != nil && currentMembership.DefaultIncidentRole.ID != nil {
			membership["default_incident_role_id"] = *currentMembership.DefaultIncidentRole.ID
		} else {
			membership["default_incident_role_id"] = ""
		}

		// Handle schedule
		if currentMembership.Schedule != nil && currentMembership.Schedule.ID != nil {
			membership["schedule_id"] = *currentMembership.Schedule.ID
		} else {
			membership["schedule_id"] = ""
		}

		// Handle user
		if currentMembership.User != nil && currentMembership.User.ID != nil {
			membership["user_id"] = *currentMembership.User.ID
		} else {
			membership["user_id"] = ""
		}

		memberships = append(memberships, membership)
	}
	attributes["memberships"] = memberships

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for team %s: %v", key, teamID, err)
		}
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Construct the create team request
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	createRequest := components.CreateTeam{
		Name:        name,
		Description: &description,
	}
	if slug, ok := d.GetOk("slug"); ok {
		slugStr := slug.(string)
		createRequest.Slug = &slugStr
	}

	// Process any optional attributes and add to the create request if necessary
	memberships := d.Get("memberships")
	for _, currentMembership := range memberships.(*schema.Set).List() {
		membership := currentMembership.(map[string]interface{})
		teamMembership := components.CreateTeamMembership{}

		if roleID, ok := membership["default_incident_role_id"].(string); ok && roleID != "" {
			teamMembership.IncidentRoleID = &roleID
		}
		if scheduleID, ok := membership["schedule_id"].(string); ok && scheduleID != "" {
			teamMembership.ScheduleID = &scheduleID
		}
		if userID, ok := membership["user_id"].(string); ok && userID != "" {
			teamMembership.UserID = &userID
		}

		createRequest.Memberships = append(createRequest.Memberships, teamMembership)
	}

	// Create the new team
	tflog.Debug(ctx, fmt.Sprintf("Create team: %s", name), map[string]interface{}{
		"name": name,
	})
	teamResponse, err := client.Sdk.Teams.CreateTeam(ctx, createRequest)
	if err != nil {
		return diag.Errorf("Error creating team %s: %v", name, err)
	}

	// Set the new team's ID in state
	d.SetId(*teamResponse.ID)

	// Update state with the latest information from the API
	return readResourceFireHydrantTeam(ctx, d, m)
}

func updateResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Construct the update team request
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	updateRequest := components.UpdateTeam{
		Name:        &name,
		Description: &description,
	}
	if slug, ok := d.GetOk("slug"); ok {
		slugStr := slug.(string)
		updateRequest.Slug = &slugStr
	}

	// Process any optional attributes and add to the update request if necessary
	memberships := d.Get("memberships")
	for _, currentMembership := range memberships.(*schema.Set).List() {
		membership := currentMembership.(map[string]interface{})
		teamMembership := components.UpdateTeamMembership{}

		if roleID, ok := membership["default_incident_role_id"].(string); ok && roleID != "" {
			teamMembership.IncidentRoleID = &roleID
		}
		if scheduleID, ok := membership["schedule_id"].(string); ok && scheduleID != "" {
			teamMembership.ScheduleID = &scheduleID
		}
		if userID, ok := membership["user_id"].(string); ok && userID != "" {
			teamMembership.UserID = &userID
		}

		updateRequest.Memberships = append(updateRequest.Memberships, teamMembership)
	}

	// Update the team
	tflog.Debug(ctx, fmt.Sprintf("Update team: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})
	_, err := client.Sdk.Teams.UpdateTeam(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating team %s: %v", d.Id(), err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantTeam(ctx, d, m)
}

func deleteResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Delete the team
	teamID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Delete team: %s", teamID), map[string]interface{}{
		"id": teamID,
	})
	err := client.Sdk.Teams.DeleteTeam(ctx, teamID)
	if err != nil {
		if sdkErr, ok := err.(*sdkerrors.SDKError); ok && sdkErr.StatusCode == 404 {
			return nil
		}
		return diag.Errorf("Error deleting team %s: %v", teamID, err)
	}

	return diag.Diagnostics{}
}
