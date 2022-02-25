---
page_title: "firehydrant_team Resource - terraform-provider-firehydrant"
subcategory: ""
description: |-
  
---

# Resource `firehydrant_team`

## Example Usage

Basic usage:

```hcl
resource "firehydrant_service" "example-service1" {
  name = "my-example-service1"
}

resource "firehydrant_service" "example-service2" {
  name = "my-example-service2"
}

resource "firehydrant_team" "example-team" {
  name        = "my-example-team"
  description = "This is an example team"
  
  service_ids = [
    firehydrant_service.example-service1.id,
    firehydrant_service.example-service2.id
  ]
}
```

## Schema

### Required

- **name** (String, Required) The name of the team.

### Optional

- **description** (String, Optional) A description for the team.
- **service_ids** (Set of String, Optional) A set of IDs of the services this team handles incident response for.
  This value _must not_ be provided if `services` is provided.
- **services** (Block List, Optional) **Deprecated** The services this team handles incident response for.
   (see [below for nested schema](#nestedblock--services)). This value _must not_ be provided if 
   `service_ids` is provided.

### Read-only

- **id** (String, Read-only) The ID of the team.

<a id="nestedblock--services"></a>
### Nested Schema for `services` (Deprecated)

Required:

- **id** (String, Required) The ID of the service.

Read-only:

- **name** (String, Read-only) The name of the service


