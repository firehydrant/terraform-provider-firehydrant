---
page_title: "FireHydrant Resource: firehydrant_rotation"
subcategory: "Signals"
---

# firehydrant_rotation Resource

FireHydrant rotations, along with on call schedules, are used to define who is on-call for a given time period.  Note that the schedule resource will create one rotation by default and additional rotations must be associated with a schedule id as well as a team id.

## Example Usage

Basic usage:
```hcl
data "firehydrant_user" "my-user-1" {
  email = "user1@example.com"
}

data "firehydrant_user" "my-user-2" {
  email = "user2@example.com"
}

data "firehydrant_team" "example-team" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

data "firehydrant_on_call_schedule" "example-schedule" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

resource "firehydrant_rotation" "new-rotation" {
  name        = "My New On-call Rotation"
  description = "This is an example on-call schedule rotation"
  team_id     = data.firehydrant_team.example-team.id
  schedule_id = data.firehydrant_on_call_schedule.example-schedule.id

  members {
    user_id = data.firehydrant_user.my-user-1.id
  }

  members {
    user_id = data.firehydrant_user.my-user-2.id
  }

  time_zone                          = "America/New_York"
  slack_user_group_id                = "S01JBG0RHUM"
  enable_slack_channel_notifications = true
  prevent_shift_deletion             = true
  coverage_gap_notification_interval = "P1H30M"

  strategy {
    type         = "weekly"
    handoff_time = "10:00:00"
    handoff_day  = "thursday"
  }

  restrictions {
    start_day = "monday"
    end_day = "monday"
    start_time = "10:00:00"
    end_time = "14:00:00"
  }

  restrictions {
    start_day = "tuesday"
    end_day = "tuesday"
    start_time = "12:00:00"
    end_time = "23:00:00"
  }

  # effective_at is required when updating rotation members
  # This will schedule the member changes to take effect at a future time
  # If not provided when updating members, an error will be returned
  # If set to athe current time or time in the past, the rotation will be effective immediately
  effective_at = "2024-12-25T10:00:00Z"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the rotation.
* `description` - (Optional) A description for the rotation.
* `team_id` - (Required) The ID of the team that the rotation belongs to.
* `schedule_id` - (Required) The ID of the on-call schedule that the rotation belongs to.
* `members` - (Optional) An ordered list of member objects that specify users on-call for the rotation. Each member object supports:
  * `user_id` - (Required) The ID of the user to add to the rotation. You can use the `firehydrant_user` data source to look up a user by email/name.
* `time_zone` - (Required) The time zone that the rotation is in.
* `slack_user_group_id` - (Optional) The ID of the Slack user group that the rotation is associated with.
* `enable_slack_channel_notifications` - (Optional, defaults to false) A boolean to define if FireHydrant should notify the team's Slack channel when handoffs occur.
* `prevent_shift_deletion` - (Optional, defaults to false) A boolean to define if FireHydrant should Prevent shifts from being deleted by users and leading to gaps in coverage.
* `coverage_gap_notification_interval` - (Optional) An [ISO8601 format](https://en.wikipedia.org/wiki/ISO_8601#Durations) (e.g. `PT8H`) duration string specifying that the team should be notified about gaps in coverage for the upcoming interval. Notifications are sent at 9am daily in the rotation's time zone via email and, if enabled, the team's Slack channel.
* `color` - (Optional) A hex color code that will be used to represent the rotation in FireHydrant's UI.
* `strategy` - (Required) A block to define the strategy for the rotation.
* `start_time` - (Optional) An ISO8601 time string specifying when the initial rotation should start. This value is only used if the rotation's strategy type is "custom".
* `restrictions` - (Optional) A block to define a restriction for the rotation.
* `effective_at` - (Optional) The date and time that the rotation becomes effective. Must be in RFC3339 format (e.g., `2024-01-15T10:00:00Z`). **Required when updating rotation members.** If not provided when updating members, an error will be returned. If set to a time in the past, the rotation will be effective immediately (the time will be automatically adjusted to the current time). This attribute is not stored in Terraform state.

The `strategy` block supports:

* `type` - (Required) The type of strategy to use for the rotation. Valid values are `weekly`, `daily`, or `custom`.
* `handoff_time` - (Required) The time of day that the rotation handoff occurs. Must be in `HH:MM:SS` format.
* `handoff_day` - (Required) The day of the week that the rotation handoff occurs. Valid values are `sunday`, `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, and `saturday`.
* `shift_duration` - (Optional) The duration of the on-call shift in [ISO8601 format](https://en.wikipedia.org/wiki/ISO_8601#Durations) (e.g. `PT8H`). Required for `custom` strategy.

The `restrictions` block supports:

* `start_day` - (Required) The day of the week that the restriction starts. Valid values are `sunday`, `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, and `saturday`.
* `end_day` - (Required) The day of the week that the restriction ends. Valid values are `sunday`, `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, and `saturday`.
* `start_time` - (Required) The time of day that the restriction starts. Must be in `HH:MM:SS` format.
* `end_time` - (Required) The time of day that the restriction ends. Must be in `HH:MM:SS` format.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the rotation.

## Import

Rotations can be imported; use `<TeamID>:<ScheduleID>:<RotationID>` as the import ID. For example:

```shell
terraform import firehydrant_rotation.example_rotation 3638b647-b99c-5051-b715-eda2c912c42e:12345678-90ab-cdef-1234-567890abcdef:3638b647-b99c-5051-b715-eda2c912c42e
```
