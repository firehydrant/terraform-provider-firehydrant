package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantTeam,
		Schema: map[string]*schema.Schema{
			// Required
			"id": {
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
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owned_service_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the team
	id := d.Get("id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read team: %s", id), map[string]interface{}{
		"id": id,
	})
	teamResponse, err := firehydrantAPIClient.Teams().Get(ctx, id)
	if err != nil {
		return diag.Errorf("Error reading team %s: %v", id, err)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"id":          teamResponse.ID,
		"name":        teamResponse.Name,
		"description": teamResponse.Description,
		"slug":        teamResponse.Slug,
	}

	// Collect mapped service IDs
	serviceIDs := make([]string, 0)
	for _, service := range teamResponse.Services {
		serviceIDs = append(serviceIDs, service.ID)
	}
	attributes["service_ids"] = serviceIDs

	// Collect mapped owned service IDs
	ownedServiceIDs := make([]string, 0)
	for _, service := range teamResponse.OwnedServices {
		ownedServiceIDs = append(ownedServiceIDs, service.ID)
	}
	attributes["owned_service_ids"] = ownedServiceIDs

	// Set the data source attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for team %s: %v", key, id, err)
		}
	}

	// Set the team's ID in state
	d.SetId(teamResponse.ID)

	return diag.Diagnostics{}
}
