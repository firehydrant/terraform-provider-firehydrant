---
page_title: "FireHydrant Resource: firehydrant_lifecycle_milestone"
---

# firehydrant_lifecycle_milestone Resource

Incident lifecycle milestones describe the current status of the incident and communicate to stakeholders the team's progress in resolving the issue.

As responders work through incidents on FireHydrant, they will typically transition the Milestone, and FireHydrant automatically logs the timestamps of these changes. This allows FireHydrant to collect data for holistic incident metrics out-of-the-box like MTT*, Impacted Infrastructure, Responder Impact, and so on.

## Example Usage

Basic usage:
```hcl
data "firehydrant_lifecycle_phase" "started" {
  name = "started"
}

resource "firehydrant_lifecycle_milestone" "new_milestone" {
  name        = "New Milestone"
  description = "This is a new lifecycle milestone"
	phase_id    = data.firehydrant_lifecycle_phase.started.id
	slug        = "new-milestone"
	position    = 2
	auto_assign_timestamp_on_create = "never_set_on_create"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the lifecycle milestone.
* `description` - (Required) A description of the lifecycle milestone.
* `phase_id` - (Required) The phase to which the lifecycle milestone belongs.
* `slug` - (Optional) The internal slug used to reference the lifecycle milestone.  If not entered, one will be created.
* `position` - (Optional) The position of the milestone within the phase. If not provided, the milestone will be added as the last milestone in the phase..
* `auto_assign_timestamp_on_create` - (Optional) The setting for auto-assigning the milestone's timestamp during incident declaration.  Must be set to one of `always_set_on_create`, `only_set_on_manual_create`, `never_set_on_create`.  Defaults to `never_set_on_create` if not set.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the lifecycle milestone.

## Import

Lifecycle milestones can be imported; use `<MILESTONE ID>` as the import ID. For example:

```shell
terraform import firehydrant_lifecycle_milestone.new_milestone 7ac44cdf-caf5-4613-b05d-c649e6de8548
```
