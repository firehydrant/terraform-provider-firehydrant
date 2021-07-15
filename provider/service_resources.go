package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceService() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantService,
		UpdateContext: updateResourceFireHydrantService,
		ReadContext:   readResourceFireHydrantService,
		DeleteContext: deleteResourceFireHydrantService,
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
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"service_tier": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func readResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	serviceID := d.Id()

	r, err := ac.Services().Get(ctx, serviceID)
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

	if err := d.Set("labels", r.Labels); err != nil {
		return diag.FromErr(err)
	}

	return ds
}

func createResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	labels := convertStringMap(d.Get("labels").(map[string]interface{}))

	r := firehydrant.CreateServiceRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		ServiceTier: d.Get("service_tier").(int),
		Labels:      labels,
	}

	newService, err := ac.Services().Create(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newService.ID)

	attributes := map[string]interface{}{
		"name":        newService.Name,
		"description": newService.Description,
		"labels":      newService.Labels,
		"service_tier": newService.ServiceTier,
	}

	if err := setAttributesFromMap(d, attributes); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func updateResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)

	r := firehydrant.UpdateServiceRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		ServiceTier: d.Get("service_tier").(int),
		Labels:      convertStringMap(d.Get("labels").(map[string]interface{})),
	}

	_, err := ac.Services().Update(ctx, d.Id(), r)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func deleteResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	serviceID := d.Id()

	err := ac.Services().Delete(ctx, serviceID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}


