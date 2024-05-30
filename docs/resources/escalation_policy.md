---
page_title: "FireHydrant Resource: firehydrant_escalation_policy"
subcategory: "Signals"
---

# firehydrant_escalation_policy Resource

FireHydrant escalation policies are used to define the order in which teams are notified of an alert.

## Example Usage

Basic usage:
```hcl
data "firehydrant_user" "my-user" {
  email = "user@example.com"
}

data "firehydrant_user" "backup-user" {
  email = "backup-user@example.com"
}

resource "firehydrant_team" "example-team" {
  name        = "example-team"
  description = "This is an example team"

  memberships {
    user_id          = data.firehydrant_user.my-user.id
  }
}

resource "firehydrant_team" "backup-team" {
  name        = "backup-team"
  description = "The backup for the example team"

  memberships {
    user_id          = data.firehydrant_user.backup-user.id
  }
}

resource "firehydrant_escalation_policy" "default_policy" {
  name = "Default Policy"
  description = "This is an example escalation policy"
  team_id = firehydrant_team.example-team.id

  step {
    timeout     = "PT1M"

    targets {
      type = "OnCallSchedule"
      id   = firehydrant_on_call_schedule.primary.id
    }

    targets {
      type = "User"
      id   = data.firehydrant_user.my-user.id
    }
  }

  repetitions = 2

  handoff_step {
    target_type = "Team"
    target_id   = firehydrant_team.backup-team.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the escalation policy.
* `description` - (Optional) A description for the escalation policy.
* `team_id` - (Required) The ID of the team to associate the escalation policy with.
* `step` - (Required) A block to define a step in the escalation policy. (a minimum of one is required)
* `handoff_step` - (Optional) A block to define a handoff step in the escalation policy.
* `default` - (Optional) Set this to `true` to mark this as the default escalation policy for this team.
* `repetitions` - (Required) The number of times to repeat the escalation policy. Defaults to 0.

The `step` block supports:

* `timeout` - (Required) The amount of time to wait before escalating to the next step. Must be in ISO 8601 duration format.
* `targets` - (Required) A block to define a target for the step. (a minimum of one is required)

The `targets` block supports:

* `id` - (Required) The ID of the target for this step.
* `type` - (Required) The type of target for this step. Must be one of `User`, `SlackChannel`, or `OnCallSchedule`.

The `handoff_step` block supports:

* `target_id` - (Required) The ID of the target to handoff to.
* `target_type` - (Required) The type of target to handoff to. Must be one of `Team` or `EscalationPolicy`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the escalation policy.
