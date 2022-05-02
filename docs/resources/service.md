---
page_title: "firehydrant_service Resource - terraform-provider-firehydrant"
subcategory: ""
description: |-

---

# Resource `firehydrant_service`

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
  name         = "my-example-service"
  alert_on_add = true
  description  = "The main service for our company"

  labels = {
    language  = "ruby",
    lifecycle = "production"
    system    = "main"
    type      = "user"
    tags      = "foo; bar; baz"
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

## Schema

### Required

- **name** (String, Required) The name of the service.

### Optional

- **alert_on_add** (Boolean, Optional) Indicates if FireHydrant should automatically create 
  an alert based on the integrations set up for this service, if this service is added to an 
  active incident. Defaults to `false`.
- **description** (String, Optional) A description for the service.
- **labels** (Map of String, Optional) Key-value pairs associated with the service. Useful for 
  supporting searching and filtering of the service catalog.
- **links** (Set of Map, Optional) Links associated with the service (see [below for nested schema](#nestedatt--links)).
- **owner_id** (String, Optional) The ID of the team that owns this service.
- **service_tier** (Integer, Optional) The service tier of this resource - between 1 - 5. 
  Lower values represent higher criticality. Defaults to `5`.
- **team_ids** (Set of String, Optional) A set of IDs of the teams responsible for this service's incident response.
<a id="nestedatt--links"></a>
### Nested Schema for `links`

- **href_url** (String, Required) The URL to use for the link.
- **name** (String, Required) The name of the link.

### Read-only

- **id** (String, Read-only) The ID of the service.
