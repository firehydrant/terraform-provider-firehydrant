package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRunbook() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantRunbook,
		UpdateContext: updateResourceFireHydrantRunbook,
		ReadContext:   readResourceFireHydrantRunbook,
		DeleteContext: deleteResourceFireHydrantRunbook,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
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

func readResourceFireHydrantRunbook(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	r, err := ac.Runbooks().Get(ctx, d.Id())
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
		"type":        r.Type,
	}

	for key, val := range svc {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return ds
}

func createResourceFireHydrantRunbook(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	name, description, typ := d.Get("name").(string), d.Get("description").(string), d.Get("type").(string)

	r := firehydrant.CreateRunbookRequest{
		Name:        name,
		Description: description,
		Type:        typ,
	}

	resource, err := ac.Runbooks().Create(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	fmt.Printf("[DEBUG] Runbook value: %#v\n", resource)

	d.SetId(resource.ID)
	d.Set("name", resource.Name)
	d.Set("type", resource.Type)
	d.Set("description", resource.Description)

	return diag.Diagnostics{}
}

func updateResourceFireHydrantRunbook(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	id := d.Id()

	r := firehydrant.UpdateRunbookRequest{
		Name:        name,
		Description: description,
	}

	_, err := ac.Runbooks().Update(ctx, id, r)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func deleteResourceFireHydrantRunbook(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)

	err := ac.Runbooks().Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}
