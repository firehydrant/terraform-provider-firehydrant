package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSignalRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantSignalRule,
		UpdateContext: updateResourceFireHydrantSignalRule,
		ReadContext:   readResourceFireHydrantSignalRule,
		DeleteContext: deleteResourceFireHydrantSignalRule,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"incident_type_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func readResourceFireHydrantSignalRule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the signal rule
	id := d.Id()
	teamID := d.Get("team_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read signal rule: %s", id), map[string]interface{}{
		"id":      id,
		"team_id": teamID,
	})

	signalRule, err := firehydrantAPIClient.SignalsRules().Get(ctx, teamID, id)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			tflog.Debug(ctx, fmt.Sprintf("Signal rule %s no longer exists", id), map[string]interface{}{
				"id":      id,
				"team_id": teamID,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading signal rule %s: %v", id, err)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"name":             signalRule.Name,
		"expression":       signalRule.Expression,
		"target_type":      signalRule.Target.Type,
		"target_id":        signalRule.Target.ID,
		"incident_type_id": signalRule.IncidentType.ID,
	}

	// Set the resource attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for signal rule %s: %v", key, id, err)
		}
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantSignalRule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the create request
	createRequest := firehydrant.CreateSignalsRuleRequest{
		Name:           d.Get("name").(string),
		Expression:     d.Get("expression").(string),
		TargetType:     d.Get("target_type").(string),
		TargetID:       d.Get("target_id").(string),
		IncidentTypeID: d.Get("incident_type_id").(string),
	}

	// Create the signal rule
	tflog.Debug(ctx, fmt.Sprintf("Create signal rule: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})
	signalRuleResponse, err := firehydrantAPIClient.SignalsRules().Create(ctx, d.Get("team_id").(string), createRequest)
	if err != nil {
		return diag.Errorf("Error creating signal rule %s: %v", d.Id(), err)
	}

	// Set the ID of the resource to the ID of the newly created signal rule
	d.SetId(signalRuleResponse.ID)

	// Update state with the latest information from the API
	return readResourceFireHydrantSignalRule(ctx, d, m)
}

func updateResourceFireHydrantSignalRule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the update request
	updateRequest := firehydrant.UpdateSignalsRuleRequest{
		Name:           d.Get("name").(string),
		Expression:     d.Get("expression").(string),
		TargetType:     d.Get("target_type").(string),
		TargetID:       d.Get("target_id").(string),
		IncidentTypeID: d.Get("incident_type_id").(string),
	}

	// Update the signal rule
	tflog.Debug(ctx, fmt.Sprintf("Update signal rule: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})
	_, err := firehydrantAPIClient.SignalsRules().Update(ctx, d.Get("team_id").(string), d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating signal rule %s: %v", d.Id(), err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantSignalRule(ctx, d, m)
}

func deleteResourceFireHydrantSignalRule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the signal rule
	tflog.Debug(ctx, fmt.Sprintf("Delete signal rule: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})
	err := firehydrantAPIClient.SignalsRules().Delete(ctx, d.Get("team_id").(string), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting signal rule %s: %v", d.Id(), err)
	}

	return diag.Diagnostics{}
}
