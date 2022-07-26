data "firehydrant_runbook_action" "slack_channel" {
  integration_slug = "slack"
  slug             = "create_incident_channel"
}

data "firehydrant_runbook_action" "notify_channel" {
  integration_slug = "slack"
  slug             = "notify_channel"
}

data "firehydrant_runbook_action" "notify_channel_custom" {
  integration_slug = "slack"
  slug             = "notify_incident_channel_custom_message"
}

data "firehydrant_runbook_action" "archive_channel" {
  integration_slug = "slack"
  slug             = "archive_incident_channel"
}

data "firehydrant_runbook_action" "email_notification" {
  integration_slug = "patchy"
  slug             = "email_notification"
}

resource "firehydrant_runbook" "default" {
  name = "Default Incident Process"
  attachment_rule = jsonencode({
    "logic" = {
      "eq" = [
        {
          "var" = "incident_current_milestone",
        },
        {
          "var" = "usr.1"
        }
      ]
    },
    "user_data" = {
      "1" = {
        "type"  = "Milestone",
        "value" = "started",
        "label" = "Started"
      }
    }
    }
  )

  steps {
    action_id = data.firehydrant_runbook_action.slack_channel.id
    name      = "Create incident channel in Slack"

    config = jsonencode({
      channel_name_format = "inc-{{ number }}"
    })

    automatic = true
    repeats   = false
  }

  steps {
    action_id = data.firehydrant_runbook_action.notify_channel_custom.id
    name      = "Incident Preamble"

    config = jsonencode({
      message = "Here's the documentation on successfully running an incident with FireHydrant's Slack bot: https://help.firehydrant.io/en/articles/3050697-incident-response-w-slack\n\nDon't worry all of your messages and actions here are tracked into your incident on the FireHydrant UI."
    })

    automatic = true
    repeats   = false
  }

  steps {
    action_id = data.firehydrant_runbook_action.email_notification.id
    name      = "Email stakeholders"

    config = jsonencode({
      email_address = "stakeholders@example.com"
      subject       = "{{ incident.severity }} - {{ incident.name }} incident has been started"
    })

    automatic = true
    repeats   = false
  }

  steps {
    action_id = data.firehydrant_runbook_action.notify_channel.id
    name      = "Notify incidents channel that a new incident has been opened"

    config = jsonencode({
      channels = "#fh-incidents"
    })

    automatic = true
    repeats   = false
  }

  steps {
    action_id = data.firehydrant_runbook_action.notify_channel_custom.id
    name      = "Remind Slack channel to update stakeholders"

    config = jsonencode({
      message = <<EOT
              Please check-in with your current status on this {{ incident.severity }} incident

              ```
              /firehydrant add note I'm calculating the power required by the flux capacitor
              ```
          EOT
    })

    automatic = true
    repeats   = false
  }

  steps {
    action_id = data.firehydrant_runbook_action.archive_channel.id
    name      = "Archive incident channel after retrospective completion"
    automatic = false
    repeats   = false
    rule = jsonencode({
      logic = {
        eq = [
          {
            var = "incident_current_milestone",
          },
          {
            var = "usr.1"
          }
        ]
      },
      user_data = {
        "1" = {
          type  = "Milestone",
          value = "resolved",
          label = "Resolved"
        }
      }
    })
  }
}
