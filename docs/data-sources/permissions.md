---
page_title: "FireHydrant Data Source: firehydrant_permissions"
subcategory: ""
---

# firehydrant_permissions Data Source

Use this data source to get information on all available permissions in FireHydrant.

FireHydrant permissions define what actions users and roles can perform within the system. This data source allows you to retrieve all available permissions for use in role definitions or policy configurations.

## Example Usage

Basic usage:
```hcl
data "firehydrant_permissions" "all_permissions" {}

output "permission_slugs" {
  value = data.firehydrant_permissions.all_permissions.permissions[*].slug
}
```

## Argument Reference

This data source does not accept any arguments.

## Attributes Reference

The following attributes are exported:

* `permissions` - A list of all available permissions in FireHydrant.

Each permission in the `permissions` list contains:

* `slug` - The unique identifier for the permission (e.g., "create_alerts", "read_teams").
* `display_name` - The human-readable name for the permission (e.g., "Create Alerts", "Read Teams").
* `description` - A description of what the permission allows.
* `available` - Whether the permission is currently available for use.
