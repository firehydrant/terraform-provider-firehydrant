---
page_title: "FireHydrant Data Source: firehydrant_runbook_action"
subcategory: ""
---

# firehydrant_runbook_action Data Source

Use this data source to get information on runbook actions.

FireHydrant runbook actions represent the different types of actions
a runbook step can take. Some runbook actions require that you have the
associated integration installed in FireHydrant.

## Example Usage

Basic usage:
```hcl
# Confluence Cloud
data "firehydrant_runbook_action" "confluence_cloud_export_retro" {
  integration_slug = "confluence_cloud"
  slug             = "export_retrospective"
}

# FireHydrant
data "firehydrant_runbook_action" "firehydrant_add_task_list" {
  integration_slug = "patchy"
  slug             = "add_task_list"
}

data "firehydrant_runbook_action" "firehydrant_assign_a_role" {
  integration_slug = "patchy"
  slug             = "assign_a_role"
}

data "firehydrant_runbook_action" "firehydrant_assign_a_team" {
  integration_slug = "patchy"
  slug             = "assign_a_team"
}

data "firehydrant_runbook_action" "firehydrant_attach_a_runbook" {
  integration_slug = "patchy"
  slug             = "attach_a_runbook"
}

data "firehydrant_runbook_action" "firehydrant_add_services_related_to_functionality" {
  integration_slug = "patchy"
  slug             = "add_services_related_to_functionality"
}

data "firehydrant_runbook_action" "firehydrant_email_notification" {
  integration_slug = "patchy"
  slug             = "email_notification"
}

data "firehydrant_runbook_action" "firehydrant_freeform_text" {
  integration_slug = "patchy"
  slug             = "freeform_text"
}

data "firehydrant_runbook_action" "firehydrant_incident_update" {
  integration_slug = "patchy"
  slug             = "incident_update"
}

data "firehydrant_runbook_action" "firehydrant_resolve_linked_alerts" {
  integration_slug = "patchy"
  slug             = "set_linked_alerts_status"
}

data "firehydrant_runbook_action" "firehydrant_script" {
  integration_slug = "patchy"
  slug             = "script"
}

data "firehydrant_runbook_action" "firehydrant_send_webhook" {
  integration_slug = "patchy"
  slug             = "send_webhook"
}

data "firehydrant_runbook_action" "firehydrant_publish_to_statuspage" {
  integration_slug = "nunc"
  slug             = "create_nunc"
}

# Giphy
data "firehydrant_runbook_action" "giphy_incident_channel_gif" {
  integration_slug = "giphy"
  slug             = "incident_channel_gif"
}

# Google Docs
data "firehydrant_runbook_action" "google_docs_export_retro" {
  integration_slug = "google_docs"
  slug             = "export_retrospective"
}

# Google Meet
data "firehydrant_runbook_action" "google_meet_create_google_meet_link" {
  integration_slug = "google_meet"
  slug             = "create_google_meet_link"
}

# Jira Cloud
data "firehydrant_runbook_action" "jira_cloud_create_incident_issue" {
  integration_slug = "jira_cloud"
  slug             = "create_incident_issue"
}

# Jira Server
data "firehydrant_runbook_action" "jira_server_create_incident_issue" {
  integration_slug = "jira_onprem"
  slug             = "create_incident_issue"
}

# Microsoft Teams
data "firehydrant_runbook_action" "microsoft_teams_create_incident_channel" {
  integration_slug = "microsoft_teams"
  slug             = "create_incident_channel"
}

data "firehydrant_runbook_action" "microsoft_teams_notify_channel" {
  integration_slug = "microsoft_teams"
  slug             = "notify_channel"
}

data "firehydrant_runbook_action" "microsoft_teams_notify_channel_custom" {
  integration_slug = "microsoft_teams"
  slug             = "notify_channel_custom_message"
}

data "firehydrant_runbook_action" "microsoft_teams_notify_incident_channel_custom" {
  integration_slug = "microsoft_teams"
  slug             = "notify_incident_channel_custom_message"
}

# OpsGenie
data "firehydrant_runbook_action" "opsgenie_create_new_incident" {
  integration_slug = "opsgenie"
  slug             = "create_new_opsgenie_incident"
}

# PagerDuty
data "firehydrant_runbook_action" "pagerduty_create_new_incident" {
  integration_slug = "pager_duty"
  slug             = "create_new_pager_duty_incident"
}

# Shortcut
data "firehydrant_runbook_action" "shortcut_create_incident_issue" {
  integration_slug = "shortcut"
  slug             = "create_incident_issue"
}

# Slack
data "firehydrant_runbook_action" "slack_add_bookmark_to_incident_channel" {
  integration_slug = "slack"
  slug             = "add_bookmark_to_incident_channel"
}

data "firehydrant_runbook_action" "slack_archive_channel" {
  integration_slug = "slack"
  slug             = "archive_incident_channel"
}

data "firehydrant_runbook_action" "slack_create_incident_channel" {
  integration_slug = "slack"
  slug             = "create_incident_channel"
}

data "firehydrant_runbook_action" "slack_notify_channel" {
  integration_slug = "slack"
  slug             = "notify_channel"
}

data "firehydrant_runbook_action" "slack_notify_channel_custom" {
  integration_slug = "slack"
  slug             = "notify_channel_custom_message"
}

data "firehydrant_runbook_action" "slack_notify_incident_channel_custom" {
  integration_slug = "slack"
  slug             = "notify_incident_channel_custom_message"
}

# Statuspage.io
data "firehydrant_runbook_action" "statuspage_create_statuspage" {
  integration_slug = "statuspage"
  slug             = "create_statuspage"
}

data "firehydrant_runbook_action" "statuspage_update_statuspage" {
  integration_slug = "statuspage"
  slug             = "update_statuspage"
}

# VictorOps
data "firehydrant_runbook_action" "victorops_create_new_incident" {
  integration_slug = "victorops"
  slug             = "create_new_victorops_incident"
}

# Webex
data "firehydrant_runbook_action" "webex_create_meeting" {
  integration_slug = "webex"
  slug             = "create_meeting"
}

# Zoom 
data "firehydrant_runbook_action" "zoom_create_meeting" {
  integration_slug = "zoom"
  slug             = "create_meeting"
}
```

## Argument Reference

The following arguments are supported:

* `integration_slug` - (Required) The slug of the integration associated with the runbook action.
* `slug` - (Required) The slug of the runbook action.
* `type` - (Optional) The type of runbook supported for the runbook action. Valid values are `incident`, 
  `infrastructure`, and `incident_role`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the runbook action.
* `name` - The name of the runbook action.
