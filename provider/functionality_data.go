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
			"functionality_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataFireHydrantFunctionality(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	id := d.Get("functionality_id").(string)

	r, err := ac.GetFunctionality(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	var ds diag.Diagnostics
	env := map[string]string{
		"name":        r.Name,
		"description": r.Description,
	}

	for key, val := range env {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(r.ID)

	return ds
}
