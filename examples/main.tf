terraform {
  required_providers {
    firehydrant = {
      versions = ["0.1.0"]
      source = "firehydrant.com/twitch/firehydrant"
    }
  }
}

provider "firehydrant" {
}

resource "firehydrant_service" "paddy-cake-paddy-cake" {
  name = "Whatever"

  labels = {
    mykey = "myvalue"
  }
}

data "firehydrant_service" "oauth2-proxy" {
  id = "2177ce81-b6b6-4063-af73-6c881c8b9899"
}

data "firehydrant_services" "logging-in-services" {
  query = "kube-system"
}

resource "firehydrant_functionality" "logging-in-2" {
  name = "Logging In (from TF) 2"

  dynamic "services" {
    for_each = [for s in data.firehydrant_services.logging-in-services.services: {
      id = s.id
    }]

    content {
      id = services.value.id
    }
  }
}

resource "firehydrant_severity" "sev1" {
  slug = "SEV1TF"
}

resource "firehydrant_runbook" "default-process-tf" {
  name = "Default IR Process (from tf)"
  type = "incident"

  severities {
    id = firehydrant_severity.sev1.slug
  }

  steps {
    name = "Create Incident Channel"
    action_id = data.firehydrant_runbook_action.create-incident-channel.id
    config = {
      channel_name_format = "-inc-123"
    }
  }
}

data "firehydrant_runbook_action" "create-incident-channel" {
  slug = "create_incident_channel"
  type = "incident"
  integration_slug = "slack"
}

