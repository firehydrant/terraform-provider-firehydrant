package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/firehydrant/firehydrant-go-sdk/models/sdkerrors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEnvironment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantEnvironment,
		Schema: map[string]*schema.Schema{
			// Required
			"environment_id": {
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
		},
	}
}

func dataFireHydrantEnvironment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Get the environment
	environmentID := d.Get("environment_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read environment: %s", environmentID), map[string]interface{}{
		"id": environmentID,
	})
	environmentResponse, err := client.Sdk.CatalogEntries.GetEnvironment(ctx, environmentID)
	if err != nil {
		if sdkErr, ok := err.(*sdkerrors.SDKError); ok && sdkErr.StatusCode == 404 {
			return diag.Errorf("Environment %s not found", environmentID)
		}
		return diag.Errorf("Error reading environment %s: %v", environmentID, err)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"description": *environmentResponse.Description,
		"name":        *environmentResponse.Name,
	}

	// Set the data source attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for environment %s: %v", key, environmentID, err)
		}
	}

	// Set the environment's ID in state
	d.SetId(*environmentResponse.ID)

	return diag.Diagnostics{}
}
