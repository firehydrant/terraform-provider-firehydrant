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
  name = "we all fall down"
}
