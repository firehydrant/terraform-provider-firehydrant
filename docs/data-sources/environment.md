---
page_title: "FireHydrant Data Source: firehydrant_environment"
subcategory: "Beta"
---

# firehydrant_environment Data Source

Use this data source to get information on environments.

FireHydrant environments let you break up your app by region (for example, "US East 1") 
or development stage (for example, "Production").

## Example Usage

Basic usage:
```hcl
data "firehydrant_environment" "example-environment" {
  environment_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `environment_id` - (Required) The ID of the environment.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the environment.
* `description` - A description of the environment.
* `name` - The name of the environment.
