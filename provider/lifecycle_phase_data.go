package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLifecyclePhase() *schema.Resource {
	return &schema.Resource{
		ReadContext: readDataLifecyclePhase,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func readDataLifecyclePhase(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	name := d.Get("name").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read lifecycle phase: %s", name), map[string]interface{}{
		"name": name,
	})

	response, err := client.Sdk.IncidentSettings.ListLifecyclePhases(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	id := ""
	var validPhases []string

	for _, phase := range response.Data {
		validPhases = append(validPhases, *phase.Name)
		if strings.EqualFold(name, *phase.Name) {
			id = *phase.ID
		}
	}
	if id == "" {
		return diag.Errorf("Lifecycle phase %s is invalid.  Valid lifecycle phases are %v", name, validPhases)
	}

	attributes := map[string]interface{}{
		"name": name,
		"id":   id,
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for lifecycle_phase %s: %v", key, id, err)
		}
	}

	d.SetId(id)

	return diag.Diagnostics{}
}
