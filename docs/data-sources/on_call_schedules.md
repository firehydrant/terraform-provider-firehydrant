---
page_title: "FireHydrant Data Source: firehydrant_schedule"
---

# firehydrant_on_call_schedules Data Source

Use this data source to get information on Signals on-call schedules for a team.

## Example Usage

Basic usage:

```hcl
data "firehydrant_on_call_schedules" "my-oncall-schedules" {
  team_id = "id-for-my-team"
  
  # optional
  query = "primary"
}
```

## Argument Reference

The following arguments are supported:

* `team_id` - (Required) The Team ID for the on-call schedule.
* `query` - (Optional) A query string for searching through the list of on-call schedules.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `on_call_schedules` - All the schedules matching the criteria specified by `query`.

The `on_call_schedules` block contains:

* `id` - The ID of the schedule.
* `name` - The name of the schedule.
* `description` - The description of the schedule.
* `time_zone` - The time zone of the schedule.
* `slack_user_group_id` - If present, the ID of the Slack user group associated with the schedule.
