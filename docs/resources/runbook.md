---
page_title: "FireHydrant Resource: firehydrant_runbook"
subcategory: "Beta"
---

# firehydrant_runbook Resource

FireHydrant runbooks allow you to configure and automate your incident response process by defining a workflow
to be followed when an incident occurs. Runbooks actually initiate actions that are fundamental steps to
resolving an incident. Such actions might be creating a Slack channel, starting a Zoom meeting, or opening
a Jira ticket. Think of a runbook as an incident response playbook that runs (or is activated) when
an incident is declared. Using runbooks, you can equip your team with updated information and best practices
for mitigating incidents.

## Example Usage

Basic usage:
```hcl
resource "firehydrant_team" "example-owner-team" {
  name        = "my-example-owner-team"
  description = "This is an example team that owns a service"
}

data "firehydrant_runbook_action" "notify-channel-action" {
  slug             = "notify_channel"
  integration_slug = "slack"
  type             = "incident"
}

resource "firehydrant_runbook" "example-runbook" {
  name        = "example-runbook"
  type        = "incident"
  description = "This is an example runbook"
  owner_id    = firehydrant_team.example-owner-team.id
  
  steps {
    name    = "Notify Channel"
    action_id = data.firehydrant_runbook_action.notify-channel-action.id
    config = {
      "channels" = "#incidents"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the runbook.
* `type` - (Required) The type of the runbook. Valid values are 
  `incident`, `general`, `infrastructure`, and `incident_role`.
* `description` - (Optional) A description of the runbook.
* `owner_id` - (Optional) The ID of the team that owns this runbook.
* `severities` - (Optional) Severities to associate with the runbook.
* `steps` - (Optional) Steps to add to the runbook.

The `severities` block supports:

* `id` - (Required) The ID of the severity.

The `steps` block supports:

* `action_id` - (Required) The ID of the runbook action for the step.
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should be automatically execute.
* `config` - (Optional) Config block for the step.
* `delation_duration` - (Optional) How long this step should wait before executing.
* `repeats` - (Optional) Whether this step should repeat.
* `repeats_duration` - (Optional) How often this step should repeat.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the runbook.

The `steps` block contains:

* `step_id` - The ID of the step.

## Import

Runbooks can be imported; use `<RUNBOOK ID>` as the import ID. For example:

```shell
terraform import firehydrant_runbook.test 3638b647-b99c-5051-b715-eda2c912c42e
```
