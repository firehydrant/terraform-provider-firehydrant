---
page_title: "FireHydrant Resource: firehydrant_team"
subcategory: ""
---

# firehydrant_team Resource

FireHydrant teams are collections of people that can be assigned to incidents 
and configured as owners of various resources, like services and runbooks.

## Example Usage

Basic usage:
```hcl
resource "firehydrant_team" "example-team" {
  name        = "example-team"
  description = "This is an example team"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the team.
* `description` - (Optional) A description for the team.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the team.

## Import

Teams can be imported; use `<TEAM ID>` as the import ID. For example:

```shell
terraform import firehydrant_team.test 3638b647-b99c-5051-b715-eda2c912c42e
```
