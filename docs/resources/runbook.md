---
page_title: "FireHydrant Resource: firehydrant_runbook"
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
  description = "This is an example team that owns a runbook"
}

data "firehydrant_runbook_action" "notify-channel-action" {
  slug             = "notify_channel"
  integration_slug = "slack"
}

resource "firehydrant_runbook" "example-runbook" {
  name        = "example-runbook"
  description = "This is an example runbook"
  owner_id    = firehydrant_team.example-owner-team.id
  attachment_rule = jsonencode({
    logic = {
      eq = [
        {
          var = "incident_current_milestone"
        },
        {
          var = "usr.1"
        }
      ]
    }
    user_data = {
      "1" = {
        type  = "Milestone"
        value = "started"
        label = "Started"
      }
    }
  })

  steps {
    name      = "Notify Channel"
    action_id = data.firehydrant_runbook_action.notify-channel-action.id

    config = jsonencode({
      channels = "#incidents"
    })

    automatic        = false
    repeats          = true
    repeats_duration = "PT15M"
    rule = jsonencode({
      logic = {
        eq = [
          {
            var = "incident_current_milestone"
          },
          {
            var = "usr.1"
          }
        ]
      }
      user_data = {
        "1" = {
          type  = "Milestone"
          value = "resolved"
          label = "Resolved"
        }
      }
    })
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the runbook.
* `steps` - (Required) Steps to add to the runbook.
* `attachment_rule` - (Optional) JSON string representing the attachment rule configuration for the runbook.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
  For more information on the conditional logic used in `attachment_rule`, see the 
  [Runbooks - Conditional Logic](../guides/runbooks_conditional_logic.md) documentation.
  Defaults to attaching manually:
  ```hcl
  attachment_rule = jsonencode({
    logic = {
      manually = [
        {
          var = "when_invoked"
        }
      ]
    }
    user_data = {}
  })
  ```
* `description` - (Optional) A description of the runbook.
* `owner_id` - (Optional) The ID of the team that owns this runbook.
* `restricted` - (Optional) Only apply this runbook to private incidents.

The `steps` block supports:

Available attributes and whether they are available and required varies depending on the specific runbook step in question.
See [Runbook Steps Configuration documentation](../guides/runbooks_steps.md) for more detailed documentation on each step. 

* `action_id` - (Required) The ID of the runbook action for the step.
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should be automatically execute.
  Defaults to `false`.
* `config` - (Optional/Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode) 
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `repeats` - (Optional) Whether this step should repeat. Defaults to `false`.
  When this value is `true`, `repeats_duration` _must_ be provided.
* `repeats_duration` - (Optional) How often this step should repeat in ISO8601. 
  Example: PT10M [Format Spec](https://www.digi.com/resources/documentation/digidocs/90001437-13/reference/r_iso_8601_duration_format.htm)
  This value _must_ be provided if `repeats` is `true`. This value _must not_ be provided if `repeats` is `false`.
* `rule` - (Optional) JSON string representing the rule configuration for the runbook step.
  For more information on the conditional logic used in `rule`, see the
  [Runbooks - Conditional Logic](../guides/runbooks_conditional_logic.md) documentation.
  The step will default to running manually if `rule` is not specified and `automatic` and `repeats` are both `false`.

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
