---
page_title: "Step Configuration - Webex"
subcategory: "Runbooks"
---

# Webex

~> **Note** You must have the Webex integration installed in FireHydrant
for any Webex runbook steps to work properly.

The FireHydrant Webex integration allows FireHydrant users to create
Webex meetings.

### Available Steps

* [Create a Webex Meeting](#create-a-webex-meeting)

## Create a Webex Meeting

The [Webex **Create a Webex Meeting** step](https://support.firehydrant.com/hc/en-us/articles/6243302737044-Integrating-with-Webex-Meetings)
allows FireHydrant users to create Webex meetings.

### Create a Webex - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "webex_create_meeting" {
  integration_slug = "webex"
  slug             = "create_meeting"
}

resource "firehydrant_runbook" "webex_create_meeting_runbook" {
  name = "webex-create-meeting-runbook"

  steps {
    name      = "Create a Webex Meeting"
    action_id = data.firehydrant_runbook_action.webex_create_meeting.id

    config = jsonencode({
      title  = "[{{incident.severity}}] {{incident.name}}"
      agenda = "Incident Description: {{ incident.description }}"

      enable_join_before_host = {
        label = "No"
        value = "false"
      }
      enable_auto_record_meeting = {
        label = "Yes"
        value = "true"
      }
    })

    automatic = true
  }
}
```

### Create a Webex - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `config` - (Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.
* `rule` - (Optional) JSON string representing the rule configuration for the runbook step.
  For more information on the conditional logic used in `rule`, see the
  [Runbooks - Conditional Logic](./runbooks_conditional_logic.md) documentation.
  The step will default to running manually if `rule` is not specified and `automatic` and `repeats` are both `false`.

The `config` block supports:

* `enable_join_before_host` - (Required) Whether to allow others to join the Webex meeting before the host.
* `enabled_auto_record_meeting` - (Required) Whether to record the Webex meeting.
* `agenda` - (Optional) The agenda that will be included in the Webex meeting.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
* `title` - (Optional) The title that will be included in the Webex meeting.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.

The `enable_join_before_host` block supports:

* `label` - (Required) The name of the enable_join_before_host option.
  Valid values are `Yes` and `No`.
* `value` - (Required) The value of the enable_join_before_host option
  Valid values are `true` and `false`.

The `enable_auto_record_meeting` block supports:

* `label` - (Required) The name of the enable_auto_record_meeting option.
  Valid values are `Yes` and `No`.
* `value` - (Required) The value of the enable_auto_record_meeting option
  Valid values are `true` and `false`.
