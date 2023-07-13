---
page_title: "Step Configuration - FireHydrant"
subcategory: "Runbooks"
---

# FireHydrant

FireHydrant supports a number of built-in steps that can be added to runbooks.

### Available Steps

* [Add Services Related to Functionality](#add-services-related-to-functionality)
* [Add Task List](#add-task-list)
* [Assign a Role](#assign-a-role)
* [Assign a Team](#assign-a-team)
* [Attach a Runbook](#attach-a-runbook)
* [Email Notification](#email-notification)
* [Freeform Text](#freeform-text)
* [Incident Update](#incident-update)
* [Publish to Status Page](#publish-to-status-page)
* [Resolve Linked Alerts](#resolve-linked-alerts)
* [Script](#script)
* [Send Webhook](#send-webhook)

## Add Services Related to Functionality

The [FireHydrant **Add Services Related to Functionality** step](https://support.firehydrant.com/hc/en-us/articles/4409237309460-Auto-Adding-Services-Related-To-Functionality-)
allows all services related to a functionality to be automatically added to an incident, 
giving you assurance that the right teams and services will be associated with your incident 
no matter what.

### Add Services Related to Functionality - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "firehydrant_add_services_related_to_functionality" {
  integration_slug = "patchy"
  slug             = "add_services_related_to_functionality"
}

resource "firehydrant_runbook" "firehydrant_add_services_related_to_functionality_runbook" {
  name        = "add-services-related-to-functionality-runbook"
  description = "This is an example configuration for a runbook that uses the FireHydrant add services related to functionality step"

  steps {
    name             = "Auto-Add Services Related To Functionality"
    action_id        = data.firehydrant_runbook_action.firehydrant_add_services_related_to_functionality.id
    automatic        = false
    repeats          = true
    repeats_duration = "PT15M"
  }
}
```

### Add Services Related to Functionality - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.
* `repeats` - (Optional) Whether this step should repeat. Defaults to `false`.
  When this value is `true`, `repeats_duration` _must_ be provided.
* `repeats_duration` - (Optional) How often this step should repeat in ISO8601.
  Example: PT10M [Format Spec](https://www.digi.com/resources/documentation/digidocs/90001488-13/reference/r_iso_8601_duration_format.htm)
  This value _must_ be provided if `repeats` is `true`. This value _must not_ be provided if `repeats` is `false`.
* `rule` - (Optional) JSON string representing the rule configuration for the runbook step.
  For more information on the conditional logic used in `rule`, see the
  [Runbooks - Conditional Logic](./runbooks_conditional_logic.md) documentation.
  The step will default to running manually if `rule` is not specified and `automatic` and `repeats` are both `false`.

## Add Task List

The FireHydrant **Add Task List** step allows 
[FireHydrant task lists](https://support.firehydrant.com/hc/en-us/articles/5505273345428-Task-lists)
to be automatically added to an incident.

### Add Task List - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "firehydrant_add_task_list" {
  integration_slug = "patchy"
  slug             = "add_task_list"
}

resource "firehydrant_incident_role" "scribe_incident_role" {
  name    = "Scribe"
  summary = "This is a scribe incident role summary"
}

resource "firehydrant_task_list" "example_task_list" {
  name        = "example-task-list"
  description = "This is an example task list"

  task_list_items {
    summary = "Example task #1"
  }

  task_list_items {
    summary     = "Example task #2"
    description = "This task is very important."
  }
}

resource "firehydrant_runbook" "firehydrant_add_task_list_runbook" {
  name        = "add-task-list-runbook"
  description = "This is an example configuration for a runbook that uses the FireHydrant add task list step"

  steps {
    name      = "Attach a Task List"
    action_id = data.firehydrant_runbook_action.firehydrant_add_task_list.id

    config = jsonencode({
      task_list = {
        label = firehydrant_task_list.example_task_list.name
        value = firehydrant_task_list.example_task_list.id
      }
      role_for_assignment = {
        label = firehydrant_incident_role.scribe_incident_role.name
        value = firehydrant_incident_role.scribe_incident_role.id
      }
    })

    automatic = true
  }
}
```

### Add Task List - Steps Argument Reference

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

* `task_list` - (Required) The task list to attach.
* `role_for_assignment` - (Optional) Which incident role to assign the tasks from the task list to.

The `task_list` block supports:

* `label` - (Required) The name of the task list.
* `value` - (Required) The ID of the task list.

The `role_for_assignment` block supports:

* `label` - (Required) The name of the incident role.
* `value` - (Required) The ID of the incident role.

## Assign a Role

The [FireHydrant **Assign a Role** step](https://support.firehydrant.com/hc/en-us/articles/360060484631-Assigning-a-role-with-a-Runbook-step) 
allows incident roles to be assigned automatically during an incident.

### Assign a Role - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "firehydrant_assign_a_role" {
  integration_slug = "patchy"
  slug             = "assign_a_role"
}

data "firehydrant_incident_role" "commander_incident_role" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

data "firehydrant_incident_role" "ops_lead_incident_role" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

resource "firehydrant_incident_role" "scribe_incident_role" {
  name    = "Scribe"
  summary = "This is a scribe incident role summary"
}

resource "firehydrant_runbook" "firehydrant_assign_a_role_runbook" {
  name        = "assign-a-role-runbook"
  description = "This is an example configuration for a runbook that uses the FireHydrant assign a role step"

  steps {
    name      = "Assign Commander Role"
    action_id = data.firehydrant_runbook_action.firehydrant_assign_a_role.id

    config = jsonencode({
      role = {
        label = data.firehydrant_incident_role.commander_incident_role.name
        value = data.firehydrant_incident_role.commander_incident_role.id
      }
      user = {
        label = "Incident Opener"
        value = jsonencode({
          type = "incident_opener"
        })
      }
    })

    automatic = true
    repeats   = false
  }

  steps {
    name      = "Assign Ops Lead Role"
    action_id = data.firehydrant_runbook_action.firehydrant_assign_a_role.id

    config = jsonencode({
      role = {
        label = data.firehydrant_incident_role.ops_lead_incident_role.name
        value = data.firehydrant_incident_role.ops_lead_incident_role.id
      }
      user = {
        label = "Example External Schedule"
        value = jsonencode({
          type                 = "external_schedule"
          external_schedule_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
        })
      }
    })

    automatic = true
    repeats   = false
  }

  steps {
    name      = "Assign Communication Role"
    action_id = data.firehydrant_runbook_action.firehydrant_assign_a_role.id

    config = jsonencode({
      role = {
        label = firehydrant_incident_role.scribe_incident_role.name
        value = firehydrant_incident_role.scribe_incident_role.id
      }
      user = {
        label = "Example User"
        value = jsonencode({
          type    = "user"
          user_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
        })
      }
    })

    automatic = true
    repeats   = false
  }
}
```

### Assign a Role - Steps Argument Reference

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

* `role` - (Required) The incident role to assign.
* `user` - (Required) The user to assign the incident role to.

The `role` block supports:

* `label` - (Required) The name of the incident role.
* `value` - (Required) The ID of the incident role.

The `user` block supports:

* `label` - (Required) The name of the user.
* `value` - (Required) JSON string representing the type and ID of the user.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).

The `value` block supports:

* `type` - (Required) The type of user. Valid values are `incident_opener`, 
  `external_schedule`, or `user`.
* `external_schedule_id` - (Optional) The ID of the external schedule type user. 
  This value _must not_ be provided if `type` is `incident_opener` or `user`.
* `user_id` - (Optional) The ID of the user type user.
  This value _must not_ be provided if `type` is `incident_opener` or `external_schedule`.

## Assign a Team

The [FireHydrant **Assign a Team** runbook step](https://support.firehydrant.com/hc/en-us/articles/360061144571-Communicating-with-teams-using-Runbook-steps)
allows teams to be assigned to an incident automatically.

### Assign a Team - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "firehydrant_assign_a_team" {
  integration_slug = "patchy"
  slug             = "assign_a_team"
}

resource "firehydrant_team" "example_team" {
  name        = "example-team"
  description = "This is an example team"
}

resource "firehydrant_runbook" "firehydrant_assign_a_team_runbook" {
  name = "assign-a-team-runbook"

  steps {
    name      = "Assign A Team"
    action_id = data.firehydrant_runbook_action.firehydrant_assign_a_team.id

    config = jsonencode({
      team = {
        label = firehydrant_team.example_team.name
        value = firehydrant_team.example_team.id
      }
    })

    automatic = true
  }
}
```

### Assign a Team - Steps Argument Reference

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

* `team` - (Required) The team to assign.

The `team` block supports:

* `label` - (Required) The name of the team.
* `value` - (Required) The ID of the team.

## Attach a Runbook

The FireHydrant **Attach a Runbook** step allows runbooks to be automatically attached to an incident.

### Attach a Runbook - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "firehydrant_attach_a_runbook" {
  integration_slug = "patchy"
  slug             = "attach_a_runbook"
}

data "firehydrant_runbook_action" "firehydrant_add_services_related_to_functionality" {
  integration_slug = "patchy"
  slug             = "add_services_related_to_functionality"
}

resource "firehydrant_runbook" "firehydrant_add_services_related_to_functionality_runbook" {
  name        = "add-services-related-to-functionality-runbook"
  description = "This is an example configuration for a runbook that uses the FireHydrant add services related to functionality step"

  steps {
    name             = "Auto-Add Services Related To Functionality"
    action_id        = data.firehydrant_runbook_action.firehydrant_add_services_related_to_functionality.id
    automatic        = false
    repeats          = true
    repeats_duration = "PT15M"
  }
}

resource "firehydrant_runbook" "firehydrant_attach_a_runbook_runbook" {
  name = "attach-a-runbook-runbook"

  steps {
    name      = "Attach A Runbook"
    action_id = data.firehydrant_runbook_action.firehydrant_attach_a_runbook.id

    config = jsonencode({
      runbook = {
        label = firehydrant_runbook.firehydrant_add_services_related_to_functionality_runbook.name
        value = firehydrant_runbook.firehydrant_add_services_related_to_functionality_runbook.id
      }
    })

    automatic = true
  }
}
```

### Attach a Runbook - Steps Argument Reference

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

* `runbook` - (Required) The runbook to attach.

The `runbook` block supports:

* `label` - (Required) The name of the runbook.
* `value` - (Required) The ID of the runbook.

## Email Notification

The [FireHydrant **Email Notification** step](https://support.firehydrant.com/hc/en-us/articles/360058202891-Sending-email-notification-via-Runbooks)
allows emails to be sent to stakeholders automatically during an incident.

### Email Notification - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "firehydrant_email_notification" {
  integration_slug = "patchy"
  slug             = "email_notification"
}

resource "firehydrant_runbook" "firehydrant_email_notification_runbook" {
  name = "email-notification-runbook"

  steps {
    name      = "Send an email notification"
    action_id = data.firehydrant_runbook_action.firehydrant_email_notification.id

    config = jsonencode({
      email_address   = "test@example.com"
      email_subject   = "Incident opened on FireHydrant"
      default_message = "Message"
    })

    automatic = true
    repeats   = false
  }
}
```

### Email Notification - Steps Argument Reference

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

* `email_address` - (Required) The email address to send a message to.
* `email_subject` - (Required) The subject line to use for the email.
* `default_message` - (Required) The body to use for the email.

## Freeform Text

The [FireHydrant **Freeform Text** step](https://support.firehydrant.com/hc/en-us/articles/360057722252-Freeform-Text)
allows custom text to be displayed during an incident.

### Freeform Text - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "firehydrant_freeform_text" {
  integration_slug = "patchy"
  slug             = "freeform_text"
}

resource "firehydrant_runbook" "firehydrant_freeform_text_runbook" {
  name = "freeform-text-runbook"

  steps {
    name      = "Freeform Text"
    action_id = data.firehydrant_runbook_action.firehydrant_freeform_text.id

    config = jsonencode({
      text = "This is example text."
    })

    repeats          = true
    repeats_duration = "PT15M"
  }
}
```

### Freeform Text - Steps Argument Reference

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
* `rule` - (Optional) JSON string representing the rule configuration for the runbook step.
  For more information on the conditional logic used in `rule`, see the
  [Runbooks - Conditional Logic](./runbooks_conditional_logic.md) documentation.
  The step will default to running manually if `rule` is not specified and `automatic` and `repeats` are both `false`.

The `config` block supports:

* `text` - (Required) The text to display when this step is executed.

## Incident Update

The [FireHydrant **Incident Update** step](https://support.firehydrant.com/hc/en-us/articles/360057721992-Updating-an-incident-using-a-Runbook-step)
allows incident details to be automatically updated during an incident.

### Incident Update - Example Usage

Basic usage:
```hcl
resource "firehydrant_priority" "example_priority" {
  slug        = "MYEXAMPLEPRIORITY"
  description = "This is an example priority"
}

resource "firehydrant_severity" "example_severity" {
  slug        = "EXAMPLESEVERITY"
  description = "This is an example severity"
}

data "firehydrant_runbook_action" "firehydrant_incident_update" {
  integration_slug = "patchy"
  slug             = "incident_update"
}

resource "firehydrant_runbook" "firehydrant_incident_update_runbook" {
  name = "incident-update-runbook"

  steps {
    name      = "Update Incident Details"
    action_id = data.firehydrant_runbook_action.firehydrant_incident_update.id

    config = jsonencode({
      milestone = {
        label = "Mitigated"
        value = "mitigated"
      }
      severity = {
        label = firehydrant_severity.example_severity.slug
        value = firehydrant_severity.example_severity.slug
      }
      priority = {
        label = firehydrant_priority.example_priority.slug
        value = firehydrant_priority.example_priority.slug
      }

      description     = "This is an example description."
      customer_impact = "Customers on the free tier are impacted."
      comment         = "Example comment"
    })
  }
}
```

### Incident Update - Steps Argument Reference

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

* `description` - (Optional) The new description for the incident.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
* `comment` - (Optional) A note to be included on the incident when it is updated.
  This field supports [FireHydrant's template variables](https://support.firehydrant.com/hc/en-us/articles/4409136426004-Using-template-variables-in-Runbooks)
  so you can automatically include details such as when the incident started, the incident summary, severity, roles involved, and much more.
* `customer_impact` - (Optional) The new customer impact summary for the incident.
* `milestone` - (Optional) The new milestone for the incident.
* `priority` - (Optional) The new priority for the incident.
* `severity` - (Optional) The new severity for the incident.

The `milestone` block supports:

* `label` - (Required) The name of the milestone.
  Valid values are:
    - `Started`
    - `Detected`
    - `Acknowledged`
    - `Investigating`
    - `Identified`
    - `Mitigated`
    - `Resolved`
    - `Retrospective Started`
    - `Retrospective Completed`
    - `Closed`
* `value` - (Required) The slug of the milestone.
  Valid values are:
  - `started`
  - `detected`
  - `acknowledged`
  - `investigating`
  - `identified`
  - `mitigated`
  - `resolved`
  - `postmortem_started`
  - `postmortem_completed`
  - `closed`

The `priority` block supports:

* `label` - (Required) The slug of the priority.
* `value` - (Required) The slug of the priority.

The `severity` block supports:

* `label` - (Required) The slug of the severity.
* `value` - (Required) The slug of the severity.

## Publish to Status Page

The FireHydrant **Publish to Status Page** step allows your 
[FireHydrant status pages](https://support.firehydrant.com/hc/en-us/articles/360057722032-Status-Pages-Overview)
to be automatically updated during an incident.

### Publish to Status Page - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "firehydrant_publish_to_status_page" {
  integration_slug = "nunc"
  slug             = "create_nunc"
}

resource "firehydrant_runbook" "firehydrant_publish_to_status_page_runbook" {
  name = "publish-to-status-page-runbook"

  steps {
    name      = "Publish to status page"
    action_id = data.firehydrant_runbook_action.firehydrant_publish_to_status_page.id

    config = jsonencode({
      title = "Example Incident #123"
      connection_id = {
        label = "status.firehydrant.com"
        value = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
      }
    })

    automatic = true
  }
}
```

### Publish to Statuspage - Steps Argument Reference

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

* `connection_id` - (Required) The status page connection to publish to.
* `title` - (Optional) The title of the incident to be used when publishing the incident to your status page.

The `connection_id` block supports:

* `label` - (Required) The name of the connection.
* `value` - (Required) The ID of the connection.

## Resolve Linked Alerts

The [FireHydrant **Resolve Linked Alerts** step](https://support.firehydrant.com/hc/en-us/articles/6689194321940-Resolve-Linked-Alerts)
allows all third party alerts that have been attached to an incident to be automatically resolved.

### Resolve Linked Alerts - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "firehydrant_resolve_linked_alerts" {
  integration_slug = "patchy"
  slug             = "set_linked_alerts_status"
}

resource "firehydrant_runbook" "firehydrant_resolve_linked_alerts_runbook" {
  name = "resolve-linked-alerts-runbook"

  steps {
    name      = "Resolve Linked Alerts"
    action_id = data.firehydrant_runbook_action.firehydrant_resolve_linked_alerts.id
    automatic = true

    rule = jsonencode({
      logic = {
        eq = [
          {
            var = "incident_current_milestone"
          },
          {
            var = "usr.1"
          }
        ]
      }
      user_data = {
        "1" = {
          type  = "Milestone"
          value = "resolved"
          label = "Resolved"
        }
      }
    })
  }
}
```

### Resolve Linked Alerts - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `name` - (Required) The name of the step.
* `automatic` - (Optional) Whether this step should automatically execute.
* `rule` - (Optional) JSON string representing the rule configuration for the runbook step.
  For more information on the conditional logic used in `rule`, see the
  [Runbooks - Conditional Logic](./runbooks_conditional_logic.md) documentation.
  The step will default to running manually if `rule` is not specified and `automatic` and `repeats` are both `false`.

## Script

The [FireHydrant **Script** step](https://support.firehydrant.com/hc/en-us/articles/360057722272-Executing-patch-scripts)
allows storing and tracking the execution of custom scripts during an incident.

### Script - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "firehydrant_script" {
  integration_slug = "patchy"
  slug             = "script"
}

resource "firehydrant_runbook" "firehydrant_script_runbook" {
  name = "script-runbook"

  steps {
    name      = "Script"
    action_id = data.firehydrant_runbook_action.firehydrant_script.id

    config = jsonencode({
      description = "This script checks the current username."
      script      = "#!/bin/bash\necho I am $USERNAME"
    })
  }
}
```

### Script - Steps Argument Reference

* `action_id` - (Required) The ID of the runbook action for the step.
* `config` - (Required) JSON string representing the configuration settings for the step.
  Use [Terraform's jsonencode](https://www.terraform.io/language/functions/jsonencode)
  function so that [Terraform can guarantee valid JSON syntax](https://www.terraform.io/language/expressions/strings#generating-json-or-yaml).
* `name` - (Required) The name of the step.
* `rule` - (Optional) JSON string representing the rule configuration for the runbook step.
  For more information on the conditional logic used in `rule`, see the
  [Runbooks - Conditional Logic](./runbooks_conditional_logic.md) documentation.
  The step will default to running manually if `rule` is not specified and `automatic` and `repeats` are both `false`.

The `config` block supports:

* `script` - (Required) The script to be executed.
* `description` - (Optional) A description of the script.

## Send Webhook

The [FireHydrant **Send Webhook** step](https://support.firehydrant.com/hc/en-us/articles/360057721712-Sending-a-webhook-from-a-Runbook)
allows webhooks to be sent automatically during an incident.

### Send Webhook - Example Usage

Basic usage:
```hcl
data "firehydrant_runbook_action" "firehydrant_send_webhook" {
  integration_slug = "patchy"
  slug             = "send_webhook"
}

resource "firehydrant_runbook" "firehydrant_send_webhook_runbook" {
  name = "webhook-runbook"

  steps {
    name      = "Send A Webhook"
    action_id = data.firehydrant_runbook_action.firehydrant_send_webhook.id

    config = jsonencode({
      endpoint = "https://example.com"
      payload  = "{\n  \"incident_id\": \"{{ incident.id }}\",\n  \"name\": \"{{ incident.name }}\",\n  \"started_at\": \"{{ incident.started_at }}\",\n  \"impacts\": {{ incident.impacts | toJSON }},\n  \"private_status_page_url\": \"{{ incident.private_status_page_url }}\"\n}\n"
    })

    automatic = true
  }
}
```

### Send Webhook - Steps Argument Reference

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

* `endpoint` - (Required) The endpoint to send the payload to.
* `payload` - (Required) The payload to send to the endpoint.
* `hmac_secret` - (Optional) An HMAC secret to use when sending the webhook
