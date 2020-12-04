---
page_title: "firehydrant_runbook Resource - terraform-provider-firehydrant"
subcategory: ""
description: |-
  
---

# Resource `firehydrant_runbook`





## Schema

### Required

- **name** (String, Required)
- **type** (String, Required)

### Optional

- **description** (String, Optional)
- **id** (String, Optional) The ID of this resource.
- **severities** (Block List) (see [below for nested schema](#nestedblock--severities))
- **steps** (Block List) (see [below for nested schema](#nestedblock--steps))

<a id="nestedblock--severities"></a>
### Nested Schema for `severities`

Required:

- **id** (String, Required) The ID of this resource.


<a id="nestedblock--steps"></a>
### Nested Schema for `steps`

Required:

- **action_id** (String, Required)
- **name** (String, Required)

Optional:

- **automatic** (Boolean, Optional)
- **config** (Map of String, Optional)
- **delation_duration** (String, Optional)
- **repeats** (Boolean, Optional)
- **repeats_duration** (String, Optional)

Read-only:

- **step_id** (String, Read-only)


