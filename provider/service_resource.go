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
			// Required
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Optional
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
			"links": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required
						"href_url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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
	serviceResponse, err := firehydrantAPIClient.Services().Get(ctx, serviceID)
	if err != nil {
		_, isNotFoundError := err.(firehydrant.NotFound)
		if isNotFoundError {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// Set values in state
	attributes := map[string]interface{}{
		"name":         serviceResponse.Name,
		"alert_on_add": serviceResponse.AlertOnAdd,
		"description":  serviceResponse.Description,
		"labels":       serviceResponse.Labels,
		"service_tier": serviceResponse.ServiceTier,
	}

	// Process any attributes that could be nil
	links := make([]map[string]interface{}, len(serviceResponse.Links))
	for index, currentLink := range serviceResponse.Links {
		links[index] = map[string]interface{}{
			"href_url": currentLink.HrefURL,
			"name":     currentLink.Name,
		}
	}
	attributes["links"] = links

	var ownerID string
	if serviceResponse.Owner != nil {
		ownerID = serviceResponse.Owner.ID
	}
	attributes["owner_id"] = ownerID

	var teamIDs []interface{}
	for _, team := range serviceResponse.Teams {
		teamIDs = append(teamIDs, team.ID)
	}
	attributes["team_ids"] = teamIDs

	// Update the resource attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get attributes from config and construct the create request
	createRequest := firehydrant.CreateServiceRequest{
		Name:        d.Get("name").(string),
		AlertOnAdd:  d.Get("alert_on_add").(bool),
		Description: d.Get("description").(string),
		Labels:      convertStringMap(d.Get("labels").(map[string]interface{})),
		ServiceTier: d.Get("service_tier").(int),
	}

	// Process any optional attributes and add to the create request if necessary
	configLinks := d.Get("links")
	for _, currentLink := range configLinks.(*schema.Set).List() {
		link := currentLink.(map[string]interface{})
		createRequest.Links = append(createRequest.Links, firehydrant.ServiceLink{
			HrefURL: link["href_url"].(string),
			Name:    link["name"].(string),
		})
	}

	if ownerID, ok := d.GetOk("owner_id"); ok && ownerID.(string) != "" {
		createRequest.Owner = &firehydrant.ServiceTeam{ID: ownerID.(string)}
	}

	teamIDs := d.Get("team_ids").(*schema.Set).List()
	for _, teamID := range teamIDs {
		createRequest.Teams = append(createRequest.Teams, firehydrant.ServiceTeam{ID: teamID.(string)})
	}

	// Create the new service
	serviceResponse, err := firehydrantAPIClient.Services().Create(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the new service's ID in state
	d.SetId(serviceResponse.ID)

	// Update state with the latest information from the API
	return readResourceFireHydrantService(ctx, d, m)
}

func updateResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the update request
	updateRequest := firehydrant.UpdateServiceRequest{
		Name:        d.Get("name").(string),
		AlertOnAdd:  d.Get("alert_on_add").(bool),
		Description: d.Get("description").(string),
		Labels:      convertStringMap(d.Get("labels").(map[string]interface{})),
		ServiceTier: d.Get("service_tier").(int),
	}

	// Process any optional attributes and add to the update request if necessary
	ownerID, ownerIDSet := d.GetOk("owner_id")
	if ownerIDSet {
		updateRequest.Owner = &firehydrant.ServiceTeam{ID: ownerID.(string)}
	} else {
		updateRequest.RemoveOwner = true
	}

	links := d.Get("links")
	for _, currentLink := range links.(*schema.Set).List() {
		link := currentLink.(map[string]interface{})

		linkAttributes := firehydrant.ServiceLink{
			HrefURL: link["href_url"].(string),
			Name:    link["name"].(string),
		}

		updateRequest.Links = append(updateRequest.Links, linkAttributes)
	}

	teamIDs := d.Get("team_ids").(*schema.Set).List()
	for _, teamID := range teamIDs {
		updateRequest.Teams = append(updateRequest.Teams, firehydrant.ServiceTeam{ID: teamID.(string)})
	}
	// This will force the update request to replace the teams with the ones we send
	updateRequest.RemoveRemainingTeams = true

	// Update the service
	_, err := firehydrantAPIClient.Services().Update(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantService(ctx, d, m)
}

func deleteResourceFireHydrantService(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the service
	serviceID := d.Id()
	err := firehydrantAPIClient.Services().Delete(ctx, serviceID)
	if err != nil {
		_, isNotFoundError := err.(firehydrant.NotFound)
		if isNotFoundError {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
