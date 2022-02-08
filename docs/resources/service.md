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

- **name** (String, Required)

### Optional

- **description** (String, Optional)
- **id** (String, Optional) The ID of this resource.
- **service_tier** (Integer, Optional) The Service Tier of this resource - between 1 - 5.
- **labels** (Map of String, Optional)
- **add_on_alert** (Boolean, Optional)
