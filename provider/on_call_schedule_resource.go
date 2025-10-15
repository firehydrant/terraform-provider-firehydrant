package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
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
		Importer: &schema.ResourceImporter{
			StateContext: importResourceFireHydrantOnCallSchedule,
		},

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
	client := m.(*firehydrant.APIClient)

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
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	timeZone := d.Get("time_zone").(string)
	startTime := d.Get("start_time").(string)
	handoffTime := d.Get("strategy.0.handoff_time").(string)
	handoffDay := d.Get("strategy.0.handoff_day").(string)
	shiftDuration := d.Get("strategy.0.shift_duration").(string)

	onCallSchedule := components.CreateTeamOnCallSchedule{
		Name:        name,
		Description: &description,
		TimeZone:    &timeZone,
		Strategy: &components.CreateTeamOnCallScheduleStrategy{
			Type:          components.CreateTeamOnCallScheduleType(d.Get("strategy.0.type").(string)),
			HandoffTime:   &handoffTime,
			HandoffDay:    (*components.CreateTeamOnCallScheduleHandoffDay)(&handoffDay),
			ShiftDuration: &shiftDuration,
		},
		StartTime:    &startTime,
		MemberIds:    memberIDs,
		Restrictions: oncallRestrictionsFromDataSDK(d),
	}

	// Get slack_user_group_id if set and non-empty
	if v, ok := d.GetOk("slack_user_group_id"); ok && v.(string) != "" {
		slackUserGroupID := v.(string)
		onCallSchedule.SlackUserGroupID = &slackUserGroupID
	}

	if onCallSchedule.Strategy.Type != "" {
		isCustomStrategy := onCallSchedule.Strategy.Type == "custom"
		if isCustomStrategy {
			if onCallSchedule.Strategy.ShiftDuration == nil || *onCallSchedule.Strategy.ShiftDuration == "" {
				return diag.Errorf("firehydrant_on_call_schedule.strategy.shift_duration is required when strategy type is 'custom'")
			}
			if onCallSchedule.StartTime == nil || *onCallSchedule.StartTime == "" {
				return diag.Errorf("firehydrant_on_call_schedule.start_time is required when strategy type is 'custom'")
			}

			// Discard unused values to avoid ambiguity.
			onCallSchedule.Strategy.HandoffTime = nil
			onCallSchedule.Strategy.HandoffDay = nil
		} else {
			if onCallSchedule.Strategy.HandoffTime == nil || *onCallSchedule.Strategy.HandoffTime == "" {
				return diag.Errorf("firehydrant_on_call_schedule.strategy.handoff_time is required when strategy type is '%s'", onCallSchedule.Strategy.Type)
			}
			if onCallSchedule.Strategy.Type == "weekly" && (onCallSchedule.Strategy.HandoffDay == nil || *onCallSchedule.Strategy.HandoffDay == "") {
				return diag.Errorf("firehydrant_on_call_schedule.strategy.handoff_day is required when strategy type is '%s'", onCallSchedule.Strategy.Type)
			}

			// Discard unused values to avoid ambiguity.
			onCallSchedule.Strategy.ShiftDuration = nil
			onCallSchedule.StartTime = nil
		}
	}

	// Create the on-call schedule
	createdOnCallSchedule, err := client.Sdk.Signals.CreateTeamOnCallSchedule(ctx, teamID, onCallSchedule)
	if err != nil {
		return diag.Errorf("Error creating on-call schedule %s: %v", teamID, err)
	}

	// Set the on-call schedule's ID in state
	d.SetId(*createdOnCallSchedule.GetID())

	return readResourceFireHydrantOnCallSchedule(ctx, d, m)
}

func readResourceFireHydrantOnCallSchedule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Get the on-call schedule
	id := d.Id()
	teamID := d.Get("team_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read on-call schedule: %s", id), map[string]interface{}{
		"id":      id,
		"team_id": teamID,
	})

	onCallSchedule, err := client.Sdk.Signals.GetTeamOnCallSchedule(ctx, teamID, id, nil, nil)
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
	memberIDs := make([]string, len(onCallSchedule.GetMembers()))
	for i, member := range onCallSchedule.GetMembers() {
		memberIDs[i] = *member.GetID()
	}

	attributes := map[string]interface{}{
		"name":         *onCallSchedule.GetName(),
		"description":  *onCallSchedule.GetDescription(),
		"time_zone":    *onCallSchedule.GetTimeZone(),
		"strategy":     strategyToMapSDK(*onCallSchedule.GetStrategy()),
		"member_ids":   memberIDs,
		"restrictions": restrictionsToDataSDK(onCallSchedule.GetRestrictions()),
	}
	if slackUserGroupID := onCallSchedule.GetSlackUserGroupID(); slackUserGroupID != nil && *slackUserGroupID != "" {
		attributes["slack_user_group_id"] = *slackUserGroupID
	}

	// Set the data source attributes to the values we got from the API
	for key, val := range attributes {
		if err := d.Set(key, val); err != nil {
			return diag.Errorf("Error setting %s for on-call schedule %s: %v", key, id, err)
		}
	}

	// Set the on-call schedule's ID in state
	d.SetId(*onCallSchedule.GetID())

	return diag.Diagnostics{}
}

func updateResourceFireHydrantOnCallSchedule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	teamID := d.Get("team_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Update on-call schedule: %s", id), map[string]interface{}{
		"id":      id,
		"team_id": teamID,
	})

	// Initialize updateRequest with basic fields
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	updateRequest := components.UpdateTeamOnCallSchedule{
		Name:        &name,
		Description: &description,
	}

	// Get slack_user_group_id if set
	if v, ok := d.GetOk("slack_user_group_id"); ok {
		slackUserGroupID := v.(string)
		updateRequest.SlackUserGroupID = &slackUserGroupID
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
			effectiveAtStr := effectiveAt.Format(time.RFC3339)
			updateRequest.EffectiveAt = &effectiveAtStr
			tflog.Debug(ctx, "Schedule update will take effect at: "+effectiveAtStr, map[string]interface{}{
				"effective_at": effectiveAtStr,
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
	updateRequest.MemberIds = memberIDs

	// Get strategy configuration
	if v, ok := d.GetOk("strategy"); ok {
		if strategies := v.([]interface{}); len(strategies) > 0 {
			strategy := strategies[0].(map[string]interface{})
			strategyType := strategy["type"].(string)
			handoffTime := strategy["handoff_time"].(string)
			handoffDay := strategy["handoff_day"].(string)

			updateRequest.Strategy = &components.UpdateTeamOnCallScheduleStrategy{
				Type:        components.UpdateTeamOnCallScheduleType(strategyType),
				HandoffTime: &handoffTime,
				HandoffDay:  (*components.UpdateTeamOnCallScheduleHandoffDay)(&handoffDay),
			}

			// Set shift duration for custom strategy
			if strategyType == "custom" {
				shiftDuration := strategy["shift_duration"].(string)
				updateRequest.Strategy.ShiftDuration = &shiftDuration
			}
		}
	}

	// Get restrictions
	restrictions := d.Get("restrictions").([]interface{})
	for _, r := range restrictions {
		restriction := r.(map[string]interface{})
		startDay := restriction["start_day"].(string)
		startTime := restriction["start_time"].(string)
		endDay := restriction["end_day"].(string)
		endTime := restriction["end_time"].(string)

		updateRequest.Restrictions = append(updateRequest.Restrictions, components.UpdateTeamOnCallScheduleRestriction{
			StartDay:  components.UpdateTeamOnCallScheduleStartDay(startDay),
			StartTime: startTime,
			EndDay:    components.UpdateTeamOnCallScheduleEndDay(endDay),
			EndTime:   endTime,
		})
	}

	// Update the on-call schedule
	_, err := client.Sdk.Signals.UpdateTeamOnCallSchedule(ctx, teamID, id, updateRequest)
	if err != nil {
		return diag.Errorf("Error updating on-call schedule %s: %v", id, err)
	}

	return readResourceFireHydrantOnCallSchedule(ctx, d, m)
}

func deleteResourceFireHydrantOnCallSchedule(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	teamID := d.Get("team_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Delete on-call schedule: %s", id), map[string]interface{}{
		"id":      id,
		"team_id": teamID,
	})

	// Delete the on-call schedule
	err := client.Sdk.Signals.DeleteTeamOnCallSchedule(ctx, teamID, id)
	if err != nil {
		return diag.Errorf("Error deleting on-call schedule %s: %v", id, err)
	}

	// Remove the on-call schedule's ID from state
	d.SetId("")

	return diag.Diagnostics{}
}

func importResourceFireHydrantOnCallSchedule(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	team_id, id, err := resourceFireHydrantOnCallScheduleParseId(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("team_id", team_id)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func resourceFireHydrantOnCallScheduleParseId(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected Team_ID:Schedule_ID", id)
	}

	return parts[0], parts[1], nil
}

func strategyToMapSDK(strategy components.NullableSignalsAPIOnCallStrategyEntity) []map[string]interface{} {
	m := map[string]interface{}{"type": *strategy.GetType()}
	if *strategy.GetType() == "custom" {
		if shiftDuration := strategy.GetShiftDuration(); shiftDuration != nil {
			m["shift_duration"] = *shiftDuration
		}
	} else {
		if handoffTime := strategy.GetHandoffTime(); handoffTime != nil {
			m["handoff_time"] = *handoffTime
		}
	}
	if *strategy.GetType() == "weekly" {
		if handoffDay := strategy.GetHandoffDay(); handoffDay != nil {
			m["handoff_day"] = *handoffDay
		}
	}
	return []map[string]interface{}{m}
}

func oncallRestrictionsFromDataSDK(d *schema.ResourceData) []components.CreateTeamOnCallScheduleRestriction {
	restrictions := make([]components.CreateTeamOnCallScheduleRestriction, 0)
	for _, restriction := range d.Get("restrictions").([]interface{}) {
		restrictionMap := restriction.(map[string]interface{})
		startDay := restrictionMap["start_day"].(string)
		startTime := restrictionMap["start_time"].(string)
		endDay := restrictionMap["end_day"].(string)
		endTime := restrictionMap["end_time"].(string)

		restrictions = append(restrictions, components.CreateTeamOnCallScheduleRestriction{
			StartDay:  components.CreateTeamOnCallScheduleStartDay(startDay),
			StartTime: startTime,
			EndDay:    components.CreateTeamOnCallScheduleEndDay(endDay),
			EndTime:   endTime,
		})
	}
	return restrictions
}

func restrictionsToDataSDK(restrictions []components.SignalsAPIOnCallRestrictionEntity) []map[string]interface{} {
	restrictionMaps := make([]map[string]interface{}, 0)
	for _, restriction := range restrictions {
		restrictionMaps = append(restrictionMaps, map[string]interface{}{
			"start_day":  *restriction.GetStartDay(),
			"start_time": *restriction.GetStartTime(),
			"end_day":    *restriction.GetEndDay(),
			"end_time":   *restriction.GetEndTime(),
		})
	}
	return restrictionMaps
}
