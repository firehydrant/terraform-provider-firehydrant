package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
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
			StateContext: schema.ImportStatePassthroughContext,
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
							Required: true,
							ForceNew: true,
						},
						"handoff_day": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"color": {
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
			Type:        d.Get("strategy.0.type").(string),
			HandoffTime: d.Get("strategy.0.handoff_time").(string),
			HandoffDay:  d.Get("strategy.0.handoff_day").(string),
		},
		MemberIDs:    memberIDs,
		Restrictions: oncallRestrictionsFromData(d),
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

	onCallSchedule := firehydrant.UpdateOnCallScheduleRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

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
	onCallSchedule.MemberIDs = memberIDs

	tflog.Debug(ctx, "Updating on-call schedule properties", map[string]interface{}{
		"id":         id,
		"properties": fmt.Sprintf("%+v", onCallSchedule),
	})

	// Update the on-call schedule
	_, err := firehydrantAPIClient.OnCallSchedules().Update(ctx, teamID, id, onCallSchedule)
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
	return []map[string]interface{}{
		{
			"type":         strategy.Type,
			"handoff_time": strategy.HandoffTime,
			"handoff_day":  strategy.HandoffDay,
		},
	}
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
