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
			"team_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func readResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the service
	serviceID := d.Id()
	r, err := firehydrantAPIClient.Services().Get(ctx, serviceID)
	if err != nil {
		return diag.FromErr(err)
	}

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

	var teamIDs []interface{}
	for _, team := range r.Teams {
		teamIDs = append(teamIDs, team.ID)
	}
	if err := d.Set("team_ids", teamIDs); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("labels", r.Labels); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

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

	teamIDs := d.Get("team_ids").(*schema.Set).List()
	for _, teamID := range teamIDs {
		r.Teams = append(r.Teams, firehydrant.ServiceTeam{ID: teamID.(string)})
	}

	// Create the new service
	newService, err := firehydrantAPIClient.Services().Create(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newService.ID)

	return readResourceFireHydrantService(ctx, d, m)
}

func updateResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

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

	teamIDs := d.Get("team_ids").(*schema.Set).List()
	for _, teamID := range teamIDs {
		r.Teams = append(r.Teams, firehydrant.ServiceTeam{ID: teamID.(string)})
	}
	// This will force the update request to replace the teams with the ones we send
	r.RemoveRemainingTeams = true

	// Update the service
	_, err := firehydrantAPIClient.Services().Update(ctx, d.Id(), r)
	if err != nil {
		return diag.FromErr(err)
	}

	return readResourceFireHydrantService(ctx, d, m)
}

func deleteResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the service
	serviceID := d.Id()
	err := firehydrantAPIClient.Services().Delete(ctx, serviceID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
