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


resource "firehydrant_environment" "production" {
  name = "FireHydrant Production"
}

resource "firehydrant_environment" "staging" {
  name = "FireHydrant Staging"
}

resource "firehydrant_service" "laddertruck-web" {
  name = "Laddertruck (web)"
  labels = {
    stack = "rails"
    version = "5.2.x"
  }
}

resource "firehydrant_service" "laddertruck-pubsub" {
  name = "Laddertruck (pubsub)"
  labels = {
    stack = "rails"
    version = "5.2.x"
  }
}

resource "firehydrant_functionality" "incident-management" {
  name = "Incident Management"

  services {
    id = firehydrant_service.laddertruck-web.id
  }

  services {
    id = firehydrant_service.laddertruck-pubsub.id
  }
}

resource "firehydrant_team" "firefighters" {
  name = "Firefighters"

  services {
    id = firehydrant_service.laddertruck-web.id
  }

  services {
    id = firehydrant_service.laddertruck-pubsub.id
  }
}

resource "firehydrant_severity" "sev1" {
  slug = "SEV1"
}

resource "firehydrant_severity" "sev2" {
  slug = "SEV2"
}

resource "firehydrant_severity" "sev3" {
  slug = "SEV3"
}

resource "firehydrant_severity" "sev4" {
  slug = "SEV4"
}

resource "firehydrant_severity" "sev5" {
  slug = "SEV5"
}

data "firehydrant_runbook_action" "email-notification" {
  slug = "email_notification"
  type = "incident"
  integration_slug = "patchy"
}

resource "firehydrant_runbook" "default-process" {
  name = "Default Incident Management"
  description = "This is the default incident management runbook"
  type = "incident"

  severities {
    id = firehydrant_severity.sev1.slug
  }

  severities {
    id = firehydrant_severity.sev2.slug
  }

  severities {
    id = firehydrant_severity.sev3.slug
  }

  steps {
    name = "Send a notification"
    action_id = data.firehydrant_runbook_action.email-notification.id
    automatic = true
    config = {
      email_address = "robert+terraform@firehydrant.io"
      subject = "An incident has been declared!"
      default_message = "A really bad incident has been declared!"
    }
  }
}
