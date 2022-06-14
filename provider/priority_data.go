package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePriority() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantPriority,
		Schema: map[string]*schema.Schema{
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataFireHydrantPriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)
	slug := d.Get("slug").(string)

	serviceResponse, err := firehydrantAPIClient.GetPriority(ctx, slug)
	if err != nil {
		return diag.FromErr(err)
	}

	env := map[string]interface{}{
		"description": serviceResponse.Description,
		"default":     serviceResponse.Default,
	}

	for key, val := range env {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(serviceResponse.Slug)

	return diag.Diagnostics{}
}
