terraform {
  required_providers {
    firehydrant = {
      source = "firehydrant/firehydrant"
      version = "0.1.5"
    }
  }
}

provider "firehydrant" {
  api_key = "fhb-00000000000000000000000000000000"
  firehydrant_base_url = "https://api.local.firehydrant.io/v1/"
}

resource "firehydrant_environment" "production" {
    name = "Production"
}

data "firehydrant_runbook_action" "slack_channel" {
  integration_slug = "slack"
  slug = "create_incident_channel"
  type = "incident"
}

data "firehydrant_runbook_action" "notify_channel" {
  integration_slug = "slack"
  slug = "notify_channel"
  type = "incident"
}
data "firehydrant_runbook_action" "notify_channel_custom" {
  integration_slug = "slack"
  slug = "notify_incident_channel_custom_message"
  type = "incident"
}

data "firehydrant_runbook_action" "create_incident_ticket" {
  integration_slug = "patchy"
  slug = "create_incident_ticket"
  type = "incident"
}


data "firehydrant_runbook_action" "archive_channel" {
  integration_slug = "slack"
  slug = "archive_incident_channel"
  type = "incident"
}

data "firehydrant_runbook_action" "email_notification" {
  integration_slug = "patchy"
  slug = "email_notification"
  type = "incident"
}



resource "firehydrant_runbook" "default" {
    name = "Default Incident Process WOOHOO"
    type = "incident"

    steps {
        action_id = data.firehydrant_runbook_action.slack_channel.id
        automatic = true
        config    = {
            "channel_name_format" = "inc-{{ number }}"
        }
        name      = "Create incident channel in Slack"
        repeats   = false
    }
    steps {
        action_id = data.firehydrant_runbook_action.notify_channel_custom.id
        automatic = true
        config    = {
            "message" = <<EOT
                Here's the documentation on successfully running an incident with FireHydrant's Slack bot: https://help.firehydrant.io/en/articles/3050697-incident-response-w-slack

                Don't worry all of your messages and actions here are tracked into your incident on the FireHydrant UI.
            EOT
        }
        name      = "Incident Preamble"
        repeats   = false
    }
    steps {
        action_id = data.firehydrant_runbook_action.email_notification.id
        automatic = true
        config    = {
            "email_address" = "stakeholders@example.com"
            "subject"       = "{{ incident.severity }} - {{ incident.name }} incident has been started"
        }
        name      = "Email stakeholders"
        repeats   = false
    }
    steps {
        action_id = data.firehydrant_runbook_action.notify_channel.id
        automatic = true
        config    = {
            "channels" = "#fh-incidents"
        }
        name      = "Notify incidents channel that a new incident has been opened"
        repeats   = false
    }
    steps {
        action_id = "b60654c8-2c15-4a25-b5e9-2ac2884f0b9b"
        automatic = true
        config    = {
            "ticket_description" = "{{ incident.description }}"
            "ticket_summary"     = "{{ incident.name }}"
        }
        name      = "Create an incident ticket in Jira"
        repeats   = false
    }
    steps {
        action_id = data.firehydrant_runbook_action.notify_channel_custom.id
        automatic = true
        config    = {
            "message" = <<EOT
                Please check-in with your current status on this {{ incident.severity }} incident

                ```
                /firehydrant add note I'm calculating the power required by the flux capacitor
                ```
            EOT
        }
        name      = "Remind Slack channel to update stakeholders"
        repeats   = false
    }

    steps {
        action_id = data.firehydrant_runbook_action.archive_channel.id
        automatic = true
        config    = {}
        name      = "Archive incident channel after retrospective completion"
        repeats   = false
    }
}
