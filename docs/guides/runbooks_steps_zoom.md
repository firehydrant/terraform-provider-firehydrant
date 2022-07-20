---
page_title: "Step Configuration - Zoom"
subcategory: "Runbooks"
---

# Zoom

~> **Note** You must have the Zoom integration installed in FireHydrant
for any Zoom runbook steps to work properly.

The FireHydrant Zoom integration allows FireHydrant users to create
Zoom meetings.

### Available Steps

* [Create a Zoom Meeting](#create-a-zoom-meeting)

## Create a Zoom Meeting

The [Zoom **Create a Zoom Meeting** step](https://support.firehydrant.com/hc/en-us/articles/360058202271-Zoom-Integration)
allows FireHydrant users to create Zoom meetings.

### Create a Zoom - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "zoom_create_meeting" {
  integration_slug = "zoom"
  slug             = "create_meeting"
}

resource "firehydrant_runbook" "zoom_create_meeting_runbook" {
  name = "zoom-create-meeting-runbook"

  steps {
    name      = "Create Zoom Meeting"
    action_id = data.firehydrant_runbook_action.zoom_create_meeting.id

    config = jsonencode({
      agenda = "Incident Description: {{ incident.description }}"
      topic  = "[{{incident.severity}}] {{incident.name}}"

      record_meeting = {
        label = "Record to cloud"
        value = "cloud"
      }
    })

    automatic = true
  }
}
```

### Create a Zoom - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `config` - (Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.

The `config` block supports:

* `record_meeting` - (Required) Whether to record the Zoom meeting.
* `agenda` - (Optional) The agenda that will be included in the Zoom meeting.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
* `topic` - (Optional) The topic that will be included in the Zoom meeting.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.

The `record_meeting` block supports:

* `label` - (Required) The name of the record_meeting option.
  Valid values are `No Recording`, `Record to cloud`, and `Record to desktop`.
* `value` - (Required) The value of the record_meeting option
  Valid values are `none`, `cloud`, and `local`.
