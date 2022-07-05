---
page_title: "FireHydrant Data Source: firehydrant_priority"
---

# firehydrant_priority Data Source

Use this data source to get information on priorities.

FireHydrant priorities define when an incident should be addressed.

## Example Usage

Basic usage:
```hcl
data "firehydrant_priority" "example-priority" {
  slug = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `slug` - (Required) The slug representing the priority.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the priority. This is the same as the slug.
* `default` - Indicates whether the priority should be the default priority for incidents.
* `description` - A description of the priority.
