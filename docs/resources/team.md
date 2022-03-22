---
page_title: "firehydrant_team Resource - terraform-provider-firehydrant"
subcategory: ""
description: |-
  
---

# Resource `firehydrant_team`

## Example Usage

Basic usage:

```hcl
resource "firehydrant_team" "example-team" {
  name        = "my-example-team"
  description = "This is an example team"
}
```

## Schema

### Required

- **name** (String, Required) The name of the team.

### Optional

- **description** (String, Optional) A description for the team.

### Read-only

- **id** (String, Read-only) The ID of the team.
