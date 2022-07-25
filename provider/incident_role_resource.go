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

func resourceIncidentRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantIncidentRole,
		UpdateContext: updateResourceFireHydrantIncidentRole,
		ReadContext:   readResourceFireHydrantIncidentRole,
		DeleteContext: deleteResourceFireHydrantIncidentRole,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"summary": {
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

func readResourceFireHydrantIncidentRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the incident role
	incidentRoleID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read incident role: %s", incidentRoleID), map[string]interface{}{
		"id": incidentRoleID,
	})
	incidentRoleResponse, err := firehydrantAPIClient.IncidentRoles().Get(ctx, incidentRoleID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			tflog.Debug(ctx, fmt.Sprintf("Incident role %s no longer exists", incidentRoleID), map[string]interface{}{
				"id": incidentRoleID,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading incident role %s: %v", incidentRoleID, err)
	}
	// Currently, the incident role API will still return deleted/archived incident roles instead of returning
	// 404. So, to check for incident roles that are deleted, we have to check for incident roles that have
	// a DiscardedAt timestamp
	if !incidentRoleResponse.DiscardedAt.IsZero() {
		tflog.Debug(ctx, fmt.Sprintf("Incident role %s has been archived", incidentRoleID), map[string]interface{}{
			"id": incidentRoleID,
		})
		d.SetId("")
		return nil
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"description": incidentRoleResponse.Description,
		"name":        incidentRoleResponse.Name,
		"summary":     incidentRoleResponse.Summary,
	}

	// Set the resource attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for incident role %s: %v", key, incidentRoleID, err)
		}
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantIncidentRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get attributes from config and construct the create request
	createRequest := firehydrant.CreateIncidentRoleRequest{
		Name:        d.Get("name").(string),
		Summary:     d.Get("summary").(string),
		Description: d.Get("description").(string),
	}

	// Create the new incident role
	tflog.Debug(ctx, fmt.Sprintf("Create incident role: %s", createRequest.Name), map[string]interface{}{
		"name": createRequest.Name,
	})
	incidentRoleResponse, err := firehydrantAPIClient.IncidentRoles().Create(ctx, createRequest)
	if err != nil {
		return diag.Errorf("Error creating incident role %s: %v", createRequest.Name, err)
	}

	// Set the new incident role's ID in state
	d.SetId(incidentRoleResponse.ID)

	// Update state with the latest information from the API
	return readResourceFireHydrantIncidentRole(ctx, d, m)
}

func updateResourceFireHydrantIncidentRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the update request
	updateRequest := firehydrant.UpdateIncidentRoleRequest{
		Name:        d.Get("name").(string),
		Summary:     d.Get("summary").(string),
		Description: d.Get("description").(string),
	}

	// Update the incident role
	tflog.Debug(ctx, fmt.Sprintf("Update incident role: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})
	_, err := firehydrantAPIClient.IncidentRoles().Update(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating incident role %s: %v", d.Id(), err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantIncidentRole(ctx, d, m)
}

func deleteResourceFireHydrantIncidentRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the incident role
	incidentRoleID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Delete incident role: %s", incidentRoleID), map[string]interface{}{
		"id": incidentRoleID,
	})
	err := firehydrantAPIClient.IncidentRoles().Delete(ctx, incidentRoleID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			return nil
		}
		return diag.Errorf("Error deleting incident role %s: %v", incidentRoleID, err)
	}

	return diag.Diagnostics{}
}
