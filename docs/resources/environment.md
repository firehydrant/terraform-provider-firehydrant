---
page_title: "FireHydrant Resource: firehydrant_environment"
---

# firehydrant_environment Resource

FireHydrant environments let you break up your app by region (for example, "US East 1")
or development stage (for example, "Production").

## Example Usage

Basic usage:
```hcl
resource "firehydrant_environment" "example-environment" {
  name        = "example-environment"
  description = "This is an example environment"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name the environment.
* `description` - (Optional) A description of the environment.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the environment.

## Import

Environments can be imported; use `<ENVIRONMENT ID>` as the import ID. For example:

```shell
terraform import firehydrant_environment.test 3638b647-b99c-5051-b715-eda2c912c42e
```
