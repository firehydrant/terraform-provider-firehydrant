---
page_title: "Step Configuration - Jira Server"
subcategory: "Runbooks"
---

# Jira Server

~> **Note** You must have the Jira Server integration installed in FireHydrant
for any Jira Server runbook steps to work properly.

The FireHydrant Jira Server integration allows FireHydrant users to track all
the actions proposed during an incident in their existing project management
workflows for estimation and scheduling.

### Available Steps

* [Create Incident Issue](#create-incident-issue)

## Create Incident Issue

The [Jira Server **Create Incident Issue** step](https://support.firehydrant.com/hc/en-us/articles/360058202631-Creating-an-Incident-Ticket)
allows FireHydrant users to create a new Jira Server incident issue ticket at the start of an incident
that will automatically sync all created action items and link them to a parent ticket.

### Create Incident Issue - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "jira_server_create_incident_issue" {
  integration_slug = "jira_server"
  slug             = "create_incident_issue"
}

data "firehydrant_incident_role" "commander_incident_role" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

resource "firehydrant_runbook" "jira_server_create_incident_issue_runbook" {
  name = "jira-server-create-incident-issue-runbook"

  steps {
    name      = "Create a Jira Server Issue"
    action_id = data.firehydrant_runbook_action.jira_server_create_incident_issue.id

    config = jsonencode({
      project = {
        label = "Platform Team Work"
        value = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
      }
      role_for_assignment = {
        label = data.firehydrant_incident_role.commander_incident_role.name
        value = data.firehydrant_incident_role.commander_incident_role.id
      }
      ticket_description = "{{ incident.description }}"
      ticket_summary     = "{{ incident.name }}"
    })
    
    automatic = true
  }
}
```

### Create Incident Issue - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `config` - (Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.

The `config` block supports:

* `project` - (Required) The FireHydrant Jira Server ticketing project to create the incident issue ticket in.
  Your FireHydrant Jira Server ticketing projects can be found at the
  [List all ticketing projects](https://developers.firehydrant.io/docs/api/5e17c443b2bc6-list-all-ticketing-projects) endpoint.
  Make sure to use the query params `configured_projects=true` and `connection_ids=YOUR_FIREHYDRANT_JIRA_SERVER_CONNECTION_ID`.
* `role_for_assignment` - (Optional) Which incident role to assign the tasks from the task list to.
* `ticket_description` - (Optional) A description of the incident for the Jira Server incident issue ticket.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
* `ticket_summary` - (Optional) A summary of the incident for the Jira Server incident issue ticket.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.

The `project` block supports:

* `label` - (Required) The name of the FireHydrant Jira Server ticketing project.
* `value` - (Required) The ID of the FireHydrant Jira Server ticketing project.

The `role_for_assignment` block supports:

* `label` - (Required) The name of the incident role.
* `value` - (Required) The ID of the incident role.