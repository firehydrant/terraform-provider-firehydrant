terraform {
  required_providers {
    firehydrant = {
      source = "firehydrant/firehydrant"
      version = "0.1.2"
    }
  }
}

provider "firehydrant" {
}

data "firehydrant_runbook_action" "email-notification" {
  slug = "email_notification"
  type = "incident"
  integration_slug = "patchy"
}

resource "firehydrant_environment" "production" {
    name = "Production"
}

resource "firehydrant_service" "heimdall" {
    name   = "Heimdall"
    labels = {}
}


# firehydrant_runbook.default:
resource "firehydrant_runbook" "default" {
    name = "Default Incident Process"
    type = "incident"

    steps {
        action_id = "a3136370-9ebd-476f-94f6-3eedf93d7bda"
        automatic = true
        config    = {
            "channel_name_format" = "inc-{{ number }}"
        }
        name      = "Create incident channel in Slack"
        repeats   = false
    }
    steps {
        action_id = "5a1b39d0-8e78-47b5-af64-947f661f9f7b"
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
        action_id = "3a59c9d4-57db-49a3-a886-01b6985785cb"
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
        action_id = "df7b7a70-037f-4d9d-84e8-00c3c8afdefe"
        automatic = true
        config    = {
            "email_address" = "stakeholders@example.com"
            "subject"       = "{{ incident.severity }} - {{ incident.name }} incident has been started"
        }
        name      = "Email stakeholders"
        repeats   = false
    }
    steps {
        action_id = "5a1b39d0-8e78-47b5-af64-947f661f9f7b"
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
        action_id = "7125e149-6ac5-4887-a11b-4d11e88902b8"
        automatic = true
        config    = {}
        name      = "Archive incident channel after retrospective completion"
        repeats   = false
    }
}
