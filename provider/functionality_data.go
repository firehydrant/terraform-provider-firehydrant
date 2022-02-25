package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

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
	r, err := firehydrantAPIClient.GetFunctionality(ctx, functionalityID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set values in state
	env := map[string]string{
		"name":        r.Name,
		"description": r.Description,
	}

	for key, val := range env {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	serviceIDs := make([]string, 0)
	for _, service := range r.Services {
		serviceIDs = append(serviceIDs, service.ID)
	}
	if err := d.Set("service_ids", serviceIDs); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(r.ID)

	return diag.Diagnostics{}
}
