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

func readResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	serviceID := d.Id()

	r, err := ac.GetService(ctx, serviceID)
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

func createResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	serviceName := d.Get("name").(string)
	serviceDescription := d.Get("description").(string)

	r := firehydrant.CreateServiceRequest{
		Name:        serviceName,
		Description: serviceDescription,
	}

	newService, err := ac.CreateService(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newService.ID)
	d.Set("name", newService.Name)
	d.Set("description", newService.Description)

	var ds diag.Diagnostics
	return ds
}

func updateResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	serviceID := d.Id()
	serviceName := d.Get("name").(string)
	serviceDescription := d.Get("description").(string)

	r := firehydrant.UpdateServiceRequest{
		Name:        serviceName,
		Description: serviceDescription,
	}

	_, err := ac.UpdateService(ctx, serviceID, r)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func deleteResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	serviceID := d.Id()

	err := ac.DeleteService(ctx, serviceID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}
