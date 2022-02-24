---
page_title: "firehydrant_service Resource - terraform-provider-firehydrant"
subcategory: ""
description: |-

---

# Resource `firehydrant_service`

## Example Usage

Basic usage:

```hcl
resource "firehydrant_service" "example-service" {
  name         = "my-example-service"
  add_on_alert = true
  description  = "The main service for our company"
  labels = {
    language  = "ruby",
    lifecycle = "production"
    system    = "main"
    type      = "user"
    tags      = "foo; bar; baz"
  }
  service_tier = 1
}
```

## Schema

### Required

- **name** (String, Required) The name of the service.

### Optional

- **add_on_alert** (Boolean, Optional) Indicates if FireHydrant should automatically create 
   an alert based on the integrations set up for this service, if this service is added to an 
   active incident. Defaults to `false`.
- **description** (String, Optional) A description for the service.
- **labels** (Map of String, Optional) Key-value pairs associated with the service. Useful for 
   supporting searching and filtering of the service catalog.
- **owner_id** (String, Optional) The ID of the team that owns this service.
- **service_tier** (Integer, Optional) The service tier of this resource - between 1 - 5. 
   Lower values represent higher criticality. Defaults to `5`.
- **team_ids** (Set of String, Optional) A set of IDs of the teams responsible for this service's incident response.

### Read-only

- **id** (String, Read-only) The ID of the service.
