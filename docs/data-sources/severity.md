---
page_title: "FireHydrant Data Source: firehydrant_severity"
---

# firehydrant_severity Data Source

Use this data source to get information on severities.

FireHydrant severities represent different levels of severity for incidents.

## Example Usage

Basic usage:
```hcl
data "firehydrant_severity" "example-severity" {
  slug = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `slug` - (Required) The slug representing the severity.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the severity. This is the same as the slug.
* `description` - A description of the severity.
* `type` - The type of the severity.
  Possible values are `gameday`, `maintenance`, and `unexpected_downtime`.
