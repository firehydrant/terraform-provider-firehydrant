---
page_title: "FireHydrant Data Source: firehydrant_rotation"
---

# firehydrant_rotation Data Source

Use this data source to get information on a Signals on-call schedule rotation.

## Example Usage

Basic usage:
```hcl
data "firehydrant_rotation" "my-rotation" {
  id      = "id-for-my-rotation"
  team_id = "id-for-my-team"
  schedule_id = "id-for-my-schedule"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The id of the rotation.
* `team_id` - (Required) The id of the team this rotation is associated with.
* `schedule_id` - (Required) The id of the schedule this rotation is associated with.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the rotation.
* `name` - The name of the rotation.
* `description` - The description of the rotation.
* `time_zone` - The time zone of the rotation.
* `slack_user_group_id` - If present, the ID of the Slack user group associated with the rotation.
