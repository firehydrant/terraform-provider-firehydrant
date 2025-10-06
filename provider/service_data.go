package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/firehydrant/firehydrant-go-sdk/models/sdkerrors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
			"auto_add_responding_team": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
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
	client := m.(*firehydrant.APIClient)

	// Get the service
	serviceID := d.Get("id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read service: %s", serviceID), map[string]interface{}{
		"id": serviceID,
	})
	serviceResponse, err := client.Sdk.CatalogEntries.GetService(ctx, serviceID)
	if err != nil {
		if sdkErr, ok := err.(*sdkerrors.SDKError); ok && sdkErr.StatusCode == 404 {
			return diag.Errorf("Service %s not found", serviceID)
		}
		return diag.Errorf("Error reading service %s: %v", serviceID, err)
	}

	// Ladder truck defines these types as `  expose :labels, documentation: {type: "object", desc: "An object of label key and values"} # rubocop:disable CustomCops/GrapeMissingType`
	// Previous implementation suggests these are always strings, adding Unmarshall into map[string]string to be defensive	labelsMap, err := unmarshalLabels(serviceResponse.Labels)
	if err != nil {
		return diag.Errorf("Error unmarshalling labels for service %s: %v", serviceID, err)
	}

	// Set values in state
	attributes := map[string]interface{}{
		"alert_on_add":             *serviceResponse.AlertOnAdd,
		"auto_add_responding_team": *serviceResponse.AutoAddRespondingTeam,
		"description":              *serviceResponse.Description,
		"labels":                   labelsMap,
		"name":                     *serviceResponse.Name,
		"service_tier":             *serviceResponse.ServiceTier,
	}

	// Process any attributes that could be nil

	// Process links
	var links []interface{}
	for _, currentLink := range serviceResponse.Links {
		links = append(links, map[string]interface{}{
			"href_url": *currentLink.HrefURL,
			"name":     *currentLink.Name,
		})
	}
	attributes["links"] = links

	// Process owner
	var ownerID string
	if serviceResponse.Owner != nil && serviceResponse.Owner.ID != nil {
		ownerID = *serviceResponse.Owner.ID
	}
	attributes["owner_id"] = ownerID

	// Process team IDs
	var teamIDs []interface{}
	for _, team := range serviceResponse.Teams {
		if team.ID != nil {
			teamIDs = append(teamIDs, *team.ID)
		}
	}
	attributes["team_ids"] = teamIDs

	// Set the data source attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for service %s: %v", key, serviceID, err)
		}
	}

	// Set the service's ID in state
	d.SetId(*serviceResponse.ID)

	return diag.Diagnostics{}
}
