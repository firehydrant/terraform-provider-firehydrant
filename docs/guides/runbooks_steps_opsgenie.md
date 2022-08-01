---
page_title: "Step Configuration - Opsgenie"
subcategory: "Runbooks"
---

# Opsgenie

~> **Note** You must have the Opsgenie integration installed in FireHydrant
for any Opsgenie runbook steps to work properly.

The FireHydrant Opsgenie integration allows FireHydrant users to link incidents 
in Opsgenie to incidents in FireHydrant. 

### Available Steps

* [Create Opsgenie Incident](#create-opsgenie-incident)

## Create Opsgenie Incident

The Opsgenie **Create Opsgenie Incident** step
allows FireHydrant users to create a new Opsgenie incident from a FireHydrant incident.

### Create Opsgenie Incident - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "opsgenie_create_new_incident" {
  integration_slug = "opsgenie"
  slug             = "create_new_opsgenie_incident"
}

resource "firehydrant_runbook" "opsgenie_create_new_incident_runbook" {
  name = "opsgenie-create-new-incident-runbook"

  steps {
    name      = "Create Opsgenie Incident"
    action_id = data.firehydrant_runbook_action.opsgenie_create_new_incident.id

    config = jsonencode({
      incident_title   = "{{incident.severity}} Incident: {{incident.name}}"
      incident_details = "There is a {{incident.severity}} incident.\n\nTriggered by {{author.name}} from FireHydrant. For more information see {{incident.incident_url}}."
      team_id = {
        label = "Platform Team"
        value = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
      }
    })

    automatic = true
  }
}
```

### Create Opsgenie Incident - Steps Argument Reference

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

* `team_id` - (Required) The Opsgenie team to assign the incident to.
* `incident_details` - (Optional) A description of the incident for Opsgenie.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
* `incident_title` - (Optional) The title of the Opsgenie incident.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.

The `team_id` block supports:

* `label` - (Required) The name of the Opsgenie team.
* `value` - (Required) The ID of the Opsgenie team.