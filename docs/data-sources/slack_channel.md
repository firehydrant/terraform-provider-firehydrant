---
page_title: "FireHydrant Data Source: firehydrant_slack_channel"
---

# firehydrant_slack_channel Data Source

Use this data source to pass Slack channel information to other resources.

## Example Usage

Basic usage:
```hcl
data "firehydrant_slack_channel" "team_rocket" {
  slack_channel_id   = "C1234567890"
  slack_channel_name = "#team-rocket"
}

resource "firehydrant_escalation_policy" "team_rocket" {
  # ...
  step {
    timeout = "PT5M"
    targets {
      type = "SlackChannel"
      id = data.firehydrant_slack_channel.team_rocket.id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

> You must provide one of these options and the `slack_channel_id` will take precedence if both are provided. 

* `slack_channel_id` - (Optional) Slack's channel ID.
* `slack_channel_name` - (Optional) Slack's channel name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The FireHydrant ID for the given Slack channel. 
