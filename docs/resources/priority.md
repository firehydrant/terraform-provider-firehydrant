---
page_title: "firehydrant_priority Resource - terraform-provider-firehydrant"
subcategory: ""
description: |-

---

# Resource `firehydrant_priority`

## Example Usage

Basic usage:

```hcl
resource "firehydrant_priority" "example-priority" {
  slug        = "MYEXAMPLEPRIORITY"
  description = "This is an example priority"
  default     = true
}
```

## Schema

### Required

- **slug** (String, Required) The slug representing the priority. It must be unique and only contain alphanumeric characters. The slug cannot be longer than 23 characters.

### Optional

- **default** (Boolean, Optional) Indicates whether the priority should be the default 
  priority for incidents. At most one resource can have default set to `true`. Setting default to `true` for multiple priority resources will result in inconsistent plans in Terraform. Defaults to `false`.
- **description** (String, Optional) A description for the priority.

### Read-only

- **id** (String, Read-only) The ID of the priority. This is the same as the slug.
