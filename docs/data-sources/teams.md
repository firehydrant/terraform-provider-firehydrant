---
page_title: "FireHydrant Data Source: firehydrant_teams"
subcategory: ""
---

# firehydrant_teams Data Source

Use this data source to get information on all teams matching the given criteria.

## Example Usage

Basic usage:
```hcl
data "firehydrant_teams" "all-teams" {
}
```

Getting all teams with `database` in the name:
```hcl
data "firehydrant_teams" "database-named-teams" {
  query = "database"
}
```

## Argument Reference

The following arguments are supported:

* `query` - (Optional) A query to search for teams by their name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `teams` - All the teams matching the criteria specified by `query`.

The `teams` block contains:

* `id` - The ID of the team.
* `name` - The name of the team.
* `description` - A description of the team.
* `slug` - The slug for the team.
* `service_ids` - A set of IDs of the services associated with this team
* `owned_service_ids` - A set of IDs of the services owned by this team
