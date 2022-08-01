---
page_title: "Step Configuration - VictorOps"
subcategory: "Runbooks"
---

# VictorOps

~> **Note** You must have the VictorOps integration installed in FireHydrant
for any VictorOps runbook steps to work properly.

The FireHydrant VictorOps integration allows FireHydrant users to link incidents
in VictorOps to incidents in FireHydrant.

### Available Steps

* [Create VictorOps Incident](#create-victorops-incident)

## Create VictorOps Incident

The VictorOps **Create VictorOps Incident** step
allows FireHydrant users to create a new VictorOps incident from a FireHydrant incident.

### Create VictorOps Incident - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "victorops_create_new_incident" {
  integration_slug = "victorops"
  slug             = "create_new_victorops_incident"
}

resource "firehydrant_runbook" "victorops_create_new_incident" {
  name = "victorops-create-new-incident-runbook"

  steps {
    name      = "Create a VictorOps Incident"
    action_id = data.firehydrant_runbook_action.victorops_create_new_incident.id

    config = jsonencode({
      incident_title   = "{{incident.severity}} Incident: {{incident.name}}"
      incident_details = "There is a {{incident.severity}} incident.\n\nTriggered by {{author.name}} from FireHydrant. For more information see {{incident.incident_url}}."

      alert_default_policy = {
        label = "Yes"
        value = "true"
      }
      routing_key = {
        label = "platform-team"
        value = "platform-team"
      }
    })

    automatic = true
    repeats   = false
  }
}
```

### Create VictorOps Incident - Steps Argument Reference

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
* `rule` - (Optional) JSON string representing the rule configuration for the runbook step.
  For more information on the conditional logic used in `rule`, see the
  [Runbooks - Conditional Logic](./runbooks_conditional_logic.md) documentation.
  The step will default to running manually if `rule` is not specified and `automatic` and `repeats` are both `false`.

The `config` block supports:

* `alert_default_policy` - (Required) Whether to alert the default escalation policy if 
  there are no impacted services linked to VictorOps and no additional Routing Key is specified.
* `incident_details` - (Optional) A description of the incident for VictorOps.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
* `incident_title` - (Optional) The title of the VictorOps incident.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
* `routing_key` - (Optional) An additional VictorOps routing key to page.

The `alert_default_policy` block supports:

* `label` - (Required) The name of the alert_default_policy option.
  Valid values are `Yes` and `No`.
* `value` - (Required) The value of the alert_default_policy option
  Valid values are `true` and `false`.

The `routing_key` block supports:

* `label` - (Required) The name of the VictorOps routing key.
  This value should match the `value` attribute.
* `value` - (Required) The name of the VictorOps routing key.
  This value should match the `label` attribute.
