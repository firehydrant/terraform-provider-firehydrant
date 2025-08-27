package provider

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOnCallSchedule() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantOnCallSchedule,
		ReadContext:   readResourceFireHydrantOnCallSchedule,
		UpdateContext: updateResourceFireHydrantOnCallSchedule,
		DeleteContext: deleteResourceFireHydrantOnCallSchedule,
		// Importer: &schema.ResourceImporter{
		// 	StateContext: schema.ImportStatePassthroughContext,
		// },

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_ids": {
				Type:          schema.TypeList,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Optional:      true, // will be required in the future once `members` has been removed.
				ConflictsWith: []string{"members"},
			},
			"members": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				// Technically, I (wilsonehusin) don't think this ever worked because it would produce HTTP 400s.
				// Documentation also always mentioned `member_ids` as the correct attribute to use.
				// Leaving this here for now to prevent potential breaking changes.
				Deprecated:    "Use member_ids to configure membership; members attribute will be removed in a future release.",
				ConflictsWith: []string{"member_ids"},
			},
			"time_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"strategy": {
				Type:     schema.TypeList, // Using TypeList to simulate a map
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"handoff_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"handoff_day": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"shift_duration": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"color": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slack_user_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"restrictions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_day": {
							Type:     schema.TypeString,
							Required: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"end_day": {
							Type:     schema.TypeString,
							Required: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"effective_at": {
				Type:     schema.TypeString,
				Optional: true,
				// Don't set computed:true since we don't want it in the state
				Description: "RFC3339 timestamp for when the schedule update should take effect. If not provided or if the time is in the past, the update will take effect immediately.",
				ValidateDiagFunc: schema.SchemaValidateDiagFunc(
					func(v interface{}, path cty.Path) diag.Diagnostics {
						timeStr := v.(string)
						_, err := time.Parse(time.RFC3339, timeStr)
						if err != nil {
							return diag.Diagnostics{
								diag.Diagnostic{
									Severity:      diag.Error,
									Summary:       "Invalid effective_at timestamp",
									Detail:        fmt.Sprintf("effective_at must be a valid RFC3339 timestamp (e.g. 2024-01-01T15:04:05Z), got: %s", timeStr),
									AttributePath: path,
								},
							}
						}

						return nil
					},
				),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
			},
		},
	}
}

func createResourceFireHydrantOnCallSchedule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Create the on-call schedule
	teamID := d.Get("team_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Create on-call schedule: %s", teamID), map[string]interface{}{
		"team_id": teamID,
	})

	inputMemberIDs := d.Get("member_ids").([]interface{})
	if len(inputMemberIDs) == 0 {
		inputMemberIDs = d.Get("members").([]interface{})
	}
	memberIDs := []string{}
	for _, memberID := range inputMemberIDs {
		if v, ok := memberID.(string); ok && v != "" {
			memberIDs = append(memberIDs, v)
		}
	}

	// Gather values from API response
	onCallSchedule := firehydrant.CreateOnCallScheduleRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		TimeZone:    d.Get("time_zone").(string),
		Strategy: firehydrant.OnCallScheduleStrategy{
			Type:          d.Get("strategy.0.type").(string),
			HandoffTime:   d.Get("strategy.0.handoff_time").(string),
			HandoffDay:    d.Get("strategy.0.handoff_day").(string),
			ShiftDuration: d.Get("strategy.0.shift_duration").(string),
		},
		StartTime:    d.Get("start_time").(string),
		MemberIDs:    memberIDs,
		Restrictions: oncallRestrictionsFromData(d),
	}

	// Get slack_user_group_id if set and non-empty
	if v, ok := d.GetOk("slack_user_group_id"); ok && v.(string) != "" {
		onCallSchedule.SlackUserGroupID = v.(string)
	}

	if onCallSchedule.Strategy.Type != "" {
		isCustomStrategy := onCallSchedule.Strategy.Type == "custom"
		if isCustomStrategy {
			if onCallSchedule.Strategy.ShiftDuration == "" {
				return diag.Errorf("firehydrant_on_call_schedule.strategy.shift_duration is required when strategy type is 'custom'")
			}
			if onCallSchedule.StartTime == "" {
				return diag.Errorf("firehydrant_on_call_schedule.start_time is required when strategy type is 'custom'")
			}

			// Discard unused values to avoid ambiguity.
			onCallSchedule.Strategy.HandoffTime = ""
			onCallSchedule.Strategy.HandoffDay = ""
		} else {
			if onCallSchedule.Strategy.HandoffTime == "" {
				return diag.Errorf("firehydrant_on_call_schedule.strategy.handoff_time is required when strategy type is '%s'", onCallSchedule.Strategy.Type)
			}
			if onCallSchedule.Strategy.Type == "weekly" && onCallSchedule.Strategy.HandoffDay == "" {
				return diag.Errorf("firehydrant_on_call_schedule.strategy.handoff_day is required when strategy type is '%s'", onCallSchedule.Strategy.Type)
			}

			// Discard unused values to avoid ambiguity.
			onCallSchedule.Strategy.ShiftDuration = ""
			onCallSchedule.StartTime = ""
		}
	}

	// Create the on-call schedule
	createdOnCallSchedule, err := firehydrantAPIClient.OnCallSchedules().Create(ctx, teamID, onCallSchedule)
	if err != nil {
		return diag.Errorf("Error creating on-call schedule %s: %v", teamID, err)
	}

	// Set the on-call schedule's ID in state
	d.SetId(createdOnCallSchedule.ID)

	return readResourceFireHydrantOnCallSchedule(ctx, d, m)
}

func readResourceFireHydrantOnCallSchedule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the signal rule
	id := d.Id()
	teamID := d.Get("team_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read signal rule: %s", id), map[string]interface{}{
		"id":      id,
		"team_id": teamID,
	})

	onCallSchedule, err := firehydrantAPIClient.OnCallSchedules().Get(ctx, teamID, id)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			tflog.Debug(ctx, fmt.Sprintf("On-call schedule %s no longer exists", id), map[string]interface{}{
				"id":      id,
				"team_id": teamID,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading on-call schedule %s: %v", id, err)
	}

	// Gather values from API response
	memberIDs := make([]string, len(onCallSchedule.Members))
	for i, member := range onCallSchedule.Members {
		memberIDs[i] = member.ID
	}

	attributes := map[string]interface{}{
		"name":         onCallSchedule.Name,
		"description":  onCallSchedule.Description,
		"time_zone":    onCallSchedule.TimeZone,
		"strategy":     strategyToMap(onCallSchedule.Strategy),
		"member_ids":   memberIDs,
		"restrictions": restrictionsToData(onCallSchedule.Restrictions),
	}
	if onCallSchedule.SlackUserGroupID != "" {
		attributes["slack_user_group_id"] = onCallSchedule.SlackUserGroupID
	}

	// Set the data source attributes to the values we got from the API
	for key, val := range attributes {
		if err := d.Set(key, val); err != nil {
			return diag.Errorf("Error setting %s for on-call schedule %s: %v", key, id, err)
		}
	}

	// Set the on-call schedule's ID in state
	d.SetId(onCallSchedule.ID)

	return diag.Diagnostics{}
}

func updateResourceFireHydrantOnCallSchedule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	id := d.Id()
	teamID := d.Get("team_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Update on-call schedule: %s", id), map[string]interface{}{
		"id":      id,
		"team_id": teamID,
	})

	// Initialize updateRequest with basic fields
	updateRequest := firehydrant.UpdateOnCallScheduleRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	// Get slack_user_group_id if set
	if v, ok := d.GetOk("slack_user_group_id"); ok {
		updateRequest.SlackUserGroupID = v.(string)
	}

	// Check if effective_at exists in raw config rather than state
	if raw := d.GetRawConfig().GetAttr("effective_at"); !raw.IsNull() {
		effectiveAtStr := raw.AsString()
		effectiveAt, err := time.Parse(time.RFC3339, effectiveAtStr)
		if err != nil {
			return diag.FromErr(err)
		}

		// Only set effective_at if it's in the future
		if effectiveAt.After(time.Now()) {
			updateRequest.EffectiveAt = effectiveAt.Format(time.RFC3339)
			tflog.Debug(ctx, "Schedule update will take effect at: "+updateRequest.EffectiveAt, map[string]interface{}{
				"effective_at": updateRequest.EffectiveAt,
			})
		} else {
			tflog.Debug(ctx, "Provided effective_at is in the past, update will take effect immediately", map[string]interface{}{
				"effective_at": effectiveAtStr,
				"now":          time.Now().Format(time.RFC3339),
			})
		}
	}

	// Get member IDs
	inputMemberIDs := d.Get("member_ids").([]interface{})
	if len(inputMemberIDs) == 0 {
		inputMemberIDs = d.Get("members").([]interface{})
	}
	memberIDs := []string{}
	for _, memberID := range inputMemberIDs {
		if v, ok := memberID.(string); ok && v != "" {
			memberIDs = append(memberIDs, v)
		}
	}
	updateRequest.MemberIDs = memberIDs

	// Get strategy configuration
	if v, ok := d.GetOk("strategy"); ok {
		if strategies := v.([]interface{}); len(strategies) > 0 {
			strategy := strategies[0].(map[string]interface{})
			updateRequest.Strategy = &firehydrant.OnCallScheduleStrategy{
				Type:        strategy["type"].(string),
				HandoffTime: strategy["handoff_time"].(string),
				HandoffDay:  strategy["handoff_day"].(string),
			}

			// Set shift duration for custom strategy
			if strategy["type"].(string) == "custom" {
				updateRequest.Strategy.ShiftDuration = strategy["shift_duration"].(string)
			}
		}
	}

	// Get restrictions
	restrictions := d.Get("restrictions").([]interface{})
	for _, r := range restrictions {
		restriction := r.(map[string]interface{})
		updateRequest.Restrictions = append(updateRequest.Restrictions, firehydrant.OnCallScheduleRestriction{
			StartDay:  restriction["start_day"].(string),
			StartTime: restriction["start_time"].(string),
			EndDay:    restriction["end_day"].(string),
			EndTime:   restriction["end_time"].(string),
		})
	}

	// Update the on-call schedule
	_, err := firehydrantAPIClient.OnCallSchedules().Update(ctx, teamID, id, updateRequest)
	if err != nil {
		return diag.Errorf("Error updating on-call schedule %s: %v", id, err)
	}

	return readResourceFireHydrantOnCallSchedule(ctx, d, m)
}

func deleteResourceFireHydrantOnCallSchedule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	id := d.Id()
	teamID := d.Get("team_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Delete on-call schedule: %s", id), map[string]interface{}{
		"id":      id,
		"team_id": teamID,
	})

	// Delete the on-call schedule
	err := firehydrantAPIClient.OnCallSchedules().Delete(ctx, teamID, id)
	if err != nil {
		return diag.Errorf("Error deleting on-call schedule %s: %v", id, err)
	}

	// Remove the on-call schedule's ID from state
	d.SetId("")

	return diag.Diagnostics{}
}

func strategyToMap(strategy firehydrant.OnCallScheduleStrategy) []map[string]interface{} {
	m := map[string]interface{}{"type": strategy.Type}
	if strategy.Type == "custom" {
		m["shift_duration"] = strategy.ShiftDuration
	} else {
		m["handoff_time"] = strategy.HandoffTime
	}
	if strategy.Type == "weekly" {
		m["handoff_day"] = strategy.HandoffDay
	}
	return []map[string]interface{}{m}
}

func oncallRestrictionsFromData(d *schema.ResourceData) []firehydrant.OnCallScheduleRestriction {
	restrictions := make([]firehydrant.OnCallScheduleRestriction, 0)
	for _, restriction := range d.Get("restrictions").([]interface{}) {
		restrictionMap := restriction.(map[string]interface{})
		restrictions = append(restrictions, firehydrant.OnCallScheduleRestriction{
			StartDay:  restrictionMap["start_day"].(string),
			StartTime: restrictionMap["start_time"].(string),
			EndDay:    restrictionMap["end_day"].(string),
			EndTime:   restrictionMap["end_time"].(string),
		})
	}
	return restrictions
}

func restrictionsToData(restrictions []firehydrant.OnCallScheduleRestriction) []map[string]interface{} {
	restrictionMaps := make([]map[string]interface{}, 0)
	for _, restriction := range restrictions {
		restrictionMaps = append(restrictionMaps, map[string]interface{}{
			"start_day":  restriction.StartDay,
			"start_time": restriction.StartTime,
			"end_day":    restriction.EndDay,
			"end_time":   restriction.EndTime,
		})
	}
	return restrictionMaps
}
