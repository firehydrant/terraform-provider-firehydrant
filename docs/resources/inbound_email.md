---
page_title: "FireHydrant Resource: firehydrant_inbound_email"
subcategory: "Signals"
---

# firehydrant_inbound_email Resource

FireHydrant inbound email resources allow you to configure email-based alerts by sending emails to the FireHydrant platform.

## Example Usage

```hcl
data "firehydrant_team" "example_team" {
  name = "Example Team"
}

resource "firehydrant_inbound_email" "example" {
  name                   = "Inbound Email Alert"
  slug                   = "inbound-email-alert"
  description            = "Inbound email alert for critical issues"
  status_cel             = "email.body.contains('has recovered') ? 'CLOSED' : 'OPEN'"
  level_cel              = "email.body.contains('panic') ? 'ERROR' : 'INFO'"
  allowed_senders        = ["@firehydrant.com", "@example.com"]
  target {
    type = "Team"
    id   = data.firehydrant_team.example_team.id
  }
  rules                  = ["email.body.contains(\"critical\")"]
  rule_matching_strategy = "all"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the inbound email resource.
* `slug` - (Required) The slug for the inbound email resource.
* `description` - (Optional) A description of the inbound email resource.
* `status_cel` - (Required) A Common Expression Language (CEL) expression to determine the status of the alert based on the email content.
* `level_cel` - (Required) A CEL expression to determine the severity level of the alert based on the email content.
* `allowed_senders` - (Required) A list of email domains or addresses allowed to send alerts.
* `target` - (Required) A block to specify the target for the alert. The block supports:
  * `type` - (Required) The type of the target. Valid values are "Team", "User", or "EscalationPolicy".
  * `id` - (Required) The ID of the target resource.
* `rules` - (Required) A list of CEL expressions that define when an alert should be triggered.
* `rule_matching_strategy` - (Required) The strategy for matching rules. Valid values are "all" or "any".

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the inbound email resource.
* `email` - The email address to send alerts to.

## Import

Inbound email resources can be imported using the resource ID, e.g.,

```
$ terraform import firehydrant_inbound_email.example 12345678-90ab-cdef-1234-567890abcdef
```
