---
page_title: "firehydrant_priority Resource - terraform-provider-firehydrant"
subcategory: ""
description: |-

---

# Resource `firehydrant_priority`

Basic usage:

```hcl
resource "firehydrant_priority" "example-priority1" {
  slug        = "my-example-priority1"
  description = "just a priority example"
  default     = true
}
```

## Schema

### Required

- **slug** (String, Required)

### Optional

- **default** (Boolean, Optional)
- **description** (String, Optional)

### Read-only

- **slug** (String, Read-only) The slug of the priority (also known as "name")
