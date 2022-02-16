package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantService,
		Schema: map[string]*schema.Schema{
			// Required
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"alert_on_add": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"links": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Computed
						"href_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_tier": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"team_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the service
	serviceID := d.Get("id").(string)
	serviceResponse, err := firehydrantAPIClient.Services().Get(ctx, serviceID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set values in state
	attributes := map[string]interface{}{
		"alert_on_add": serviceResponse.AlertOnAdd,
		"description":  serviceResponse.Description,
		"name":         serviceResponse.Name,
		"service_tier": serviceResponse.ServiceTier,
	}

	// Process any attributes that could be nil
	var links []interface{}
	for _, currentLink := range serviceResponse.Links {
		links = append(links, map[string]interface{}{
			"href_url": currentLink.HrefURL,
			"name":     currentLink.Name,
		})
	}
	attributes["links"] = links

	if serviceResponse.Owner != nil {
		attributes["owner_id"] = serviceResponse.Owner.ID
	}

	var teamIDs []interface{}
	for _, team := range serviceResponse.Teams {
		teamIDs = append(teamIDs, team.ID)
	}
	attributes["team_ids"] = teamIDs

	// Set the data source attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
	}

	// Set the service's ID in state
	d.SetId(serviceResponse.ID)

	return diag.Diagnostics{}
}
