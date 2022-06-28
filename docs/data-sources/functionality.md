---
page_title: "FireHydrant Data Source: firehydrant_functionality"
subcategory: ""
---

# firehydrant_functionality Data Source

Use this data source to get information on functionalities.

A functionality (function) is a programming construct that performs a specific task. 
FireHydrant functionalities let you associate backend services with the features your 
end users interact with.

## Example Usage

Basic usage:
```hcl
data "firehydrant_functionality" "example-functionality" {
  functionality_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `functionality_id` - (Required) The ID of the functionality.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the functionality.
* `description` - A description of the functionality.
* `name` - The name of the functionality.
* `service_ids` - A set of IDs of the services this functionality is associated with.
