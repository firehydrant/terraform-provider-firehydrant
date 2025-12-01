package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/firehydrant-go-sdk/models/sdkerrors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFunctionality() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantFunctionality,
		UpdateContext: updateResourceFireHydrantFunctionality,
		ReadContext:   readResourceFireHydrantFunctionality,
		DeleteContext: deleteResourceFireHydrantFunctionality,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Optional
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"team_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"auto_add_responding_team": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func readResourceFireHydrantFunctionality(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Get the functionality
	functionalityID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read functionality: %s", functionalityID), map[string]interface{}{
		"id": functionalityID,
	})
	functionalityResponse, err := client.Sdk.CatalogEntries.GetFunctionality(ctx, functionalityID)
	if err != nil {
		if sdkErr, ok := err.(*sdkerrors.SDKError); ok && sdkErr.StatusCode == 404 {
			tflog.Debug(ctx, fmt.Sprintf("Functionality %s no longer exists", functionalityID), map[string]interface{}{
				"id": functionalityID,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading functionality %s: %v", functionalityID, err)
	}

	// Ladder truck defines these types as `  expose :labels, documentation: {type: "object", desc: "An object of label key and values"} # rubocop:disable CustomCops/GrapeMissingType`
	// Previous implementation suggests these are always strings, adding Unmarshall into map[string]string to be defensive
	labelsMap, err := unmarshalLabels(functionalityResponse.Labels)
	if err != nil {
		return diag.Errorf("Error unmarshalling labels for functionality %s: %v", functionalityID, err)
	}

	description := ""
	if functionalityResponse.Description != nil {
		description = *functionalityResponse.Description
	}

	autoAddRespondingTeam := false
	if functionalityResponse.AutoAddRespondingTeam != nil {
		autoAddRespondingTeam = *functionalityResponse.AutoAddRespondingTeam
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"name":                     *functionalityResponse.Name,
		"description":              description,
		"auto_add_responding_team": autoAddRespondingTeam,
		"labels":                   labelsMap,
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

	// Set the resource attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for functionality %s: %v", key, functionalityID, err)
		}
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantFunctionality(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Get attributes from config and construct the create request
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	autoAddRespondingTeam := d.Get("auto_add_responding_team").(bool)
	labels := convertStringMap(d.Get("labels").(map[string]interface{}))

	createRequest := components.CreateFunctionality{
		Name:                  name,
		Description:           &description,
		Labels:                labels,
		AutoAddRespondingTeam: &autoAddRespondingTeam,
	}

	// Process services
	serviceIDs := d.Get("service_ids").(*schema.Set).List()
	for _, serviceID := range serviceIDs {
		createRequest.Services = append(createRequest.Services, components.CreateFunctionalityService{
			ID: serviceID.(string),
		})
	}

	// Process owner if set
	if ownerID, ok := d.GetOk("owner_id"); ok && ownerID.(string) != "" {
		createRequest.Owner = &components.CreateFunctionalityOwner{
			ID: ownerID.(string),
		}
	}

	// Process team IDs if set
	teamIDs := d.Get("team_ids").(*schema.Set).List()
	for _, teamID := range teamIDs {
		createRequest.Teams = append(createRequest.Teams, components.CreateFunctionalityTeam{
			ID: teamID.(string),
		})
	}

	// Create the new functionality
	tflog.Debug(ctx, fmt.Sprintf("Create functionality: %s", name), map[string]interface{}{
		"name": name,
	})
	functionalityResponse, err := client.Sdk.CatalogEntries.CreateFunctionality(ctx, createRequest)
	if err != nil {
		return diag.Errorf("Error creating functionality %s: %v", name, err)
	}

	// Set the new functionality's ID in state
	d.SetId(*functionalityResponse.ID)

	// Update state with the latest information from the API
	return readResourceFireHydrantFunctionality(ctx, d, m)
}

func updateResourceFireHydrantFunctionality(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Construct the update request
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	autoAddRespondingTeam := d.Get("auto_add_responding_team").(bool)
	labels := convertStringMap(d.Get("labels").(map[string]interface{}))

	removeRemainingServices := true
	updateRequest := components.UpdateFunctionality{
		Name:                    &name,
		Description:             &description,
		Labels:                  labels,
		AutoAddRespondingTeam:   &autoAddRespondingTeam,
		RemoveRemainingServices: &removeRemainingServices,
	}

	// Process services
	// Always initialize Services as empty slice (not nil) so empty array is sent to clear services
	updateRequest.Services = []components.UpdateFunctionalityService{}
	serviceIDs := d.Get("service_ids").(*schema.Set).List()
	for _, serviceID := range serviceIDs {
		updateRequest.Services = append(updateRequest.Services, components.UpdateFunctionalityService{
			ID: serviceID.(string),
		})
	}

	// Process owner - set or remove
	ownerID, ownerIDSet := d.GetOk("owner_id")
	if ownerIDSet && ownerID.(string) != "" {
		updateRequest.Owner = &components.UpdateFunctionalityOwner{
			ID: ownerID.(string),
		}
	} else {
		removeOwner := true
		updateRequest.RemoveOwner = &removeOwner
	}

	// Process team IDs
	// Always initialize Teams as empty slice (not nil) so empty array is sent to clear teams
	updateRequest.Teams = []components.UpdateFunctionalityTeam{}
	teamIDs := d.Get("team_ids").(*schema.Set).List()
	for _, teamID := range teamIDs {
		updateRequest.Teams = append(updateRequest.Teams, components.UpdateFunctionalityTeam{
			ID: teamID.(string),
		})
	}
	// Force replacement of teams with the ones we send
	removeRemainingTeams := true
	updateRequest.RemoveRemainingTeams = &removeRemainingTeams

	// Update the functionality
	tflog.Debug(ctx, fmt.Sprintf("Update functionality: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})
	_, err := client.Sdk.CatalogEntries.UpdateFunctionality(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating functionality %s: %v", d.Id(), err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantFunctionality(ctx, d, m)
}

func deleteResourceFireHydrantFunctionality(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Delete the functionality
	functionalityID := d.Id()
	err := client.Sdk.CatalogEntries.DeleteFunctionality(ctx, functionalityID)
	if err != nil {
		if sdkErr, ok := err.(*sdkerrors.SDKError); ok && sdkErr.StatusCode == 404 {
			return nil
		}
		return diag.Errorf("Error deleting functionality %s: %v", functionalityID, err)
	}

	return diag.Diagnostics{}
}
