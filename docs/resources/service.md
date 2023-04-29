---
page_title: "FireHydrant Resource: firehydrant_service"
subcategory: ""
---

# firehydrant_service Resource

FireHydrant services are collections of functions that perform a targeted business operation.
A service can be a microservice, a mono-repository, or a library that you maintain.

## Example Usage

Basic usage:
```hcl
resource "firehydrant_team" "example-owner-team" {
  name        = "my-example-owner-team"
  description = "This is an example team that owns a service"
}

resource "firehydrant_team" "example-responder-team1" {
  name        = "my-example-responder-team1"
  description = "This is an example team that is responsible for responding to incidents for a service"
}

resource "firehydrant_team" "example-responder-team2" {
  name        = "my-example-responder-team2"
  description = "This is an example team that is responsible for responding to incidents for a service"
}

resource "firehydrant_service" "example-service" {
  name         = "example-service"
  alert_on_add = true
  description  = "This is an example service"

  labels = {
    language  = "go",
    lifecycle = "production"
  }

  links {
    href_url = "https://example.com/internal-dashboard"
    name     = "Internal Dashboard"
  }

  owner_id     = firehydrant_team.example-owner-team.id
  service_tier = 1

  team_ids = [
    firehydrant_team.example-responder-team1.id,
    firehydrant_team.example-responder-team2.id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the service.
* `alert_on_add` - (Optional) Indicates if FireHydrant should automatically create
  an alert based on the integrations set up for this service, if this service is added to an
  active incident. Defaults to `false`.
* `auto_add_responding_team` - (Optional) Indicates if FireHydrant should automatically add
  the responding team if this service is added to an active incident. Defaults to `false`.
* `external_resources` - (Optional) External resources associated with the service
* `description` - (Optional) A description for the service.
* `labels` - (Optional) Key-value pairs associated with the service. Useful for
  supporting searching and filtering of the service catalog.
* `links` - (Optional) Links associated with the service
* `owner_id` - (Optional) The ID of the team that owns this service.
* `service_tier` - (Optional) The service tier of this resource - between 1 - 5.
  Lower values represent higher criticality. Defaults to `5`.
* `team_ids` - (Optional) A set of IDs of the teams responsible for this service's incident
  response.

The `links` block supports:

* `href_url` - (Required) The URL to use for the link.
* `name` - (Required) The name of the link.

The `exteral_resources` block supports:

* `remote_id` - (Required) The ID of the resource in the remote provider
* `connection_type` - (Required) The connection type configured in FireHydrant (`pager_duty` for example)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the service.

## Import

Services can be imported; use `<SERVICE ID>` as the import ID. For example:

```shell
terraform import firehydrant_service.test 3638b647-b99c-5051-b715-eda2c912c42e
```
