package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
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
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"rule_matching_strategy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceInboundEmailCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(firehydrant.Client)

	createReq := firehydrant.CreateInboundEmailRequest{
		Name:                 d.Get("name").(string),
		Slug:                 d.Get("slug").(string),
		Description:          d.Get("description").(string),
		StatusCEL:            d.Get("status_cel").(string),
		LevelCEL:             d.Get("level_cel").(string),
		AllowedSenders:       expandStringSet(d.Get("allowed_senders").(*schema.Set)),
		Target:               expandTarget(d.Get("target").([]interface{})[0].(map[string]interface{})),
		Rules:                expandStringSet(d.Get("rules").(*schema.Set)),
		RuleMatchingStrategy: d.Get("rule_matching_strategy").(string),
	}

	inboundEmail, err := c.InboundEmails().Create(ctx, createReq)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(inboundEmail.ID)

	return resourceInboundEmailRead(ctx, d, m)
}

func resourceInboundEmailRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(firehydrant.Client)

	inboundEmail, err := c.InboundEmails().Get(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", inboundEmail.Name)
	d.Set("slug", inboundEmail.Slug)
	d.Set("description", inboundEmail.Description)
	d.Set("status_cel", inboundEmail.StatusCEL)
	d.Set("level_cel", inboundEmail.LevelCEL)
	d.Set("allowed_senders", inboundEmail.AllowedSenders)
	d.Set("target", []interface{}{flattenTarget(inboundEmail.Target)})
	d.Set("rules", inboundEmail.Rules)
	d.Set("rule_matching_strategy", inboundEmail.RuleMatchingStrategy)
	d.Set("email", inboundEmail.Email)

	return nil
}

func resourceInboundEmailUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(firehydrant.Client)

	updateReq := firehydrant.UpdateInboundEmailRequest{
		Name:                 d.Get("name").(string),
		Slug:                 d.Get("slug").(string),
		Description:          d.Get("description").(string),
		StatusCEL:            d.Get("status_cel").(string),
		LevelCEL:             d.Get("level_cel").(string),
		AllowedSenders:       expandStringSet(d.Get("allowed_senders").(*schema.Set)),
		Target:               expandTarget(d.Get("target").([]interface{})[0].(map[string]interface{})),
		Rules:                expandStringSet(d.Get("rules").(*schema.Set)),
		RuleMatchingStrategy: d.Get("rule_matching_strategy").(string),
	}

	_, err := c.InboundEmails().Update(ctx, d.Id(), updateReq)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceInboundEmailRead(ctx, d, m)
}

func resourceInboundEmailDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(firehydrant.Client)

	err := c.InboundEmails().Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
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

func expandTarget(m map[string]interface{}) firehydrant.Target {
	return firehydrant.Target{
		Type: m["type"].(string),
		ID:   m["id"].(string),
	}
}

func flattenTarget(target firehydrant.Target) map[string]interface{} {
	return map[string]interface{}{
		"type": target.Type,
		"id":   target.ID,
	}
}
