package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSlackChannel() *schema.Resource {
	return &schema.Resource{
		Description: "The `firehydrant_slack_channel` data source allows access to the details of a Slack channel.",
		ReadContext: dataFireHydrantSlackChannelRead,
		Schema: map[string]*schema.Schema{
			"slack_channel_id": {
				Description: "ID of the channel, provided by Slack.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"id": {
				Description: "FireHydrant internal ID for the Slack channel.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataFireHydrantSlackChannelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the Slack channel
	channelID := d.Get("slack_channel_id").(string)
	slackChannel, err := firehydrantAPIClient.SlackChannels().Get(ctx, channelID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the ID
	d.SetId(slackChannel.ID)

	return diag.Diagnostics{}
}
