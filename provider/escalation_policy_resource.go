package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceEscalationPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantEscalationPolicy,
		UpdateContext: updateResourceFireHydrantEscalationPolicy,
		ReadContext:   readResourceFireHydrantEscalationPolicy,
		DeleteContext: deleteResourceFireHydrantEscalationPolicy,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repetitions": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"step": {
				Type:     schema.TypeList, // or TypeSet if ordering is not important
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timeout": {
							Type:     schema.TypeString,
							Required: true,
						},
						"targets": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"handoff_step": {
				Type:     schema.TypeList,
				Optional: true, // or Required based on your API
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"step_strategy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The strategy for handling steps in the escalation policy. Can be 'static' or 'dynamic_by_priority'.",
			},
			"notification_priority_policies": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Priority-specific policies for dynamic escalation policies",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(firehydrant.NotificationPriorityHigh),
								string(firehydrant.NotificationPriorityMedium),
								string(firehydrant.NotificationPriorityLow),
							}, false),
						},
						"repetitions": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of repetitions for this priority level",
						},
						"handoff_step": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Handoff step for this priority level",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"target_id": {
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

func readResourceFireHydrantEscalationPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Get the escalation policy
	id := d.Id()
	teamID := d.Get("team_id").(string)
	tflog.Debug(ctx, "Read escalation policy", map[string]interface{}{
		"id":      id,
		"team_id": teamID,
	})

	escalationPolicy, err := client.Sdk.Signals.GetTeamEscalationPolicy(ctx, teamID, id)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			tflog.Debug(ctx, fmt.Sprintf("Escalation Policy %s no longer exists", id), map[string]interface{}{
				"id": id,
			})
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading escalation policy %s: %v", id, err)
	}

	// Update the resource data
	d.Set("name", *escalationPolicy.GetName())
	d.Set("description", *escalationPolicy.GetDescription())
	d.Set("default", *escalationPolicy.GetDefault())
	d.Set("repetitions", *escalationPolicy.GetRepetitions())

	// Set step strategy
	if stepStrategy := escalationPolicy.GetStepStrategy(); stepStrategy != nil {
		d.Set("step_strategy", *stepStrategy)
	}

	// Set notification priority policies
	var policies []map[string]interface{}
	if priorityPolicies := escalationPolicy.GetNotificationPriorityPolicies(); priorityPolicies != nil {
		tflog.Debug(ctx, fmt.Sprintf("Found %d notification priority policies", len(priorityPolicies)), map[string]interface{}{
			"count": len(priorityPolicies),
		})
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
	} else {
		tflog.Debug(ctx, "No notification priority policies found in API response", map[string]interface{}{})
	}
	// Always set notification_priority_policies, even if empty, to prevent "attribute not found" errors
	tflog.Debug(ctx, fmt.Sprintf("Setting notification_priority_policies with %d policies", len(policies)), map[string]interface{}{
		"policies": policies,
	})
	d.Set("notification_priority_policies", policies)

	// Set the steps
	var steps []map[string]interface{}
	for _, step := range escalationPolicy.GetSteps() {
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

		steps = append(steps, stepMap)
	}

	d.Set("step", steps)

	// Set the handoff step
	if handoffStep := escalationPolicy.GetHandoffStep(); handoffStep != nil {
		if target := handoffStep.GetTarget(); target != nil {
			handoffStepMap := map[string]interface{}{
				"target_type": *target.GetType(),
				"target_id":   *target.GetID(),
			}

			d.Set("handoff_step", []map[string]interface{}{handoffStepMap})
		}
	}

	return diag.Diagnostics{}
}

// creates an escalation policy for a team using the firehydrant api client
func createResourceFireHydrantEscalationPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	// Create the escalation policy
	teamID := d.Get("team_id").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	defaultVal := d.Get("default").(bool)
	repetitions := d.Get("repetitions").(int)

	createReq := components.CreateTeamEscalationPolicy{
		Name:        name,
		Description: &description,
		Default:     &defaultVal,
		Repetitions: &repetitions,
		Steps:       getStepsFromResourceData(d),
		HandoffStep: getHandoffStepFromResourceData(d),
	}

	// Handle step strategy
	if stepStrategy := d.Get("step_strategy").(string); stepStrategy != "" {
		createReq.StepStrategy = &stepStrategy
	}

	// Handle notification priority policies
	if priorityPolicies := d.Get("notification_priority_policies").([]interface{}); len(priorityPolicies) > 0 {
		createReq.PrioritizedSettings = getNotificationPriorityPoliciesFromResourceData(d)
	}

	tflog.Debug(ctx, "Creating escalation policy", map[string]interface{}{
		"team_id": teamID,
		"request": createReq,
	})

	escalationPolicy, err := client.Sdk.Signals.CreateTeamEscalationPolicy(ctx, teamID, createReq)
	if err != nil {
		return diag.Errorf("Error creating escalation policy %s: %v", d.Id(), err)
	}

	// Set the ID of the escalation policy
	d.SetId(*escalationPolicy.GetID())
	return readResourceFireHydrantEscalationPolicy(ctx, d, m)
}

func updateResourceFireHydrantEscalationPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	// Update the escalation policy
	teamID := d.Get("team_id").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	defaultVal := d.Get("default").(bool)
	repetitions := d.Get("repetitions").(int)

	updateReq := components.UpdateTeamEscalationPolicy{
		Name:        &name,
		Description: &description,
		Default:     &defaultVal,
		Repetitions: &repetitions,
		Steps:       getStepsFromResourceDataUpdateSDK(d),
		HandoffStep: getHandoffStepFromResourceDataUpdateSDK(d),
	}

	// Handle step strategy
	if stepStrategy := d.Get("step_strategy").(string); stepStrategy != "" {
		updateReq.StepStrategy = &stepStrategy
	}

	// Handle notification priority policies
	if priorityPolicies := d.Get("notification_priority_policies").([]interface{}); len(priorityPolicies) > 0 {
		updateReq.PrioritizedSettings = getNotificationPriorityPoliciesFromResourceDataUpdate(d)
	}

	tflog.Debug(ctx, "Updating escalation policy", map[string]interface{}{
		"team_id": teamID,
		"request": spew.Sdump(updateReq),
	})

	_, err := client.Sdk.Signals.UpdateTeamEscalationPolicy(ctx, teamID, d.Id(), updateReq)
	if err != nil {
		return diag.Errorf("Error updating escalation policy %s: %v", d.Id(), err)
	}

	return readResourceFireHydrantEscalationPolicy(ctx, d, m)
}

func deleteResourceFireHydrantEscalationPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	// Delete the escalation policy
	teamID := d.Get("team_id").(string)
	tflog.Debug(ctx, "Deleting escalation policy", map[string]interface{}{
		"team_id": teamID,
		"id":      d.Id(),
	})

	err := client.Sdk.Signals.DeleteTeamEscalationPolicy(ctx, teamID, d.Id())
	if err != nil {
		return diag.Errorf("Error deleting escalation policy %s: %v", d.Id(), err)
	}

	return nil
}

func getHandoffStepFromResourceData(d *schema.ResourceData) *components.CreateTeamEscalationPolicyHandoffStep {
	if v, ok := d.GetOk("handoff_step"); ok {
		handoffStepList := v.([]interface{})

		if len(handoffStepList) > 0 {
			handoffStepMap := handoffStepList[0].(map[string]interface{})

			handoffStep := &components.CreateTeamEscalationPolicyHandoffStep{
				TargetType: components.CreateTeamEscalationPolicyTargetType(handoffStepMap["target_type"].(string)),
				TargetID:   handoffStepMap["target_id"].(string),
			}

			return handoffStep
		}
	}

	return nil
}

func getHandoffStepFromResourceDataUpdateSDK(d *schema.ResourceData) *components.UpdateTeamEscalationPolicyHandoffStep {
	if v, ok := d.GetOk("handoff_step"); ok {
		handoffStepList := v.([]interface{})

		if len(handoffStepList) > 0 {
			handoffStepMap := handoffStepList[0].(map[string]interface{})

			handoffStep := &components.UpdateTeamEscalationPolicyHandoffStep{
				TargetType: components.UpdateTeamEscalationPolicyTargetType(handoffStepMap["target_type"].(string)),
				TargetID:   handoffStepMap["target_id"].(string),
			}

			return handoffStep
		}
	}

	return nil
}

func getStepsFromResourceData(d *schema.ResourceData) []components.CreateTeamEscalationPolicyStep {
	var steps []components.CreateTeamEscalationPolicyStep

	if v, ok := d.GetOk("step"); ok {
		stepList := v.([]interface{})
		for _, stepItem := range stepList {
			stepMap := stepItem.(map[string]interface{})
			targets := stepMap["targets"].([]interface{})

			var stepTargets []components.CreateTeamEscalationPolicyTarget
			for _, targetItem := range targets {
				targetMap := targetItem.(map[string]interface{})
				target := components.CreateTeamEscalationPolicyTarget{
					Type: components.CreateTeamEscalationPolicyType(targetMap["type"].(string)),
					ID:   targetMap["id"].(string),
				}

				stepTargets = append(stepTargets, target)
			}

			step := components.CreateTeamEscalationPolicyStep{
				Timeout: stepMap["timeout"].(string),
				Targets: stepTargets,
			}

			steps = append(steps, step)
		}
	}

	return steps
}

func getNotificationPriorityPoliciesFromResourceDataUpdate(d *schema.ResourceData) *components.UpdateTeamEscalationPolicyPrioritizedSettings {
	priorityPolicies := d.Get("notification_priority_policies").([]interface{})
	if len(priorityPolicies) == 0 {
		return nil
	}

	settings := &components.UpdateTeamEscalationPolicyPrioritizedSettings{}

	for _, policyItem := range priorityPolicies {
		policyMap := policyItem.(map[string]interface{})
		priority := policyMap["priority"].(string)

		// Get repetitions if specified
		var repetitions *int
		if reps, ok := policyMap["repetitions"]; ok && reps != nil {
			if repsInt, ok := reps.(int); ok {
				repetitions = &repsInt
			}
		}

		// Set the appropriate priority policy
		switch firehydrant.NotificationPriority(priority) {
		case firehydrant.NotificationPriorityHigh:
			var handoffStep *components.UpdateTeamEscalationPolicyHighHandoffStep
			if handoff, ok := policyMap["handoff_step"]; ok && handoff != nil {
				handoffList := handoff.([]interface{})
				if len(handoffList) > 0 {
					handoffMap := handoffList[0].(map[string]interface{})
					handoffStep = &components.UpdateTeamEscalationPolicyHighHandoffStep{
						TargetType: components.UpdateTeamEscalationPolicyHighTargetType(handoffMap["target_type"].(string)),
						TargetID:   handoffMap["target_id"].(string),
					}
				}
			}
			settings.High = &components.UpdateTeamEscalationPolicyHigh{
				Repetitions: repetitions,
				HandoffStep: handoffStep,
			}
		case firehydrant.NotificationPriorityMedium:
			var handoffStep *components.UpdateTeamEscalationPolicyMediumHandoffStep
			if handoff, ok := policyMap["handoff_step"]; ok && handoff != nil {
				handoffList := handoff.([]interface{})
				if len(handoffList) > 0 {
					handoffMap := handoffList[0].(map[string]interface{})
					handoffStep = &components.UpdateTeamEscalationPolicyMediumHandoffStep{
						TargetType: components.UpdateTeamEscalationPolicyMediumTargetType(handoffMap["target_type"].(string)),
						TargetID:   handoffMap["target_id"].(string),
					}
				}
			}
			settings.Medium = &components.UpdateTeamEscalationPolicyMedium{
				Repetitions: repetitions,
				HandoffStep: handoffStep,
			}
		case firehydrant.NotificationPriorityLow:
			var handoffStep *components.UpdateTeamEscalationPolicyLowHandoffStep
			if handoff, ok := policyMap["handoff_step"]; ok && handoff != nil {
				handoffList := handoff.([]interface{})
				if len(handoffList) > 0 {
					handoffMap := handoffList[0].(map[string]interface{})
					handoffStep = &components.UpdateTeamEscalationPolicyLowHandoffStep{
						TargetType: components.UpdateTeamEscalationPolicyLowTargetType(handoffMap["target_type"].(string)),
						TargetID:   handoffMap["target_id"].(string),
					}
				}
			}
			settings.Low = &components.UpdateTeamEscalationPolicyLow{
				Repetitions: repetitions,
				HandoffStep: handoffStep,
			}
		}
	}

	return settings
}

func getNotificationPriorityPoliciesFromResourceData(d *schema.ResourceData) *components.CreateTeamEscalationPolicyPrioritizedSettings {
	priorityPolicies := d.Get("notification_priority_policies").([]interface{})
	if len(priorityPolicies) == 0 {
		return nil
	}

	settings := &components.CreateTeamEscalationPolicyPrioritizedSettings{}

	for _, policyItem := range priorityPolicies {
		policyMap := policyItem.(map[string]interface{})
		priority := policyMap["priority"].(string)

		// For dynamic escalation policies, the steps are defined at the main policy level
		// Priority-specific policies mainly handle repetitions and handoff steps
		// The steps from the policy are used for all priorities unless overridden

		// Get repetitions if specified
		var repetitions *int
		if reps, ok := policyMap["repetitions"]; ok && reps != nil {
			if repsInt, ok := reps.(int); ok {
				repetitions = &repsInt
			}
		}

		// Set the appropriate priority policy
		switch firehydrant.NotificationPriority(priority) {
		case firehydrant.NotificationPriorityHigh:
			var handoffStep *components.CreateTeamEscalationPolicyHighHandoffStep
			if handoff, ok := policyMap["handoff_step"]; ok && handoff != nil {
				handoffList := handoff.([]interface{})
				if len(handoffList) > 0 {
					handoffMap := handoffList[0].(map[string]interface{})
					handoffStep = &components.CreateTeamEscalationPolicyHighHandoffStep{
						TargetType: components.CreateTeamEscalationPolicyHighTargetType(handoffMap["target_type"].(string)),
						TargetID:   handoffMap["target_id"].(string),
					}
				}
			}
			settings.High = &components.CreateTeamEscalationPolicyHigh{
				Repetitions: repetitions,
				HandoffStep: handoffStep,
			}
		case firehydrant.NotificationPriorityMedium:
			var handoffStep *components.CreateTeamEscalationPolicyMediumHandoffStep
			if handoff, ok := policyMap["handoff_step"]; ok && handoff != nil {
				handoffList := handoff.([]interface{})
				if len(handoffList) > 0 {
					handoffMap := handoffList[0].(map[string]interface{})
					handoffStep = &components.CreateTeamEscalationPolicyMediumHandoffStep{
						TargetType: components.CreateTeamEscalationPolicyMediumTargetType(handoffMap["target_type"].(string)),
						TargetID:   handoffMap["target_id"].(string),
					}
				}
			}
			settings.Medium = &components.CreateTeamEscalationPolicyMedium{
				Repetitions: repetitions,
				HandoffStep: handoffStep,
			}
		case firehydrant.NotificationPriorityLow:
			var handoffStep *components.CreateTeamEscalationPolicyLowHandoffStep
			if handoff, ok := policyMap["handoff_step"]; ok && handoff != nil {
				handoffList := handoff.([]interface{})
				if len(handoffList) > 0 {
					handoffMap := handoffList[0].(map[string]interface{})
					handoffStep = &components.CreateTeamEscalationPolicyLowHandoffStep{
						TargetType: components.CreateTeamEscalationPolicyLowTargetType(handoffMap["target_type"].(string)),
						TargetID:   handoffMap["target_id"].(string),
					}
				}
			}
			settings.Low = &components.CreateTeamEscalationPolicyLow{
				Repetitions: repetitions,
				HandoffStep: handoffStep,
			}
		}
	}

	return settings
}

func getStepsFromResourceDataUpdateSDK(d *schema.ResourceData) []components.UpdateTeamEscalationPolicyStep {
	var steps []components.UpdateTeamEscalationPolicyStep

	if v, ok := d.GetOk("step"); ok {
		stepList := v.([]interface{})
		for _, stepItem := range stepList {
			stepMap := stepItem.(map[string]interface{})
			targets := stepMap["targets"].([]interface{})

			var stepTargets []components.UpdateTeamEscalationPolicyTarget
			for _, targetItem := range targets {
				targetMap := targetItem.(map[string]interface{})
				target := components.UpdateTeamEscalationPolicyTarget{
					Type: components.UpdateTeamEscalationPolicyType(targetMap["type"].(string)),
					ID:   targetMap["id"].(string),
				}

				stepTargets = append(stepTargets, target)
			}

			step := components.UpdateTeamEscalationPolicyStep{
				Timeout: stepMap["timeout"].(string),
				Targets: stepTargets,
			}

			steps = append(steps, step)
		}
	}

	return steps
}
