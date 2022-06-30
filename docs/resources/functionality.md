---
page_title: "FireHydrant Resource: firehydrant_functionality"
subcategory: "Beta"
---

# firehydrant_functionality Resource

A functionality (function) is a programming construct that performs a specific task.
FireHydrant functionalities let you associate backend services with the features your
end users interact with.

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

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the functionality.
* `description` - (Optional) A description of the functionality.
* `service_ids` - (Optional) A set of IDs of the services this functionality is associated with.
  This value _must not_ be provided if `services` is provided.
* `services` - **Deprecated** Use `service_ids` instead. The services this functionality is associated with. 
  This value _must not_ be provided if `service_ids` is provided.

**Deprecated** The `services` block supports:

* `id` - (Required) The ID of the service.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the functionality.

**Deprecated** The `services` block contains:

* `id` - The ID of the service.

## Import

Functionalities can be imported; use `<FUNCTIONALITY ID>` as the import ID. For example:

```shell
terraform import firehydrant_functionality.test 3638b647-b99c-5051-b715-eda2c912c42e
```
