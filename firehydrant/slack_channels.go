package firehydrant

import (
	"context"
	"fmt"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// SlackChannelsClient is an interface for interacting with Slack channels
type SlackChannelsClient interface {
	Get(ctx context.Context, slackID string) (*SlackChannelResponse, error)
}

// RESTSlackChannelsClient implements the SlackChannelClient interface
type RESTSlackChannelsClient struct {
	client *APIClient
}

var _ SlackChannelsClient = &RESTSlackChannelsClient{}

func (c *RESTSlackChannelsClient) restClient() *sling.Sling {
	return c.client.client()
}

// Get retrieves a Slack channel from FireHydrant using Slack ID. This is useful for looking up
// a Slack channel's internal ID.
func (c *RESTSlackChannelsClient) Get(ctx context.Context, slackID string) (*SlackChannelResponse, error) {
	channels := &SlackChannelsResponse{}
	apiError := &APIError{}
	response, err := c.restClient().Get("integrations/slack/channels?slack_channel_id="+slackID).Receive(channels, apiError)
	if err != nil {
		return nil, fmt.Errorf("could not get slack channel: %w", err)
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	if channels.Channels == nil || len(channels.Channels) == 0 {
		return nil, fmt.Errorf("no slack channel found with id '%s'", slackID)
	}
	if channelCount := len(channels.Channels); channelCount > 1 {
		// "at least" because it may paginate.
		tflog.Error(ctx, "found more than one Slack channel", map[string]interface{}{
			"id":    slackID,
			"found": channelCount,
		})
		for _, channel := range channels.Channels {
			tflog.Error(ctx, "found Slack channel", map[string]interface{}{
				"id":               channel.ID,
				"slack_channel_id": channel.SlackChannelID,
				"name":             channel.Name,
			})
		}
		return nil, fmt.Errorf("more than one Slack channel found: see Terraform logs for more information.")
	}

	tflog.Info(ctx, "found Slack channel", map[string]interface{}{
		"id":               channels.Channels[0].ID,
		"slack_channel_id": channels.Channels[0].SlackChannelID,
		"name":             channels.Channels[0].Name,
	})

	return channels.Channels[0], nil
}
