terraform {
  required_providers {
    firehydrant = {
      source = "firehydrant/firehydrant"
      version = "0.1.3"
    }
  }
}

provider "firehydrant" {
  api_key = "fhb-00000000000000000000000000000000"
  firehydrant_base_url = "https://api.local.firehydrant.io/v1/"
}

resource "firehydrant_service" "web" {
    name   = "web"
    description = "The main web application for our company"
    labels = {
      language = "ruby",
      lifecycle = "production"
      system = "main"
      type = "user"
      tags = "foo; bar; baz"
    }
    service_tier = 1
}

resource "firehydrant_service" "frontend" {
    name   = "web"
    description = "The main web UI for our company"
    labels = {
      language = "javascript",
      lifecycle = "production"
      system = "main"
      type = "user"
      tags = "foo; bar; baz"
    }
    service_tier = 1
}
resource "firehydrant_service" "async" {
    name   = "async"
    description = "Sidekiq worker queue"
    labels = {
      language = "ruby"
      lifecycle = "production"
      system = "main"
      type = "user"
      tags = "foo"
    }
    service_tier = 1
}

resource "firehydrant_service" "cron" {
    name   = "cron"
    description = "Scheduled tasks using clockwork"
    labels = {
      language = "ruby"
      lifecycle = "production"
      system = "main"
      type = "user"
      tags = "foo"
    }
    service_tier = 4
}

resource "firehydrant_service" "db" {
    name   = "db"
    description = "Primary PG database"
    labels = {
      database = "postgresql"
      lifecycle = "production"
      system = "main"
      type = "ops"
      tags = "foo"
    }
    service_tier = 1
}

resource "firehydrant_service" "logging" {
    name   = "logging"
    description = "Datadog"
    labels = {
      lifecycle = "production"
      system = "main"
      type = "ops"
      tags = "foo"
    }
    service_tier = 3
}

resource "firehydrant_service" "uploads" {
    name   = "uploads"
    description = "AWS S3 bucket"
    labels = {
      lifecycle = "production"
      system = "main"
      type = "ops"
      tags = "foo"
    }
    service_tier = 1
}
resource "firehydrant_service" "backups" {
    name   = "backups"
    description = "Another AWS S3 bucket"
    labels = {
      lifecycle = "production"
      system = "main"
      type = "ops"
      tags = "foo"
    }
    service_tier = 2
}
