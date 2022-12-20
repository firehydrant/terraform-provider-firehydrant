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
						"incident_role_id": {
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
		},
	}
}

func readResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the team
	teamID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read team: %s", teamID), map[string]interface{}{
		"id": teamID,
	})
	teamResponse, err := firehydrantAPIClient.GetTeam(ctx, teamID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
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
		"name":        teamResponse.Name,
		"description": teamResponse.Description,
	}

	// Process any attributes that could be nil
	memberships := make([]map[string]interface{}, len(teamResponse.Memberships))
	for index, currentMembership := range teamResponse.Memberships {
		memberships[index] = map[string]interface{}{
			"incident_role_id": currentMembership.DefaultIncidentRole.ID,
			"schedule_id":      currentMembership.Schedule.ID,
			"user_id":          currentMembership.User.ID,
		}
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
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the create team request
	createRequest := firehydrant.CreateTeamRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Process any optional attributes and add to the create request if necessary
	memberships := d.Get("memberships")
	for _, currentMembership := range memberships.(*schema.Set).List() {
		membership := currentMembership.(map[string]interface{})
		createRequest.Memberships = append(createRequest.Memberships, firehydrant.Membership{
			IncidentRoleId: membership["incident_role_id"].(string),
			ScheduleId:     membership["schedule_id"].(string),
			UserId:         membership["user_id"].(string),
		})
	}

	// Create the new team
	tflog.Debug(ctx, fmt.Sprintf("Create team: %s", createRequest.Name), map[string]interface{}{
		"name": createRequest.Name,
	})
	teamResponse, err := firehydrantAPIClient.CreateTeam(ctx, createRequest)
	if err != nil {
		return diag.Errorf("Error creating team %s: %v", createRequest.Name, err)
	}

	// Set the new team's ID in state
	d.SetId(teamResponse.ID)

	// Update state with the latest information from the API
	return readResourceFireHydrantTeam(ctx, d, m)
}

func updateResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the update team request
	updateRequest := firehydrant.UpdateTeamRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Process any optional attributes and add to the update request if necessary
	memberships := d.Get("memberships")
	for _, currentMembership := range memberships.(*schema.Set).List() {
		membership := currentMembership.(map[string]interface{})
		updateRequest.Memberships = append(updateRequest.Memberships, firehydrant.Membership{
			IncidentRoleId: membership["incident_role_id"].(string),
			ScheduleId:     membership["schedule_id"].(string),
			UserId:         membership["user_id"].(string),
		})
	}

	// Update the team
	tflog.Debug(ctx, fmt.Sprintf("Update team: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})
	_, err := firehydrantAPIClient.UpdateTeam(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating team %s: %v", d.Id(), err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantTeam(ctx, d, m)
}

func deleteResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the team
	teamID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Delete team: %s", teamID), map[string]interface{}{
		"id": teamID,
	})
	err := firehydrantAPIClient.DeleteTeam(ctx, teamID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			return nil
		}
		return diag.Errorf("Error deleting team %s: %v", teamID, err)
	}

	return diag.Diagnostics{}
}
