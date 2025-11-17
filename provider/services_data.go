package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/firehydrant/firehydrant-go-sdk/models/operations"
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
	client := m.(*firehydrant.APIClient)

	// Build the list services request
	query := d.Get("query").(string)
	labels := d.Get("labels").(map[string]interface{})

	// Convert labels map to comma-separated string format: key=value,key2=value2
	var labelsStr string
	if len(labels) > 0 {
		labelPairs := make([]string, 0, len(labels))
		for key, value := range labels {
			labelPairs = append(labelPairs, fmt.Sprintf("%s=%s", key, value.(string)))
		}
		labelsStr = strings.Join(labelPairs, ",")
	}

	tflog.Debug(ctx, "Read services", map[string]interface{}{
		"query":  query,
		"labels": labelsStr,
	})

	request := operations.ListServicesRequest{}
	if query != "" {
		request.Query = &query
	}
	if labelsStr != "" {
		request.Labels = &labelsStr
	}

	servicesResponse, err := client.Sdk.CatalogEntries.ListServices(ctx, request)
	if err != nil {
		return diag.Errorf("Error reading services: %v", err)
	}

	// Set the data source attributes to the values we got from the API
	services := make([]interface{}, 0)
	for _, service := range servicesResponse.Data {
		// Unmarshal labels from SDK struct to map[string]string
		labelsMap, err := unmarshalLabels(service.Labels)
		if err != nil {
			return diag.Errorf("Error unmarshalling labels for service %s: %v", *service.ID, err)
		}

		// Safely dereference pointers with default values
		description := ""
		if service.Description != nil {
			description = *service.Description
		}

		alertOnAdd := false
		if service.AlertOnAdd != nil {
			alertOnAdd = *service.AlertOnAdd
		}

		autoAddRespondingTeam := false
		if service.AutoAddRespondingTeam != nil {
			autoAddRespondingTeam = *service.AutoAddRespondingTeam
		}

		serviceTier := 0
		if service.ServiceTier != nil {
			serviceTier = *service.ServiceTier
		}

		attributes := map[string]interface{}{
			"id":                       *service.ID,
			"alert_on_add":             alertOnAdd,
			"auto_add_responding_team": autoAddRespondingTeam,
			"description":              description,
			"labels":                   labelsMap,
			"name":                     *service.Name,
			"service_tier":             serviceTier,
		}

		// Process any attributes that could be nil

		// Process links
		var links []interface{}
		for _, currentLink := range service.Links {
			links = append(links, map[string]interface{}{
				"href_url": *currentLink.HrefURL,
				"name":     *currentLink.Name,
			})
		}
		attributes["links"] = links

		// Process owner
		var ownerID string
		if service.Owner != nil && service.Owner.ID != nil {
			ownerID = *service.Owner.ID
		}
		attributes["owner_id"] = ownerID

		// Process team IDs
		var teamIDs []interface{}
		for _, team := range service.Teams {
			if team.ID != nil {
				teamIDs = append(teamIDs, *team.ID)
			}
		}
		attributes["team_ids"] = teamIDs

		services = append(services, attributes)
	}
	if err := d.Set("services", services); err != nil {
		return diag.Errorf("Error setting services: %v", err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}
