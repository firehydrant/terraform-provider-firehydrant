---
page_title: "firehydrant_services Data Source - terraform-provider-firehydrant"
subcategory: ""
description: |-

---

# Data Source `firehydrant_services`

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

## Schema

### Optional

- **labels** (Map of String, Optional)
- **query** (String, Optional)

### Read-only

- **services** (List of Object, Read-only) (see [below for nested schema](#nestedatt--services))
<a id="nestedatt--services"></a>
### Nested Schema for `services`

- **id** (String, Read-only) The ID of the service.
- **add_on_alert** (Boolean, Read-only) Indicates if FireHydrant should automatically create
  an alert based on the integrations set up for this service, if this service is added to an
  active incident. Defaults to `false`.
- **description** (String, Read-only) A description for the service.
- **labels** (Map of String, Read-only) Key-value pairs associated with the service. Useful for
  supporting searching and filtering of the service catalog.
- **links** (Set of Map, Read-only) Links associated with the service (see [below for nested schema](#nestedatt--links)).
- **name** (String, Read-only) The name of the service.
- **owner_id** (String, Read-only) The ID of the team that owns this service.
- **service_tier** (Integer, Read-only) The service tier of this resource - between 1 - 5.
  Lower values represent higher criticality. Defaults to `5`.
- **team_ids** (Set of String, Optional) A set of IDs of the teams responsible for this service's incident response.
<a id="nestedatt--links"></a>
### Nested Schema for `links`

- **href_url** (String, Read-only) The URL to use for the link.
- **name** (String, Read-only) The name of the link.
