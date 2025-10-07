package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIncidentType() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataIncidentType,
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
			"template": {
				Type:     schema.TypeList, // Using TypeList to simulate a map
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"customer_impact_summary": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"severity_slug": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"priority_slug": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"private_incident": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						// "labels": {
						// 	Type:     schema.TypeMap,
						// 	Optional: true,
						// },
						"tags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"runbook_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"team_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"impacts": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"impact_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"condition_id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func readDataIncidentType(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read incident type: %s", id), map[string]interface{}{
		"id": id,
	})

	response, err := client.Sdk.IncidentSettings.GetIncidentType(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	var templateSlice []map[string]interface{}

	if response.Template != nil {
		var description, customerImpactSummary, severity, priority string
		var privateIncident bool

		if response.Template.Description != nil {
			description = *response.Template.Description
		}
		if response.Template.CustomerImpactSummary != nil {
			customerImpactSummary = *response.Template.CustomerImpactSummary
		}
		if response.Template.Severity != nil {
			severity = *response.Template.Severity
		}
		if response.Template.Priority != nil {
			priority = *response.Template.Priority
		}
		if response.Template.PrivateIncident != nil {
			privateIncident = *response.Template.PrivateIncident
		}

		template := map[string]interface{}{
			"description":             description,
			"customer_impact_summary": customerImpactSummary,
			"severity_slug":           severity,
			"priority_slug":           priority,
			"private_incident":        privateIncident,
		}

		// labels is in the sdk as an empty struct, which seems... wrong.  I'm going to implement the rest of this without it
		// (because I can only hold so much complexity in my head), and then investigate this from the API side to see if
		// this is being generated correctly.

		var tags []interface{}
		for _, tag := range response.Template.TagList {
			tags = append(tags, tag)
		}
		template["tags"] = tags

		var runbookIDs []interface{}
		for _, r := range response.Template.RunbookIds {
			runbookIDs = append(runbookIDs, r)
		}
		template["runbook_ids"] = runbookIDs

		var teamIDs []interface{}
		for _, team := range response.Template.TeamIds {
			teamIDs = append(teamIDs, team)
		}
		template["team_ids"] = teamIDs

		var impacts []map[string]interface{}
		for _, im := range response.Template.Impacts {
			impact := map[string]interface{}{
				"impact_id":    im.ID,
				"condition_id": im.ConditionID,
			}
			impacts = append(impacts, impact)
		}
		template["impacts"] = impacts

		templateSlice = append(templateSlice, template)
	}
	var name, description string
	if response.Name != nil {
		name = *response.Name
	}
	if response.Description != nil {
		description = *response.Description
	}

	attributes := map[string]interface{}{
		"name":        name,
		"description": description,
		"template":    templateSlice,
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for incident_type %s: %v", key, id, err)
		}
	}

	if response.ID != nil {
		d.SetId(*response.ID)
	}

	return diag.Diagnostics{}
}
