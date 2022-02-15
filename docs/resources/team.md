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
  services {
    id = firehydrant_service.example-service1.id
  }
  services {
    id = firehydrant_service.example-service2.id
  }
}
```




## Schema

### Required

- **name** (String, Required)

### Optional

- **description** (String, Optional)
- **services** (Block List) (see [below for nested schema](#nestedblock--services))

### Read-only

- **id** (String, Read-only) The ID of the team.

<a id="nestedblock--services"></a>
### Nested Schema for `services`

Required:

- **id** (String, Required) The ID of the service.

Read-only:

- **name** (String, Read-only)


