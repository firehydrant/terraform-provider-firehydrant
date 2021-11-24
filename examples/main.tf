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

resource "firehydrant_team" "firefighters" {
  name = "Firefighters"
}

data "firehydrant_runbook_action" "slack_channel" {
  integration_slug = "slack"
  slug = "create_incident_channel"
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
        action_id = data.firehydrant_runbook_action.archive_channel.id
        automatic = true
        config    = {}
        name      = "Archive incident channel after retrospective completion"
        repeats   = false
    }
}
