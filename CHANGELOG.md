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
