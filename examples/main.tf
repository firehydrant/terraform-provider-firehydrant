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

data "firehydrant_service" "chow-hall" {
  id = "chow-hall"
}

output "chow-hall-name" {
  value = data.firehydrant_service.chow-hall.name
}

resource "firehydrant_service" "paddy-cake-paddy-cake" {
  name = data.firehydrant_environment.production.name
}

data "firehydrant_environment" "production" {
  environment_id = "ef2fe740-df17-4d57-9cc0-ecae7c8e9594"
}

output "production-name" {
  value = data.firehydrant_environment.production.name
}

data "firehydrant_functionality" "cloud" {
  functionality_id = "211cb0c3-cf50-4723-a7e3-7f87f7883f88"
}

output "func-cloud-name" {
  value = data.firehydrant_functionality.cloud.name
}
