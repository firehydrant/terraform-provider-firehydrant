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
				Type:          schema.TypeSet,
				ConflictsWith: []string{"services"},
				Optional:      true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"services": {
				Type:          schema.TypeList,
				ConflictsWith: []string{"service_ids"},
				Deprecated:    "Use service_ids instead. The services attribute will be removed in the future. See the CHANGELOG to learn more: https://github.com/firehydrant/terraform-provider-firehydrant/blob/v0.2.0/CHANGELOG.md",
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
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

	// Set the resource attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for functionality %s: %v", key, functionalityID, err)
		}
	}

	// TODO: refactor this once deprecated attribute is removed
	// Update service IDs in state
	_, servicesSet := d.GetOk("services")
	if servicesSet {
		// If the config is using the services attribute, update the services attribute
		// in state with the information we got from the API
		var services []interface{}
		for _, service := range functionalityResponse.Services {
			services = append(services, map[string]interface{}{
				"id":   service.ID,
				"name": service.Name,
			})
		}
		if err := d.Set("services", services); err != nil {
			return diag.FromErr(err)
		}
	} else {
		// Otherwise, default to the preferred service_ids attribute and update the
		// service_ids attribute in state with the information we got from the API
		serviceIDs := make([]string, 0)
		for _, service := range functionalityResponse.Services {
			serviceIDs = append(serviceIDs, service.ID)
		}
		if err := d.Set("service_ids", serviceIDs); err != nil {
			return diag.FromErr(err)
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
	// TODO: refactor this once deprecated attribute is removed
	// Add service IDs to the create request
	services, servicesSet := d.GetOk("services")
	serviceIDs, serviceIDsSet := d.GetOk("service_ids")
	if servicesSet {
		// If the services attribute is set, use the service IDs from that attribute
		// to set the service IDs for the create functionality request
		for _, service := range services.([]interface{}) {
			serviceAttributes := service.(map[string]interface{})
			createRequest.Services = append(createRequest.Services, firehydrant.FunctionalityService{
				ID: serviceAttributes["id"].(string),
			})
		}
	} else if serviceIDsSet {
		// If the service_ids attribute is set, use the service IDs from that attribute
		// to set the service IDs for the create functionality request
		for _, serviceID := range serviceIDs.(*schema.Set).List() {
			createRequest.Services = append(createRequest.Services, firehydrant.FunctionalityService{
				ID: serviceID.(string),
			})
		}
	}
	// Otherwise, don't send any service IDs in the create functionality request,
	// which will create a functionality with no services

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
	// TODO: refactor this once deprecated attribute is removed
	// Add service IDs to the update request
	services, servicesSet := d.GetOk("services")
	serviceIDs, serviceIDsSet := d.GetOk("service_ids")
	updatedServices := make([]firehydrant.FunctionalityService, 0)
	if servicesSet {
		// If the services attribute is set, use the service IDs from that attribute
		// to populate the list of service IDs for the update functionality request
		for _, service := range services.([]interface{}) {
			serviceAttributes := service.(map[string]interface{})
			updatedServices = append(updatedServices, firehydrant.FunctionalityService{
				ID: serviceAttributes["id"].(string),
			})
		}
	} else if serviceIDsSet {
		// If the service_ids attribute is set, use the service IDs from that attribute
		// to populate the list of for service IDs for the update functionality request
		for _, serviceID := range serviceIDs.(*schema.Set).List() {
			updatedServices = append(updatedServices, firehydrant.FunctionalityService{
				ID: serviceID.(string),
			})
		}
	}
	// Otherwise, neither attribute is set, so updatedServiceIDs remains empty,
	// which will allow us to remove services from a functionality if either attribute
	// has been removed from the config

	// Set the service IDs for the update functionality request
	updateRequest.Services = updatedServices
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
