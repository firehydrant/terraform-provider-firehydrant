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
data "firehydrant_user" "my-user" {
  email = "example@firehydrant.io"
}

data "firehydrant_incident_role" "ops-lead" {
  id = "1c679abe-3060-47d4-ab5e-e1ecbd5ce55f"
}

resource "firehydrant_team" "example-team" {
  name        = "example-team"
  description = "This is an example team"

  memberships {
    user_id          = data.firehydrant_user.my-user.id
    incident_role_id = data.firehydrant_incident_role.ops-lead.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the team.
* `description` - (Optional) A description for the team.
* `memberships` - (Optional) A resource to tie a schedule or user to a team via a incident role.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the team.

## Import

Teams can be imported; use `<TEAM ID>` as the import ID. For example:

```shell
terraform import firehydrant_team.test 3638b647-b99c-5051-b715-eda2c912c42e
```
