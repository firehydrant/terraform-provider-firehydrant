package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		},
	}
}

func readResourceFireHydrantEscalationPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the escalation policy
	id := d.Id()
	teamID := d.Get("team_id").(string)
	tflog.Debug(ctx, "Read escalation policy", map[string]interface{}{
		"id":      id,
		"team_id": teamID,
	})

	escalationPolicy, err := firehydrantAPIClient.EscalationPolicies().Get(ctx, teamID, id)
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
	d.Set("name", escalationPolicy.Name)
	d.Set("description", escalationPolicy.Description)
	d.Set("default", escalationPolicy.Default)
	d.Set("repetitions", escalationPolicy.Repetitions)

	// Set the steps
	var steps []map[string]interface{}
	for _, step := range escalationPolicy.Steps {
		targets := []map[string]interface{}{}
		for _, target := range step.Targets {
			targetMap := map[string]interface{}{
				"type": target.Type,
				"id":   target.ID,
			}

			targets = append(targets, targetMap)
		}

		stepMap := map[string]interface{}{
			"timeout": step.Timeout,
			"targets": targets,
		}

		steps = append(steps, stepMap)
	}

	d.Set("step", steps)

	// Set the handoff step
	if escalationPolicy.HandoffStep != nil {
		handoffStepMap := map[string]interface{}{
			"target_type": escalationPolicy.HandoffStep.Target.Type,
			"target_id":   escalationPolicy.HandoffStep.Target.ID,
		}

		d.Set("handoff_step", []map[string]interface{}{handoffStepMap})
	}

	return diag.Diagnostics{}
}

// creates an escalation policy for a team using the firehydrant api client
func createResourceFireHydrantEscalationPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	// Create the escalation policy
	teamID := d.Get("team_id").(string)
	createReq := firehydrant.CreateEscalationPolicyRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Default:     d.Get("default").(bool),
		Repetitions: d.Get("repetitions").(int),
		Steps:       getStepsFromResourceData(d),
		HandoffStep: getHandoffStepFromResourceData(d),
	}

	tflog.Debug(ctx, "Creating escalation policy", map[string]interface{}{
		"team_id": teamID,
		"request": createReq,
	})

	escalationPolicy, err := firehydrantAPIClient.EscalationPolicies().Create(ctx, teamID, createReq)
	if err != nil {
		return diag.Errorf("Error creating signal rule %s: %v", d.Id(), err)
	}

	// Set the ID of the escalation policy
	d.SetId(escalationPolicy.ID)
	return readResourceFireHydrantEscalationPolicy(ctx, d, m)
}

func updateResourceFireHydrantEscalationPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	// Update the escalation policy
	teamID := d.Get("team_id").(string)
	updateReq := firehydrant.UpdateEscalationPolicyRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Default:     d.Get("default").(bool),
		Repetitions: d.Get("repetitions").(int),
		Steps:       getStepsFromResourceData(d),
		HandoffStep: getHandoffStepFromResourceData(d),
	}

	tflog.Debug(ctx, "Updating escalation policy", map[string]interface{}{
		"team_id": teamID,
		"request": spew.Sdump(updateReq),
	})

	_, err := firehydrantAPIClient.EscalationPolicies().Update(ctx, teamID, d.Id(), updateReq)
	if err != nil {
		return diag.Errorf("Error updating escalation policy %s: %v", d.Id(), err)
	}

	return readResourceFireHydrantEscalationPolicy(ctx, d, m)
}

func deleteResourceFireHydrantEscalationPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the escalation policy
	teamID := d.Get("team_id").(string)
	tflog.Debug(ctx, "Deleting escalation policy", map[string]interface{}{
		"team_id": teamID,
		"id":      d.Id(),
	})

	err := firehydrantAPIClient.EscalationPolicies().Delete(ctx, teamID, d.Id())
	if err != nil {
		return diag.Errorf("Error deleting escalation policy %s: %v", d.Id(), err)
	}

	return nil
}

func getHandoffStepFromResourceData(d *schema.ResourceData) *firehydrant.CreateEscalationPolicyHandoffStep {
	if v, ok := d.GetOk("handoff_step"); ok {
		handoffStepList := v.([]interface{})

		if len(handoffStepList) > 0 {
			handoffStepMap := handoffStepList[0].(map[string]interface{})

			handoffStep := &firehydrant.CreateEscalationPolicyHandoffStep{
				Type: handoffStepMap["target_type"].(string),
				ID:   handoffStepMap["target_id"].(string),
			}

			return handoffStep
		}
	}

	return nil
}

func getStepsFromResourceData(d *schema.ResourceData) []firehydrant.EscalationPolicyStep {
	var steps []firehydrant.EscalationPolicyStep

	if v, ok := d.GetOk("step"); ok {
		stepList := v.([]interface{})
		for position, stepItem := range stepList {
			stepMap := stepItem.(map[string]interface{})
			targets := stepMap["targets"].([]interface{})

			var stepTargets []firehydrant.EscalationPolicyTarget
			for _, targetItem := range targets {
				targetMap := targetItem.(map[string]interface{})
				target := firehydrant.EscalationPolicyTarget{
					// ID is not provided in resource data, so it's not set here
					Type: targetMap["type"].(string),
					ID:   targetMap["id"].(string),
				}

				stepTargets = append(stepTargets, target)
			}

			step := firehydrant.EscalationPolicyStep{
				Position: position,
				Timeout:  stepMap["timeout"].(string),
				Targets:  stepTargets,
			}

			steps = append(steps, step)
		}
	}

	return steps
}
