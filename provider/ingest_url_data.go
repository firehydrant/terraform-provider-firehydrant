package provider

import (
	"context"
	"strconv"
	"time"

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

	if team_id == "" && (escalation_policy_id != "" || schedule_id != "") {
		return diag.Errorf("`team_id` must be set if either `escalation_policy_id` or `on_call_schedule_id` is set")
	}

	ingestURL := ""
	if transposer == "" {
		// If no transposer is requested, we use the ingest URL API endpoint.  Otherwise, we use the transposers endpoint.
		params := firehydrant.IngestURLParams{
			TeamID:             team_id,
			UserID:             user_id,
			EscalationPolicyID: escalation_policy_id,
			OnCallScheduleID:   schedule_id,
		}

		url, err := firehydrantAPIClient.IngestURL().Get(ctx, params)
		if err != nil {
			return diag.FromErr(err)
		}
		ingestURL = url.URL
	} else {
		params := firehydrant.TransposersParams{
			TeamID:             team_id,
			UserID:             user_id,
			EscalationPolicyID: escalation_policy_id,
			OnCallScheduleID:   schedule_id,
		}

		ts, err := firehydrantAPIClient.Transposers().Get(ctx, params)
		if err != nil {
			return diag.FromErr(err)
		}
		for _, t := range ts.Transposers {
			if t.Slug == transposer {
				ingestURL = t.IngestURL
			}
		}
		if ingestURL == "" {
			return diag.Errorf("No transposer found with slug %s", transposer)
		}
	}

	// Set the ID
	if err := d.Set("url", ingestURL); err != nil {
		return diag.Errorf("Error setting url: %v", err)
	}
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}
