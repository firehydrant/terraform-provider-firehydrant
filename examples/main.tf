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
