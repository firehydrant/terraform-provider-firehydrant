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

- **id** (String, Required) The ID of this resource.

### Optional

- **service_tier** (Integer, Optional) The Service Tier of this resource - between 1 - 5.
- **add_on_alert** (Boolean, Optional)
-
### Read-only

- **description** (String, Read-only)
- **name** (String, Read-only)
- **add_on_alert** (Boolean, Optional)
