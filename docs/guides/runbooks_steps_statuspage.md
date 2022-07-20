---
page_title: "Step Configuration - Statuspage"
subcategory: "Runbooks"
---

# Statuspage

~> **Note** You must have the Statuspage integration installed in FireHydrant
for any Statuspage runbook steps to work properly.

The FireHydrant Statuspage integration allows FireHydrant users to integrate their 
public status page updates into their incident response process

### Available Steps

* [Create Statuspage Incident](#create-statuspage-incident)
* [Update Statuspage Incident](#update-statuspage-incident)

## Create Statuspage Incident

The [Statuspage **Create Statuspage Incident** step](https://support.firehydrant.com/hc/en-us/articles/360058202851-Create-an-incident-on-your-Atlassian-Statuspage)
allows FireHydrant users to create a new Statuspage incident from a FireHydrant incident.

### Create Statuspage Incident - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "statuspage_create_statuspage" {
  integration_slug = "statuspage"
  slug             = "create_statuspage"
}

resource "firehydrant_runbook" "statuspage_create_statuspage_runbook" {
  name = "statuspage-create-statuspage-runbook"

  steps {
    name      = "Creates a Statuspage.io Incident"
    action_id = data.firehydrant_runbook_action.statuspage_create_statuspage.id

    config = jsonencode({
      title   = "{{ incident.name }}"
      message = "{{ incident.description }}"

      connection_id = {
        label = "Acme, Inc Status Page"
        value = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
      }
    })

    automatic = true
  }
}
```

### Create Statuspage Incident - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `config` - (Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.

The `config` block supports:

* `connection_id` - (Required) The FireHydrant Statuspage connection representing the status page to publish to.
  Your FireHydrant Statuspage connections can be found at the 
  [List Statuspage connections](https://developers.firehydrant.io/docs/api/48069b4939db5-list-statuspage-connections) endpoint. 
* `message` - (Optional) A description of the incident for Statuspage.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
* `title` - (Optional) The title of the Statuspage incident.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.

The `connection_id` block supports:

* `label` - (Required) The page name of the FireHydrant Statuspage connection.
* `value` - (Required) The ID of the FireHydrant Statuspage connection.

## Update Statuspage Incident

The Statuspage **Update Statuspage Incident** step
allows FireHydrant users to update a Statuspage incident from a FireHydrant incident.

### Update Statuspage Incident - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "statuspage_create_statuspage" {
  integration_slug = "statuspage"
  slug             = "create_statuspage"
}

data "firehydrant_runbook_action" "statuspage_update_statuspage" {
  integration_slug = "statuspage"
  slug             = "update_statuspage"
}

resource "firehydrant_runbook" "statuspage_update_statuspage_runbook" {
  name = "statuspage-update-statuspage-runbook"

  steps {
    name      = "Creates a Statuspage.io Incident"
    action_id = data.firehydrant_runbook_action.statuspage_create_statuspage.id

    config = jsonencode({
      title   = "{{ incident.name }}"
      message = "{{ incident.description }}"

      connection_id = {
        label = "Acme, Inc Status Page"
        value = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
      }
    })
  }

  steps {
    name      = "Update a Statuspage.io Incident"
    action_id = data.firehydrant_runbook_action.statuspage_update_statuspage.id

    config = jsonencode({
      message = "This incident has been mitigated and we are continuing to monitor."
    })

    automatic        = true
    repeats          = true
    repeats_duration = "PT15M"
  }
}
```

### Update Statuspage Incident - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `config` - (Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `name` - (Required) The name of the step.
* `repeats` - (Optional) Whether this step should repeat. Defaults to `false`.
  When this value is `true`, `repeats_duration` _must_ be provided.
* `repeats_duration` - (Optional) How often this step should repeat in ISO8601.
  Example: PT10M [Format Spec](https://www.digi.com/resources/documentation/digidocs/90001437-13/reference/r_iso_8601_duration_format.htm)
  This value _must_ be provided if `repeats` is `true`. This value _must not_ be provided if `repeats` is `false`.

The `config` block supports:

* `message` - (Required) A status update message for the Statuspage incident.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
