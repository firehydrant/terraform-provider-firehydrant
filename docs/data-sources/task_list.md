---
page_title: "FireHydrant Data Source: firehydrant_task_list"
---

# firehydrant_task_list Data Source

Use this data source to get information on task lists.

FireHydrant task lists help ensure integrity in your incident response.
With tasks lists, you can create predefined tasks for your responders to
reduce cognitive load during an incident. When used during an incident,
each task list item will be assigned as a separate task.


## Example Usage

Basic usage:
```hcl
data "firehydrant_task_list" "example-task-list" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the task list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the task list.
* `description` - A description of the task list.
* `name` - The name of the task list.
* `task_list_items` - A list of tasks in the task list.

The `task_list_items` block contains:

* `summary` - (Required) A summary of the task.
* `description` - (Optional) A description of the task list.
