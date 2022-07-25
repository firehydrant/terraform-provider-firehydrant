package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSeverity() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantSeverity,
		Schema: map[string]*schema.Schema{
			// Required
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataFireHydrantSeverity(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the severity
	slug := d.Get("slug").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read severity: %s", slug), map[string]interface{}{
		"id": slug,
	})
	severityResponse, err := firehydrantAPIClient.Severities().Get(ctx, slug)
	if err != nil {
		return diag.Errorf("Error reading severity %s: %v", slug, err)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"slug":        severityResponse.Slug,
		"description": severityResponse.Description,
	}

	// Set the data source attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for severity %s: %v", key, slug, err)
		}
	}

	// Set the severity's ID in state
	d.SetId(severityResponse.Slug)

	return diag.Diagnostics{}
}
