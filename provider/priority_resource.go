package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePriority() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantPriority,
		UpdateContext: updateResourceFireHydrantPriority,
		ReadContext:   readResourceFireHydrantPriority,
		DeleteContext: deleteResourceFireHydrantPriority,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k string, oldValue string, newValue string, d *schema.ResourceData) bool {
					// Slug is case-insensitive, so don't show a diff if the string are the same when compared
					// in all lowercase
					if strings.ToLower(oldValue) == strings.ToLower(newValue) {
						return true
					}
					return false
				},
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.All(
						validation.StringLenBetween(0, 23),
						validation.StringMatch(regexp.MustCompile(`\A[[:alnum:]]+\z`), "must only include letters and numbers"),
					),
				),
			},

			// Optional
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func readResourceFireHydrantPriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the priority
	priorityID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read priority: %s", priorityID), map[string]interface{}{
		"id": priorityID,
	})
	priorityResponse, err := firehydrantAPIClient.Priorities().Get(ctx, priorityID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			tflog.Debug(ctx, fmt.Sprintf("Priority %s no longer exists", priorityID), map[string]interface{}{
				"id": priorityID,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading priority %s: %v", priorityID, err)
	}

	// Gather values from the API response
	attributes := map[string]interface{}{
		"slug":        priorityResponse.Slug,
		"default":     priorityResponse.Default,
		"description": priorityResponse.Description,
	}

	// Set the resource attributes to the values we got from the API
	for key, val := range attributes {
		if err := d.Set(key, val); err != nil {
			return diag.Errorf("Error setting %s for priority %s: %v", key, priorityID, err)
		}
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantPriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get attributes from config and construct the create request
	createRequest := firehydrant.CreatePriorityRequest{
		Slug:        d.Get("slug").(string),
		Default:     d.Get("default").(bool),
		Description: d.Get("description").(string),
	}

	// Create the new priority
	tflog.Debug(ctx, fmt.Sprintf("Create priority: %s", createRequest.Slug), map[string]interface{}{
		"slug": createRequest.Slug,
	})
	priorityResponse, err := firehydrantAPIClient.Priorities().Create(ctx, createRequest)
	if err != nil {
		return diag.Errorf("Error creating priority %s: %v", createRequest.Slug, err)
	}

	// Set the new priority's ID in state
	d.SetId(priorityResponse.Slug)

	// Update state with the latest information from the API
	return readResourceFireHydrantPriority(ctx, d, m)
}

func updateResourceFireHydrantPriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the update request
	updateRequest := firehydrant.UpdatePriorityRequest{
		Slug:        d.Get("slug").(string),
		Description: d.Get("description").(string),
		Default:     d.Get("default").(bool),
	}

	// Update the priority
	tflog.Debug(ctx, fmt.Sprintf("Update priority: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})
	_, err := firehydrantAPIClient.Priorities().Update(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating priority %s: %v", d.Id(), err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantPriority(ctx, d, m)
}

func deleteResourceFireHydrantPriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the priority
	priorityID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Delete priority: %s", priorityID), map[string]interface{}{
		"id": priorityID,
	})
	err := firehydrantAPIClient.Priorities().Delete(ctx, priorityID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			return nil
		}
		return diag.Errorf("Error deleting priority %s: %v", priorityID, err)
	}

	return diag.Diagnostics{}
}
