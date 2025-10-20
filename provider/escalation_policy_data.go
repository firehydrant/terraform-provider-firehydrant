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

func dataSourceEscalationPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantEscalationPolicy,
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
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"repetitions": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"step": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timeout": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"targets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"priorities": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"handoff_step": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"step_strategy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notification_priority_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repetitions": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"handoff_step": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_id": {
										Type:     schema.TypeString,
										Computed: true,
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

func dataFireHydrantEscalationPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)
	var policy *components.SignalsAPIEscalationPolicyEntity

	name := d.Get("name").(string)
	teamID := d.Get("team_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Fetch escalation policy: %s", name), map[string]interface{}{
		"name":    name,
		"team_id": teamID,
	})

	policyResponse, err := client.Sdk.Signals.ListTeamEscalationPolicies(ctx, teamID, &name, nil, nil)
	if err != nil {
		return diag.Errorf("Error fetching escalation policy '%s': %v", name, err)
	}

	policies := policyResponse.GetData()
	if len(policies) == 0 {
		return diag.Errorf("Did not find escalation policy matching '%s'", name)
	}
	if len(policies) > 1 {
		for _, p := range policies {
			if p.GetName() != nil && *p.GetName() == name {
				// we do not allow multiple escalation policies with the same name so we can return the first match
				policy = &p
				break
			}
		}
		if policy == nil {
			return diag.Errorf("Did not find escalation policy matching '%s'", name)
		}
	} else {
		policy = &policies[0]
	}

	attributes := map[string]interface{}{
		"id":          *policy.GetID(),
		"name":        *policy.GetName(),
		"description": *policy.GetDescription(),
		"default":     *policy.GetDefault(),
		"repetitions": *policy.GetRepetitions(),
	}

	// Set step strategy
	if stepStrategy := policy.GetStepStrategy(); stepStrategy != nil {
		attributes["step_strategy"] = *stepStrategy
	}

	// Set notification priority policies
	if priorityPolicies := policy.GetNotificationPriorityPolicies(); priorityPolicies != nil {
		var policies []map[string]interface{}
		for _, policy := range priorityPolicies {
			policyMap := map[string]interface{}{
				"priority": *policy.GetNotificationPriority(),
			}

			// Set repetitions if available
			if repetitions := policy.GetRepetitions(); repetitions != nil {
				policyMap["repetitions"] = *repetitions
			}

			// Set handoff step if available
			if handoffStep := policy.GetHandoffStep(); handoffStep != nil {
				if target := handoffStep.GetTarget(); target != nil {
					handoffStepMap := map[string]interface{}{
						"target_type": *target.GetType(),
						"target_id":   *target.GetID(),
					}
					policyMap["handoff_step"] = []map[string]interface{}{handoffStepMap}
				}
			}

			policies = append(policies, policyMap)
		}
		attributes["notification_priority_policies"] = policies
	}

	var steps []map[string]interface{}
	for _, step := range policy.GetSteps() {
		targets := []map[string]interface{}{}
		for _, target := range step.GetTargets() {
			targetMap := map[string]interface{}{
				"type": *target.GetType(),
				"id":   *target.GetID(),
			}
			targets = append(targets, targetMap)
		}

		stepMap := map[string]interface{}{
			"timeout": *step.GetTimeout(),
			"targets": targets,
		}

		// Add priorities if available
		if priorities := step.GetPriorities(); priorities != nil {
			stepMap["priorities"] = priorities
		}
		steps = append(steps, stepMap)
	}
	attributes["step"] = steps

	if handoffStep := policy.GetHandoffStep(); handoffStep != nil {
		if target := handoffStep.GetTarget(); target != nil {
			handoffStepMap := map[string]interface{}{
				"target_type": *target.GetType(),
				"target_id":   *target.GetID(),
			}
			attributes["handoff_step"] = []map[string]interface{}{handoffStepMap}
		}
	}

	for key, val := range attributes {
		if err := d.Set(key, val); err != nil {
			return diag.Errorf("Error setting %s for escalation policy %s: %v", key, name, err)
		}
	}

	d.SetId(*policy.GetID())

	return diag.Diagnostics{}
}
