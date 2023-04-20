package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTeams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantTeams,
		Schema: map[string]*schema.Schema{
			// Optional
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed
			"teams": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceTeam(),
			},
		},
	}
}

func dataFireHydrantTeams(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the teams
	query := d.Get("query").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read teams"), map[string]interface{}{
		"query": query,
	})
	teamsResponse, err := firehydrantAPIClient.Teams().List(ctx, &firehydrant.TeamQuery{
		Query: query,
	})
	if err != nil {
		return diag.Errorf("Error reading teams: %v", err)
	}

	// Set the data source attributes to the values we got from the API
	teams := make([]interface{}, 0)
	for _, team := range teamsResponse.Teams {
		attributes := map[string]interface{}{
			"id":          team.ID,
			"name":        team.Name,
			"description": team.Description,
			"slug":        team.Slug,
		}

		// Collect mapped service IDs
		serviceIDs := make([]string, 0)
		for _, service := range team.Services {
			serviceIDs = append(serviceIDs, service.ID)
		}
		attributes["service_ids"] = serviceIDs

		// Collect mapped owned service IDs
		ownedServiceIDs := make([]string, 0)
		for _, ownedService := range team.OwnedServices {
			ownedServiceIDs = append(ownedServiceIDs, ownedService.ID)
		}
		attributes["owned_service_ids"] = ownedServiceIDs

		teams = append(teams, attributes)
	}
	if err := d.Set("teams", teams); err != nil {
		return diag.Errorf("Error setting teams: %v", err)
	}

	d.SetId("does-not-matter")

	return diag.Diagnostics{}
}
