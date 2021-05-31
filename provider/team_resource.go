package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantTeam,
		UpdateContext: updateResourceFireHydrantTeam,
		ReadContext:   readResourceFireHydrantTeam,
		DeleteContext: deleteResourceFireHydrantTeam,
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

func readResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	r, err := ac.GetTeam(ctx, d.Id())
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
		svcs[index] = map[string]interface{}{
			"id":   s.ID,
			"name": s.Name,
		}
	}

	if err := d.Set("services", svcs); err != nil {
		return diag.FromErr(err)
	}

	return ds
}

func createResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	name, description := d.Get("name").(string), d.Get("description").(string)

	r := firehydrant.CreateTeamRequest{
		Name:        name,
		Description: description,
		ServiceIDs:  []string{},
	}

	services := d.Get("services").([]interface{})
	for _, svc := range services {
		data := svc.(map[string]interface{})
		r.ServiceIDs = append(r.ServiceIDs, data["id"].(string))
	}

	resource, err := ac.CreateTeam(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.ID)

	attributes := map[string]interface{}{
		"name":        resource.Name,
		"description": resource.Description,
	}

	if err := setAttributesFromMap(d, attributes); err != nil {
		return diag.FromErr(err)
	}

	svcs := make([]interface{}, len(resource.Services))
	for index, s := range resource.Services {
		svcs[index] = map[string]interface{}{
			"id":   s.ID,
			"name": s.Name,
		}
	}

	if err := d.Set("services", svcs); err != nil {
		return diag.FromErr(err)
	}

	var ds diag.Diagnostics
	return ds
}

func updateResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	r := firehydrant.UpdateTeamRequest{
		Name:        name,
		Description: description,
	}

	services := d.Get("services").([]interface{})
	for _, svc := range services {
		data := svc.(map[string]interface{})
		r.ServiceIDs = append(r.ServiceIDs, data["id"].(string))
	}

	functionality, err := ac.UpdateTeam(ctx, id, r)
	if err != nil {
		return diag.FromErr(err)
	}

	svcs := make([]interface{}, len(functionality.Services))
	for index, s := range functionality.Services {
		svcs[index] = map[string]interface{}{
			"id":   s.ID,
			"name": s.Name,
		}
	}

	if err := d.Set("services", svcs); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func deleteResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	teamID := d.Id()

	err := ac.DeleteTeam(ctx, teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}
