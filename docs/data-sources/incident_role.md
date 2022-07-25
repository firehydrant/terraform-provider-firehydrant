---
page_title: "FireHydrant Data Source: firehydrant_incident_role"
---

# firehydrant_incident_role Data Source

Use this data source to get information on incident roles.

FireHydrant incident roles allow you to use predefined roles during your 
incidents. These predefined roles help responders know exactly what their 
responsibilities are as soon as they drop into an incident.


## Example Usage

Basic usage:
```hcl
data "firehydrant_incident_role" "example-incident-role" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the incident role.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the incident role.
* `description` - A description of the incident role.
* `name` - The name of the incident role.
* `summary` - A summary of the incident role.