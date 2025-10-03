package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTeam() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantTeam,
		Schema: map[string]*schema.Schema{
			// Required
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"owned_service_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"memberships": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_incident_role_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schedule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataFireHydrantTeam(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Get the team
	id := d.Get("id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read team: %s", id), map[string]interface{}{
		"id": id,
	})
	teamResponse, err := client.Sdk.Teams.GetTeam(ctx, id, nil)
	if err != nil {
		return diag.Errorf("Error reading team %s: %v", id, err)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"id":          *teamResponse.ID,
		"name":        *teamResponse.Name,
		"description": *teamResponse.Description,
		"slug":        *teamResponse.Slug,
	}

	// Collect mapped service IDs
	serviceIDs := make([]string, 0)
	for _, service := range teamResponse.Services {
		if service.ID != nil {
			serviceIDs = append(serviceIDs, *service.ID)
		}
	}
	attributes["service_ids"] = serviceIDs

	// Collect mapped owned service IDs
	ownedServiceIDs := make([]string, 0)
	for _, service := range teamResponse.OwnedServices {
		if service.ID != nil {
			ownedServiceIDs = append(ownedServiceIDs, *service.ID)
		}
	}
	attributes["owned_service_ids"] = ownedServiceIDs

	// Process memberships
	memberships := make([]map[string]interface{}, 0)
	for _, currentMembership := range teamResponse.Memberships {
		membership := map[string]interface{}{}

		// Handle default incident role
		if currentMembership.DefaultIncidentRole != nil && currentMembership.DefaultIncidentRole.ID != nil {
			membership["default_incident_role_id"] = *currentMembership.DefaultIncidentRole.ID
		} else {
			membership["default_incident_role_id"] = ""
		}

		// Handle schedule
		if currentMembership.Schedule != nil && currentMembership.Schedule.ID != nil {
			membership["schedule_id"] = *currentMembership.Schedule.ID
		} else {
			membership["schedule_id"] = ""
		}

		// Handle user
		if currentMembership.User != nil && currentMembership.User.ID != nil {
			membership["user_id"] = *currentMembership.User.ID
		} else {
			membership["user_id"] = ""
		}

		memberships = append(memberships, membership)
	}
	attributes["memberships"] = memberships

	// Set the data source attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for team %s: %v", key, id, err)
		}
	}

	// Set the team's ID in state
	d.SetId(*teamResponse.ID)

	return diag.Diagnostics{}
}
