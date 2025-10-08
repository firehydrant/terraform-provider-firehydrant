package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/firehydrant-go-sdk/models/sdkerrors"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIncidentType() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceIncidentType,
		ReadContext:   readResourceIncidentType,
		UpdateContext: updateResourceIncidentType,
		DeleteContext: deleteResourceIncidentType,
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
			"template": {
				Type:     schema.TypeList, // Using TypeList to simulate a map
				Required: true,
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

func createResourceIncidentType(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	description := d.Get("description").(string)
	templateDescription := d.Get("template.0.description").(string)
	cis := d.Get("template.0.customer_impact_summary").(string)
	severity_id := d.Get("template.0.severity_slug").(string)
	priority_id := d.Get("template.0.priority_slug").(string)
	//Seriously?!?  A pointer to a boolean?  The pointer takes up more space that the actual value.  Ugh.
	is_private := d.Get("template.0.private_incident").(bool)

	inputTags := d.Get("template.0.tags").([]interface{})
	tags := []string{}
	for _, tag := range inputTags {
		if v, ok := tag.(string); ok && v != "" {
			tags = append(tags, v)
		}
	}

	inputRunbooks := d.Get("template.0.runbook_ids").([]interface{})
	runbooks := []string{}
	for _, runbook := range inputRunbooks {
		if v, ok := runbook.(string); ok && v != "" {
			runbooks = append(runbooks, v)
		}
	}

	inputTeams := d.Get("template.0.team_ids").([]interface{})
	teams := []string{}
	for _, team := range inputTeams {
		if v, ok := team.(string); ok && v != "" {
			teams = append(teams, v)
		}
	}

	inputImpacts := d.Get("template.0.impacts").([]interface{})
	impacts := []components.CreateIncidentTypeImpact{}
	for _, impact := range inputImpacts {
		impactMap := impact.(map[string]interface{})
		impacts = append(impacts, components.CreateIncidentTypeImpact{
			ID:          impactMap["impact_id"].(string),
			ConditionID: impactMap["condition_id"].(string),
		})
	}

	request := components.CreateIncidentType{
		Name:        d.Get("name").(string),
		Description: &description,
		Template: components.CreateIncidentTypeTemplate{
			Description:           &templateDescription,
			CustomerImpactSummary: &cis,
			Severity:              &severity_id,
			Priority:              &priority_id,
			PrivateIncident:       &is_private,
			TagList:               tags,
			RunbookIds:            runbooks,
			TeamIds:               teams,
			Impacts:               impacts,
		},
	}

	tflog.Debug(ctx, "Create new Incident Type")
	response, err := client.Sdk.IncidentSettings.CreateIncidentType(ctx, request)
	if err != nil {
		return diag.Errorf("Error creating new Incident Type: %v", err)
	}

	d.SetId(*response.ID)

	return readResourceIncidentType(ctx, d, m)
}

func readResourceIncidentType(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		impacts = append(impacts, map[string]interface{}{
			"impact_id":    im.ID,
			"condition_id": im.ConditionID,
		})
	}
	template["impacts"] = impacts

	templateSlice := []map[string]interface{}{template}

	attributes := map[string]interface{}{
		"name":        *response.Name,
		"description": *response.Description,
		"template":    templateSlice,
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for incident_type %s: %v", key, id, err)
		}
	}

	d.SetId(*response.ID)

	return diag.Diagnostics{}
}

func updateResourceIncidentType(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	description := d.Get("description").(string)
	templateDescription := d.Get("template.0.description").(string)
	cis := d.Get("template.0.customer_impact_summary").(string)
	severity_id := d.Get("template.0.severity_slug").(string)
	priority_id := d.Get("template.0.priority_slug").(string)
	//Seriously?!?  A pointer to a boolean?  The pointer takes up more space that the actual value.  Ugh.
	is_private := d.Get("template.0.private_incident").(bool)

	inputTags := d.Get("template.0.tags").([]interface{})
	tags := []string{}
	for _, tag := range inputTags {
		if v, ok := tag.(string); ok && v != "" {
			tags = append(tags, v)
		}
	}

	inputRunbooks := d.Get("template.0.runbook_ids").([]interface{})
	runbooks := []string{}
	for _, runbook := range inputRunbooks {
		if v, ok := runbook.(string); ok && v != "" {
			runbooks = append(runbooks, v)
		}
	}

	inputTeams := d.Get("template.0.team_ids").([]interface{})
	teams := []string{}
	for _, team := range inputTeams {
		if v, ok := team.(string); ok && v != "" {
			teams = append(teams, v)
		}
	}

	inputImpacts := d.Get("template.0.impacts").([]interface{})
	impacts := []components.UpdateIncidentTypeImpact{}
	for _, impact := range inputImpacts {
		impactMap := impact.(map[string]interface{})
		impacts = append(impacts, components.UpdateIncidentTypeImpact{
			ID:          impactMap["impact_id"].(string),
			ConditionID: impactMap["condition_id"].(string),
		})
	}

	request := components.UpdateIncidentType{
		Name:        d.Get("name").(string),
		Description: &description,
		Template: components.UpdateIncidentTypeTemplate{
			Description:           &templateDescription,
			CustomerImpactSummary: &cis,
			Severity:              &severity_id,
			Priority:              &priority_id,
			PrivateIncident:       &is_private,
			TagList:               tags,
			RunbookIds:            runbooks,
			TeamIds:               teams,
			Impacts:               impacts,
		},
	}

	tflog.Debug(ctx, "Update Incident Type")
	_, err := client.Sdk.IncidentSettings.UpdateIncidentType(ctx, id, request)
	if err != nil {
		return diag.Errorf("Error updating Incident Type: %v", err)
	}

	return readResourceIncidentType(ctx, d, m)
}

func deleteResourceIncidentType(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Delete incident type: %s", id), map[string]interface{}{
		"ID": id,
	})
	err := client.Sdk.IncidentSettings.DeleteIncidentType(ctx, id)
	if err != nil {
		if err.(*sdkerrors.SDKError).StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting incident type %s: %v", id, err)
	}

	return diag.Diagnostics{}
}
