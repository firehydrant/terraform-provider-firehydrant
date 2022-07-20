---
page_title: "Step Configuration - Google Meet"
subcategory: "Runbooks"
---

# Google Meet

~> **Note** You must have the Google Meet integration installed in FireHydrant
for any Google Meet runbook steps to work properly.

The FireHydrant Google Meet integration allows FireHydrant users to create Google 
Meet meeting rooms.

### Available Steps

* [Create a Google Meet](#create-a-google-meet)

## Create a Google Meet

The [Google Meet **Create a Google Meet** step](https://support.firehydrant.com/hc/en-us/articles/360061049852-Integrating-with-Google-Meet)
allows FireHydrant users to create Google Meet meeting rooms.

### Create a Google Meet - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "google_meet_create_google_meet_link" {
  integration_slug = "google_meet"
  slug             = "create_google_meet_link"
}

resource "firehydrant_runbook" "runbook17" {
  name = "google-meet-runbook"

  steps {
    name      = "Create a Google Meet"
    action_id = data.firehydrant_runbook_action.google_meet_create_google_meet_link.id

    config = jsonencode({
      topic = "[{{incident.severity}}] {{incident.name}}"
    })
  }
}
```

### Create a Google Meet - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `config` - (Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.

The `config` block supports:

* `body_template` - (Optional) The topic that will be included in the Google Meet meeting
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
