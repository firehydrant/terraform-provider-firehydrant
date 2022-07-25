package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Singular services data source
func dataSourceIncidentRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantIncidentRole,
		Schema: map[string]*schema.Schema{
			// Required
			"id": {
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
			"summary": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataFireHydrantIncidentRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the incident role
	incidentRoleID := d.Get("id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read incident role: %s", incidentRoleID), map[string]interface{}{
		"id": incidentRoleID,
	})
	incidentRoleResponse, err := firehydrantAPIClient.IncidentRoles().Get(ctx, incidentRoleID)
	if err != nil {
		return diag.Errorf("Error reading incident role %s: %v", incidentRoleID, err)
	}
	// Currently, the incident role API will still return deleted/archived incident roles instead of returning
	// 404. So, to check for incident roles that are deleted, we have to check for incident roles that have
	// a DiscardedAt timestamp
	if !incidentRoleResponse.DiscardedAt.IsZero() {
		return diag.Errorf("Error reading incident role %s: Incident role has been archived", incidentRoleID)
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

	// Set the incident role's ID in state
	d.SetId(incidentRoleResponse.ID)

	return diag.Diagnostics{}
}
