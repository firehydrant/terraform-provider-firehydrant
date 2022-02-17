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
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alert_on_add": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
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
	r, err := firehydrantAPIClient.Services().Get(ctx, serviceID)
	if err != nil {
		return diag.FromErr(err)
	}

	svc := map[string]interface{}{
		"alert_on_add": r.AlertOnAdd,
		"description":  r.Description,
		"name":         r.Name,
		"service_tier": r.ServiceTier,
	}

	// Process any attributes that could be nil
	if r.Owner != nil {
		svc["owner_id"] = r.Owner.ID
	}

	var teamIDs []interface{}
	for _, team := range r.Teams {
		teamIDs = append(teamIDs, team.ID)
	}
	svc["team_ids"] = teamIDs

	// Set the data source attributes to the values we got from the API
	for key, val := range svc {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(r.ID)

	return diag.Diagnostics{}
}
