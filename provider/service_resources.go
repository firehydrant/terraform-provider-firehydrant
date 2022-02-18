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
			"alert_on_add": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_tier": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
		},
	}
}

func readResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)

	// Get the service
	serviceID := d.Id()
	r, err := ac.Services().Get(ctx, serviceID)
	if err != nil {
		return diag.FromErr(err)
	}

	var ds diag.Diagnostics

	svc := map[string]interface{}{
		"name":         r.Name,
		"alert_on_add": r.AlertOnAdd,
		"description":  r.Description,
		"service_tier": r.ServiceTier,
	}

	// Process any attributes that could be nil
	if r.Owner != nil {
		svc["owner_id"] = r.Owner.ID
	}

	// Update the resource attributes to the values we got from the API
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

	// Get attributes from config and construct the create request
	labels := convertStringMap(d.Get("labels").(map[string]interface{}))
	r := firehydrant.CreateServiceRequest{
		Name:        d.Get("name").(string),
		AlertOnAdd:  d.Get("alert_on_add").(bool),
		Description: d.Get("description").(string),
		Labels:      labels,
		ServiceTier: d.Get("service_tier").(int),
	}

	// Process any optional attributes and add to the create request if necessary
	if ownerID, ok := d.GetOk("owner_id"); ok && ownerID.(string) != "" {
		r.Owner = &firehydrant.ServiceTeam{ID: ownerID.(string)}
	}

	// Create the new service
	newService, err := ac.Services().Create(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newService.ID)

	// TODO: Replace this whole section with a call to readResource
	// Update resource attributes
	attributes := map[string]interface{}{
		"name":         newService.Name,
		"alert_on_add": newService.AlertOnAdd,
		"description":  newService.Description,
		"labels":       newService.Labels,
		"service_tier": newService.ServiceTier,
	}

	// Process any attributes that could be nil
	if newService.Owner != nil {
		attributes["owner_id"] = newService.Owner.ID
	}

	if err := setAttributesFromMap(d, attributes); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func updateResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)

	// Construct the update request
	r := firehydrant.UpdateServiceRequest{
		Name:        d.Get("name").(string),
		AlertOnAdd:  d.Get("alert_on_add").(bool),
		Description: d.Get("description").(string),
		Labels:      convertStringMap(d.Get("labels").(map[string]interface{})),
		ServiceTier: d.Get("service_tier").(int),
	}

	// Process any optional attributes and add to the create request if necessary
	// Only set ownerID if it has actually been changed
	if d.HasChange("owner_id") {
		ownerID, ownerIDSet := d.GetOk("owner_id")
		if ownerIDSet {
			r.Owner = &firehydrant.ServiceTeam{ID: ownerID.(string)}
		} else {
			r.RemoveOwner = true
		}
	}

	// Update the service
	_, err := ac.Services().Update(ctx, d.Id(), r)
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO: Add a call to readResource to update attribute from API

	return diag.Diagnostics{}
}

func deleteResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)

	// Delete the service
	serviceID := d.Id()
	err := ac.Services().Delete(ctx, serviceID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diag.Diagnostics{}
}
