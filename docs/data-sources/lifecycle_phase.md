---
page_title: "FireHydrant Data Source: firehydrant_lifecycle_phase"
subcategory: ""
---

# firehydrant_lifecycle_phase Data Source

Use this data source to get the phase ID for a lifecycle phase, needed for creating or updating lifecycle milestones.

Lifecycle phases are identified by name and are not editable as of this version.  The existing lifecycle phase names are (in order) `started`, `active`, `post-incident`, and `closed`.

## Example Usage

Basic usage:
```hcl
data "firehydrant_lifecycle_phase" "started" {
  name = "started"
}

resource "firehydrant_lifecycle_milestone" "new_milestone" {
  name        = "New Milestone"
  description = "A new lifecycle milestone"
	phase_id    = data.firehydrant_lifecycle_phase.started.id
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the lifecycle milestone.The existing lifecycle phase names are (in order) `started`, `active`, `post-incident`, and `closed`.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the lifecycle phase.
