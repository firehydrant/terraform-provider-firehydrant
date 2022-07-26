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

func resourceServiceDependency() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantServiceDependency,
		UpdateContext: updateResourceFireHydrantServiceDependency,
		ReadContext:   readResourceFireHydrantServiceDependency,
		DeleteContext: deleteResourceFireHydrantServiceDependency,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"connected_service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// Optional
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func readResourceFireHydrantServiceDependency(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the service dependency
	serviceDependencyID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read service dependency: %s", serviceDependencyID), map[string]interface{}{
		"id": serviceDependencyID,
	})
	serviceDependencyResponse, err := firehydrantAPIClient.ServiceDependencies().Get(ctx, serviceDependencyID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			tflog.Debug(ctx, fmt.Sprintf("Service dependency %s no longer exists", serviceDependencyID), map[string]interface{}{
				"id": serviceDependencyID,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading service dependency %s: %v", serviceDependencyID, err)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"connected_service_id": serviceDependencyResponse.ConnectedService.ID,
		"service_id":           serviceDependencyResponse.Service.ID,
		"notes":                serviceDependencyResponse.Notes,
	}

	// Set the resource attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for service dependency %s: %v", key, serviceDependencyID, err)
		}
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantServiceDependency(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get attributes from config and construct the create request
	createRequest := firehydrant.CreateServiceDependencyRequest{
		ConnectedServiceID: d.Get("connected_service_id").(string),
		ServiceID:          d.Get("service_id").(string),
		Notes:              d.Get("notes").(string),
	}

	// Create the new service dependency
	tflog.Debug(ctx, fmt.Sprintf("Create service dependency: %s:%s", createRequest.ServiceID, createRequest.ConnectedServiceID), map[string]interface{}{
		"connected_service_id": createRequest.ConnectedServiceID,
		"service_id":           createRequest.ServiceID,
	})
	serviceDependencyResponse, err := firehydrantAPIClient.ServiceDependencies().Create(ctx, createRequest)
	if err != nil {
		return diag.Errorf("Error creating service dependency %s:%s: %v", createRequest.ServiceID, createRequest.ConnectedServiceID, err)
	}

	// Set the new service dependency's ID in state
	d.SetId(serviceDependencyResponse.ID)

	// Update state with the latest information from the API
	return readResourceFireHydrantServiceDependency(ctx, d, m)
}

func updateResourceFireHydrantServiceDependency(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the update request
	updateRequest := firehydrant.UpdateServiceDependencyRequest{
		Notes: d.Get("notes").(string),
	}

	// Update the service dependency
	tflog.Debug(ctx, fmt.Sprintf("Update service dependency: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})
	_, err := firehydrantAPIClient.ServiceDependencies().Update(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating service dependency %s: %v", d.Id(), err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantServiceDependency(ctx, d, m)
}

func deleteResourceFireHydrantServiceDependency(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the service dependency
	serviceDependencyID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Delete service dependency: %s", serviceDependencyID), map[string]interface{}{
		"id": serviceDependencyID,
	})
	err := firehydrantAPIClient.ServiceDependencies().Delete(ctx, serviceDependencyID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			return nil
		}
		return diag.Errorf("Error deleting service dependency %s: %v", serviceDependencyID, err)
	}

	return diag.Diagnostics{}
}
