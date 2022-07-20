---
page_title: "Step Configuration - Slack"
subcategory: "Runbooks"
---

# Slack

~> **Note** You must have the Slack integration installed in FireHydrant
for any Slack runbook steps to work properly.

The FireHydrant Slack integration allows users to interact with FireHydrant
through Slack. This allows your engineers to stay in Slack while
still leveraging all the automation FireHydrant provides.

### Available Steps

* [Add a Bookmark to Incident Channel](#add-a-bookmark-to-incident-channel)
* [Archive Incident Channel](#archive-incident-channel)
* [Create Incident Channel](#create-incident-channel)
* [Notify Channel](#notify-channel)
* [Notify Channel With a Custom Message](#notify-channel-with-a-custom-message)
* [Notify Incident Channel With a Custom Message](#notify-incident-channel-with-a-custom-message)

## Add a Bookmark to Incident Channel

The Slack **Add a Bookmark to Incident Channel** step
allows FireHydrant users to automatically add a bookmark to the incident channel.

### Add a Bookmark to Incident Channel - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "slack_add_bookmark_to_incident_channel" {
  integration_slug = "slack"
  slug             = "add_bookmark_to_incident_channel"
}

resource "firehydrant_runbook" "slack_add_bookmark_to_incident_channel_runbook" {
  name = "slack-add-bookmark-to-incident-channel-runbook"

  steps {
    name      = "Add Bookmark to Slack Incident Channel"
    action_id = data.firehydrant_runbook_action.slack_add_bookmark_to_incident_channel.id

    config = jsonencode({
      bookmark_link  = "https://example.com"
      bookmark_title = "Service Dashboard"
    })
    
    automatic = true
  }
}
```

### Add a Bookmark to Incident Channel - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `config` - (Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.

The `config` block supports:

* `bookmark_link` - (optional) The bookmark link.
* `bookmark_title` - (optional) A title for the bookmark.

## Archive Incident Channel

The Slack **Archive Incident Channel** step
allows FireHydrant users to automatically archive the incident channel
after an incident is over.

### Archive Incident Channel - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "slack_create_incident_channel" {
  integration_slug = "slack"
  slug             = "create_incident_channel"
}

data "firehydrant_runbook_action" "slack_archive_channel" {
  integration_slug = "slack"
  slug             = "archive_incident_channel"
}

resource "firehydrant_runbook" "slack_archive_channel_runbook" {
  name = "slack-archive-incident-channel-runbook"

  steps {
    name      = "Create a Slack Incident Channel"
    action_id = data.firehydrant_runbook_action.slack_create_incident_channel.id

    config = jsonencode({
      channel_name_format = "incident-{{ number }}"
      channel_visibility = {
        label = "Private"
        value = "private"
      }
    })

    automatic = true
  }

  steps {
    name      = "Archive Slack Incident Channel"
    action_id = data.firehydrant_runbook_action.slack_archive_channel.id
    automatic = false
  }
}
```

### Archive Incident Channel - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.

## Create Incident Channel

The [Slack **Create Incident Channel** step](https://support.firehydrant.com/hc/en-us/articles/360058202871)
allows FireHydrant users to automatically create a centralized channel for an incident.

### Create Incident Channel - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "slack_create_incident_channel" {
  integration_slug = "slack"
  slug             = "create_incident_channel"
}

resource "firehydrant_runbook" "slack_create_incident_channel_runbook" {
  name = "slack-create-incident-channel-runbook"

  steps {
    name      = "Create a Slack Incident Channel"
    action_id = data.firehydrant_runbook_action.slack_create_incident_channel.id

    config = jsonencode({
      channel_name_format = "incident-{{ number }}"
      channel_visibility = {
        label = "Private"
        value = "private"
      }
    })

    automatic = true
  }
}
```

### Create Incident Channel - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `config` - (Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.

The `config` block supports:

* `channel_visibility` - (Required) Whether the channel should be public or private.
* `channel_name_format` - (Optional) The format to use for the channel's name.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.

The `channel_visibility` block supports:

* `label` - (Required) The name of the channel_visibility option.
  Valid values are `Public` and `Private`.
* `value` - (Required) The value of the channel_visibility option
  Valid values are `public` and `private`.

## Notify Channel

The [Slack **Notify Channel** step](https://support.firehydrant.com/hc/en-us/articles/360057722292)
allows FireHydrant users to automatically notify different channels during an incident.

### Notify Channel - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "slack_notify_channel" {
  integration_slug = "slack"
  slug             = "notify_channel"
}

resource "firehydrant_runbook" "slack_notify_channel_runbook" {
  name = "slack-notify-channel-runbook"

  steps {
    name      = "Notify a Slack Channel"
    action_id = data.firehydrant_runbook_action.slack_notify_channel.id

    config = jsonencode({
      channels = "#general, #incidents"
    })

    repeats          = true
    repeats_duration = "PT15M"
  }
}
```

### Notify Channel - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `config` - (Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.
* `repeats` - (Optional) Whether this step should repeat. Defaults to `false`.
  When this value is `true`, `repeats_duration` _must_ be provided.
* `repeats_duration` - (Optional) How often this step should repeat in ISO8601.
  Example: PT10M [Format Spec](https://www.digi.com/resources/documentation/digidocs/90001437-13/reference/r_iso_8601_duration_format.htm)
  This value _must_ be provided if `repeats` is `true`. This value _must not_ be provided if `repeats` is `false`.

The `config` block supports:

* `channels` - (Required) A comma separated list of channels to notify.
  Each channel should start with the `#` character.

## Notify Channel With a Custom Message

The [Slack **Notify Channel With a Custom Message** step](https://support.firehydrant.com/hc/en-us/articles/360058202811)
allows FireHydrant users to automatically notify different channels with a custom message during an incident.

### Notify Channel With a Custom Message - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "slack_notify_channel_custom" {
  integration_slug = "slack"
  slug             = "notify_channel_custom_message"
}

resource "firehydrant_runbook" "slack_notify_channel_custom_runbook" {
  name = "slack-notify-channel-custom-runbook"

  steps {
    name      = "Notify a Slack Channel with a Custom Message"
    action_id = data.firehydrant_runbook_action.slack_notify_channel_custom.id

    config = jsonencode({
      channels = "#general, #incidents"
      message  = "This is a {{ incident.severity }} incident"
    })

    repeats   = false
    automatic = true
  }
}
```

### Notify Channel With a Custom Message - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `config` - (Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.
* `repeats` - (Optional) Whether this step should repeat. Defaults to `false`.
  When this value is `true`, `repeats_duration` _must_ be provided.
* `repeats_duration` - (Optional) How often this step should repeat in ISO8601.
  Example: PT10M [Format Spec](https://www.digi.com/resources/documentation/digidocs/90001437-13/reference/r_iso_8601_duration_format.htm)
  This value _must_ be provided if `repeats` is `true`. This value _must not_ be provided if `repeats` is `false`.

The `config` block supports:

* `channels` - (Required) A comma separated list of channels to notify.
  Each channel should start with the `#` character.
* `message` - (Required) The custom message to send to the list of channels.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.

## Notify Incident Channel With a Custom Message

The [Slack **Notify Incident Channel With a Custom Message** step](https://support.firehydrant.com/hc/en-us/articles/360058202811)
allows FireHydrant users to automatically notify the incident channel with a custom message.

### Notify Incident Channel With a Custom Message - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "slack_create_incident_channel" {
  integration_slug = "slack"
  slug             = "create_incident_channel"
}

data "firehydrant_runbook_action" "slack_notify_incident_channel_custom" {
  integration_slug = "slack"
  slug             = "notify_incident_channel_custom_message"
}

resource "firehydrant_runbook" "slack_notify_incident_channel_custom_runbook" {
  name = "slack-notify-incident-channel-custom-runbook"

  steps {
    name      = "Create a Slack Incident Channel"
    action_id = data.firehydrant_runbook_action.slack_create_incident_channel.id

    config = jsonencode({
      channel_name_format = "incident-{{ number }}"
      channel_visibility = {
        label = "Private"
        value = "private"
      }
    })

    automatic = true
  }

  steps {
    name      = "Notify Slack Incident Channel with a Custom Message"
    action_id = data.firehydrant_runbook_action.slack_notify_incident_channel_custom.id

    config = jsonencode({
      message = "This is a {{ incident.severity }} incident"
      action_button = {
        label = "Add a note"
        value = "new_note"
      }
    })

    automatic = true
  }
}
```

### Notify Incident Channel With a Custom Message - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `config` - (Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.
* `repeats` - (Optional) Whether this step should repeat. Defaults to `false`.
  When this value is `true`, `repeats_duration` _must_ be provided.
* `repeats_duration` - (Optional) How often this step should repeat in ISO8601.
  Example: PT10M [Format Spec](https://www.digi.com/resources/documentation/digidocs/90001437-13/reference/r_iso_8601_duration_format.htm)
  This value _must_ be provided if `repeats` is `true`. This value _must not_ be provided if `repeats` is `false`.

The `config` block supports:

* `message` - (Required) The custom message to send to the incident channel.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
* `action_button` - (Optional) An action button to include in the custom message.

The `action_button` block supports:

* `label` - (Required) The name of the action button.
  Valid values are:
    - `Assign team`
    - `Assign role`
    - `Post an update`
    - `Add a note`
    - `Edit incident details`
    - `View tasks`
    - `View service info`
    - `See who's on call`
    - `Page a service`
    - `Resolve incident`
* `value` - (Required) The slug of the action button.
  Valid values are:
    - `assign_team`
    - `assign_role`
    - `update_impact`
    - `new_note`
    - `edit_incident`
    - `view_all_tasks`
    - `service_info`
    - `on_call`
    - `page_service`
    - `resolve_incident`
