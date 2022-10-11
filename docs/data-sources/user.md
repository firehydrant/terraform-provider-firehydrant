---
page_title: "FireHydrant Data Source: firehydrant_user"
---

# firehydrant_user Data Source

Use this data source to get information on users.

## Example Usage

Basic usage:
```hcl
data "firehydrant_user" "my-user" {
  email = "my-user@firehydrant.io"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) The user's email address.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the user. 
