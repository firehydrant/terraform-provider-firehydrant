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

func dataSourceFunctionality() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantFunctionality,
		Schema: map[string]*schema.Schema{
			// Required
			"functionality_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"service_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"team_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"auto_add_responding_team": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataFireHydrantFunctionality(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Get the functionality
	functionalityID := d.Get("functionality_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read functionality: %s", functionalityID), map[string]interface{}{
		"id": functionalityID,
	})
	functionalityResponse, err := client.Sdk.CatalogEntries.GetFunctionality(ctx, functionalityID)
	if err != nil {
		if sdkErr, ok := err.(*sdkerrors.SDKError); ok && sdkErr.StatusCode == 404 {
			return diag.Errorf("Functionality %s not found", functionalityID)
		}
		return diag.Errorf("Error reading functionality %s: %v", functionalityID, err)
	}

	// Ladder truck defines these types as `  expose :labels, documentation: {type: "object", desc: "An object of label key and values"} # rubocop:disable CustomCops/GrapeMissingType`
	// Previous implementation suggests these are always strings, adding Unmarshall into map[string]string to be defensive
	labelsMap, err := unmarshalLabels(functionalityResponse.Labels)
	if err != nil {
		return diag.Errorf("Error unmarshalling labels for functionality %s: %v", functionalityID, err)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"name":                     *functionalityResponse.Name,
		"description":              *functionalityResponse.Description,
		"labels":                   labelsMap,
		"auto_add_responding_team": *functionalityResponse.AutoAddRespondingTeam,
	}

	// Process service IDs
	serviceIDs := make([]string, 0)
	for _, service := range functionalityResponse.Services {
		if service.ID != nil {
			serviceIDs = append(serviceIDs, *service.ID)
		}
	}
	attributes["service_ids"] = serviceIDs

	// Process owner
	var ownerID string
	if functionalityResponse.Owner != nil && functionalityResponse.Owner.ID != nil {
		ownerID = *functionalityResponse.Owner.ID
	}
	attributes["owner_id"] = ownerID

	// Process team IDs
	var teamIDs []interface{}
	for _, team := range functionalityResponse.Teams {
		if team.ID != nil {
			teamIDs = append(teamIDs, *team.ID)
		}
	}
	attributes["team_ids"] = teamIDs

	// Set the data source attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for functionality %s: %v", key, functionalityID, err)
		}
	}

	// Set the functionality's ID in state
	d.SetId(*functionalityResponse.ID)

	return diag.Diagnostics{}
}
