package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func readResourceFireHydrantPriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)
	serviceResponse, err := firehydrantAPIClient.GetPriority(ctx, d.Get("slug").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	var ds diag.Diagnostics
	svc := map[string]interface{}{
		"slug":        serviceResponse.Slug,
		"description": serviceResponse.Description,
		"default":     serviceResponse.Default,
	}

	for key, val := range svc {
		if err := d.Set(key, val); err != nil {
			return diag.FromErr(err)
		}
	}

	return ds
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

	// Create the new service
	priorityResponse, err := firehydrantAPIClient.CreatePriority(ctx, createRequest)
	if err != nil {
		return diag.FromErr(err)
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
	_, err := firehydrantAPIClient.UpdatePriority(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantPriority(ctx, d, m)
}

func deleteResourceFireHydrantPriority(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the priority
	priorityID := d.Id()
	err := firehydrantAPIClient.DeletePriority(ctx, priorityID)
	if err != nil {
		_, isNotFoundError := err.(firehydrant.NotFound)
		if isNotFoundError {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
