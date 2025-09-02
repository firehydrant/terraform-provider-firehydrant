---
page_title: "FireHydrant Data Source: firehydrant_user_permissions"
subcategory: ""
---

# firehydrant_user_permissions Data Source

Use this data source to get information on the current user's permissions in FireHydrant.

This data source retrieves all permissions that the currently authenticated user has access to. This is useful for understanding what actions the current user can perform and for building user-specific role configurations.

## Example Usage

Basic usage:
```hcl
data "firehydrant_user_permissions" "mine" {}

output "my_permissions" {
  value = data.firehydrant_user_permissions.mine.permissions[*].slug
}
```

## Argument Reference

This data source does not accept any arguments.

## Attributes Reference

The following attributes are exported:

* `permissions` - A list of permissions that the current user has access to.

Each permission in the `permissions` list contains:

* `slug` - The unique identifier for the permission (e.g., "create_alerts", "read_teams").
* `display_name` - The human-readable name for the permission (e.g., "Create Alerts", "Read Teams").
* `description` - A description of what the permission allows.
* `available` - Whether the permission is currently available for use.
