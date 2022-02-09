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
			"alert_on_add": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_id": {
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
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceService(),
			},
		},
	}
}

func dataFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)

	// Get the service
	serviceID := d.Get("id").(string)
	r, err := ac.Services().Get(ctx, serviceID)
	if err != nil {
		return diag.FromErr(err)
	}

	var ds diag.Diagnostics

	svc := map[string]interface{}{
		"alert_on_add": r.AlertOnAdd,
		"description":  r.Description,
		"name":         r.Name,
		"service_tier": r.ServiceTier,
	}

	// Process any attributes that could be nil
	if r.Owner != nil {
		svc["owner_id"] = r.Owner.ID
	}

	// Set the data source attributes to the values we got from the API
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

	// Get the services
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

	// Set the data source attributes to the values we got from the API
	for _, svc := range r.Services {
		values := map[string]interface{}{
			"id":           svc.ID,
			"alert_on_add": svc.AlertOnAdd,
			"description":  svc.Description,
			"name":         svc.Name,
			"service_tier": svc.ServiceTier,
		}

		// Process any attributes that could be nil
		if svc.Owner != nil {
			values["owner_id"] = svc.Owner.ID
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
