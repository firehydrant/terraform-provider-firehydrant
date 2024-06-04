---
page_title: "FireHydrant Resource: firehydrant_status_update_template"
---

# firehydrant_status_update_template Resource

FireHydrant status update templates are used to create preset templates for status page updates.

## Example Usage

Basic usage:
```hcl
resource "firehydrant_status_update_template" "my_template" {
  name = "My Template"
  body = "We are investigating reports of **ADD_WHAT_YOU_ARE_INVESTIGATING**"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the status update template.
* `body` - (Required) The body text of the template.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the status update template.
