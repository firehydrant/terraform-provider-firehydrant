package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Singular services data source
func dataSourceService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantService,
		Schema: map[string]*schema.Schema{
			"id": {
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
			"service_tier": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantServices,
		Schema: map[string]*schema.Schema{
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_tier": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	serviceID := d.Get("id").(string)

	r, err := ac.Services().Get(ctx, serviceID)
	if err != nil {
		return diag.FromErr(err)
	}

	var ds diag.Diagnostics
	svc := map[string]interface{}{
		"name":         r.Name,
		"description":  r.Description,
		"service_tier": r.ServiceTier,
	}

	for key, val := range svc {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(r.ID)

	return ds
}

// Multiple services data source
func dataFireHydrantServices(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)

	query := d.Get("query").(string)
	labels := d.Get("labels").(map[string]interface{})

	ls := firehydrant.LabelsSelector{}
	for k, v := range labels {
		ls[k] = v.(string)
	}

	r, err := ac.Services().List(ctx, &firehydrant.ServiceQuery{
		Query:          query,
		LabelsSelector: ls,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	services := make([]interface{}, 0)

	for _, svc := range r.Services {
		values := map[string]interface{}{
			"id":           svc.ID,
			"name":         svc.Name,
			"description":  svc.Description,
			"service_tier": svc.ServiceTier,
		}
		services = append(services, values)
	}

	if err := d.Set("services", services); err != nil {
		return diag.FromErr(err)
	}

	var ds diag.Diagnostics
	d.SetId("does-not-matter")

	return ds
}
