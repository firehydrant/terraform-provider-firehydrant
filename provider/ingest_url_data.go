package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIngestURL() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantIngestURL,
		Schema: map[string]*schema.Schema{
			// Optional
			"transposer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"team_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"user_id"},
			},
			"user_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"team_id"},
			},
			"escalation_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"on_call_schedule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataFireHydrantIngestURL(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	team_id := d.Get("team_id").(string)
	user_id := d.Get("user_id").(string)
	escalation_policy_id := d.Get("escalation_policy_id").(string)
	schedule_id := d.Get("on_call_schedule_id").(string)
	transposer := d.Get("transposer").(string)

	params := firehydrant.IngestURLParams{
		TeamID:             team_id,
		UserID:             user_id,
		EscalationPolicyID: escalation_policy_id,
		OnCallScheduleID:   schedule_id,
	}

	if team_id == "" && (escalation_policy_id != "" || schedule_id != "") {
		return diag.Errorf("`team_id` must be set if either `escalation_policy_id` or `on_call_schedule_id` is set")
	}
	url, err := firehydrantAPIClient.IngestURL().Get(ctx, params)
	if err != nil {
		return diag.FromErr(err)
	}

	// Not a huge fan of the URL hacking here, but we can't get this directly from the API.  If we want, we can validate transposer
	// name against list of slugs from `curl https://api.firehydrant.io/v1/signals/transposers | jq -r '.data | .[] | .slug'`
	// but even then the composition here bothers me.

	finalURL := url.URL
	if transposer != "" {
		finalURL = strings.Replace(url.URL, "/process", fmt.Sprintf("/transpose/%s", transposer), 1)
	}

	// Set the ID
	if err := d.Set("url", finalURL); err != nil {
		return diag.Errorf("Error setting url: %v", err)
	}

	return diag.Diagnostics{}
}
