package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/firehydrant-go-sdk/models/sdkerrors"
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
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the user to add to the rotation. You can use the firehydrant_user data source to look up a user by email/name. Leave empty to create an unassigned slot in the rotation.",
						},
					},
				},
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
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	teamID := d.Get("team_id").(string)
	scheduleID := d.Get("schedule_id").(string)
	tflog.Debug(ctx, "Read rotation", map[string]interface{}{
		"id":          id,
		"team_id":     teamID,
		"schedule_id": scheduleID,
	})

	rotation, err := client.Sdk.Signals.GetOnCallScheduleRotation(ctx, id, teamID, scheduleID)
	if err != nil {
		if sdkErr, ok := err.(*sdkerrors.SDKError); ok && sdkErr.StatusCode == 404 {
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

	members := make([]map[string]interface{}, 0)
	rotationMembers := rotation.GetMembers()
	for _, member := range rotationMembers {
		if userID := member.GetID(); userID != nil && *userID != "" {
			members = append(members, map[string]interface{}{
				"user_id": *userID,
			})
		} else {
			// Include unassigned slots as empty string to preserve rotation order
			members = append(members, map[string]interface{}{
				"user_id": "",
			})
		}
	}

	attributes := map[string]interface{}{
		"name":                               *rotation.GetName(),
		"time_zone":                          *rotation.GetTimeZone(),
		"enable_slack_channel_notifications": *rotation.GetEnableSlackChannelNotifications(),
		"prevent_shift_deletion":             *rotation.GetPreventShiftDeletion(),
		"members":                            members,
	}

	// Handle optional description field
	if description := rotation.GetDescription(); description != nil {
		attributes["description"] = *description
	}

	// Handle optional color field
	if color := rotation.GetColor(); color != nil {
		attributes["color"] = *color
	}

	// Handle optional slack_user_group_id field
	if slackUserGroupID := rotation.GetSlackUserGroupID(); slackUserGroupID != nil && *slackUserGroupID != "" {
		attributes["slack_user_group_id"] = *slackUserGroupID
	}

	// Handle optional coverage_gap_notification_interval field
	if coverageGapNotificationInterval := rotation.GetCoverageGapNotificationInterval(); coverageGapNotificationInterval != nil && *coverageGapNotificationInterval != "" {
		attributes["coverage_gap_notification_interval"] = *coverageGapNotificationInterval
	}

	// Note: start_time is not returned in the API response, it's only used during creation for custom strategies

	// Handle strategy
	if strategy := rotation.GetStrategy(); strategy != nil {
		attributes["strategy"] = rotationStrategyToMapSDK(*strategy)
	}

	// Handle restrictions
	if restrictions := rotation.GetRestrictions(); restrictions != nil {
		attributes["restrictions"] = rotationRestrictionsToDataSDK(restrictions)
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for rotation %s: %v", key, id, err)
		}
	}

	d.SetId(*rotation.GetID())

	return diag.Diagnostics{}
}

func createResourceFireHydrantRotation(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	teamID := d.Get("team_id").(string)
	scheduleID := d.Get("schedule_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Create rotation: %s", teamID), map[string]interface{}{
		"team_id":     teamID,
		"schedule_id": scheduleID,
	})

	inputMembers := d.Get("members").([]interface{})
	members := []components.CreateOnCallScheduleRotationMember{}
	for _, member := range inputMembers {
		var userID *string

		if member == nil {
			// Member block is nil - treat as unassigned slot
			members = append(members, components.CreateOnCallScheduleRotationMember{UserID: nil})
			continue
		}

		memberMap, ok := member.(map[string]interface{})
		if !ok {
			return diag.Errorf("Invalid member format: expected map, got %T", member)
		}

		// user_id field is required in member block
		userIDVal, exists := memberMap["user_id"]
		if !exists {
			return diag.Errorf("member block must include user_id field (use empty string for unassigned slots)")
		}

		// Check if user_id is non-empty
		if userIDStr, ok := userIDVal.(string); ok && userIDStr != "" {
			userID = &userIDStr
		}
		// If userID is empty string or not a string, userID stays nil (unassigned slot)

		members = append(members, components.CreateOnCallScheduleRotationMember{UserID: userID})
	}

	// Gather values from schema
	name := d.Get("name").(string)
	timeZone := d.Get("time_zone").(string)
	strategyType := d.Get("strategy.0.type").(string)
	handoffTime := d.Get("strategy.0.handoff_time").(string)
	handoffDay := d.Get("strategy.0.handoff_day").(string)
	shiftDuration := d.Get("strategy.0.shift_duration").(string)

	rotation := components.CreateOnCallScheduleRotation{
		Name:     name,
		TimeZone: timeZone,
		Members:  members,
		Strategy: components.CreateOnCallScheduleRotationStrategy{
			Type: components.CreateOnCallScheduleRotationType(strategyType),
		},
		Restrictions: rotationRestrictionsFromDataSDK(d),
	}

	// Handle optional description field
	if desc := d.Get("description").(string); desc != "" {
		rotation.Description = &desc
	}

	// Handle optional slack_user_group_id field
	if v, ok := d.GetOk("slack_user_group_id"); ok && v.(string) != "" {
		slackUserGroupID := v.(string)
		rotation.SlackUserGroupID = &slackUserGroupID
	}

	// Handle optional enable_slack_channel_notifications field
	if v, ok := d.GetOk("enable_slack_channel_notifications"); ok {
		enableSlackChannelNotifications := v.(bool)
		rotation.EnableSlackChannelNotifications = &enableSlackChannelNotifications
	}

	// Handle optional prevent_shift_deletion field
	if v, ok := d.GetOk("prevent_shift_deletion"); ok {
		preventShiftDeletion := v.(bool)
		rotation.PreventShiftDeletion = &preventShiftDeletion
	}

	// Handle optional coverage_gap_notification_interval field
	if v, ok := d.GetOk("coverage_gap_notification_interval"); ok && v.(string) != "" {
		coverageGapNotificationInterval := v.(string)
		rotation.CoverageGapNotificationInterval = &coverageGapNotificationInterval
	}

	// Handle optional start_time field
	if v, ok := d.GetOk("start_time"); ok && v.(string) != "" {
		startTime := v.(string)
		rotation.StartTime = &startTime
	}

	// Handle optional color field
	if v, ok := d.GetOk("color"); ok && v.(string) != "" {
		color := v.(string)
		rotation.Color = &color
	}

	// Handle strategy fields
	if strategyType != "" {
		isCustomStrategy := strategyType == "custom"
		if isCustomStrategy {
			if shiftDuration == "" {
				return diag.Errorf("firehydrant_rotation.strategy.shift_duration is required when strategy type is 'custom'")
			}
			if rotation.StartTime == nil || *rotation.StartTime == "" {
				return diag.Errorf("firehydrant_rotation.start_time is required when strategy type is 'custom'")
			}

			shiftDurationPtr := &shiftDuration
			rotation.Strategy.ShiftDuration = shiftDurationPtr
		} else {
			if handoffTime == "" {
				return diag.Errorf("firehydrant_rotation.strategy.handoff_time is required when strategy type is '%s'", strategyType)
			}
			if strategyType == "weekly" && handoffDay == "" {
				return diag.Errorf("firehydrant_rotation.strategy.handoff_day is required when strategy type is '%s'", strategyType)
			}

			handoffTimePtr := &handoffTime
			rotation.Strategy.HandoffTime = handoffTimePtr
			if handoffDay != "" {
				handoffDayPtr := components.CreateOnCallScheduleRotationHandoffDay(handoffDay)
				rotation.Strategy.HandoffDay = &handoffDayPtr
			}
		}
	}

	// Create the rotation
	createdRotation, err := client.Sdk.Signals.CreateOnCallScheduleRotation(ctx, teamID, scheduleID, rotation)
	if err != nil {
		return diag.Errorf("Error creating rotation %s: %v", teamID, err)
	}

	// Set the rotation's ID in state
	d.SetId(*createdRotation.GetID())

	return readResourceFireHydrantRotation(ctx, d, m)
}

func updateResourceFireHydrantRotation(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	teamID := d.Get("team_id").(string)
	scheduleID := d.Get("schedule_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Update rotation: %s", id), map[string]interface{}{
		"id":          id,
		"team_id":     teamID,
		"schedule_id": scheduleID,
	})

	// Initialize updateRequest
	updateRequest := components.UpdateOnCallScheduleRotation{}

	// Set name
	name := d.Get("name").(string)
	updateRequest.Name = &name

	// Handle optional description field
	if desc := d.Get("description").(string); desc != "" {
		updateRequest.Description = &desc
	}

	// Handle optional slack_user_group_id field
	if v, ok := d.GetOk("slack_user_group_id"); ok {
		slackUserGroupID := v.(string)
		updateRequest.SlackUserGroupID = &slackUserGroupID
	}

	// Handle optional enable_slack_channel_notifications field
	if v, ok := d.GetOk("enable_slack_channel_notifications"); ok {
		enableSlackChannelNotifications := v.(bool)
		updateRequest.EnableSlackChannelNotifications = &enableSlackChannelNotifications
	}

	// Handle optional prevent_shift_deletion field
	if v, ok := d.GetOk("prevent_shift_deletion"); ok {
		preventShiftDeletion := v.(bool)
		updateRequest.PreventShiftDeletion = &preventShiftDeletion
	}

	// Handle optional color field
	if v, ok := d.GetOk("color"); ok && v.(string) != "" {
		color := v.(string)
		updateRequest.Color = &color
	}

	// Handle effective_at - always set it to ensure API gets a valid timestamp
	if raw := d.GetRawConfig().GetAttr("effective_at"); !raw.IsNull() {
		effectiveAtStr := raw.AsString()
		if effectiveAtStr != "" {
			// Validate the timestamp format
			effectiveAt, err := time.Parse(time.RFC3339, effectiveAtStr)
			if err != nil {
				return diag.FromErr(err)
			}
			// If it's in the past, use current time instead
			if effectiveAt.After(time.Now()) {
				// Send the timestamp as-is to the API
				updateRequest.EffectiveAt = &effectiveAtStr
				tflog.Debug(ctx, "Rotation update will take effect at: "+effectiveAtStr, map[string]interface{}{
					"effective_at": effectiveAtStr,
				})
			} else {
				// effective_at is in the past, update to now
				effectiveAtFormatted := time.Now().Format(time.RFC3339)
				updateRequest.EffectiveAt = &effectiveAtFormatted
				tflog.Info(ctx, "Provided effective_at is in the past, update will take effect immediately", map[string]interface{}{
					"provided_effective_at": effectiveAtStr,
					"effective_at":          effectiveAtFormatted,
				})
			}
		} else {
			// If effective_at is provided but empty, use current time
			now := time.Now()
			effectiveAtFormatted := now.Format(time.RFC3339)
			updateRequest.EffectiveAt = &effectiveAtFormatted
			tflog.Debug(ctx, "effective_at is empty, using current time for immediate effect", map[string]interface{}{
				"current_time": effectiveAtFormatted,
			})
		}
	} else {
		// If effective_at is not provided at all, use current time for immediate effect
		now := time.Now()
		effectiveAtFormatted := now.Format(time.RFC3339)
		updateRequest.EffectiveAt = &effectiveAtFormatted
		tflog.Debug(ctx, "effective_at not provided, using current time for immediate effect", map[string]interface{}{
			"current_time": effectiveAtFormatted,
		})
	}

	inputMembers := d.Get("members").([]interface{})
	members := []components.UpdateOnCallScheduleRotationMember{}
	for _, member := range inputMembers {
		var userID *string

		if member == nil {
			// Member block is nil - treat as unassigned slot
			members = append(members, components.UpdateOnCallScheduleRotationMember{UserID: nil})
			continue
		}

		memberMap, ok := member.(map[string]interface{})
		if !ok {
			return diag.Errorf("Invalid member format: expected map, got %T", member)
		}

		// user_id field is required in member block
		userIDVal, exists := memberMap["user_id"]
		if !exists {
			return diag.Errorf("member block must include user_id field (use empty string for unassigned slots)")
		}

		// Check if user_id is non-empty
		if userIDStr, ok := userIDVal.(string); ok && userIDStr != "" {
			userID = &userIDStr
		}
		// If userID is empty string or not a string, userID stays nil (unassigned slot)

		members = append(members, components.UpdateOnCallScheduleRotationMember{UserID: userID})
	}
	// Always set members, even if empty, to allow clearing members
	updateRequest.Members = members

	// Get strategy configuration
	if v, ok := d.GetOk("strategy"); ok {
		if strategies := v.([]interface{}); len(strategies) > 0 {
			strategy := strategies[0].(map[string]interface{})
			strategyType := strategy["type"].(string)
			handoffTime := strategy["handoff_time"].(string)
			handoffDay := strategy["handoff_day"].(string)
			shiftDuration := strategy["shift_duration"].(string)

			updateStrategy := &components.UpdateOnCallScheduleRotationStrategy{
				Type: components.UpdateOnCallScheduleRotationType(strategyType),
			}

			if strategyType == "custom" {
				if shiftDuration != "" {
					updateStrategy.ShiftDuration = &shiftDuration
				}
			} else {
				if handoffTime != "" {
					updateStrategy.HandoffTime = &handoffTime
				}
				if strategyType == "weekly" && handoffDay != "" {
					handoffDayPtr := components.UpdateOnCallScheduleRotationHandoffDay(handoffDay)
					updateStrategy.HandoffDay = &handoffDayPtr
				}
			}

			updateRequest.Strategy = updateStrategy
		}
	}

	// Get restrictions - always set this field, even if empty, to allow clearing restrictions
	restrictions := d.Get("restrictions").([]interface{})
	updateRequest.Restrictions = make([]components.UpdateOnCallScheduleRotationRestriction, 0, len(restrictions))
	for _, r := range restrictions {
		restriction := r.(map[string]interface{})
		startDay := restriction["start_day"].(string)
		startTime := restriction["start_time"].(string)
		endDay := restriction["end_day"].(string)
		endTime := restriction["end_time"].(string)

		updateRequest.Restrictions = append(updateRequest.Restrictions, components.UpdateOnCallScheduleRotationRestriction{
			StartDay:  components.UpdateOnCallScheduleRotationStartDay(startDay),
			StartTime: startTime,
			EndDay:    components.UpdateOnCallScheduleRotationEndDay(endDay),
			EndTime:   endTime,
		})
	}

	_, err := client.Sdk.Signals.UpdateOnCallScheduleRotation(ctx, id, teamID, scheduleID, updateRequest)
	if err != nil {
		return diag.Errorf("Error updating rotation %s: %v", id, err)
	}

	return readResourceFireHydrantRotation(ctx, d, m)
}

func deleteResourceFireHydrantRotation(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	teamID := d.Get("team_id").(string)
	scheduleID := d.Get("schedule_id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Delete rotation: %s", id), map[string]interface{}{
		"id":          id,
		"team_id":     teamID,
		"schedule_id": scheduleID,
	})

	err := client.Sdk.Signals.DeleteOnCallScheduleRotation(ctx, id, teamID, scheduleID)
	if err != nil {
		// If the resource is already deleted (404), treat as success
		if sdkErr, ok := err.(*sdkerrors.SDKError); ok && sdkErr.StatusCode == 404 {
			return nil
		}
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

func rotationStrategyToMapSDK(strategy components.NullableSignalsAPIOnCallStrategyEntity) []map[string]interface{} {
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

func rotationRestrictionsFromDataSDK(d *schema.ResourceData) []components.CreateOnCallScheduleRotationRestriction {
	restrictions := make([]components.CreateOnCallScheduleRotationRestriction, 0)
	for _, restriction := range d.Get("restrictions").([]interface{}) {
		restrictionMap := restriction.(map[string]interface{})
		startDay := restrictionMap["start_day"].(string)
		startTime := restrictionMap["start_time"].(string)
		endDay := restrictionMap["end_day"].(string)
		endTime := restrictionMap["end_time"].(string)

		restrictions = append(restrictions, components.CreateOnCallScheduleRotationRestriction{
			StartDay:  components.CreateOnCallScheduleRotationStartDay(startDay),
			StartTime: startTime,
			EndDay:    components.CreateOnCallScheduleRotationEndDay(endDay),
			EndTime:   endTime,
		})
	}
	return restrictions
}

func rotationRestrictionsToDataSDK(restrictions []components.SignalsAPIOnCallRestrictionEntity) []map[string]interface{} {
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
