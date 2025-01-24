---
page_title: "FireHydrant Data Source: firehydrant_schedule"
---

# firehydrant_on_call_schedule Data Source

Use this data source to get information on a Signals on-call schedule.

## Example Usage

Basic usage:
```hcl
data "firehydrant_on_call_schedule" "my-oncall-schedule" {
  id      = "id-for-my-schedule"
  team_id = "id-for-my-team"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The name of the oncall schedule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the schedule.
* `name` - The name of the schedule.
* `description` - The description of the schedule.
* `time_zone` - The time zone of the schedule.
* `slack_user_group_id` - If present, the ID of the Slack user group associated with the schedule.
