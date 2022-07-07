---
page_title: "FireHydrant Data Source: firehydrant_runbook"
subcategory: "Beta"
---

# firehydrant_runbook Data Source

Use this data source to get information on runbooks.

FireHydrant runbooks allow you to configure and automate your incident response process by defining a workflow 
to be followed when an incident occurs. Runbooks actually initiate actions that are fundamental steps to 
resolving an incident. Such actions might be creating a Slack channel, starting a Zoom meeting, or opening 
a Jira ticket. Think of a runbook as an incident response playbook that runs (or is activated) when
an incident is declared. Using runbooks, you can equip your team with updated information and best practices 
for mitigating incidents.

## Example Usage

Basic usage:
```hcl
data "firehydrant_runbook" "example-runbook" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the runbook.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the runbook.
* `description` - A description of the runbook.
* `owner_id` - The ID of the team that owns this runbook.
* `name` - The name of the runbook.
