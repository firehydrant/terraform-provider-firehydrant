resource "firehydrant_inbound_email" "example" {
  name                   = "Inbound Email"
  slug                   = "inbound-email"
  description            = "Description here"
  status_cel             = "email.body.contains('has recovered') ? 'CLOSED' : 'OPEN'"
  level_cel              = "email.body.contains('panic') ? 'ERROR' : 'INFO'"
  allowed_senders        = ["@firehydrant.com"]
  target {
    type = "Team"
    id   = "30715262-7fa9-4328-95e2-7cd45433d518"
  }
  rules                  = ["email.body.contains(\"hello\")"]
  rule_matching_strategy = "all"
}
