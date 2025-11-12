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

func resourceInboundEmail() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInboundEmailCreate,
		ReadContext:   resourceInboundEmailRead,
		UpdateContext: resourceInboundEmailUpdate,
		DeleteContext: resourceInboundEmailDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status_cel": {
				Type:     schema.TypeString,
				Required: true,
			},
			"level_cel": {
				Type:     schema.TypeString,
				Required: true,
			},
			"allowed_senders": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"target": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
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
			"rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rule_matching_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceInboundEmailCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	slug := d.Get("slug").(string)
	statusCel := d.Get("status_cel").(string)
	levelCel := d.Get("level_cel").(string)

	createReq := components.CreateSignalsEmailTarget{
		Name:           d.Get("name").(string),
		Slug:           &slug,
		StatusCel:      &statusCel,
		LevelCel:       &levelCel,
		AllowedSenders: expandStringSet(d.Get("allowed_senders").(*schema.Set)),
	}

	// Handle optional description
	if desc := d.Get("description").(string); desc != "" {
		createReq.Description = &desc
	}

	// Handle optional target
	if target := targetFromResourceData(d); target != nil {
		createReq.Target = target
	}

	// Handle optional rules - provide empty array if not specified
	if rulesSet := d.Get("rules").(*schema.Set); rulesSet.Len() > 0 {
		createReq.Rules = expandStringSet(rulesSet)
	} else {
		createReq.Rules = []string{}
	}

	// Handle optional rule_matching_strategy
	// The api defaults rule matching strategy to "all" for inbound emails so we need to handle
	// this field the same way within the provider otherwise, terrform state will always have unexpected diffs
	if ruleMatchingStrategy := d.Get("rule_matching_strategy").(string); ruleMatchingStrategy != "" {
		strategy := components.CreateSignalsEmailTargetRuleMatchingStrategy(ruleMatchingStrategy)
		createReq.RuleMatchingStrategy = &strategy
	} else {
		// Always set default strategy to match API behavior
		defaultStrategy := components.CreateSignalsEmailTargetRuleMatchingStrategyAll // api defaults to all
		createReq.RuleMatchingStrategy = &defaultStrategy
	}

	tflog.Debug(ctx, fmt.Sprintf("Create inbound email: %s", createReq.Name), map[string]interface{}{
		"name": createReq.Name,
	})

	inboundEmail, err := client.Sdk.Signals.CreateSignalsEmailTarget(ctx, createReq)
	if err != nil {
		return diag.Errorf("Error creating inbound email %s: %v", createReq.Name, err)
	}

	d.SetId(*inboundEmail.GetID())

	return resourceInboundEmailRead(ctx, d, m)
}

func resourceInboundEmailRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	inboundEmail, err := client.Sdk.Signals.GetSignalsEmailTarget(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var targetResourceData []interface{}
	if target := inboundEmail.GetTarget(); target != nil {
		targetResourceData = []interface{}{flattenTarget(target)}
	}

	d.Set("name", *inboundEmail.GetName())
	d.Set("slug", *inboundEmail.GetSlug())

	// Handle optional description
	if desc := inboundEmail.GetDescription(); desc != nil {
		d.Set("description", *desc)
	}

	d.Set("status_cel", *inboundEmail.GetStatusCel())
	d.Set("level_cel", *inboundEmail.GetLevelCel())
	d.Set("allowed_senders", inboundEmail.GetAllowedSenders())
	d.Set("target", targetResourceData)
	d.Set("rules", inboundEmail.GetRules())

	// Handle rule_matching_strategy - API always returns this field with default "all"
	if strategy := inboundEmail.GetRuleMatchingStrategy(); strategy != nil {
		d.Set("rule_matching_strategy", *strategy)
	} else {
		// Fallback to default if somehow not present
		d.Set("rule_matching_strategy", "all")
	}

	// Handle computed email field
	if email := inboundEmail.GetEmail(); email != nil {
		d.Set("email", *email)
	}

	return nil
}

func resourceInboundEmailUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	updateReq := components.UpdateSignalsEmailTarget{}

	// Set required fields
	name := d.Get("name").(string)
	updateReq.Name = &name
	slug := d.Get("slug").(string)
	updateReq.Slug = &slug
	statusCel := d.Get("status_cel").(string)
	updateReq.StatusCel = &statusCel
	levelCel := d.Get("level_cel").(string)
	updateReq.LevelCel = &levelCel
	updateReq.AllowedSenders = expandStringSet(d.Get("allowed_senders").(*schema.Set))

	// Handle optional description
	if desc := d.Get("description").(string); desc != "" {
		updateReq.Description = &desc
	}

	// Handle optional target
	if target := targetFromResourceDataForUpdate(d); target != nil {
		updateReq.Target = target
	}

	// Handle optional rules - provide empty array if not specified
	if rulesSet := d.Get("rules").(*schema.Set); rulesSet.Len() > 0 {
		updateReq.Rules = expandStringSet(rulesSet)
	} else {
		updateReq.Rules = []string{}
	}

	// Handle optional rule_matching_strategy - always set a value to match API defaults
	if ruleMatchingStrategy := d.Get("rule_matching_strategy").(string); ruleMatchingStrategy != "" {
		strategy := components.UpdateSignalsEmailTargetRuleMatchingStrategy(ruleMatchingStrategy)
		updateReq.RuleMatchingStrategy = &strategy
	} else {
		// Always set default strategy to match API behavior
		defaultStrategy := components.UpdateSignalsEmailTargetRuleMatchingStrategyAll
		updateReq.RuleMatchingStrategy = &defaultStrategy
	}

	tflog.Debug(ctx, fmt.Sprintf("Update inbound email: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})

	_, err := client.Sdk.Signals.UpdateSignalsEmailTarget(ctx, d.Id(), updateReq)
	if err != nil {
		return diag.Errorf("Error updating inbound email %s: %v", d.Id(), err)
	}

	return resourceInboundEmailRead(ctx, d, m)
}

func resourceInboundEmailDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	tflog.Debug(ctx, fmt.Sprintf("Delete inbound email: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})

	err := client.Sdk.Signals.DeleteSignalsEmailTarget(ctx, d.Id())
	if err != nil {
		return diag.Errorf("Error deleting inbound email %s: %v", d.Id(), err)
	}

	d.SetId("")

	return nil
}

func expandStringSet(set *schema.Set) []string {
	s := make([]string, 0, set.Len())
	for _, v := range set.List() {
		s = append(s, v.(string))
	}
	return s
}

func targetFromResourceData(d *schema.ResourceData) *components.CreateSignalsEmailTargetTarget {
	if len(d.Get("target").([]interface{})) == 0 {
		return nil
	}

	t := d.Get("target").([]interface{})[0].(map[string]interface{})
	return &components.CreateSignalsEmailTargetTarget{
		Type: components.CreateSignalsEmailTargetType(t["type"].(string)),
		ID:   t["id"].(string),
	}
}

func targetFromResourceDataForUpdate(d *schema.ResourceData) *components.UpdateSignalsEmailTargetTarget {
	if len(d.Get("target").([]interface{})) == 0 {
		return nil
	}

	t := d.Get("target").([]interface{})[0].(map[string]interface{})
	return &components.UpdateSignalsEmailTargetTarget{
		Type: components.UpdateSignalsEmailTargetType(t["type"].(string)),
		ID:   t["id"].(string),
	}
}

func flattenTarget(target *components.NullableSignalsAPITargetEntity) map[string]interface{} {
	if target == nil {
		return nil
	}

	id := target.GetID()
	typeVal := target.GetType()

	if id == nil || typeVal == nil {
		return nil
	}

	return map[string]interface{}{
		"type": *typeVal,
		"id":   *id,
	}
}
