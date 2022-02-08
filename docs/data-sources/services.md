---
page_title: "firehydrant_services Data Source - terraform-provider-firehydrant"
subcategory: ""
description: |-

---

# Data Source `firehydrant_services`

## Example Usage

Basic usage:
```hcl
data "firehydrant_services" "all-services" {
}
```

Getting all services with `database` in the name:
```hcl
data "firehydrant_services" "database-named-services" {
  query = "database"
}
```

Getting all services with the label `managed: true`:
```hcl
data "firehydrant_services" "managed-true-labeled-services" {
  labels = {
    managed = "true"
  }
}
```

## Schema

### Optional

- **id** (String, Optional) The ID of this resource.
- **labels** (Map of String, Optional)
- **query** (String, Optional)

### Read-only

- **services** (List of Object, Read-only) (see [below for nested schema](#nestedatt--services))

<a id="nestedatt--services"></a>
### Nested Schema for `services`

- **description** (String)
- **id** (String)
- **name** (String)
