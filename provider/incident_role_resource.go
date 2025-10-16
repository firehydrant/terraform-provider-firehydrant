package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
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
	client := m.(*firehydrant.APIClient)

	// Get the incident role
	incidentRoleID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read incident role: %s", incidentRoleID), map[string]interface{}{
		"id": incidentRoleID,
	})
	incidentRole, err := client.Sdk.IncidentSettings.GetIncidentRole(ctx, incidentRoleID)
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
	if incidentRole.GetDiscardedAt() != nil && !incidentRole.GetDiscardedAt().IsZero() {
		tflog.Debug(ctx, fmt.Sprintf("Incident role %s has been archived", incidentRoleID), map[string]interface{}{
			"id": incidentRoleID,
		})
		d.SetId("")
		return nil
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"name":    *incidentRole.GetName(),
		"summary": *incidentRole.GetSummary(),
	}

	// Handle optional description field
	if description := incidentRole.GetDescription(); description != nil {
		attributes["description"] = *description
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
	client := m.(*firehydrant.APIClient)

	// Get attributes from config and construct the create request
	createReq := components.CreateIncidentRole{
		Name:    d.Get("name").(string),
		Summary: d.Get("summary").(string),
	}

	// Handle optional description field
	if desc := d.Get("description").(string); desc != "" {
		createReq.Description = &desc
	}

	// Create the new incident role
	tflog.Debug(ctx, fmt.Sprintf("Create incident role: %s", createReq.Name), map[string]interface{}{
		"name": createReq.Name,
	})
	incidentRole, err := client.Sdk.IncidentSettings.CreateIncidentRole(ctx, createReq)
	if err != nil {
		return diag.Errorf("Error creating incident role %s: %v", createReq.Name, err)
	}

	// Set the new incident role's ID in state
	d.SetId(*incidentRole.GetID())

	// Update state with the latest information from the API
	return readResourceFireHydrantIncidentRole(ctx, d, m)
}

func updateResourceFireHydrantIncidentRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Construct the update request (all fields optional with pointers)
	updateReq := components.UpdateIncidentRole{}
	name := d.Get("name").(string)
	updateReq.Name = &name
	summary := d.Get("summary").(string)
	updateReq.Summary = &summary

	// Handle optional description field
	// TODO: The Go SDK uses omitempty, so we cannot send empty strings to clear fields.
	// This isn't a big deal since it's just the description field, but we should fix this in the future
	if desc := d.Get("description").(string); desc != "" {
		updateReq.Description = &desc
	}

	// Update the incident role
	tflog.Debug(ctx, fmt.Sprintf("Update incident role: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})
	_, err := client.Sdk.IncidentSettings.UpdateIncidentRole(ctx, d.Id(), updateReq)
	if err != nil {
		return diag.Errorf("Error updating incident role %s: %v", d.Id(), err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantIncidentRole(ctx, d, m)
}

func deleteResourceFireHydrantIncidentRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Delete the incident role
	incidentRoleID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Delete incident role: %s", incidentRoleID), map[string]interface{}{
		"id": incidentRoleID,
	})
	err := client.Sdk.IncidentSettings.DeleteIncidentRole(ctx, incidentRoleID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			return nil
		}
		return diag.Errorf("Error deleting incident role %s: %v", incidentRoleID, err)
	}

	return diag.Diagnostics{}
}
