package provider

import (
	"context"
	"errors"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceStatusUpdateTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantStatusUpdateTemplate,
		UpdateContext: updateResourceFireHydrantStatusUpdateTemplate,
		ReadContext:   readResourceFireHydrantStatusUpdateTemplate,
		DeleteContext: deleteResourceFireHydrantStatusUpdateTemplate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"body": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func readResourceFireHydrantStatusUpdateTemplate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	id := d.Id()
	tflog.Debug(ctx, "Read status update template", map[string]interface{}{"id": id})

	template, err := firehydrantAPIClient.StatusUpdateTemplates().Get(ctx, id)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			tflog.Debug(ctx, "Status update template %s does not exist", map[string]interface{}{"id": id})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading status update template %s: %v", id, err)
	}

	attributes := map[string]interface{}{
		"name": template.Name,
		"body": template.Body,
		"id":   template.ID,
	}
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for status update template %s: %v", key, id, err)
		}
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantStatusUpdateTemplate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	createRequest := firehydrant.CreateStatusUpdateTemplateRequest{
		Name: d.Get("name").(string),
		Body: d.Get("body").(string),
	}

	tflog.Debug(ctx, "Create status update template", map[string]interface{}{"id": d.Id()})
	statusUpdateTemplateResponse, err := firehydrantAPIClient.StatusUpdateTemplates().Create(ctx, createRequest)
	if err != nil {
		return diag.Errorf("Error creating status update template %s: %v", d.Id(), err)
	}

	d.SetId(statusUpdateTemplateResponse.ID)

	return readResourceFireHydrantStatusUpdateTemplate(ctx, d, m)
}

func updateResourceFireHydrantStatusUpdateTemplate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	updateRequest := firehydrant.UpdateStatusUpdateTemplateRequest{
		Name: d.Get("name").(string),
		Body: d.Get("body").(string),
	}

	tflog.Debug(ctx, "Update status update template", map[string]interface{}{"id": d.Id()})
	_, err := firehydrantAPIClient.StatusUpdateTemplates().Update(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating status update template %s: %v", d.Id(), err)
	}

	return readResourceFireHydrantStatusUpdateTemplate(ctx, d, m)
}

func deleteResourceFireHydrantStatusUpdateTemplate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	tflog.Debug(ctx, "Delete status update template", map[string]interface{}{"id": d.Id()})
	err := firehydrantAPIClient.StatusUpdateTemplates().Delete(ctx, d.Id())
	if err != nil {
		return diag.Errorf("Error deleting status update template %s: %v", d.Id(), err)
	}

	return diag.Diagnostics{}
}
