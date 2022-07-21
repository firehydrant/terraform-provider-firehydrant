---
page_title: "FireHydrant Resource: firehydrant_task_list"
---

# firehydrant_task_list Resource

FireHydrant task lists help ensure integrity in your incident response. 
With tasks lists, you can create predefined tasks for your responders to 
reduce cognitive load during an incident. When used during an incident, 
each task list item will be assigned as a separate task.

## Example Usage

Basic usage:
```hcl
resource "firehydrant_task_list" "example-task-list" {
  name        = "example-task-list"
  description = "This is an example task list"

  task_list_items {
    summary = "Example task #1"
  }

  task_list_items {
    summary     = "Example task #2"
    description = "This task is very important."
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the task list.
* `task_list_items` - (Required) A list of tasks to include in the task list.
* `description` - (Optional) A description of the task list.

The `task_list_items` block supports:

* `summary` - (Required) A summary of the task.
* `description` - (Optional) A description of the task list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the task list.

## Import

Task Lists can be imported; use `<TASK LIST ID>` as the import ID. For example:

```shell
terraform import firehydrant_task_list.test 3638b647-b99c-5051-b715-eda2c912c42e
```
