---
page_title: "Step Configuration - Microsoft Teams"
subcategory: "Runbooks"
---

# Microsoft Teams

~> **Note** You must have the Microsoft Teams integration installed in FireHydrant
for any Microsoft Teams runbook steps to work properly.

The FireHydrant Microsoft Teams integration allows users to interact with FireHydrant
through Microsoft Teams. This allows your engineers to stay in Microsoft Teams while 
still leveraging all the automation FireHydrant provides.

### Available Steps

* [Create Incident Channel](#create-incident-channel)
* [Notify Channel](#notify-channel)
* [Notify Channel With a Custom Message](#notify-channel-with-a-custom-message)
* [Notify Incident Channel With a Custom Message](#notify-incident-channel-with-a-custom-message)

## Create Incident Channel

The [Microsoft Teams **Create Incident Channel** step](https://support.firehydrant.com/hc/en-us/articles/360058202871)
allows FireHydrant users to automatically create a centralized channel for an incident. 

### Create Incident Channel - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "microsoft_teams_create_incident_channel" {
  integration_slug = "microsoft_teams"
  slug             = "create_incident_channel"
}

resource "firehydrant_runbook" "rmicrosoft_teams_create_incident_channel_runbook" {
  name = "microsoft-teams-create-incident-channel-runbook"

  steps {
    name      = "Create Microsoft Teams Incident Channel"
    action_id = data.firehydrant_runbook_action.microsoft_teams_create_incident_channel.id

    config = jsonencode({
      channel_name_format = "incident-{{ number }}"
    })

    automatic = true
  }
}
```

### Create Incident Channel - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.
* `config` - (Optional) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).

The `config` block supports:

* `channel_name_format` - (Optional) The format to use for the channel's name.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.

## Notify Channel

The [Microsoft Teams **Notify Channel** step](https://support.firehydrant.com/hc/en-us/articles/360057722292)
allows FireHydrant users to automatically notify different channels during an incident.

### Notify Channel - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "microsoft_teams_notify_channel" {
  integration_slug = "microsoft_teams"
  slug             = "notify_channel"
}

resource "firehydrant_runbook" "microsoft_teams_notify_channel_runbook" {
  name = "microsoft-teams-notify-channel-runbook"

  steps {
    name      = "Notify a Microsoft Teams Channel"
    action_id = data.firehydrant_runbook_action.microsoft_teams_notify_channel.id

    config = jsonencode({
      channels = "General, Incidents"
    })

    automatic        = true
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

## Notify Channel With a Custom Message

The [Microsoft Teams **Notify Channel With a Custom Message** step](https://support.firehydrant.com/hc/en-us/articles/360058202811)
allows FireHydrant users to automatically notify different channels with a custom message during an incident.

### Notify Channel With a Custom Message - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "microsoft_teams_notify_channel_custom" {
  integration_slug = "microsoft_teams"
  slug             = "notify_channel_custom_message"
}

resource "firehydrant_runbook" "microsoft_teams_notify_channel_custom_runbook" {
  name = "microsoft-teams-notify-channel-custom-runbook"

  steps {
    name      = "Notify a Microsoft Teams Channel with a Custom Message"
    action_id = data.firehydrant_runbook_action.microsoft_teams_notify_channel_custom.id

    config = jsonencode({
      channels = "General, Incidents"
      message  = "Please check-in with your current status on this {{ incident.severity }} incident\n\n```\n/firehydrant add note I'm calculating the power required by the flux capacitor\n```\n"
    })

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
* `message` - (Required) The custom message to send to the list of channels.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.

## Notify Incident Channel With a Custom Message

The [Microsoft Teams **Notify Incident Channel With a Custom Message** step](https://support.firehydrant.com/hc/en-us/articles/360058202811)
allows FireHydrant users to automatically notify the incident channel with a custom message.

### Notify Incident Channel With a Custom Message - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "microsoft_teams_create_incident_channel" {
  integration_slug = "microsoft_teams"
  slug             = "create_incident_channel"
}

data "firehydrant_runbook_action" "microsoft_teams_notify_incident_channel_custom" {
  integration_slug = "microsoft_teams"
  slug             = "notify_incident_channel_custom_message"
}

resource "firehydrant_runbook" "microsoft_teams_notify_incident_channel_custom_runbook" {
  name = "microsoft-teams-notify-incident-channel-custom-runbook"

  steps {
    name      = "Create Microsoft Teams Incident Channel"
    action_id = data.firehydrant_runbook_action.microsoft_teams_create_incident_channel.id

    config = jsonencode({
      channel_name_format = "incident-{{ number }}"
    })

    automatic = true
  }

  steps {
    name      = "Notify the Microsoft Teams Incident Channel with a Custom Message"
    action_id = data.firehydrant_runbook_action.microsoft_teams_notify_incident_channel_custom.id

    config = jsonencode({
      message = "Please check-in with your current status on this {{ incident.severity }} incident\n\n```\n/firehydrant add note I'm calculating the power required by the flux capacitor\n```\n"
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
