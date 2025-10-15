package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/firehydrant-go-sdk/models/operations"
	"github.com/firehydrant/firehydrant-go-sdk/models/sdkerrors"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLifecycleMilestone() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceLifecycleMilestone,
		ReadContext:   readResourceLifecycleMilestone,
		UpdateContext: updateResourceLifecycleMilestone,
		DeleteContext: deleteResourceLifecycleMilestone,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"phase_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"position": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"auto_assign_timestamp_on_create": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func createResourceLifecycleMilestone(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	name := d.Get("name").(string)
	phase_id := d.Get("phase_id").(string)
	slug := d.Get("slug").(string)
	position := d.Get("position").(int)

	var assign_timestamp operations.CreateLifecycleMilestoneAutoAssignTimestampOnCreate
	desired_assign_timestamp := d.Get("auto_assign_timestamp_on_create").(string)
	switch desired_assign_timestamp {
	case "always_set_on_create":
		fallthrough
	case "only_set_on_manual_create":
		fallthrough
	case "never_set_on_create":
		assign_timestamp = operations.CreateLifecycleMilestoneAutoAssignTimestampOnCreate(desired_assign_timestamp)
	default:
		return diag.Errorf("invalid value for auto_assign_timestamp_on_create: %v", desired_assign_timestamp)
	}

	request := operations.CreateLifecycleMilestoneRequest{
		Name:                        name,
		Description:                 d.Get("description").(string),
		PhaseID:                     phase_id,
		Slug:                        &slug,
		Position:                    &position,
		AutoAssignTimestampOnCreate: &assign_timestamp,
	}

	tflog.Debug(ctx, "Create new Lifecycle Milestone")
	response, err := client.Sdk.IncidentSettings.CreateLifecycleMilestone(ctx, request)
	if err != nil {
		return diag.Errorf("Error creating new Lifecycle Milestone: %v", err)
	}

	// we get back the entire list of phases and milestones, so we need to dig through it to get the id of the one we just created
	milestoneID := ""
	for _, phase := range response.Data {
		if *phase.ID == phase_id {
			for _, milestone := range phase.Milestones {
				if *milestone.Name == name {
					milestoneID = *milestone.ID
				}
			}
		}
	}
	if milestoneID == "" {
		return diag.Errorf("Lifecycle Milestone %v not found in response", milestoneID)
	}
	d.SetId(milestoneID)

	return readResourceIncidentType(ctx, d, m)
}

func readResourceLifecycleMilestone(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read lifecycle milestone: %s", id), map[string]interface{}{
		"id": id,
	})

	// There is no call to get a specific milestone.  So we get the list of phases and milestones and search through them for the id of the one we want.
	// We're also pulling the phase ID from here as we search.
	response, err := client.Sdk.IncidentSettings.ListLifecyclePhases(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var desired_milestone *components.LifecyclesMilestoneEntity
	phase_id := ""
	for _, phase := range response.Data {
		for _, milestone := range phase.Milestones {
			if *milestone.ID == id {
				desired_milestone = &milestone
				phase_id = *phase.ID
			}
		}
	}
	if desired_milestone == nil {
		return diag.Errorf("Lifecycle Milestone %v not found in response", id)
	}

	attributes := map[string]interface{}{
		"name":                            *desired_milestone.Name,
		"description":                     *desired_milestone.Description,
		"slug":                            *desired_milestone.Slug,
		"phase_id":                        phase_id,
		"position":                        *desired_milestone.Position,
		"auto_assign_timestamp_on_create": *desired_milestone.AutoAssignTimestampOnCreate,
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for lifecycle_milestone %s: %v", key, id, err)
		}
	}

	d.SetId(*desired_milestone.ID)

	return diag.Diagnostics{}
}

func updateResourceLifecycleMilestone(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	slug := d.Get("slug").(string)
	position := d.Get("position").(int)

	var assign_timestamp operations.UpdateLifecycleMilestoneAutoAssignTimestampOnCreate
	desired_assign_timestamp := d.Get("auto_assign_timestamp_on_create").(string)
	switch desired_assign_timestamp {
	case "always_set_on_create":
		fallthrough
	case "only_set_on_manual_create":
		fallthrough
	case "never_set_on_create":
		assign_timestamp = operations.UpdateLifecycleMilestoneAutoAssignTimestampOnCreate(desired_assign_timestamp)
	default:
		return diag.Errorf("invalid value for auto_assign_timestamp_on_create: %v", desired_assign_timestamp)
	}

	request := operations.UpdateLifecycleMilestoneRequestBody{
		Name:                        &name,
		Description:                 &description,
		Slug:                        &slug,
		Position:                    &position,
		AutoAssignTimestampOnCreate: &assign_timestamp,
	}

	tflog.Debug(ctx, "Update Lifecycle Milestone")
	_, err := client.Sdk.IncidentSettings.UpdateLifecycleMilestone(ctx, id, &request)
	if err != nil {
		return diag.Errorf("Error updating lifecycle milestone: %v", err)
	}

	return readResourceLifecycleMilestone(ctx, d, m)
}

func deleteResourceLifecycleMilestone(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Delete lifecycle milestone: %s", id), map[string]interface{}{
		"ID": id,
	})
	err := client.Sdk.IncidentSettings.DeleteLifecycleMilestone(ctx, id)
	if err != nil {
		if err.(*sdkerrors.SDKError).StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting lifecycle milestone %s: %v", id, err)
	}

	return diag.Diagnostics{}
}
