package provider

import (
	"context"
	"strconv"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePriority() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantPriority,
		UpdateContext: updateResourceFireHydrantPriority,
		ReadContext:   readResourceFireHydrantPriority,
		DeleteContext: deleteResourceFireHydrantPriority,
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
			"default": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func readResourceFireHydrantPriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	r, err := ac.GetPriority(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var ds diag.Diagnostics
	svc := map[string]string{
		"slug":        r.Slug,
		"description": r.Description,
		"default":     strconv.FormatBool(r.Default),
	}

	for key, val := range svc {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return ds
}

func createResourceFireHydrantPriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	slug, description, defaultV := d.Get("slug").(string), d.Get("description").(string), d.Get("default").(bool)

	r := firehydrant.CreatePriorityRequest{
		Slug:        slug,
		Description: description,
		Default:     defaultV,
	}

	resource, err := ac.CreatePriority(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.Slug)

	if err := d.Set("description", resource.Description); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func updateResourceFireHydrantPriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	description := d.Get("description").(string)
	defaultV := d.Get("default").(bool)
	id := d.Id()
	r := firehydrant.UpdatePriorityRequest{
		Slug:        id,
		Description: description,
		Default:     defaultV,
	}

	_, err := ac.UpdatePriority(ctx, id, r)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func deleteResourceFireHydrantPriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	priorityID := d.Id()

	err := ac.DeletePriority(ctx, priorityID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}
