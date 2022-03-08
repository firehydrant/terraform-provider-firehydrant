---
page_title: "firehydrant_functionality Resource - terraform-provider-firehydrant"
subcategory: ""
description: |-
  
---

# Resource `firehydrant_functionality`

## Example Usage

Basic usage:

```hcl
resource "firehydrant_service" "example-service1" {
  name = "my-example-service1"
}

resource "firehydrant_service" "example-service2" {
  name = "my-example-service2"
}

resource "firehydrant_functionality" "example-functionality" {
  name        = "my-example-functionality"
  description = "This is an example functionality"
  
  service_ids = [
    firehydrant_service.example-service1.id,
    firehydrant_service.example-service2.id
  ]
}
```

## Schema

### Required

- **name** (String, Required) The name of the functionality.

### Optional

- **description** (String, Optional) A description for the functionality.
- **service_ids** (Set of String, Optional) A set of IDs of the services this functionality is associated with.
  This value _must not_ be provided if `services` is provided.
- **services** (Block List, Optional) **Deprecated** The services this functionality is associated with 
   (see [below for nested schema](#nestedblock--services)). This value _must not_ be provided if
  `service_ids` is provided.

### Read-only

- **id** (String, Read-only) The ID of the functionality.

<a id="nestedblock--services"></a>
### Nested Schema for `services` (Deprecated)

Required:

- **id** (String, Required) The ID of the service.

Read-only:

- **name** (String, Read-only) The name of the service.
