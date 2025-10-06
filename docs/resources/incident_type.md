---
page_title: "FireHydrant Resource: firehydrant_incident_type"
---

# firehydrant_incident_type Resource

FireHydrant incident types allow you to use predefined types during your
incidents. These predefined settings let you control the default severity, priority, 
runbooks, and other aspects of an incident based on a single setting.

## Example Usage

Basic usage:
```hcl
resource "firehydrant_incident_type" "test_incident_type" {
  name        = "Test incident type"
  description = "This is a test"

	template {
    description = "The template for a test incident"
  }
}
```

Advanced usage with full template:
```hcl
resource "firehydrant_incident_type" "test_incident_type" {
  name        = "Test incident type"
  description = "This is a test"
	
  template {
	  description = "The template for a test incident"
		customer_impact_summary = "This is a test.  Customers will be unaffected."
		severity_slug = "SEV5"
		priority_slug = "TESTPRIORITY"
		private_incident = false

		tags = [ "foo", "bar" ]
		runbook_ids = [ "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" ]
		team_ids = [ "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" ]
		
		impacts {
		    impact_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
				condition_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
		}
		
		impacts {
			  condition_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
				impact_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    }
	}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the incident type.
* `description` - (Optional) A description of the incident type.
* `template` - (Required) The template used for the incident type.

The `template` block supports:

* `description` - (Optional) The description for the template used for the incident type.
* `customer_impact_summary` - (Optional) A brief statement describing how this incident type impacts customers.
* `severity_slug` - (Optional) The slug for the severity to be used for incidents of this type.
* `priority_slug` - (Optional) The slug for the priority to be used for incidents of this type.
* `private_incident` - (Optional) A boolean to indicate if incidents of this type should be made private.

* `tags` - (Optional) A list of tags to be applied to incidents of this type.
* `runbook_ids` - (Optional) A list of runbook ids for the runbooks to be attached to incidents of this type.
* `team_ids` - (Optional) A list of team ids for the teams to be added to incidents of this type.
* `impacts` - (Optional) A block indicating which service catalog items would be impacted by incidents of this type and the condition they should be set to.

The `impacts` block supports: 

* `impact_id` - (Required) The id of the service, functionality, or environment that would be impacted by incidents of this type.
* `condition_id` - (Required) The id of the pre-defined condition to set for the above impact when creating incidents of this type.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the incident type.

## Import

Incident types can be imported; use `<INCIDENT TYPE ID>` as the import ID. For example:

```shell
terraform import firehydrant_incident_type.test 3638b647-b99c-5051-b715-eda2c912c42e
```