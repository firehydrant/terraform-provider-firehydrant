package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRotation() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantRotation,
		ReadContext:   readResourceFireHydrantRotation,
		UpdateContext: updateResourceFireHydrantRotation,
		DeleteContext: deleteResourceFireHydrantRotation,
		Importer: &schema.ResourceImporter{
			StateContext: importResourceFireHydrantRotation,
		},

		Schema: map[string]*schema.Schema{
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"schedule_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slack_user_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_slack_channel_notifications": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"prevent_shift_deletion": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"coverage_gap_notification_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"color": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"members": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
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
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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

func readResourceFireHydrantRotation(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	id := d.Id()
	teamID := d.Get("team_id").(string)
	scheduleID := d.Get("schedule_id").(string)
	tflog.Debug(ctx, "Read rotation", map[string]interface{}{
		"id":          id,
		"team_id":     teamID,
		"schedule_id": scheduleID,
	})

	rotation, err := firehydrantAPIClient.Rotations().Get(ctx, teamID, scheduleID, id)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			tflog.Debug(ctx, "Rotation %s does not exist", map[string]interface{}{
				"id":          id,
				"team_id":     teamID,
				"schedule_id": scheduleID,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading rotation %s: %v", id, err)
	}

	outputMemberIDs := []string{}
	memberIDs := rotation.Members
	for _, memberID := range memberIDs {
		if v := memberID.ID; v != "" {
			outputMemberIDs = append(outputMemberIDs, memberID.ID)
		}
	}

	attributes := map[string]interface{}{
		"name":                               rotation.Name,
		"time_zone":                          rotation.TimeZone,
		"description":                        rotation.Description,
		"enable_slack_channel_notifications": rotation.EnableSlackChannelNotifications,
		"prevent_shift_deletion":             rotation.PreventShiftDeletion,
		"color":                              rotation.Color,
		"members":                            outputMemberIDs,
		"strategy":                           rotationStrategyToMap(rotation.Strategy),
		"restrictions":                       rotationRestrictionsToData(rotation.Restrictions),
	}
	if rotation.SlackUserGroupID != "" {
		attributes["slack_user_group_id"] = rotation.SlackUserGroupID
	}
	if rotation.CoverageGapNotificationInterval != "" {
		attributes["coverage_gap_notification_interval"] = rotation.CoverageGapNotificationInterval
	}
	if rotation.StartTime != "" {
		attributes["start_time"] = rotation.StartTime
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for rotation %s: %v", key, id, err)
		}
	}

	d.SetId(rotation.ID)

	return diag.Diagnostics{}
}

func createResourceFireHydrantRotation(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	teamID := d.Get("team_id").(string)
	scheduleID := d.Get("schedule_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Create rotation: %s", teamID), map[string]interface{}{
		"team_id":     teamID,
		"schedule_id": scheduleID,
	})

	inputMemberIDs := d.Get("members").([]interface{})
	memberIDs := []firehydrant.RotationMember{}
	for _, memberID := range inputMemberIDs {
		if v, ok := memberID.(string); ok && v != "" {
			memberIDs = append(memberIDs, firehydrant.RotationMember{ID: v})
		}
	}

	// Gather values from API response
	rotation := firehydrant.CreateRotationRequest{
		Name:        d.Get("name").(string),
		TimeZone:    d.Get("time_zone").(string),
		Description: d.Get("description").(string),
		Members:     memberIDs,
		Strategy: firehydrant.RotationStrategy{
			Type:          d.Get("strategy.0.type").(string),
			HandoffTime:   d.Get("strategy.0.handoff_time").(string),
			HandoffDay:    d.Get("strategy.0.handoff_day").(string),
			ShiftDuration: d.Get("strategy.0.shift_duration").(string),
		},
		Restrictions: rotationRestrictionsFromData(d),
	}

	// Get slack_user_group_id if set and non-empty
	if v, ok := d.GetOk("slack_user_group_id"); ok && v.(string) != "" {
		rotation.SlackUserGroupID = v.(string)
	}
	if v, ok := d.GetOk("enable_slack_channel_notifications"); ok {
		rotation.EnableSlackChannelNotifications = v.(bool)
	}
	if v, ok := d.GetOk("prevent_shift_deletion"); ok {
		rotation.PreventShiftDeletion = v.(bool)
	}
	if v, ok := d.GetOk("coverage_gap_notification_interval"); ok && v.(string) != "" {
		rotation.CoverageGapNotificationInterval = v.(string)
	}
	if v, ok := d.GetOk("start_time"); ok && v.(string) != "" {
		rotation.StartTime = v.(string)
	}
	if v, ok := d.GetOk("color"); ok && v.(string) != "" {
		rotation.Color = v.(string)
	}

	if rotation.Strategy.Type != "" {
		isCustomStrategy := rotation.Strategy.Type == "custom"
		if isCustomStrategy {
			if rotation.Strategy.ShiftDuration == "" {
				return diag.Errorf("firehydrant_rotation.strategy.shift_duration is required when strategy type is 'custom'")
			}
			if rotation.StartTime == "" {
				return diag.Errorf("firehydrant_rotation.start_time is required when strategy type is 'custom'")
			}

			// Discard unused values to avoid ambiguity.
			rotation.Strategy.HandoffTime = ""
			rotation.Strategy.HandoffDay = ""
		} else {
			if rotation.Strategy.HandoffTime == "" {
				return diag.Errorf("firehydrant_rotation.strategy.handoff_time is required when strategy type is '%s'", rotation.Strategy.Type)
			}
			if rotation.Strategy.Type == "weekly" && rotation.Strategy.HandoffDay == "" {
				return diag.Errorf("firehydrant_rotation.strategy.handoff_day is required when strategy type is '%s'", rotation.Strategy.Type)
			}

			// Discard unused values to avoid ambiguity.
			rotation.Strategy.ShiftDuration = ""
			rotation.StartTime = ""
		}
	}

	// Create the rotation
	createdRotation, err := firehydrantAPIClient.Rotations().Create(ctx, teamID, scheduleID, rotation)
	if err != nil {
		return diag.Errorf("Error creating rotation %s: %v", teamID, err)
	}

	// Set the rotation's ID in state
	d.SetId(createdRotation.ID)

	return readResourceFireHydrantRotation(ctx, d, m)
}

func updateResourceFireHydrantRotation(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	id := d.Id()
	teamID := d.Get("team_id").(string)
	scheduleID := d.Get("schedule_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Update rotation: %s", id), map[string]interface{}{
		"id":          id,
		"team_id":     teamID,
		"schedule_id": scheduleID,
	})

	// Initialize updateRequest with basic fields
	updateRequest := firehydrant.UpdateRotationRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if v, ok := d.GetOk("slack_user_group_id"); ok {
		updateRequest.SlackUserGroupID = v.(string)
	}
	if v, ok := d.GetOk("enable_slack_channel_notifications"); ok {
		updateRequest.EnableSlackChannelNotifications = v.(bool)
	}
	if v, ok := d.GetOk("prevent_shift_deletion"); ok {
		updateRequest.PreventShiftDeletion = v.(bool)
	}
	if v, ok := d.GetOk("color"); ok && v.(string) != "" {
		updateRequest.Color = v.(string)
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
			tflog.Debug(ctx, "Rotation update will take effect at: "+updateRequest.EffectiveAt, map[string]interface{}{
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
	inputMemberIDs := d.Get("members").([]interface{})
	members := []firehydrant.RotationMember{}
	for _, memberID := range inputMemberIDs {
		if v, ok := memberID.(string); ok && v != "" {
			members = append(members, firehydrant.RotationMember{ID: v})
		}
	}
	updateRequest.Members = members

	// Get strategy configuration
	if v, ok := d.GetOk("strategy"); ok {
		if strategies := v.([]interface{}); len(strategies) > 0 {
			strategy := strategies[0].(map[string]interface{})
			updateRequest.Strategy = &firehydrant.RotationStrategy{
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
		updateRequest.Restrictions = append(updateRequest.Restrictions, firehydrant.RotationRestriction{
			StartDay:  restriction["start_day"].(string),
			StartTime: restriction["start_time"].(string),
			EndDay:    restriction["end_day"].(string),
			EndTime:   restriction["end_time"].(string),
		})
	}

	_, err := firehydrantAPIClient.Rotations().Update(ctx, teamID, scheduleID, id, updateRequest)
	if err != nil {
		return diag.Errorf("Error updating rotation %s: %v", id, err)
	}

	return readResourceFireHydrantRotation(ctx, d, m)
}

func deleteResourceFireHydrantRotation(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	id := d.Id()
	teamID := d.Get("team_id").(string)
	scheduleID := d.Get("schedule_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Delete rotation: %s", id), map[string]interface{}{
		"id":          id,
		"team_id":     teamID,
		"schedule_id": scheduleID,
	})

	err := firehydrantAPIClient.Rotations().Delete(ctx, teamID, scheduleID, id)
	if err != nil {
		return diag.Errorf("Error deleting rotation %s: %v", id, err)
	}

	d.SetId("")

	return diag.Diagnostics{}
}

func importResourceFireHydrantRotation(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	team_id, schedule_id, id, err := resourceFireHydrantRotationParseId(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("team_id", team_id)
	d.Set("schedule_id", schedule_id)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func resourceFireHydrantRotationParseId(id string) (string, string, string, error) {
	parts := strings.SplitN(id, ":", 3)

	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("unexpected format of ID (%s), expected Team_ID:Schedule_ID:Rotation_ID", id)
	}

	return parts[0], parts[1], parts[2], nil
}

func rotationStrategyToMap(strategy firehydrant.RotationStrategy) []map[string]interface{} {
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

func rotationRestrictionsFromData(d *schema.ResourceData) []firehydrant.RotationRestriction {
	restrictions := make([]firehydrant.RotationRestriction, 0)
	for _, restriction := range d.Get("restrictions").([]interface{}) {
		restrictionMap := restriction.(map[string]interface{})
		restrictions = append(restrictions, firehydrant.RotationRestriction{
			StartDay:  restrictionMap["start_day"].(string),
			StartTime: restrictionMap["start_time"].(string),
			EndDay:    restrictionMap["end_day"].(string),
			EndTime:   restrictionMap["end_time"].(string),
		})
	}
	return restrictions
}

func rotationRestrictionsToData(restrictions []firehydrant.RotationRestriction) []map[string]interface{} {
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
