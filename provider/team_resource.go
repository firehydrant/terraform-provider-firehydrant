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
			"service_ids": {
				Type:          schema.TypeSet,
				ConflictsWith: []string{"services"},
				Optional:      true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"services": {
				Type:          schema.TypeList,
				ConflictsWith: []string{"service_ids"},
				Deprecated:    "Use service_ids instead. The services attribute will be removed in the future. See the CHANGELOG to learn more: https://github.com/firehydrant/terraform-provider-firehydrant/blob/v0.2.0/CHANGELOG.md",
				Optional:      true,
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
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the team
	r, err := firehydrantAPIClient.GetTeam(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Set values in state
	svc := map[string]string{
		"name":        r.Name,
		"description": r.Description,
	}
	for key, val := range svc {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	// TODO: refactor this once deprecated attribute is removed
	// Update service IDs in state
	_, servicesSet := d.GetOk("services")
	if servicesSet {
		// If the config is using the services attribute, update the services attribute
		// in state with the information we got from the API
		var services []interface{}
		for _, service := range r.Services {
			services = append(services, map[string]interface{}{
				"id":   service.ID,
				"name": service.Name,
			})
		}
		if err := d.Set("services", services); err != nil {
			return diag.FromErr(err)
		}
	} else {
		// Otherwise, default to the preferred service_ids attribute and update the
		// service_ids attribute in state with the information we got from the API
		serviceIDs := make([]string, 0)
		for _, service := range r.Services {
			serviceIDs = append(serviceIDs, service.ID)
		}
		if err := d.Set("service_ids", serviceIDs); err != nil {
			return diag.FromErr(err)
		}
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the create team request
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	r := firehydrant.CreateTeamRequest{
		Name:        name,
		Description: description,
		ServiceIDs:  []string{},
	}

	// TODO: refactor this once deprecated attribute is removed
	// Add service IDs to the create request
	services, servicesSet := d.GetOk("services")
	serviceIDs, serviceIDsSet := d.GetOk("service_ids")
	if servicesSet {
		// If the services attribute is set, use the service IDs from that attribute
		// to set the service IDs for the create team request
		for _, service := range services.([]interface{}) {
			serviceAttributes := service.(map[string]interface{})
			r.ServiceIDs = append(r.ServiceIDs, serviceAttributes["id"].(string))
		}
	} else if serviceIDsSet {
		// If the service_ids attribute is set, use the service IDs from that attribute
		// to set the service IDs for the create team request
		for _, serviceID := range serviceIDs.(*schema.Set).List() {
			r.ServiceIDs = append(r.ServiceIDs, serviceID.(string))
		}
	}
	// Otherwise, don't send any service IDs in the create team request,
	// which will create a team with no services

	// Create the new team
	teamResponse, err := firehydrantAPIClient.CreateTeam(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the new team's ID in state
	d.SetId(teamResponse.ID)

	// Update state with the latest information from the API
	return readResourceFireHydrantTeam(ctx, d, m)
}

func updateResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the update team request
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	r := firehydrant.UpdateTeamRequest{
		Name:        name,
		Description: description,
	}

	// TODO: refactor this once deprecated attribute is removed
	// Add service IDs to the update request
	services, servicesSet := d.GetOk("services")
	serviceIDs, serviceIDsSet := d.GetOk("service_ids")
	updatedServiceIDs := make([]string, 0)
	if servicesSet {
		// If the services attribute is set, use the service IDs from that attribute
		// to populate the list of service IDs for the update team request
		for _, service := range services.([]interface{}) {
			serviceAttributes := service.(map[string]interface{})
			updatedServiceIDs = append(updatedServiceIDs, serviceAttributes["id"].(string))
		}
	} else if serviceIDsSet {
		// If the service_ids attribute is set, use the service IDs from that attribute
		// to populate the list of for service IDs for the update team request
		for _, serviceID := range serviceIDs.(*schema.Set).List() {
			updatedServiceIDs = append(updatedServiceIDs, serviceID.(string))
		}
	}
	// Otherwise, neither attribute is set, so updatedServiceIDs remains empty,
	// which will allow us to remove services from a team if either attribute
	// has been removed from the config

	// Set the service IDs for the update team request
	r.ServiceIDs = updatedServiceIDs

	// Update the team
	teamID := d.Id()
	_, err := firehydrantAPIClient.UpdateTeam(ctx, teamID, r)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantTeam(ctx, d, m)
}

func deleteResourceFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the team
	teamID := d.Id()
	err := firehydrantAPIClient.DeleteTeam(ctx, teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
