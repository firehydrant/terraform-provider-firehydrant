---
page_title: "FireHydrant Resource: firehydrant_severity"
subcategory: "Beta"
---

# firehydrant_severity Resource

FireHydrant severities represent different levels of severity for incidents.

## Example Usage

Basic usage:
```hcl
resource "firehydrant_severity" "example-severity" {
  slug        = "EXAMPLESEVERITY"
  description = "This is an example severity"
}
```

## Argument Reference

The following arguments are supported:

* `slug` - (Required) The slug representing the priority. It must be unique and only contain 
  alphanumeric characters. The slug cannot be longer than 23 characters.
* `description` - (Optional) A description for the severity.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the severity.

## Import

Severities can be imported; use `<SEVERITY ID>` as the import ID. For example:

```shell
terraform import firehydrant_severity.test 3638b647-b99c-5051-b715-eda2c912c42e
```
