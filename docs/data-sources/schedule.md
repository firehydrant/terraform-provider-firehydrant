---
page_title: "FireHydrant Data Source: firehydrant_schedule"
---

# firehydrant_schedule Data Source

Use this data source to get information on schedules.

## Example Usage

Basic usage:
```hcl
data "firehydrant_schedule" "my-oncall-schedule" {
  email = "Main Oncall Schedule"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the oncall schedule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the schedule. 
