package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	fhsdk "github.com/firehydrant/firehydrant-go-sdk"
	"github.com/firehydrant/firehydrant-go-sdk/models/components"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCustomEventSourceResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckCustomEventSourceResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomEventSourceResourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCustomEventSourceResourceExistsWithAttributes_basic("firehydrant_custom_event_source.foo_transposer"),
					resource.TestCheckResourceAttrSet("firehydrant_custom_event_source.foo_transposer", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "name", "The Foo Transposer"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "slug", "foo"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "description", "This is the foo transposer"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "javascript", "function transpose(input) {\n  return input.data;\n}"),
				),
			},
		},
	})
}

func TestAccCustomEventSourceResource_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckCustomEventSourceResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomEventSourceResourceConfig_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCustomEventSourceResourceExistsWithAttributes_basic("firehydrant_custom_event_source.foo_transposer"),
					resource.TestCheckResourceAttrSet("firehydrant_custom_event_source.foo_transposer", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "name", "The Foo Transposer"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "slug", "foo"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "description", "This is the foo transposer"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "javascript", "function transpose(input) {\n  return input.data;\n}"),
				),
			},
			{
				Config: testAccCustomEventSourceResourceConfig_update(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCustomEventSourceResourceExistsWithAttributes_basic("firehydrant_custom_event_source.foo_transposer"),
					resource.TestCheckResourceAttrSet("firehydrant_custom_event_source.foo_transposer", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "name", "The Foo Transposer"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "slug", "foo"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "description", "A new foo transposer description"),
					resource.TestCheckResourceAttr(
						"firehydrant_custom_event_source.foo_transposer", "javascript", "function transpose(input) {\n  return input.foo;\n}"),
				),
			},
		},
	})
}

func testAccCheckCustomEventSourceResourceExistsWithAttributes_basic(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		customEventSourceResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if customEventSourceResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client := fhsdk.New(fhsdk.WithSecurity(components.Security{APIKey: os.Getenv("FIREHYDRANT_API_KEY")}))

		response, err := client.Signals.GetSignalsEventSource(context.TODO(), customEventSourceResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := customEventSourceResource.Primary.Attributes["name"], *response.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = customEventSourceResource.Primary.Attributes["slug"], *response.Slug
		if expected != got {
			return fmt.Errorf("Unexpected max_delay. Expected: %s, got: %s", expected, got)
		}

		expected, got = customEventSourceResource.Primary.Attributes["description"], *response.Description
		if expected != got {
			return fmt.Errorf("Unexpected priority. Expected: %s, got: %s", expected, got)
		}

		expected, got = customEventSourceResource.Primary.Attributes["javascript"], *response.Expression
		if expected != got {
			return fmt.Errorf("Unexpected priority. Expected: %s, got: %s", expected, got)
		}

		return nil
	}
}

func testAccCheckCustomEventSourceResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := fhsdk.New(fhsdk.WithSecurity(components.Security{APIKey: os.Getenv("FIREHYDRANT_API_KEY")}))

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_custom_event_source" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			_, err := client.Signals.GetSignalsEventSource(context.TODO(), stateResource.Primary.ID)
			if err == nil {
				return fmt.Errorf("Custom Event Source %s still exists", stateResource.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCustomEventSourceResourceConfig_basic() string {
	return `
resource "firehydrant_custom_event_source" "foo_transposer" {
  name = "The Foo Transposer"
	slug = "foo"
	description = "This is the foo transposer"
	javascript = "function transpose(input) {\n  return input.data;\n}"
}`
}

func testAccCustomEventSourceResourceConfig_update() string {
	return `
resource "firehydrant_custom_event_source" "foo_transposer" {
  name = "The Foo Transposer"
	slug = "foo"
	description = "A new foo transposer description"
	javascript = "function transpose(input) {\n  return input.foo;\n}"
}`
}
