---
page_title: "firehydrant_team Resource - terraform-provider-firehydrant"
subcategory: ""
description: |-
  
---

# Resource `firehydrant_team`





## Schema

### Required

- **name** (String, Required)

### Optional

- **description** (String, Optional)
- **id** (String, Optional) The ID of this resource.
- **services** (Block List) (see [below for nested schema](#nestedblock--services))

<a id="nestedblock--services"></a>
### Nested Schema for `services`

Required:

- **id** (String, Required) The ID of this resource.

Read-only:

- **name** (String, Read-only)


