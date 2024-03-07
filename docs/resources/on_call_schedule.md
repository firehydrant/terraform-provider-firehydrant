---
page_title: "FireHydrant Resource: firehydrant_on_call_schedule"
subcategory: "Signals"
---

# firehydrant_on_call_schedule Resource

FireHydrant on-call schedules are used to define who is on-call for a given time period.

## Example Usage

Basic usage:
```hcl
data "firehydrant_user" "my-user" {
  email = "user@example.com"
}

resource "firehydrant_team" "example-team" {
  name        = "example-team"
  description = "This is an example team"

  memberships {
    user_id          = data.firehydrant_user.my-user.id
  }
}

resource "firehydrant_on_call_schedule" "primary" {
  name        = "Primary On-Call Schedule"
  description = "This is an example on-call schedule"
  team_id     = firehydrant_team.example-team.id

  member_ids = [
    data.firehydrant_user.my-user.id,
  ]

  time_zone   = "America/New_York"

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
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the on-call schedule.
* `description` - (Optional) A description for the on-call schedule.
* `team_id` - (Required) The ID of the team that the on-call schedule belongs to.
* `member_ids` - (Required) A list of user IDs that are on-call for the on-call schedule.
* `members` - (Deprecated) use `member_ids` instead.
* `time_zone` - (Required) The time zone that the on-call schedule is in.
* `strategy` - (Required) A block to define the strategy for the on-call schedule.
* `restrictions` - (Optional) A block to define a restriction for the on-call schedule.

The `strategy` block supports:

* `type` - (Required) The type of strategy to use for the on-call schedule. Valid values are `weekly` and `daily`.
* `handoff_time` - (Required) The time of day that the on-call schedule handoff occurs. Must be in `HH:MM:SS` format.
* `handoff_day` - (Required) The day of the week that the on-call schedule handoff occurs. Valid values are `sunday`, `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, and `saturday`.

The `restrictions` block supports:

* `start_day` - (Required) The day of the week that the restriction starts. Valid values are `sunday`, `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, and `saturday`.
* `end_day` - (Required) The day of the week that the restriction ends. Valid values are `sunday`, `monday`, `tuesday`, `wednesday`, `thursday`, `friday`, and `saturday`.
* `start_time` - (Required) The time of day that the restriction starts. Must be in `HH:MM:SS` format.
* `end_time` - (Required) The time of day that the restriction ends. Must be in `HH:MM:SS` format.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the on-call schedule.
