package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/firehydrant-go-sdk/models/sdkerrors"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCustomEventSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceCustomEventSource,
		ReadContext:   readResourceCustomEventSource,
		UpdateContext: updateResourceCustomEventSource,
		DeleteContext: deleteResourceCustomEventSource,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"javascript": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ingest_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func createResourceCustomEventSource(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	description := d.Get("description").(string)

	request := components.CreateSignalsEventSource{
		Name:        d.Get("name").(string),
		Slug:        d.Get("slug").(string),
		Description: &description,
		Javascript:  d.Get("javascript").(string),
	}

	tflog.Debug(ctx, "Create new Custom Event Source")
	response, err := client.Sdk.Signals.CreateSignalsEventSource(ctx, request)
	if err != nil {
		return diag.Errorf("Error creating new Custom Event Source: %v", err)
	}

	d.SetId(*response.Slug)

	return readResourceCustomEventSource(ctx, d, m)
}

func readResourceCustomEventSource(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	slug := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read custom event source: %s", slug), map[string]interface{}{
		"slug": slug,
	})

	response, err := client.Sdk.Signals.GetSignalsEventSource(ctx, slug)
	if err != nil {
		return diag.FromErr(err)
	}

	attributes := map[string]interface{}{
		"name":        *response.Name,
		"slug":        *response.Slug,
		"description": *response.Description,
		"javascript":  *response.Expression,
		"ingest_url":  *response.IngestURL,
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for custom_event_source %s: %v", key, slug, err)
		}
	}

	return diag.Diagnostics{}
}

func updateResourceCustomEventSource(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	description := d.Get("description").(string)
	slug := d.Get("slug").(string)

	//we don't have an update api call in the client, but the API says that we create via a PUT, so I'm guessing this is correct?
	request := components.CreateSignalsEventSource{
		Name:        d.Get("name").(string),
		Slug:        slug,
		Description: &description,
		Javascript:  d.Get("javascript").(string),
	}

	tflog.Debug(ctx, "Update Custom Event Source")
	response, err := client.Sdk.Signals.CreateSignalsEventSource(ctx, request)
	if err != nil {
		return diag.Errorf("Error updating Custom Event Source %s: %v", slug, err)
	}

	d.SetId(*response.Slug)

	return readResourceCustomEventSource(ctx, d, m)
}

func deleteResourceCustomEventSource(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	slug := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Delete custom event source: %s", slug), map[string]interface{}{
		"slug": slug,
	})
	err := client.Sdk.Signals.DeleteSignalsEventSource(ctx, slug)
	if err != nil {
		if err.(*sdkerrors.SDKError).StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting custom event source %s: %v", slug, err)
	}

	return diag.Diagnostics{}
}
