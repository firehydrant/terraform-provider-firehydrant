terraform {
  required_providers {
    firehydrant = {
      source  = "firehydrant/firehydrant"
      version = "0.1.4"
    }
  }
}

provider "firehydrant" {
  api_key              = "fhb-00000000000000000000000000000000"
  firehydrant_base_url = "https://api.production.firehydrant.io/v1/"
}

resource "firehydrant_environment" "production" {
  name = "Production"
}
