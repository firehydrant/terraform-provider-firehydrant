package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantServices,
		Schema: map[string]*schema.Schema{
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceService(),
			},
		},
	}
}

func dataFireHydrantServices(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the services
	query := d.Get("query").(string)
	labels := d.Get("labels").(map[string]interface{})
	ls := firehydrant.LabelsSelector{}
	for k, v := range labels {
		ls[k] = v.(string)
	}
	r, err := firehydrantAPIClient.Services().List(ctx, &firehydrant.ServiceQuery{
		Query:          query,
		LabelsSelector: ls,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the data source attributes to the values we got from the API
	services := make([]interface{}, 0)
	for _, svc := range r.Services {
		values := map[string]interface{}{
			"id":           svc.ID,
			"alert_on_add": svc.AlertOnAdd,
			"description":  svc.Description,
			"name":         svc.Name,
			"service_tier": svc.ServiceTier,
		}

		// Process any attributes that could be nil
		if svc.Owner != nil {
			values["owner_id"] = svc.Owner.ID
		}

		var teamIDs []interface{}
		for _, team := range svc.Teams {
			teamIDs = append(teamIDs, team.ID)
		}
		values["team_ids"] = teamIDs

		services = append(services, values)
	}
	if err := d.Set("services", services); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("does-not-matter")

	return diag.Diagnostics{}
}
