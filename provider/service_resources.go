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
				Default:  5,
			},
			"functionalities": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"summary": {
							Type:     schema.TypeString,
							Required: false,
						},
						"id": {
							Type:     schema.TypeString,
							Required: false,
						},
					},
				},
			},
			"links": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: false,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"href_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"icon_url": {
							Type:     schema.TypeString,
							Required: false,
						},
					},
				},
			},
			"teams": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"alert_on_add": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
	svc := map[string]interface{}{
		"name":         r.Name,
		"description":  r.Description,
		"service_tier": r.ServiceTier,
		"alert_on_add": r.AlertOnAdd,
	}

	for key, val := range svc {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	functionalities := make([]interface{}, len(*r.Functionalities))
	for index, f := range *r.Functionalities {

		functionalities[index] = map[string]interface{}{
			"id":          f.ID,
			"name":        f.Name,
			"summary":     f.Summary,
			"description": f.Description,
		}
	}

	if err := d.Set("functionalities", functionalities); err != nil {
		return diag.FromErr(err)
	}

	links := make([]interface{}, len(*r.Links))
	for index, l := range *r.Links {

		functionalities[index] = map[string]interface{}{
			"id":       l.ID,
			"name":     l.Name,
			"href_url": l.HrefURL,
			"icon_url": l.IconURL,
		}
	}

	if err := d.Set("links", links); err != nil {
		return diag.FromErr(err)
	}

	teams := make([]interface{}, len(*r.Teams))
	for index, l := range *r.Teams {
		functionalities[index] = map[string]interface{}{
			"id": l.ID,
		}
	}

	if err := d.Set("teams", teams); err != nil {
		return diag.FromErr(err)
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
		AlertOnAdd:  d.Get("alert_on_add").(bool),
	}

	functionalities := d.Get("functionalities").([]interface{})
	for _, functionality := range functionalities {

		f := functionality.(map[string]interface{})
		*r.Functionalities = append(*r.Functionalities, firehydrant.ServiceFunctionalities{
			Summary:     f["summary"].(string),
			ID:          f["id"].(string),
			Name:        f["name"].(string),
			Description: f["description"].(string),
		})
	}

	links := d.Get("links").([]interface{})
	for _, link := range links {

		l := link.(map[string]interface{})
		*r.Links = append(*r.Links, firehydrant.ServiceLinks{
			ID:      l["id"].(string),
			Name:    l["name"].(string),
			HrefURL: l["href_url"].(string),
			IconURL: l["icon_url"].(string),
		})
	}

	teams := d.Get("teams").([]interface{})
	for _, team := range teams {

		t := team.(map[string]interface{})
		*r.Teams = append(*r.Teams, firehydrant.ServiceTeams{
			ID: t["id"].(string),
		})
	}

	newService, err := ac.Services().Create(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newService.ID)

	attributes := map[string]interface{}{
		"name":            newService.Name,
		"description":     newService.Description,
		"labels":          newService.Labels,
		"service_tier":    newService.ServiceTier,
		"functionalities": newService.Functionalities,
		"links":           newService.Links,
		"teams":           newService.Teams,
		"alert_on_add":    newService.AlertOnAdd,
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
		AlertOnAdd:  d.Get("alert_on_add").(bool),
	}

	functionalities := d.Get("functionalities").([]interface{})
	for _, functionality := range functionalities {
		f := functionality.(map[string]interface{})
		*r.Functionalities = append(*r.Functionalities, firehydrant.ServiceFunctionalities{
			Summary:     f["summary"].(string),
			ID:          f["id"].(string),
			Name:        f["name"].(string),
			Description: f["description"].(string),
		})
	}

	links := d.Get("links").([]interface{})
	for _, link := range links {
		l := link.(map[string]interface{})
		*r.Links = append(*r.Links, firehydrant.ServiceLinks{
			ID:      l["id"].(string),
			Name:    l["name"].(string),
			HrefURL: l["href_url"].(string),
			IconURL: l["icon_url"].(string),
		})
	}

	teams := d.Get("teams").([]interface{})
	for _, team := range teams {
		t := team.(map[string]interface{})
		*r.Teams = append(*r.Teams, firehydrant.ServiceTeams{
			ID: t["id"].(string),
		})
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
