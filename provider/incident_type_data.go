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
				ForceNew: true,
				MaxItems: 1,
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
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"runbook_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"team_ids": {
							Type:     schema.TypeSet,
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

	template := map[string]interface{}{
		"description":             *response.Template.Description,
		"customer_impact_summary": *response.Template.CustomerImpactSummary,
		"severity_slug":           *response.Template.Severity,
		"priority_slug":           *response.Template.Priority,
		"private_incident":        *response.Template.PrivateIncident,
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

	attributes := map[string]interface{}{
		"name":        *response.Name,
		"description": *response.Description,
		"template":    template,
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for incident_type %s: %v", key, id, err)
		}
	}

	return diag.Diagnostics{}
}
