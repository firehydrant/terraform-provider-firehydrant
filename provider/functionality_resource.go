package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFunctionality() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantFunctionality,
		UpdateContext: updateResourceFireHydrantFunctionality,
		ReadContext:   readResourceFireHydrantFunctionality,
		DeleteContext: deleteResourceFireHydrantFunctionality,
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

func readResourceFireHydrantFunctionality(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	r, err := ac.GetFunctionality(ctx, d.Id())
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

func createResourceFireHydrantFunctionality(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	name, description := d.Get("name").(string), d.Get("description").(string)

	r := firehydrant.CreateFunctionalityRequest{
		Name:        name,
		Description: description,
	}

	resource, err := ac.CreateFunctionality(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.ID)
	d.Set("name", resource.Name)
	d.Set("description", resource.Description)

	var ds diag.Diagnostics
	return ds
}

func updateResourceFireHydrantFunctionality(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	r := firehydrant.UpdateFunctionalityRequest{
		Name:        name,
		Description: description,
	}

	_, err := ac.UpdateFunctionality(ctx, id, r)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func deleteResourceFireHydrantFunctionality(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	FunctionalityID := d.Id()

	err := ac.DeleteFunctionality(ctx, FunctionalityID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}
