---
page_title: "Step Configuration - PagerDuty"
subcategory: "Runbooks"
---

# PagerDuty

~> **Note** You must have the PagerDuty integration installed in FireHydrant
for any PagerDuty runbook steps to work properly.

The FireHydrant PagerDuty integration allows FireHydrant users to link incidents
in PagerDuty to incidents in FireHydrant.

### Available Steps

* [Create PagerDuty Incident](#create-pagerduty-incident)

## Create PagerDuty Incident

The [PagerDuty **Create PagerDuty Incident** step](https://support.firehydrant.com/hc/en-us/articles/360057722212-Starting-a-PagerDuty-Incident)
allows FireHydrant users to create a new PagerDuty incident from a FireHydrant incident.

### Create PagerDuty Incident - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "pagerduty_create_new_incident" {
  integration_slug = "pager_duty"
  slug             = "create_new_pager_duty_incident"
}

resource "firehydrant_runbook" "pagerduty_create_new_incident_runbook" {
  name = "pagerduty-create-new-incident-runbook"

  steps {
    name      = "Create PagerDuty Incident"
    action_id = data.firehydrant_runbook_action.pagerduty_create_new_incident.id

    config = jsonencode({
      incident_title   = "{{incident.severity}} Incident: {{incident.name}}"
      incident_details = "There is a {{incident.severity}} incident.\n\nTriggered by {{author.name}} from FireHydrant. For more information see {{incident.incident_url}}."

      escalation_policy_id = {
        label = "Registration Team Escalation Policy"
        value = "xxxxxxx"
      }
      incident_creator = {
        label = "user@example.com"
        value = "user@example.com"
      }
      service_id = {
        label = "Registration Service"
        value = "xxxxxxx"
      }
    })

    automatic = true
  }
}
```

### Create PagerDuty Incident - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `config` - (Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.

The `config` block supports:

* `escalation_policy_id` - (Required) The PagerDuty escalation policy for determining the users to be alerted.
* `incident_creator` - (Required) The PagerDuty user to use for creating the incident.
* `incident_details` - (Required) A description of the incident for PagerDuty.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
* `incident_title` - (Required) The title of the PagerDuty incident.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
* `service_id` - (Required) The PagerDuty service to include in the incident.

The `escalation_policy_id` block supports:

* `label` - (Required) The name of the PagerDuty escalation policy.
* `value` - (Required) The ID of the PagerDuty escalation policy.

The `incident_creator` block supports:

* `label` - (Required) The email of the PagerDuty user.
  This value should match the `value` attribute.
* `value` - (Required) The email of the PagerDuty user.
  This value should match the `label` attribute.

The `service_id` block supports:

* `label` - (Required) The name of the PagerDuty service.
* `value` - (Required) The ID of the PagerDuty service.