package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

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
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataFireHydrantFunctionality(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the functionality
	functionalityID := d.Get("functionality_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read functionality: %s", functionalityID), map[string]interface{}{
		"id": functionalityID,
	})
	functionalityResponse, err := firehydrantAPIClient.Functionalities().Get(ctx, functionalityID)
	if err != nil {
		return diag.Errorf("Error reading functionality %s: %v", functionalityID, err)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"name":        functionalityResponse.Name,
		"description": functionalityResponse.Description,
		"labels":      functionalityResponse.Labels,
		"Owner": functionalityResponse.Owner,
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

	// Set the functionality's ID in state
	d.SetId(functionalityResponse.ID)

	return diag.Diagnostics{}
}
