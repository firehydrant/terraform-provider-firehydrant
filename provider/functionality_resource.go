package provider

import (
	"context"
	"log"

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
			"services": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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

	svcs := make([]interface{}, len(r.Services))
	for index, s := range r.Services {
		log.Println("[DEBUG] (Read) Service Name: " + s.Name)
		svcs[index] = map[string]interface{}{
			"id":   s.ID,
			"name": s.Name,
		}
	}
	d.Set("services", svcs)

	return ds
}

func createResourceFireHydrantFunctionality(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	name, description := d.Get("name").(string), d.Get("description").(string)

	r := firehydrant.CreateFunctionalityRequest{
		Name:        name,
		Description: description,
		Services:    []firehydrant.FunctionalityService{},
	}

	services := d.Get("services").([]interface{})
	for _, svc := range services {
		data := svc.(map[string]interface{})
		r.Services = append(r.Services, firehydrant.FunctionalityService{
			ID: data["id"].(string),
		})
	}

	resource, err := ac.CreateFunctionality(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.ID)
	d.Set("name", resource.Name)
	d.Set("description", resource.Description)

	svcs := make([]interface{}, len(resource.Services))
	for index, s := range resource.Services {
		svcs[index] = map[string]interface{}{
			"id":   s.ID,
			"name": s.Name,
		}
		log.Println("[DEBUG] (Create) Service Name: " + s.Name)
	}
	d.Set("services", svcs)

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

	services := d.Get("services").([]interface{})
	for _, svc := range services {
		data := svc.(map[string]interface{})
		r.Services = append(r.Services, firehydrant.FunctionalityService{
			ID: data["id"].(string),
		})
	}

	functionality, err := ac.UpdateFunctionality(ctx, id, r)
	if err != nil {
		return diag.FromErr(err)
	}

	svcs := make([]interface{}, len(functionality.Services))
	for index, s := range functionality.Services {
		svcs[index] = map[string]interface{}{
			"id":   s.ID,
			"name": s.Name,
		}
		log.Println("[DEBUG] (Update) Service Name: " + s.Name)
	}
	d.Set("services", svcs)

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
