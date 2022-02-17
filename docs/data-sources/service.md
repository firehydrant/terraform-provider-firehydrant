---
page_title: "firehydrant_service Data Source - terraform-provider-firehydrant"
subcategory: ""
description: |-

---

# Data Source `firehydrant_service`

## Example Usage

Basic usage:
```hcl
data "firehydrant_service" "example-service" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- **id** (String, Required) The ID of the service.

### Read-only

- **add_on_alert** (Boolean, Read-only) Indicates if FireHydrant should automatically create
  an alert based on the integrations set up for this service, if this service is added to an
  active incident. Defaults to `false`.
- **description** (String, Read-only) A description for the service.
- **name** (String, Read-only) The name of the service.
- **owner_id** (String, Read-only) The ID of the team that owns this service.
- **service_tier** (Integer, Read-only) The service tier of this resource - between 1 - 5.
  Lower values represent higher criticality. Defaults to `5`.
- **team_ids** (Set of String, Optional) A set of IDs of the teams responsible for this service's incident response.

