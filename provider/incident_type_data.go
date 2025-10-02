package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataIncidentType() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataIncidentType,
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

func readDataIncidentType(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
