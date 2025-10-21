package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/firehydrant/firehydrant-go-sdk/models/operations"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOnCallSchedules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantOnCallSchedules,
		Schema: map[string]*schema.Schema{
			// Required
			"team_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Optional
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed
			"on_call_schedules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceOnCallSchedule(),
			},
		},
	}
}

func dataFireHydrantOnCallSchedules(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Get the on-call schedule
	teamID := d.Get("team_id").(string)
	query := d.Get("query").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read on-call schedules: %s", teamID), map[string]interface{}{
		"team_id": teamID,
		"query":   query,
	})
	request := operations.ListTeamOnCallSchedulesRequest{
		TeamID: teamID,
	}
	if query != "" {
		request.Query = &query
	}
	schedulesResponse, err := client.Sdk.Signals.ListTeamOnCallSchedules(ctx, request)

	if err != nil {
		return diag.Errorf("Error reading on-call schedules: %v", err)
	}

	// Set the data source attributes to the values we got from the API
	schedules := make([]interface{}, 0)
	for _, schedule := range schedulesResponse.GetData() {
		schedules = append(schedules, dataFireHydrantOnCallScheduleToAttributesMap(teamID, schedule))
	}
	if err = d.Set("on_call_schedules", schedules); err != nil {
		return diag.Errorf("Error setting on-call schedules: %v", err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diag.Diagnostics{}
}
