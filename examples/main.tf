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

data "firehydrant_service" "monolith" {
	id = "2177ce81-b6b6-4063-af73-6c881c8b9899"
}

output "monolith_name" {
  value = data.firehydrant_service.monolith
}
