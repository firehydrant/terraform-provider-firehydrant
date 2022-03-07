---
page_title: "firehydrant_service_dependency Resource - terraform-provider-firehydrant"
subcategory: ""
description: |-
FireHydrant service_dependencies are used to describe a link between two services.
---

# Resource `firehydrant_service_dependency`

FireHydrant environments are used to tag incidents with where they are occurring.



## Schema

### Required

- **service_id** (String, Required)
- **connecting_service_id** (String, Required)

### Optional

- **notes** (String, Optional)
