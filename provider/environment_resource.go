package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Description:   "FireHydrant environments are used to tag incidents with where they are occurring.",
		CreateContext: createResourceFireHydrantEnvironment,
		UpdateContext: updateResourceFireHydrantEnvironment,
		ReadContext:   readResourceFireHydrantEnvironment,
		DeleteContext: deleteResourceFireHydrantEnvironment,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
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

func readResourceFireHydrantEnvironment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	r, err := ac.GetEnvironment(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	var ds diag.Diagnostics
	svc := map[string]string{
		"name":        r.Name,
		"description": r.Description,
	}

	for key, val := range svc {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return ds
}

func createResourceFireHydrantEnvironment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	name, description := d.Get("name").(string), d.Get("description").(string)

	r := firehydrant.CreateEnvironmentRequest{
		Name:        name,
		Description: description,
	}

	resource, err := ac.CreateEnvironment(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.ID)

	attributes := map[string]interface{}{
		"name":        resource.Name,
		"description": resource.Description,
	}
	if err := setAttributesFromMap(d, attributes); err != nil {
		return diag.FromErr(fmt.Errorf("could not set attributes: %w", err))
	}

	return diag.Diagnostics{}
}

func updateResourceFireHydrantEnvironment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	r := firehydrant.UpdateEnvironmentRequest{
		Name:        name,
		Description: description,
	}

	_, err := ac.UpdateEnvironment(ctx, id, r)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func deleteResourceFireHydrantEnvironment(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	EnvironmentID := d.Id()

	err := ac.DeleteEnvironment(ctx, EnvironmentID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}
