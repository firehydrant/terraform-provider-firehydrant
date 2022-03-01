package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
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
		},
	}
}

func readResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the team
	teamResponse, err := firehydrantAPIClient.GetTeam(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Set values in state
	attributes := map[string]interface{}{
		"name":        teamResponse.Name,
		"description": teamResponse.Description,
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the create team request
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	createRequest := firehydrant.CreateTeamRequest{
		Name:        name,
		Description: description,
	}

	// Create the new team
	teamResponse, err := firehydrantAPIClient.CreateTeam(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
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
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	updateRequest := firehydrant.UpdateTeamRequest{
		Name:        name,
		Description: description,
	}

	// Update the team
	teamID := d.Id()
	_, err := firehydrantAPIClient.UpdateTeam(ctx, teamID, updateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantTeam(ctx, d, m)
}

func deleteResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the team
	teamID := d.Id()
	err := firehydrantAPIClient.DeleteTeam(ctx, teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
