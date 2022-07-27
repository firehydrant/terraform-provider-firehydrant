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

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Description:   "FireHydrant environments are used to tag incidents with where they are occurring.",
		CreateContext: createResourceFireHydrantEnvironment,
		UpdateContext: updateResourceFireHydrantEnvironment,
		ReadContext:   readResourceFireHydrantEnvironment,
		DeleteContext: deleteResourceFireHydrantEnvironment,
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
		},
	}
}

func readResourceFireHydrantEnvironment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the environment
	environmentID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read environment: %s", environmentID), map[string]interface{}{
		"id": environmentID,
	})
	environmentResponse, err := firehydrantAPIClient.Environments().Get(ctx, environmentID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			tflog.Debug(ctx, fmt.Sprintf("Environment %s no longer exists", environmentID), map[string]interface{}{
				"id": environmentID,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading environment %s: %v", environmentID, err)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"name":        environmentResponse.Name,
		"description": environmentResponse.Description,
	}

	// Set the resource attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for environment %s: %v", key, environmentID, err)
		}
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantEnvironment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get attributes from config and construct the create request
	createRequest := firehydrant.CreateEnvironmentRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Create the new environment
	tflog.Debug(ctx, fmt.Sprintf("Create environment: %s", createRequest.Name), map[string]interface{}{
		"name": createRequest.Name,
	})
	environmentResponse, err := firehydrantAPIClient.Environments().Create(ctx, createRequest)
	if err != nil {
		return diag.Errorf("Error creating environment %s: %v", createRequest.Name, err)
	}

	// Set the new environment's ID in state
	d.SetId(environmentResponse.ID)

	// Update state with the latest information from the API
	return readResourceFireHydrantEnvironment(ctx, d, m)
}

func updateResourceFireHydrantEnvironment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the update request
	updateRequest := firehydrant.UpdateEnvironmentRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Update the environment
	tflog.Debug(ctx, fmt.Sprintf("Update environment: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})
	_, err := firehydrantAPIClient.Environments().Update(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating environment %s: %v", d.Id(), err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantEnvironment(ctx, d, m)
}

func deleteResourceFireHydrantEnvironment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the environment
	environmentID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Delete environment: %s", environmentID), map[string]interface{}{
		"id": environmentID,
	})
	err := firehydrantAPIClient.Environments().Delete(ctx, environmentID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			return nil
		}
		return diag.Errorf("Error deleting environment %s: %v", environmentID, err)
	}

	return diag.Diagnostics{}
}
