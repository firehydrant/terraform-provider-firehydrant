package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantServices,
		Schema: map[string]*schema.Schema{
			// Optional
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed
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
	labelsSelector := firehydrant.LabelsSelector{}
	for key, value := range labels {
		labelsSelector[key] = value.(string)
	}
	tflog.Debug(ctx, fmt.Sprintf("Read services"), map[string]interface{}{
		"query":  query,
		"labels": labels,
	})
	servicesResponse, err := firehydrantAPIClient.Services().List(ctx, &firehydrant.ServiceQuery{
		Query:          query,
		LabelsSelector: labelsSelector,
	})
	if err != nil {
		return diag.Errorf("Error reading services: %v", err)
	}

	// Set the data source attributes to the values we got from the API
	services := make([]interface{}, 0)
	for _, service := range servicesResponse.Services {
		attributes := map[string]interface{}{
			"id":           service.ID,
			"alert_on_add": service.AlertOnAdd,
			"description":  service.Description,
			"labels":       service.Labels,
			"name":         service.Name,
			"service_tier": service.ServiceTier,
		}

		// Process any attributes that could be nil
		var links []interface{}
		for _, currentLink := range service.Links {
			links = append(links, map[string]interface{}{
				"href_url": currentLink.HrefURL,
				"name":     currentLink.Name,
			})
		}
		attributes["links"] = links

		if service.Owner != nil {
			attributes["owner_id"] = service.Owner.ID
		}

		var teamIDs []interface{}
		for _, team := range service.Teams {
			teamIDs = append(teamIDs, team.ID)
		}
		attributes["team_ids"] = teamIDs

		services = append(services, attributes)
	}
	if err := d.Set("services", services); err != nil {
		return diag.Errorf("Error setting services: %v", err)
	}

	d.SetId("does-not-matter")

	return diag.Diagnostics{}
}
