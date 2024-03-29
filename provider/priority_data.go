package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePriority() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantPriority,
		Schema: map[string]*schema.Schema{
			// Required
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataFireHydrantPriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the priority
	slug := d.Get("slug").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read priority: %s", slug), map[string]interface{}{
		"id": slug,
	})
	priorityResponse, err := firehydrantAPIClient.Priorities().Get(ctx, slug)
	if err != nil {
		return diag.Errorf("Error reading priority %s: %v", slug, err)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"description": priorityResponse.Description,
		"default":     priorityResponse.Default,
	}

	// Set the data source attributes to the values we got from the API
	for key, val := range attributes {
		if err := d.Set(key, val); err != nil {
			return diag.Errorf("Error setting %s for priority %s: %v", key, slug, err)
		}
	}

	// Set the priority's ID in state
	d.SetId(priorityResponse.Slug)

	return diag.Diagnostics{}
}
