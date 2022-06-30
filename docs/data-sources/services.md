---
page_title: "FireHydrant Data Source: firehydrant_services"
subcategory: ""
---

# firehydrant_services Data Source

Use this data source to get information on all services matching the given criteria.

FireHydrant services are collections of functions that perform a targeted business operation.
A service can be a microservice, a mono-repository, or a library that you maintain.

## Example Usage

Basic usage:
```hcl
data "firehydrant_services" "all-services" {
}
```

Getting all services with `database` in the name:
```hcl
data "firehydrant_services" "database-named-services" {
  query = "database"
}
```

Getting all services with the label `managed: true`:
```hcl
data "firehydrant_services" "managed-true-labeled-services" {
  labels = {
    managed = "true"
  }
}
```

## Argument Reference

The following arguments are supported:

* `labels` - (Optional) Labels on the runbooks being searched for.
* `query` - (Optional) A query to search for runbooks by their name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `services` - All the services matching the criteria specified by `labels` and `query`.

The `services` block contains:

* `id` - The ID of the service.
* `alert_on_add` - Indicates if FireHydrant should automatically create an alert
  based on the integrations set up for this service, if this service is added to
  an active incident.
* `description` - A description of the service.
* `labels` - Key-value pairs associated with the service. Useful for supporting
  searching and filtering of the service catalog.
* `links` - Links associated with the service
* `name` - The name of the service.
* `owner_id` - The ID of the team that owns this service.
* `service_tier` - The service tier of this resource - between 1 - 5.
  Lower values represent higher criticality.
* `team_ids` - A set of IDs of the teams responsible for this service's incident
  response.

The `links` block contains:

* `href_url` - The URL for the link.
* `name` - The name of the link.
