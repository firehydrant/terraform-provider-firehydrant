---
page_title: "FireHydrant Data Source: firehydrant_ingest_url"
---

# firehydrant_ingest_url Data Source

Use this data source to get the ingest URL for signals, optionally targeting a specific user, team, escalation policy, or on call schedule.  A transposer can also be used with or without any of the above targets

## Example Usage

Basic usage:
```hcl
data "firehydrant_ingest_url" "team_rocket" {
  team_id    = "team_rocket"
  transposer = "datadog"
}

resource "datadog_webhook" "firehydrant-myteam" {
  name           = "firehydrant-myteam"
  url            = firehydrant_ingest_url.team_rocket.url
  encode_as      = "json"
  # ...
}
```

## Argument Reference

The following arguments are supported:

> You must provide one of these options.

* `user_id` - (Optional) The ID for an existing firehydrant user.
* `team_id` - (Optional) The ID for an existing firehydrant team.
* `escalation_policy_id` - (Optional) The ID for an escalation policy belonging to a given team.  If this is used, `team_id` must be provided.
* `on_call_schedule_id` - (Optional) The ID for an on call schedule belonging to a given team.  If this is used, `team_id` must be provided.

* `transposer` - (Optional) A transposer to use when ingesting data.  See the list of valid transposers on the Event Sources tab of the Signals page in the FireHydrant UI.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `url` - The URL that will receive signals data, based on the provided attributes. 
