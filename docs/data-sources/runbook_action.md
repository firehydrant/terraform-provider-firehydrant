---
page_title: "firehydrant_runbook_action Data Source - terraform-provider-firehydrant"
subcategory: ""
description: |-
  
---

# Data Source `firehydrant_runbook_action`

## Example Usage

Basic usage:
```hcl
# Confluence Cloud
data "firehydrant_runbook_action" "confluence_cloud_export_retro" {
  integration_slug = "confluence_cloud"
  slug             = "export_retrospective"
  type             = "incident"
}

# FireHydrant
data "firehydrant_runbook_action" "firehydrant_assign_a_role" {
  integration_slug = "patchy"
  slug             = "assign_a_role"
  type             = "incident"
}

data "firehydrant_runbook_action" "firehydrant_assign_a_team" {
  integration_slug = "patchy"
  slug             = "assign_a_team"
  type             = "incident"
}

data "firehydrant_runbook_action" "firehydrant_attach_a_runbook" {
  integration_slug = "patchy"
  slug             = "attach_a_runbook"
  type             = "incident"
}

data "firehydrant_runbook_action" "firehydrant_add_services_related_to_functionality" {
  integration_slug = "patchy"
  slug             = "add_services_related_to_functionality"
  type             = "incident"
}

data "firehydrant_runbook_action" "firehydrant_create_incident_ticket" {
  integration_slug = "patchy"
  slug             = "create_incident_ticket"
  type             = "incident"
}

data "firehydrant_runbook_action" "firehydrant_email_notification" {
  integration_slug = "patchy"
  slug             = "email_notification"
  type             = "incident"
}

data "firehydrant_runbook_action" "firehydrant_freeform_text" {
  integration_slug = "patchy"
  slug             = "freeform_text"
  type             = "incident"
}

data "firehydrant_runbook_action" "firehydrant_incident_update" {
  integration_slug = "patchy"
  slug             = "incident_update"
  type             = "incident"
}

data "firehydrant_runbook_action" "firehydrant_resolve_linked_alerts" {
  integration_slug = "patchy"
  slug             = "set_linked_alerts_status"
  type             = "incident"
}

data "firehydrant_runbook_action" "firehydrant_script" {
  integration_slug = "patchy"
  slug             = "script"
  type             = "incident"
}

data "firehydrant_runbook_action" "firehydrant_send_webhook" {
  integration_slug = "patchy"
  slug             = "send_webhook"
  type             = "incident"
}

data "firehydrant_runbook_action" "firehydrant_publish_to_statuspage" {
  integration_slug = "nunc"
  slug             = "create_nunc"
  type             = "incident"
}

# Giphy
data "firehydrant_runbook_action" "giphy_incident_channel_gif" {
  integration_slug = "giphy"
  slug             = "incident_channel_gif"
  type             = "incident"
}

# Google
data "firehydrant_runbook_action" "google_create_google_meet_link" {
  integration_slug = "google_meet"
  slug             = "create_google_meet_link"
  type             = "incident"
}

# Jira Cloud
data "firehydrant_runbook_action" "jira_cloud_create_incident_issue" {
  integration_slug = "jira_cloud"
  slug             = "create_incident_issue"
  type             = "incident"
}

# Jira Server
data "firehydrant_runbook_action" "jira_server_create_incident_issue" {
  integration_slug = "jira_onprem"
  slug             = "create_incident_issue"
  type             = "incident"
}

# OpsGenie
data "firehydrant_runbook_action" "opsgenie_create_new_incident" {
  integration_slug = "opsgenie"
  slug             = "create_new_opsgenie_incident"
  type             = "incident"
}

# PagerDuty
data "firehydrant_runbook_action" "pagerduty_create_new_incident" {
  integration_slug = "pager_duty"
  slug             = "create_new_pager_duty_incident"
  type             = "incident"
}

# Shortcut
data "firehydrant_runbook_action" "shortcut_create_incident_issue" {
  integration_slug = "shortcut"
  slug             = "create_incident_issue"
  type             = "incident"
}

# Slack
data "firehydrant_runbook_action" "slack_archive_channel" {
  integration_slug = "slack"
  slug             = "archive_incident_channel"
  type             = "incident"
}

data "firehydrant_runbook_action" "slack_create_incident_channel" {
  integration_slug = "slack"
  slug             = "create_incident_channel"
  type             = "incident"
}

data "firehydrant_runbook_action" "slack_notify_channel" {
  integration_slug = "slack"
  slug             = "notify_channel"
  type             = "incident"
}

data "firehydrant_runbook_action" "slack_notify_channel_custom" {
  integration_slug = "slack"
  slug             = "notify_channel_custom_message"
  type             = "incident"
}

data "firehydrant_runbook_action" "slack_notify_incident_channel_custom" {
  integration_slug = "slack"
  slug             = "notify_incident_channel_custom_message"
  type             = "incident"
}

# Statuspage.io
data "firehydrant_runbook_action" "statuspage_create_statuspage" {
  integration_slug = "statuspage"
  slug             = "create_statuspage"
  type             = "incident"
}

data "firehydrant_runbook_action" "statuspage_update_statuspage" {
  integration_slug = "statuspage"
  slug             = "update_statuspage"
  type             = "incident"
}

# VictorOps
data "firehydrant_runbook_action" "victorops_create_new_incident" {
  integration_slug = "victorops"
  slug             = "create_new_victorops_incident"
  type             = "incident"
}

# Zoom 
data "firehydrant_runbook_action" "zoom_create_meeting" {
  integration_slug = "zoom"
  slug             = "create_meeting"
  type             = "incident"
}
```

## Schema

### Required

- **integration_slug** (String, Required) The slug of the integration associated with the runbook action.
- **slug** (String, Required) The slug of the runbook action.
- **type** (String, Required) The type of runbook supported for the runbook action. 
   Valid values are `incident`, `infrastructure`, and `incident_role`.

### Read-only

- **id** (String, Read-only) The ID of the runbook action.
- **name** (String, Read-only) The name of the runbook action.
