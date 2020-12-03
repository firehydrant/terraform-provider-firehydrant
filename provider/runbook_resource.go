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
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
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
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"step_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"config": {
							Type:     schema.TypeMap,
							Optional: true,
						},
						"automatic": {
							Type:     schema.TypeBool,
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
						"delation_duration": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func readResourceFireHydrantRunbook(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	r, err := ac.Runbooks().Get(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	var ds diag.Diagnostics
	svc := map[string]string{
		"name":        r.Name,
		"description": r.Description,
		"type":        r.Type,
	}

	for key, val := range svc {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := convertRunbookToState(r, d); err != nil {
		return diag.FromErr(err)
	}

	return ds
}

func createResourceFireHydrantRunbook(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	name, description, typ := d.Get("name").(string), d.Get("description").(string), d.Get("type").(string)

	r := firehydrant.CreateRunbookRequest{
		Name:        name,
		Description: description,
		Type:        typ,
	}

	steps := d.Get("steps").([]interface{})
	for _, step := range steps {
		s := step.(map[string]interface{})

		r.Steps = append(r.Steps, firehydrant.RunbookStep{
			Name:      s["name"].(string),
			ActionID:  s["action_id"].(string),
			Automatic: s["automatic"].(bool),
			Config:    convertStringMap(s["config"].(map[string]interface{})),
		})
	}

	severities := d.Get("severities").([]interface{})
	for _, sev := range severities {
		s := sev.(map[string]interface{})

		r.Severities = append(r.Severities, firehydrant.RunbookRelation{
			ID: s["id"].(string),
		})
	}

	resource, err := ac.Runbooks().Create(ctx, r)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.ID)

	attributes := map[string]interface{}{
		"name":        resource.Name,
		"type":        resource.Type,
		"description": resource.Description,
	}

	if err := setAttributesFromMap(d, attributes); err != nil {
		return diag.FromErr(err)
	}

	if err := convertRunbookToState(resource, d); err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func updateResourceFireHydrantRunbook(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	id := d.Id()

	r := firehydrant.UpdateRunbookRequest{
		Name:        name,
		Description: description,
	}

	steps := d.Get("steps").([]interface{})
	for _, step := range steps {
		s := step.(map[string]interface{})

		r.Steps = append(r.Steps, firehydrant.RunbookStep{
			Name:      s["name"].(string),
			ActionID:  s["action_id"].(string),
			Automatic: s["automatic"].(bool),
			Config:    convertStringMap(s["config"].(map[string]interface{})),
		})
	}

	severities := d.Get("severities").([]interface{})
	for _, sev := range severities {
		s := sev.(map[string]interface{})

		r.Severities = append(r.Severities, firehydrant.RunbookRelation{
			ID: s["id"].(string),
		})
	}

	_, err := ac.Runbooks().Update(ctx, id, r)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func deleteResourceFireHydrantRunbook(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	ac := m.(firehydrant.Client)

	err := ac.Runbooks().Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}

func convertRunbookToState(runbook *firehydrant.RunbookResponse, d *schema.ResourceData) error {
	resourceSteps := make([]interface{}, len(runbook.Steps))
	for index, s := range runbook.Steps {
		stepConfig := map[string]interface{}{}
		for k, v := range s.Config {
			stepConfig[k] = v
		}

		resourceSteps[index] = map[string]interface{}{
			"step_id":   s.StepID,
			"name":      s.Name,
			"action_id": s.ActionID,
			"config":    stepConfig,
			"automatic": s.Automatic,
		}
	}

	if err := d.Set("steps", resourceSteps); err != nil {
		return err
	}

	sevs := make([]interface{}, len(runbook.Severities))
	for index, s := range runbook.Severities {
		sevs[index] = map[string]interface{}{
			"id": s.ID,
		}
	}
	if err := d.Set("severities", sevs); err != nil {
		return err
	}

	return nil
}
