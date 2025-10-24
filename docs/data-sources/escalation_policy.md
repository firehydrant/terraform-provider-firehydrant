---
page_title: "FireHydrant Data Source: firehydrant_escalation_policy"
subcategory: "Signals"
---

# firehydrant_escalation_policy Data Source

Use this data source to get information on an escalation policy matching the given criteria.

## Example Usage

Basic usage:
```hcl
data "firehydrant_escalation_policy" "example-policy" {
  team_id = "3638b647-b99c-5051-b715-eda2c912c42e"
  name    = "Default Policy"
}
```

## Argument Reference

The following arguments are supported:

* `team_id` - (Required) The ID of the team that owns the escalation policy.
* `name` - (Required) The name of the escalation policy to retrieve.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the escalation policy.
* `description` - A description of the escalation policy.
* `default` - Whether this is the default escalation policy for the team.
* `repetitions` - The number of times to repeat the escalation policy.
* `step_strategy` - The strategy for handling steps in the escalation policy. Can be 'static' or 'dynamic_by_priority'.
* `step` - A list of steps in the escalation policy.
* `handoff_step` - A handoff step for the escalation policy.
* `notification_priority_policies` - Priority-specific policies for dynamic escalation policies.

The `step` block supports:

* `timeout` - The amount of time to wait before escalating to the next step. Must be in ISO 8601 duration format.
* `targets` - A list of targets for the step.
* `priorities` - A list of priorities this step applies to (for dynamic escalation policies).

The `targets` block supports:

* `id` - The ID of the target for this step.
* `type` - The type of target for this step. Must be one of `User`, `SlackChannel`, or `OnCallSchedule`.

The `handoff_step` block supports:

* `target_id` - The ID of the target to handoff to.
* `target_type` - The type of target to handoff to. Must be one of `Team` or `EscalationPolicy`.

The `notification_priority_policies` block supports:

* `priority` - The priority level. Must be one of `HIGH`, `MEDIUM`, or `LOW`.
* `repetitions` - The number of repetitions for this priority level.
* `handoff_step` - A handoff step for this priority level.

The `notification_priority_policies.handoff_step` block supports:

* `target_id` - The ID of the target to handoff to for this priority level.
* `target_type` - The type of target to handoff to for this priority level. Must be one of `Team` or `EscalationPolicy`.
