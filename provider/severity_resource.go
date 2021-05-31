package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSeverity() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantSeverity,
		UpdateContext: updateResourceFireHydrantSeverity,
		ReadContext:   readResourceFireHydrantSeverity,
		DeleteContext: deleteResourceFireHydrantSeverity,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func readResourceFireHydrantSeverity(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	r, err := ac.GetSeverity(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	var ds diag.Diagnostics
	svc := map[string]string{
		"slug":        r.Slug,
		"description": r.Description,
	}

	for key, val := range svc {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return ds
}

func createResourceFireHydrantSeverity(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	slug, description := d.Get("slug").(string), d.Get("description").(string)

	r := firehydrant.CreateSeverityRequest{
		Slug:        slug,
		Description: description,
	}

	resource, err := ac.CreateSeverity(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.Slug)

	if err := d.Set("description", resource.Description); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func updateResourceFireHydrantSeverity(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	description := d.Get("description").(string)
	id := d.Id()
	r := firehydrant.UpdateSeverityRequest{
		Slug:        id,
		Description: description,
	}

	_, err := ac.UpdateSeverity(ctx, id, r)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func deleteResourceFireHydrantSeverity(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	severityID := d.Id()

	err := ac.DeleteSeverity(ctx, severityID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}
