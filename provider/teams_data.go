package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/firehydrant-go-sdk/models/operations"
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
	client := m.(*firehydrant.APIClient)

	// Build the list teams request
	query := d.Get("query").(string)
	tflog.Debug(ctx, "Read teams", map[string]interface{}{
		"query": query,
	})

	request := operations.ListTeamsRequest{}
	if query != "" {
		request.Query = &query
	}

	allTeams, err := fetchAllPages(ctx,
		func(ctx context.Context, req operations.ListTeamsRequest) ([]components.TeamEntity, *components.NullablePaginationEntity, error) {
			resp, err := client.Sdk.Teams.ListTeams(ctx, req)
			if err != nil {
				return nil, nil, err
			}
			return resp.GetData(), resp.GetPagination(), nil
		},
		func(req *operations.ListTeamsRequest, page int) { req.Page = &page },
		request,
	)
	if err != nil {
		return diag.Errorf("Error reading teams: %v", err)
	}

	// Set the data source attributes to the values we got from the API
	teams := make([]interface{}, 0)
	for _, team := range allTeams {
		attributes := map[string]interface{}{
			"id":          *team.GetID(),
			"name":        *team.GetName(),
			"description": *team.GetDescription(),
			"slug":        *team.GetSlug(),
		}

		// Collect mapped service IDs
		serviceIDs := make([]string, 0)
		for _, service := range team.GetServices() {
			if service.GetID() != nil {
				serviceIDs = append(serviceIDs, *service.GetID())
			}
		}
		attributes["service_ids"] = serviceIDs

		// Collect mapped owned service IDs
		ownedServiceIDs := make([]string, 0)
		for _, ownedService := range team.GetOwnedServices() {
			if ownedService.GetID() != nil {
				ownedServiceIDs = append(ownedServiceIDs, *ownedService.GetID())
			}
		}
		attributes["owned_service_ids"] = ownedServiceIDs

		teams = append(teams, attributes)
	}
	if err := d.Set("teams", teams); err != nil {
		return diag.Errorf("Error setting teams: %v", err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}
