package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLifecycleMilestoneResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckLifecycleMilestoneResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLifecycleMilestoneResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLifecycleMilestoneResourceExistsWithAttributes_basic("firehydrant_lifecycle_milestone.new_milestone"),
					resource.TestCheckResourceAttrSet("firehydrant_lifecycle_milestone.new_milestone", "id"),
					resource.TestCheckResourceAttrSet("firehydrant_lifecycle_milestone.new_milestone", "phase_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_lifecycle_milestone.new_milestone", "name", fmt.Sprintf("Test Milestone %s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_lifecycle_milestone.new_milestone", "description", fmt.Sprintf("test description %s", rName)),
				),
			},
		},
	})
}

func TestAccLifecycleMilestoneResource_update(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckLifecycleMilestoneResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLifecycleMilestoneResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLifecycleMilestoneResourceExistsWithAttributes_basic("firehydrant_lifecycle_milestone.new_milestone"),
					resource.TestCheckResourceAttrSet("firehydrant_lifecycle_milestone.new_milestone", "id"),
					resource.TestCheckResourceAttrSet("firehydrant_lifecycle_milestone.new_milestone", "phase_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_lifecycle_milestone.new_milestone", "name", fmt.Sprintf("Test Milestone %s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_lifecycle_milestone.new_milestone", "description", fmt.Sprintf("test description %s", rName)),
				),
			},
			{
				Config: testAccLifecycleMilestoneResourceConfig_update(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLifecycleMilestoneResourceExistsWithAttributes_basic("firehydrant_lifecycle_milestone.new_milestone"),
					resource.TestCheckResourceAttrSet("firehydrant_lifecycle_milestone.new_milestone", "id"),
					resource.TestCheckResourceAttrSet("firehydrant_lifecycle_milestone.new_milestone", "phase_id"),
					resource.TestCheckResourceAttr(
						"firehydrant_lifecycle_milestone.new_milestone", "name", fmt.Sprintf("Test Milestone %s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_lifecycle_milestone.new_milestone", "description", fmt.Sprintf("test description %s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_lifecycle_milestone.new_milestone", "slug", "test-milestone"),
					resource.TestCheckResourceAttr(
						"firehydrant_lifecycle_milestone.new_milestone", "position", "2"),
					resource.TestCheckResourceAttr(
						"firehydrant_lifecycle_milestone.new_milestone", "auto_assign_timestamp_on_create", "never_set_on_create"),
				),
			},
		},
	})
}

func TestAccLifecycleMilestoneResourceImport_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLifecycleMilestoneResourceConfig_basic(rName),
			},
			{
				ResourceName:      "firehydrant_lifecycle_milestone.new_milestone",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLifecycleMilestoneResourceExistsWithAttributes_basic(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		milestoneResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if milestoneResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		response, err := client.Sdk.IncidentSettings.ListLifecyclePhases(context.TODO())
		if err != nil {
			return fmt.Errorf("Unable to get Lifecycle milestones: %v", err)
		}

		var desired_milestone *components.LifecyclesMilestoneEntity
		for _, phase := range response.Data {
			for _, milestone := range phase.Milestones {
				if *milestone.ID == milestoneResource.Primary.ID {
					desired_milestone = &milestone
				}
			}
		}
		if desired_milestone == nil {
			return fmt.Errorf("Lifecycle milestone %s still exists", milestoneResource.Primary.ID)
		}

		expected, got := milestoneResource.Primary.Attributes["name"], *desired_milestone.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = milestoneResource.Primary.Attributes["description"], *desired_milestone.Description
		if expected != got {
			return fmt.Errorf("Unexpected summary. Expected: %s, got: %s", expected, got)
		}

		return nil
	}
}

func testAccCheckLifecycleMilestoneResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_lifecycle_milestone" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			response, err := client.Sdk.IncidentSettings.ListLifecyclePhases(context.TODO())
			if err != nil {
				return fmt.Errorf("Unable to get Lifecycle milestones: %v", err)
			}

			var desired_milestone *components.LifecyclesMilestoneEntity
			for _, phase := range response.Data {
				for _, milestone := range phase.Milestones {
					if *milestone.ID == stateResource.Primary.ID {
						desired_milestone = &milestone
					}
				}
			}
			if desired_milestone != nil {
				return fmt.Errorf("Lifecycle milestone %s still exists", stateResource.Primary.ID)
			}
		}

		return nil
	}
}

func testAccLifecycleMilestoneResourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "firehydrant_lifecycle_phase" "started" {
  name = "started"
}

resource "firehydrant_lifecycle_milestone" "new_milestone" {
  name        = "Test Milestone %s"
  description = "test description %s"
	phase_id    = data.firehydrant_lifecycle_phase.started.id
}`, rName, rName)
}

func testAccLifecycleMilestoneResourceConfig_update(rName string) string {
	return fmt.Sprintf(`
data "firehydrant_lifecycle_phase" "started" {
  name = "started"
}

resource "firehydrant_lifecycle_milestone" "new_milestone" {
  name        = "Test Milestone %s"
  description = "test description %s"
	phase_id    = data.firehydrant_lifecycle_phase.started.id
	slug        = "test-milestone"
	position    = 2
	auto_assign_timestamp_on_create = "never_set_on_create"
}`, rName, rName)
}
