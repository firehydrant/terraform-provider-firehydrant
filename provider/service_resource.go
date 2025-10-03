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

func resourceService() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantService,
		UpdateContext: updateResourceFireHydrantService,
		ReadContext:   readResourceFireHydrantService,
		DeleteContext: deleteResourceFireHydrantService,
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
			"alert_on_add": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"auto_add_responding_team": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"links": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required
						"href_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_tier": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"team_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"external_resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required
						"remote_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"connection_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func readResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Get the service
	serviceID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read service: %s", serviceID), map[string]interface{}{
		"id": serviceID,
	})
	serviceResponse, err := client.Sdk.CatalogEntries.GetService(ctx, serviceID)
	if err != nil {
		if sdkErr, ok := err.(*sdkerrors.SDKError); ok && sdkErr.StatusCode == 404 {
			tflog.Debug(ctx, fmt.Sprintf("Service %s no longer exists", serviceID), map[string]interface{}{
				"id": serviceID,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading service %s: %v", serviceID, err)
	}

	// Set values in state
	attributes := map[string]interface{}{
		"name":                     *serviceResponse.Name,
		"alert_on_add":             *serviceResponse.AlertOnAdd,
		"auto_add_responding_team": *serviceResponse.AutoAddRespondingTeam,
		"description":              *serviceResponse.Description,
		"labels":                   make(map[string]interface{}), // Labels handling TBD - empty for now
		"service_tier":             *serviceResponse.ServiceTier,
	}

	// Process links
	links := make([]map[string]interface{}, len(serviceResponse.Links))
	for index, currentLink := range serviceResponse.Links {
		links[index] = map[string]interface{}{
			"href_url": *currentLink.HrefURL,
			"name":     *currentLink.Name,
		}
	}
	attributes["links"] = links

	// Process external resources
	ers := make([]map[string]interface{}, len(serviceResponse.ExternalResources))
	for index, currentER := range serviceResponse.ExternalResources {
		ers[index] = map[string]interface{}{
			"remote_id":       *currentER.RemoteID,
			"connection_type": *currentER.ConnectionType,
		}
	}
	attributes["external_resources"] = ers

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

	// Update the resource attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for service %s: %v", key, serviceID, err)
		}
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Get attributes from config and construct the create request
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	alertOnAdd := d.Get("alert_on_add").(bool)
	autoAddRespondingTeam := d.Get("auto_add_responding_team").(bool)
	serviceTier := d.Get("service_tier").(int)
	labels := convertStringMap(d.Get("labels").(map[string]interface{}))

	createRequest := components.CreateService{
		Name:                  name,
		Description:           &description,
		Labels:                labels,
		ServiceTier:           (*components.CreateServiceServiceTier)(&serviceTier),
		AlertOnAdd:            &alertOnAdd,
		AutoAddRespondingTeam: &autoAddRespondingTeam,
	}

	// Process links
	configLinks := d.Get("links").(*schema.Set).List()
	for _, currentLink := range configLinks {
		link := currentLink.(map[string]interface{})
		hrefURL := link["href_url"].(string)
		linkName := link["name"].(string)
		createRequest.Links = append(createRequest.Links, components.CreateServiceLink{
			HrefURL: hrefURL,
			Name:    linkName,
		})
	}

	// Process owner if set
	if ownerID, ok := d.GetOk("owner_id"); ok && ownerID.(string) != "" {
		createRequest.Owner = &components.CreateServiceOwner{
			ID: ownerID.(string),
		}
	}

	// Process team IDs if set
	teamIDs := d.Get("team_ids").(*schema.Set).List()
	for _, teamID := range teamIDs {
		createRequest.Teams = append(createRequest.Teams, components.CreateServiceTeam{
			ID: teamID.(string),
		})
	}

	// Process external resources
	externalResources := d.Get("external_resources").(*schema.Set).List()
	for _, currentER := range externalResources {
		er := currentER.(map[string]interface{})
		remoteID := er["remote_id"].(string)
		connectionType := er["connection_type"].(string)
		createRequest.ExternalResources = append(createRequest.ExternalResources, components.CreateServiceExternalResource{
			RemoteID:       remoteID,
			ConnectionType: &connectionType,
		})
	}

	// Create the new service
	tflog.Debug(ctx, fmt.Sprintf("Create service: %s", name), map[string]interface{}{
		"name": name,
	})
	serviceResponse, err := client.Sdk.CatalogEntries.CreateService(ctx, createRequest)
	if err != nil {
		return diag.Errorf("Error creating service %s: %v", name, err)
	}

	// Set the new service's ID in state
	d.SetId(*serviceResponse.ID)

	// Update state with the latest information from the API
	return readResourceFireHydrantService(ctx, d, m)
}

func updateResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Construct the update request
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	alertOnAdd := d.Get("alert_on_add").(bool)
	autoAddRespondingTeam := d.Get("auto_add_responding_team").(bool)
	serviceTier := d.Get("service_tier").(int)
	labels := convertStringMap(d.Get("labels").(map[string]interface{}))

		// Process any optional attributes and add to the update request if necessary

	updateRequest := components.UpdateService{
		Name:                  &name,
		Description:           &description,
		Labels:                labels,
		ServiceTier:           (*components.UpdateServiceServiceTier)(&serviceTier),
		AlertOnAdd:            &alertOnAdd,
		AutoAddRespondingTeam: &autoAddRespondingTeam,
	}

	// Process owner - set or remove
	ownerID, ownerIDSet := d.GetOk("owner_id")
	if ownerIDSet && ownerID.(string) != "" {
		updateRequest.Owner = &components.UpdateServiceOwner{
			ID: ownerID.(string),
		}
	} else {
		removeOwner := true
		updateRequest.RemoveOwner = &removeOwner
	}

	// Process links
	links := d.Get("links").(*schema.Set).List()
	for _, currentLink := range links {
		link := currentLink.(map[string]interface{})
		hrefURL := link["href_url"].(string)
		linkName := link["name"].(string)
		updateRequest.Links = append(updateRequest.Links, components.UpdateServiceLink{
			HrefURL: hrefURL,
			Name:    linkName,
		})
	}

	// Process team IDs
	teamIDs := d.Get("team_ids").(*schema.Set).List()
	for _, teamID := range teamIDs {
		updateRequest.Teams = append(updateRequest.Teams, components.UpdateServiceTeam{
			ID: teamID.(string),
		})
	}
	// This will force the update request to replace the teams with the ones we send
	removeRemainingTeams := true
	updateRequest.RemoveRemainingTeams = &removeRemainingTeams

	// Process external resources
	externalResources := d.Get("external_resources").(*schema.Set).List()
	for _, currentER := range externalResources {
		er := currentER.(map[string]interface{})
		remoteID := er["remote_id"].(string)
		connectionType := er["connection_type"].(string)
		updateRequest.ExternalResources = append(updateRequest.ExternalResources, components.UpdateServiceExternalResource{
			RemoteID:       remoteID,
			ConnectionType: &connectionType,
		})
	}
	removeRemainingExternalResources := true
	updateRequest.RemoveRemainingExternalResources = &removeRemainingExternalResources

	// Update the service
	tflog.Debug(ctx, fmt.Sprintf("Update service: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})
	_, err := client.Sdk.CatalogEntries.UpdateService(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating service %s: %v", d.Id(), err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantService(ctx, d, m)
}

func deleteResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Delete the service
	serviceID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Delete service: %s", serviceID), map[string]interface{}{
		"id": serviceID,
	})
	err := client.Sdk.CatalogEntries.DeleteService(ctx, serviceID)
	if err != nil {
		if sdkErr, ok := err.(*sdkerrors.SDKError); ok && sdkErr.StatusCode == 404 {
			return nil
		}
		return diag.Errorf("Error deleting service %s: %v", serviceID, err)
	}

	return diag.Diagnostics{}
}
