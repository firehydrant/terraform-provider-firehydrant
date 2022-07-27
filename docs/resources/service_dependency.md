---
page_title: "FireHydrant Resource: firehydrant_service_dependency"
---

# firehydrant_service_dependency Resource

Adding [FireHydrant service dependencies](https://support.firehydrant.com/hc/en-us/articles/5347782635924-Service-Dependencies) 
to respective services can help establish the network of relationships across your FireHydrant 
service catalog. Any time a service is impacted by an incident, you can quickly discover dependencies 
to accurately understand the scope of impact from an incident. This can allow you to pull in all needed 
responding teams to quickly resolve the incident.

## Example Usage

Basic usage:
```hcl
resource "firehydrant_service" "example-service1" {
  name = "Example Service 1"
}

resource "firehydrant_service" "example-service2" {
  name = "Example Service 2"
}

resource "firehydrant_service_dependency" "example-service-dependency" {
  service_id = firehydrant_service.example-service1.id
  connected_service_id = firehydrant_service.example-service2.id
  notes = "These are example service dependency notes"
}
```

## Argument Reference

The following arguments are supported:

* `connected_service_id` - (Required) The ID of the service to define a dependency for.
  The `connected_service_id` represents a service that is a downstream dependency of the 
  service represented by `service_id`.
* `service_id` - (Required) The ID of the service to define a dependency for.
  The `service_id` represents a service that is an upstream dependency of the
  service represented by `connected_service_id`.
* `notes` - (Optional) Any notes to add to the service dependency.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the service dependency.

## Import

Service dependencies can be imported; use `<SERVICE DEPENDENCY ID>` as the import ID. For example:

```shell
terraform import firehydrant_service_dependency.test 3638b647-b99c-5051-b715-eda2c912c42e
```
