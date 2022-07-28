package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

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
		},
	}
}

func readResourceFireHydrantFunctionality(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the functionality
	functionalityID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read functionality: %s", functionalityID), map[string]interface{}{
		"id": functionalityID,
	})
	functionalityResponse, err := firehydrantAPIClient.Functionalities().Get(ctx, functionalityID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			tflog.Debug(ctx, fmt.Sprintf("Functionality %s no longer exists", functionalityID), map[string]interface{}{
				"id": functionalityID,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading functionality %s: %v", functionalityID, err)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"name":        functionalityResponse.Name,
		"description": functionalityResponse.Description,
	}

	serviceIDs := make([]string, 0)
	for _, service := range functionalityResponse.Services {
		serviceIDs = append(serviceIDs, service.ID)
	}
	attributes["service_ids"] = serviceIDs

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
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get attributes from config and construct the create request
	createRequest := firehydrant.CreateFunctionalityRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Process any optional attributes and add to the create request if necessary
	serviceIDs := d.Get("service_ids")
	for _, serviceID := range serviceIDs.(*schema.Set).List() {
		createRequest.Services = append(createRequest.Services, firehydrant.FunctionalityService{
			ID: serviceID.(string),
		})
	}

	// Create the new functionality
	tflog.Debug(ctx, fmt.Sprintf("Create functionality: %s", createRequest.Name), map[string]interface{}{
		"name": createRequest.Name,
	})
	functionalityResponse, err := firehydrantAPIClient.Functionalities().Create(ctx, createRequest)
	if err != nil {
		return diag.Errorf("Error creating functionality %s: %v", createRequest.Name, err)
	}

	// Set the new functionality's ID in state
	d.SetId(functionalityResponse.ID)

	// Update state with the latest information from the API
	return readResourceFireHydrantFunctionality(ctx, d, m)
}

func updateResourceFireHydrantFunctionality(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the update request
	updateRequest := firehydrant.UpdateFunctionalityRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Process any optional attributes and add to the update request if necessary
	serviceIDs := d.Get("service_ids")
	for _, serviceID := range serviceIDs.(*schema.Set).List() {
		updateRequest.Services = append(updateRequest.Services, firehydrant.FunctionalityService{
			ID: serviceID.(string),
		})
	}
	// Otherwise, neither attribute is set, so updatedServiceIDs remains empty,
	// which will allow us to remove services from a functionality if either attribute
	// has been removed from the config

	// This will force the update request to replace the services with the ones we send
	updateRequest.RemoveRemainingServices = true

	// Update the functionality
	tflog.Debug(ctx, fmt.Sprintf("Update functionality: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})
	_, err := firehydrantAPIClient.Functionalities().Update(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating functionality %s: %v", d.Id(), err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantFunctionality(ctx, d, m)
}

func deleteResourceFireHydrantFunctionality(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the functionality
	functionalityID := d.Id()
	err := firehydrantAPIClient.Functionalities().Delete(ctx, functionalityID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
