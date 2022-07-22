---
page_title: "FireHydrant Resource: firehydrant_incident_role"
---

# firehydrant_incident_role Resource

FireHydrant incident roles allow you to use predefined roles during your
incidents. These predefined roles help responders know exactly what their
responsibilities are as soon as they drop into an incident.

## Example Usage

Basic usage:
```hcl
resource "firehydrant_incident_role" "example-incident-role" {
  name        = "Example Incident Role"
  description = "This is an example incident role description"
  summary     = "This is an example incident role summary"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the incident role.
* `summary` - (Required) A summary of the incident role.
* `description` - (Optional) A description of the incident role.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the incident role.

## Import

Incident roles can be imported; use `<INCIDENT ROLE ID>` as the import ID. For example:

```shell
terraform import firehydrant_incident_role.test 3638b647-b99c-5051-b715-eda2c912c42e
```