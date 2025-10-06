---
page_title: "FireHydrant Data Source: firehydrant_incident_type"
---

# firehydrant_incident_type Data Source

Use this data source to get information on incident types.

FireHydrant incident types allow you to use predefined incident settings during your 
incidents.  These predefined settings let you control the default severity, priority, 
runbooks, and other aspects of an incident based on a single setting.


## Example Usage

Basic usage:
```hcl
data "firehydrant_incident_type" "example-incident-type" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the incident type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the incident type.
* `name` - The name of the incident type.
* `description` - A description of the incident type.
* `template` - The template used for the incident type.

* `template.0.description` - The description for the template used for the incident type.
* `template.0.customer_impact_summary` - A brief statement describing how this incident type impacts customers.
* `template.0.severity_slug` - The slug for the severity to be used for incidents of this type.
* `template.0.priority_slug` - The slug for the priority to be used for incidents of this type.
* `template.0.private_incident` - A boolean to indicate if incidents of this type should be made private.
* `template.0.tags` - The tags to be applied to incidents of this type.
* `template.0.runbook_ids` - The runbooks to be attached to incidents of this type.
* `template.0.team_ids` - The teams to be added to incidents of this type.
