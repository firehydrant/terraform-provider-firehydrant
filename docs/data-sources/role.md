---
page_title: "FireHydrant Data Source: firehydrant_role"
subcategory: ""
---

# firehydrant_role Data Source

Use this data source to get information on a specific role in FireHydrant.

FireHydrant roles define sets of permissions that can be assigned to users and teams. This data source allows you to retrieve information about existing roles, including built-in roles and custom roles.

## Example Usage

Basic usage by slug:
```hcl
data "firehydrant_role" "member" {
  slug = "member"
}

output "role_permissions" {
  value = data.firehydrant_role.member.permissions[*].slug
}
```

Lookup by ID:
```hcl
data "firehydrant_role" "custom_role" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

output "role_name" {
  value = data.firehydrant_role.custom_role.name
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the role. Either `id` or `slug` must be specified.
* `slug` - (Optional) The slug of the role. Either `id` or `slug` must be specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the role.
* `name` - The name of the role.
* `slug` - The slug of the role.
* `description` - A description of the role.
* `organization_id` - The ID of the organization this role belongs to.
* `built_in` - Whether this is a built-in role provided by FireHydrant.
* `read_only` - Whether this role can be modified.
* `permissions` - A list of permissions assigned to this role.
* `created_at` - When the role was created.
* `updated_at` - When the role was last updated.

Each permission in the `permissions` list contains:

* `slug` - The unique identifier for the permission.
* `display_name` - The human-readable name for the permission.
* `description` - A description of what the permission allows.
* `available` - Whether the permission is currently available for use.