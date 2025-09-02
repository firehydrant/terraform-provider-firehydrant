---
page_title: "FireHydrant Data Source: firehydrant_team_permissions"
subcategory: ""
---

# firehydrant_team_permissions Data Source

Use this data source to get information on a team's permissions in FireHydrant.

This data source retrieves all permissions that a specific team has access to. This is useful for understanding what actions team members can perform and for building team-specific role configurations.

## Example Usage

Basic usage:
```hcl
data "firehydrant_team_permissions" "example_team" {
  team_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

output "team_permissions" {
  value = data.firehydrant_team_permissions.example_team.permissions[*].slug
}
```

## Argument Reference

The following arguments are supported:

* `team_id` - (Required) The ID of the team to retrieve permissions for.

## Attributes Reference

The following attributes are exported:

* `permissions` - A list of permissions that the specified team has access to.

Each permission in the `permissions` list contains:

* `slug` - The unique identifier for the permission (e.g., "create_alerts", "read_teams").
* `display_name` - The human-readable name for the permission (e.g., "Create Alerts", "Read Teams").
* `description` - A description of what the permission allows.
* `available` - Whether the permission is currently available for use.
