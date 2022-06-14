resource "firehydrant_priority" "low" {
  slug        = "LOW"
  description = "low priority"
  default     = true
}

resource "firehydrant_priority" "medium" {
  slug        = "MEDIUM"
  description = "medium priority"
}

resource "firehydrant_priority" "high" {
  slug        = "HIGH"
  description = "high priority"
}

data "firehydrant_priority" "P1" {
  slug        = "P1"
}