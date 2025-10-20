package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSignalRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantSignalRule,
		Schema: map[string]*schema.Schema{
			// Required
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expression": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"incident_type_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notification_priority_override": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_incident_condition_when": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deduplication_expiry": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Target fields for additional information
			"target_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_team_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_is_pageable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataFireHydrantSignalRule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)
	var rule *components.SignalsAPIRuleEntity

	// Get the signal rule name
	name := d.Get("name").(string)
	teamID := d.Get("team_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Fetch signal rule: %s", name), map[string]interface{}{
		"name":    name,
		"team_id": teamID,
	})

	ruleResponse, err := client.Sdk.Signals.ListTeamSignalRules(ctx, teamID, &name, nil, nil)
	if err != nil {
		return diag.Errorf("Error fetching signal rule '%s': %v", name, err)
	}

	rules := ruleResponse.GetData()
	if len(rules) == 0 {
		return diag.Errorf("Did not find signal rule matching '%s'", name)
	}
	if len(rules) > 1 {
		for _, r := range rules {
			if r.GetName() != nil && *r.GetName() == name {
				// we do not allow multiple signal rules with the same name so we can return the first match
				rule = &r
				break
			}
		}
		if rule == nil {
			return diag.Errorf("Did not find signal rule matching '%s'", name)
		}
	} else {
		rule = &rules[0]
	}

	attributes := map[string]interface{}{
		"id":          *rule.GetID(),
		"name":        *rule.GetName(),
		"expression":  *rule.GetExpression(),
		"target_type": *rule.GetTarget().GetType(),
		"target_id":   *rule.GetTarget().GetID(),
	}

	// Handle target additional fields
	if target := rule.GetTarget(); target != nil {
		if target.GetName() != nil {
			attributes["target_name"] = *target.GetName()
		}
		if target.GetTeamID() != nil {
			attributes["target_team_id"] = *target.GetTeamID()
		}
		if target.GetIsPageable() != nil {
			attributes["target_is_pageable"] = *target.GetIsPageable()
		}
	}

	if incidentType := rule.GetIncidentType(); incidentType != nil && incidentType.GetID() != nil {
		attributes["incident_type_id"] = *incidentType.GetID()
	}

	if priority := rule.GetNotificationPriorityOverride(); priority != nil {
		attributes["notification_priority_override"] = string(*priority)
	}

	if condition := rule.GetCreateIncidentConditionWhen(); condition != nil {
		attributes["create_incident_condition_when"] = string(*condition)
	}

	if expiry := rule.GetDeduplicationExpiry(); expiry != nil {
		attributes["deduplication_expiry"] = *expiry
	}

	for key, val := range attributes {
		if err := d.Set(key, val); err != nil {
			return diag.Errorf("Error setting %s for signal rule %s: %v", key, name, err)
		}
	}

	// Set the signal rule's ID in state
	d.SetId(*rule.GetID())

	return diag.Diagnostics{}
}
