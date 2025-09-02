---
page_title: "FireHydrant Resource: firehydrant_role"
subcategory: ""
---

# firehydrant_role Resource

FireHydrant roles define sets of permissions that can be assigned to users and teams. This resource allows you to create, update, and manage custom roles with specific permission sets.

## Example Usage

Basic usage:
```hcl
resource "firehydrant_role" "example_role" {
  name        = "example-role"
  description = "A custom role with basic permissions"
  permissions = [
    "read_teams",
    "read_users"
  ]
}
```

Advanced usage with comprehensive permissions:
```hcl
resource "firehydrant_role" "incident_manager" {
  name        = "incident-manager"
  description = "Role for managing incidents and alerts"
  permissions = [
    "read_alerts",
    "create_alerts",
    "read_escalation_policies",
    "read_on_call_schedules",
    "read_teams",
    "read_users",
    "read_incident_settings",
    "read_integrations",
    "read_incidents",
    "read_webhooks",
    "read_runbooks",
    "read_status_templates",
    "read_audiences",
    "read_change_events",
    "read_organization_settings",
    "read_service_catalog",
    "read_analytics",
    "read_alert_rules",
    "read_call_routes",
    "read_support_hours"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the role.
* `description` - (Optional) A description of the role.
* `permissions` - (Required) A list of permission slugs to assign to this role. Note that some permissions may have dependencies that require other permissions to be included.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the role.
* `slug` - The slug of the role (auto-generated from the name).
* `organization_id` - The ID of the organization this role belongs to.
* `built_in` - Whether this is a built-in role provided by FireHydrant (always `false` for custom roles).
* `read_only` - Whether this role can be modified (always `false` for custom roles).
* `created_at` - When the role was created.
* `updated_at` - When the role was last updated.

## Import

Roles can be imported; use `<ROLE ID>` as the import ID. For example:

```shell
terraform import firehydrant_role.example 3638b647-b99c-5051-b715-eda2c912c42e
```

## Notes

* **Permission Dependencies**: Some permissions have dependencies on other permissions. For example, `create_alerts` requires `read_alerts` and several other read permissions. The provider will validate these dependencies when creating or updating roles.
* **Built-in Roles**: Built-in roles provided by FireHydrant (like `member`, `admin`) cannot be modified or deleted through this resource.
* **Permission Validation**: The provider validates that all specified permissions exist and are available before creating or updating a role.
