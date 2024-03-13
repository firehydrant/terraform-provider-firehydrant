## 0.6.0

ENHANCEMENTS:

* **New Data Source**: `firehydrant_slack_channel` ([#147](https://github.com/firehydrant/terraform-provider-firehydrant/pull/147)).
* Improved "resource not found error" with more details of URL.

## 0.5.0

* `firehydrant_team` now supports `slug` attribute ([#145](https://github.com/firehydrant/terraform-provider-firehydrant/pull/145)).

## 0.4.2

ENHANCEMENTS:

* Bump Go version to 1.22 ([#144](https://github.com/firehydrant/terraform-provider-firehydrant/pull/144))

BUG FIXES:

* Resource `firehydrant_on_call_schedule` now handles `member_ids` correctly. Field `members` has been deprecated. ([#143](https://github.com/firehydrant/terraform-provider-firehydrant/pull/143))

## 0.4.1

BUG FIXES:

* Fixes a bug where handoff steps in escalation policies were assigning to the incorrect schema and causing a panic.

## 0.4.0

* Add labels to functionalities
* Add initial Signals resources for beta period

## 0.3.6

* Fix client versioning

## 0.3.5

* Use timestamps for data source IDs
* Better document how to have runbooks automatically attach

## 0.3.4

* New Data Source: `firehydrant_team` ([#121](https://github.com/firehydrant/terraform-provider-firehydrant/pull/121))
* New Data Source: `firehydrant_teams` ([#121](https://github.com/firehydrant/terraform-provider-firehydrant/pull/121))

## 0.3.3

* Bump golang from 1.16 to 1.18
* resource/service: Added `external_resources` attribute to service ([#123](https://github.com/firehydrant/terraform-provider-firehydrant/pull/123))

## 0.3.2

ENHANCEMENTS:

* resource/service: Added `auto_add_responding_team` attribute to service ([#117](https://github.com/firehydrant/terraform-provider-firehydrant/pull/117))
* data_source/service: Added `auto_add_responding_team` attribute to service ([#117](https://github.com/firehydrant/terraform-provider-firehydrant/pull/117))
* data_source/services: Added `auto_add_responding_team` attribute to service ([#117](https://github.com/firehydrant/terraform-provider-firehydrant/pull/117))
* resource/team: Add the ability to attach memberships to teams ([#116](https://github.com/firehydrant/terraform-provider-firehydrant/pull/116))

## 0.3.1

BUG FIXES:

* documentation: Fixed broken links in runbook resource documentation ([#101](https://github.com/firehydrant/terraform-provider-firehydrant/pull/101))

## 0.3.0

BREAKING CHANGES:

* resource/functionality: Removed deprecated `services` attribute, preferring `service_ids` instead ([#94](https://github.com/firehydrant/terraform-provider-firehydrant/pull/94))
* resource/runbook: Changed the type of the `steps` `config` attribute to a JSON string ([#79](https://github.com/firehydrant/terraform-provider-firehydrant/pull/79))
* resource/runbook: Removed deprecated `type` attribute ([#80](https://github.com/firehydrant/terraform-provider-firehydrant/pull/80))
* resource/runbook: Removed deprecated `severities` attribute ([#80](https://github.com/firehydrant/terraform-provider-firehydrant/pull/80))
* resource/runbook: Changed `steps` attribute to be required ([#80](https://github.com/firehydrant/terraform-provider-firehydrant/pull/80))

BUG FIXES:

* resource/environment: Fixed bug that prevented the `description` attribute from being unset ([#93](https://github.com/firehydrant/terraform-provider-firehydrant/pull/93))
* resource/runbook: Fixed bug that prevented the `description` attribute from being unset ([#80](https://github.com/firehydrant/terraform-provider-firehydrant/pull/80))
* resource/severity: Fixed bug that prevented the `description` attribute from being unset ([#88](https://github.com/firehydrant/terraform-provider-firehydrant/pull/88))
* data_source/runbook_action: Fixed bug that prevented the Slack `add_bookmark_to_incident_channel` action from working by making `type` optional ([#92](https://github.com/firehydrant/terraform-provider-firehydrant/pull/92))

FEATURES:

* **New Resource:** `firehydrant_incident_role` ([#87](https://github.com/firehydrant/terraform-provider-firehydrant/pull/87))
* **New Resource:** `firehydrant_priority` ([#65](https://github.com/firehydrant/terraform-provider-firehydrant/pull/65))
* **New Resource:** `firehydrant_service_dependency` ([#89](https://github.com/firehydrant/terraform-provider-firehydrant/pull/89))
* **New Resource:** `firehydrant_task_list` ([#85](https://github.com/firehydrant/terraform-provider-firehydrant/pull/85))
* **New Data Source:** `firehydrant_incident_role` ([#87](https://github.com/firehydrant/terraform-provider-firehydrant/pull/87))
* **New Data Source:** `firehydrant_priority` ([#65](https://github.com/firehydrant/terraform-provider-firehydrant/pull/65))
* **New Data Source:** `firehydrant_severity` ([#88](https://github.com/firehydrant/terraform-provider-firehydrant/pull/88))
* **New Data Source:** `firehydrant_task_list` ([#85](https://github.com/firehydrant/terraform-provider-firehydrant/pull/85))

ENHANCEMENTS:

* provider: Improved error messages by adding details from the API response ([#75](https://github.com/firehydrant/terraform-provider-firehydrant/pull/75))
* provider: Improved documentation for configuring the provider ([#95](https://github.com/firehydrant/terraform-provider-firehydrant/pull/95))
* resource/environment: Added logging ([#93](https://github.com/firehydrant/terraform-provider-firehydrant/pull/93))
* resource/functionality: Added logging ([#94](https://github.com/firehydrant/terraform-provider-firehydrant/pull/94))
* resource/runbook: Added logging ([#74](https://github.com/firehydrant/terraform-provider-firehydrant/pull/74))
* resource/runbook: Added the `owner_id` attribute to runbook ([#76](https://github.com/firehydrant/terraform-provider-firehydrant/pull/76))
* resource/runbook: Added the `repeats` and `repeat_duration` attribute to runbook step ([#78](https://github.com/firehydrant/terraform-provider-firehydrant/pull/78))
* resource/runbook: Added the `attachment_rule` attribute to runbook ([#82](https://github.com/firehydrant/terraform-provider-firehydrant/pull/82))
* resource/runbook: Added default value of `false` to the steps `automatic` attribute ([#83](https://github.com/firehydrant/terraform-provider-firehydrant/pull/83))
* resource/runbook: Added the `rule` attribute to runbook steps ([#84](https://github.com/firehydrant/terraform-provider-firehydrant/pull/84))
* resource/runbook: Added documentation for step configuration for every runbook action ([#81](https://github.com/firehydrant/terraform-provider-firehydrant/pull/81))
* resource/service: Added logging ([#96](https://github.com/firehydrant/terraform-provider-firehydrant/pull/96))
* resource/severity: Added logging to the resource and validation to the `slug` attribute ([#88](https://github.com/firehydrant/terraform-provider-firehydrant/pull/88))
* resource/severity: Added support for the `type` attribute ([#88](https://github.com/firehydrant/terraform-provider-firehydrant/pull/88))
* resource/team: Added logging ([#96](https://github.com/firehydrant/terraform-provider-firehydrant/pull/96))
* data_source/environment: Added logging ([#93](https://github.com/firehydrant/terraform-provider-firehydrant/pull/93))
* data_source/functionality: Added logging ([#94](https://github.com/firehydrant/terraform-provider-firehydrant/pull/94))
* data_source/runbook: Added logging ([#74](https://github.com/firehydrant/terraform-provider-firehydrant/pull/74))
* data_source/runbook: Added the `owner_id` attribute to runbook ([#76](https://github.com/firehydrant/terraform-provider-firehydrant/pull/76))
* data_source/runbook: Added the `attachment_rule` attribute to runbook ([#82](https://github.com/firehydrant/terraform-provider-firehydrant/pull/82))
* data_source/runbook_action: Added logging ([#74](https://github.com/firehydrant/terraform-provider-firehydrant/pull/74))
* data_source/service: Added logging ([#96](https://github.com/firehydrant/terraform-provider-firehydrant/pull/96))
* data_source/services: Added logging ([#96](https://github.com/firehydrant/terraform-provider-firehydrant/pull/96))

NOTES:

* resource/functionality: The deprecated `services` attribute has been removed. See the ["Notes" section in 0.2.0](#020)
  or the [original deprecation PR](https://github.com/firehydrant/terraform-provider-firehydrant/pull/49) for more information.
* resource/runbook: There are a number of breaking changes for the runbook resource. The `steps` attribute is now required,
  the `steps` `config` attribute is now a JSON string, and the `type` and `severities` attribute have been removed. In order
  to upgrade to 0.3.0, you will need to destroy your existing runbooks and recreate them after changing your configuration to
  account for the breaking changes.
  As an example, the configuration below was valid in 0.2.1
   ```hcl
   # An example of a valid 0.2.1 configuration
   resource "firehydrant_runbook" "example-runbook" {
     name = "example-runbook"
     type = "incident"

     steps {
       name    = "Send me an email"
       action_id = data.firehydrant_runbook_action.firehydrant_email_notification.id

       config = {
         email_address   = "test@example.com"
         email_subject   = "Incident opened on FireHydrant"
         default_message = "Message"
       }
     }
   }
   ```
  To upgrade to 0.3.0, that configuration would have to change to the following:
   ```hcl
   # The same configuration as above, updated to be valid for 0.3.0
   resource "firehydrant_runbook" "example-runbook" {
     name = "example-runbook"

     steps {
       name    = "Send me an email"
       action_id = data.firehydrant_runbook_action.firehydrant_email_notification.id

       config = jsonencode({
         email_address   = "test@example.com"
         email_subject   = "Incident opened on FireHydrant"
         default_message = "Message"
       })
     }
   }
   ```

## 0.2.1

ENHANCEMENTS:

* documentation: Refactored documentation to follow best practices, added descriptions for all arguments, and added configuration examples for all resources and data sources ([#69](https://github.com/firehydrant/terraform-provider-firehydrant/pull/69))

## 0.2.0

BREAKING CHANGES:

* resource/team: Removed `services` attribute. Use resource/service to associate teams with a service. See "Notes" for more information ([#54](https://github.com/firehydrant/terraform-provider-firehydrant/pull/54))

BUG FIXES:

* provider: Fixed bug where errors weren't being checked or handled for various API requests ([#58](https://github.com/firehydrant/terraform-provider-firehydrant/pull/58))
* resource/functionality: Fixed bug that prevented the `description` attribute from being unset ([#49](https://github.com/firehydrant/terraform-provider-firehydrant/pull/49))
* resource/functionality: Fixed bug that prevented the `services` attribute from being unset ([#49](https://github.com/firehydrant/terraform-provider-firehydrant/pull/49))
* resource/functionality: Fixed bug that prevented functionalities from being removed from state when they had been deleted outside of Terraform ([#58](https://github.com/firehydrant/terraform-provider-firehydrant/pull/58))
* resource/service: Fixed bug that prevented the `description` attribute from being unset ([#51](https://github.com/firehydrant/terraform-provider-firehydrant/pull/51))
* resource/service: Fixed bug that prevented `labels` from being removed from services ([#52](https://github.com/firehydrant/terraform-provider-firehydrant/pull/52))
* resource/service: Fixed bug that prevented services from being removed from state when they had been deleted outside of Terraform ([#58](https://github.com/firehydrant/terraform-provider-firehydrant/pull/58))
* resource/team: Fixed bug that prevented teams from being removed from state when they had been deleted outside of Terraform ([#58](https://github.com/firehydrant/terraform-provider-firehydrant/pull/58))
* data_source/runbook_action: Fixed bug that caused the wrong action to be returned when multiple actions existed for the same slug ([#56](https://github.com/firehydrant/terraform-provider-firehydrant/pull/56))

ENHANCEMENTS:

* provider: Added Terraform version to the user agent header ([#24](https://github.com/firehydrant/terraform-provider-firehydrant/pull/24))
* resource/functionality: Added deprecation warning to the `services` attribute, preferring `service_ids` instead ([#49](https://github.com/firehydrant/terraform-provider-firehydrant/pull/49))
* resource/service: Added the `alert_on_add` attribute to services ([#24](https://github.com/firehydrant/terraform-provider-firehydrant/pull/24))
* resource/service: Added the `owner_id` attribute to services ([#23](https://github.com/firehydrant/terraform-provider-firehydrant/pull/23))
* resource/service: Added the `team_ids` attribute to services ([#54](https://github.com/firehydrant/terraform-provider-firehydrant/pull/54))
* resource/service: Added the `links` attribute to services ([#30](https://github.com/firehydrant/terraform-provider-firehydrant/pull/30))
* data_source/functionality: Added the `service_ids` attribute to functionality ([#49](https://github.com/firehydrant/terraform-provider-firehydrant/pull/49))
* data_source/service: Added the `alert_on_add` attribute to services ([#24](https://github.com/firehydrant/terraform-provider-firehydrant/pull/24))
* data_source/service: Added the `owner_id` attribute to services ([#23](https://github.com/firehydrant/terraform-provider-firehydrant/pull/23))
* data_source/service: Added the `team_ids` attribute to services ([#54](https://github.com/firehydrant/terraform-provider-firehydrant/pull/54))
* data_source/service: Added the `links` and `labels` attributes to services ([#30](https://github.com/firehydrant/terraform-provider-firehydrant/pull/30))
* data_source/services: Added the `alert_on_add` attribute to services ([#24](https://github.com/firehydrant/terraform-provider-firehydrant/pull/24))
* data_source/services: Added the `owner_id` attribute to services ([#23](https://github.com/firehydrant/terraform-provider-firehydrant/pull/23))
* data_source/services: Added the `team_ids` attribute to services ([#54](https://github.com/firehydrant/terraform-provider-firehydrant/pull/54))
* data_source/services: Added the `links` and `labels` attributes to services ([#30](https://github.com/firehydrant/terraform-provider-firehydrant/pull/30))

NOTES:

* The services attribute has been removed from resource/team. If you need to add a team as an owner or responder to a service, use resource/service and specify `owner_id` or `team_ids`.
   When upgrading to 0.2.0, you should remove the `services` attribute from any team resources and instead add `team_ids` to each service resource.
   As an example, the configuration below was valid in 0.1.4.
   ```hcl
   # An example of a valid 0.1.4 configuration
   resource "firehydrant_service" "service1" {
     name        = "service1"
     description = "description1"

     labels = {
       language  = "ruby",
       lifecycle = "production"
     }

     service_tier = 1
   }

   resource "firehydrant_team" "team1" {
     name = "team1"

     services {
       id = firehydrant_service.service1.id
     }
   }

   resource "firehydrant_team" "team2" {
     name        = "team2"
     description = "description2"

     services {
       id = firehydrant_service.service1.id
     }
   }
   ```
   To upgrade to 0.2.0, that configuration would have to change to the following:
   ```hcl
   # The same configuration as above, updated to be valid for 0.2.0
   resource "firehydrant_service" "service1" {
     name        = "service1"
     description = "description1"

     labels = {
       language  = "ruby",
       lifecycle = "production"
     }

     service_tier = 1

     team_ids = [
       firehydrant_team.team1.id,
       firehydrant_team.team2.id
     ]
   }


   resource "firehydrant_team" "team1" {
     name = "team1"
   }

   resource "firehydrant_team" "team2" {
     name        = "team2"
     description = "description2"
   }
   ```
* The deprecated attribute `services` will be removed from resource/functionality 3 months after the release of v0.2.0. You will have until June 30, 2022 to migrate to the preferred attribute.
  More information about this deprecation can be found in the description of ([#49](https://github.com/firehydrant/terraform-provider-firehydrant/pull/49))

## 0.1.4

BUG FIXES:
* Only one page of results was returned when using`data "firehydrant_runbook_action"` to lookup runbook steps. This meant some steps were completely inaccessible through the provider. This patch will return all steps.

FEATURES:
* CHANGELOG.md added to track contents of each release

IMPROVEMENTS
* The `examples` directory now is separated by resource and includes more and extended examples of `service` resources

## 0.1.3
FEATURES:
* Automated release process using GitHub Actions
* Added `service_tier` to service resource
## 0.1.2 (Dec 3, 2020)

NO CHANGES
## 0.1.1 (Dec 3, 2020)

IMPROVEMENTS
* Added autogenerated documentation

## 0.1.0 (Dec 3, 2020)

INITIAL RELEASE
