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

data "firehydrant_service" "lake-service" {
	id = "lake-service"
}

output "monolith_name" {
  value = data.firehydrant_service.monolith.name
}

output "lake_name" {
  value = data.firehydrant_service.monolith.name
}

resource "firehydrant_service" "my-new-service" {
  name = "Pond Service"
  description = "Lakes provide ecological benefits to a complex system"
}
