---
page_title: "firehydrant_service_dependency Data Source - terraform-provider-firehydrant"
subcategory: ""
description: |-

---

# Data Source `firehydrant_service_dependency`

## Example Usage
Basic usage:
```hcl
data "firehydrant_service_dependency" "example-service-dependency" {
  id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```
## Schema

### Required

- **service_dependency_id** (String, Required) The UUID of this resource.

### Optional

### Read-only
- **service_id** (String, Read-only) The UUID of a service resource.
- **connected_service_id** (String, Read-only) The UUID of a service resource.
- **notes** (String, Read-only)

