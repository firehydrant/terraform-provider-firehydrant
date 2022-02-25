---
page_title: "firehydrant_functionality Data Source - terraform-provider-firehydrant"
subcategory: ""
description: |-
  
---

# Data Source `firehydrant_functionality`

## Example Usage

Basic usage:

```hcl
data "firehydrant_functionality" "example-functionality" {
  functionality_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

## Schema

### Required

- **functionality_id** (String, Required) The ID of the functionality.

### Read-only

- **id** (String, Read-only) The ID of the functionality.
- **description** (String, Read-only)
- **name** (String, Read-only) The name of the functionality.
- **description** (String, Read-only) A description for the functionality.
- **service_ids** (Set of String, Read-only) A set of IDs of the services this functionality is associated with.
