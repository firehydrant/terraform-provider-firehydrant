package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRunbook() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantRunbook,
		UpdateContext: updateResourceFireHydrantRunbook,
		ReadContext:   readResourceFireHydrantRunbook,
		DeleteContext: deleteResourceFireHydrantRunbook,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Optional
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"severities": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"steps": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Required
						"action_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						// Optional
						"automatic": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"config": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"delation_duration": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"repeats": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"repeats_duration": {
							Type:     schema.TypeString,
							Optional: true,
						},

						// Computed
						"step_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func readResourceFireHydrantRunbook(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the runbook
	runbookID := d.Id()
	runbookResponse, err := firehydrantAPIClient.Runbooks().Get(ctx, runbookID)
	if err != nil {
		_, isNotFoundError := err.(firehydrant.NotFound)
		if isNotFoundError {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"name":        runbookResponse.Name,
		"description": runbookResponse.Description,
		"type":        runbookResponse.Type,
	}

	steps := make([]interface{}, len(runbookResponse.Steps))
	for index, currentStep := range runbookResponse.Steps {
		stepConfig := map[string]interface{}{}
		for key, value := range currentStep.Config {
			stepConfig[key] = value
		}

		steps[index] = map[string]interface{}{
			"step_id":   currentStep.StepID,
			"name":      currentStep.Name,
			"action_id": currentStep.ActionID,
			"config":    stepConfig,
			"automatic": currentStep.Automatic,
		}
	}
	attributes["steps"] = steps

	severities := make([]interface{}, len(runbookResponse.Severities))
	for index, currentSeverity := range runbookResponse.Severities {
		severities[index] = map[string]interface{}{
			"id": currentSeverity.ID,
		}
	}
	attributes["severities"] = severities

	// Set the resource attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantRunbook(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get attributes from config and construct the create request
	createRequest := firehydrant.CreateRunbookRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
	}

	// Process any optional attributes and add to the create request if necessary
	steps := d.Get("steps").([]interface{})
	for _, currentStep := range steps {
		step := currentStep.(map[string]interface{})

		createRequest.Steps = append(createRequest.Steps, firehydrant.RunbookStep{
			Name:      step["name"].(string),
			ActionID:  step["action_id"].(string),
			Automatic: step["automatic"].(bool),
			Config:    convertStringMap(step["config"].(map[string]interface{})),
		})
	}

	severities := d.Get("severities").([]interface{})
	for _, severity := range severities {
		currentSeverity := severity.(map[string]interface{})

		createRequest.Severities = append(createRequest.Severities, firehydrant.RunbookRelation{
			ID: currentSeverity["id"].(string),
		})
	}

	// Create the new runbook
	runbookResponse, err := firehydrantAPIClient.Runbooks().Create(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the new runbook's ID in state
	d.SetId(runbookResponse.ID)

	// Update state with the latest information from the API
	return readResourceFireHydrantRunbook(ctx, d, m)
}

func updateResourceFireHydrantRunbook(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the update request
	updateRequest := firehydrant.UpdateRunbookRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Process any optional attributes and add to the update request if necessary
	steps := d.Get("steps").([]interface{})
	for _, currentStep := range steps {
		step := currentStep.(map[string]interface{})

		updateRequest.Steps = append(updateRequest.Steps, firehydrant.RunbookStep{
			Name:      step["name"].(string),
			ActionID:  step["action_id"].(string),
			Automatic: step["automatic"].(bool),
			Config:    convertStringMap(step["config"].(map[string]interface{})),
		})
	}

	severities := d.Get("severities").([]interface{})
	for _, currentSeverity := range severities {
		severity := currentSeverity.(map[string]interface{})

		updateRequest.Severities = append(updateRequest.Severities, firehydrant.RunbookRelation{
			ID: severity["id"].(string),
		})
	}

	// Update the runbook
	_, err := firehydrantAPIClient.Runbooks().Update(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantRunbook(ctx, d, m)
}

func deleteResourceFireHydrantRunbook(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the runbook
	runbookID := d.Id()
	err := firehydrantAPIClient.Runbooks().Delete(ctx, runbookID)
	if err != nil {
		_, isNotFoundError := err.(firehydrant.NotFound)
		if isNotFoundError {
			return nil
		}
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
