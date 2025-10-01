---
page_title: "FireHydrant Resource: firehydrant_custom_event_source"
subcategory: "Signals"
---

# firehydrant_custom_event_source Resource

FireHydrant custom event sources allow users to submit javascript to transpose signals from their custom event source to a format compatible with Signals.

## Example Usage

Basic usage:
```hcl
resource "firehydrant_custom_event_source" "example_transposer" {
  name = "Example Transposer"
	slug = "example-transposer"
	description = "Maps alerts produced by the Example service to FireHydrant Signals"
	javascript = "function transpose(input) {\n  return input.foo;\n}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the custom event source.
* `slug` - (Required) A unique identifier for the custom event source.  Only lowercase letters and `-` are allowed.
* `description` - (Optional) A description of the custom event source.
* `javascript` - (Required) The javascript expression that handles the transposition to a Firehydrant Signal.  See https://docs.firehydrant.com/docs/custom-event-source#getting-started for details and example javascript.


## Import

Custom Event Sources can be imported; use `slug` as the import ID. For example:

```shell
terraform import firehydrant_custom_event_source.example_transposer example-transposer
```
