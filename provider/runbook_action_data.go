package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Singular services data source
func dataSourceRunbookAction() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantRunbookAction,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"integration_slug": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataFireHydrantRunbookAction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	typ, slug, integrationSlug := d.Get("type").(string), d.Get("slug").(string), d.Get("integration_slug").(string)

	r, err := ac.RunbookActions().Get(ctx, typ, fmt.Sprintf("%s.%s", integrationSlug, slug))
	if err != nil {
		return diag.FromErr(err)
	}

	var ds diag.Diagnostics
	svc := map[string]string{
		"name":             r.Name,
		"slug":             r.Slug,
		"integration_slug": "bunk",
	}

	for key, val := range svc {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(r.ID)

	return ds
}
