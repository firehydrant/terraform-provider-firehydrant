# Templates for adding new resources

The following are basic example templates for creating new resources and the tests
that go with them.

For more complex examples, see the task list or service resources in this provider.

## Resource template

```go
// example_model_resource.go
package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceExampleModel() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantExampleModel,
		UpdateContext: updateResourceFireHydrantExampleModel,
		ReadContext:   readResourceFireHydrantExampleModel,
		DeleteContext: deleteResourceFireHydrantExampleModel,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"required_attribute": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Optional
			"optional_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func readResourceFireHydrantExampleModel(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the example model
	exampleModelID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read example model: %s", exampleModelID), map[string]interface{}{
		"id": exampleModelID,
	})
	exampleModelResponse, err := firehydrantAPIClient.ExampleModels().Get(ctx, exampleModelID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			tflog.Debug(ctx, fmt.Sprintf("Example model %s no longer exists", exampleModelID), map[string]interface{}{
				"id": exampleModelID,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading example model %s: %v", exampleModelID, err)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"required_attribute": exampleModelResponse.RequiredAttribute,
		"optional_attribute": exampleModelResponse.OptionalAttribute,
	}

	// Process any data that could be nil or has a more complex structure
	// and add it to the attributes block.
	// Good examples of this can be found in the task list and service
	// resources.

	// Set the resource attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for example model %s: %v", key, exampleModelID, err)
		}
	}

	return diag.Diagnostics{}
}

func createResourceFireHydrantExampleModel(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get attributes from config and construct the create request
	createRequest := firehydrant.CreateExampleModelRequest{
		RequiredAttribute: d.Get("required_attribute").(string),
		OptionalAttribute: d.Get("optional_attribute").(string),
	}

	// Process any optional attributes and add to the create request if necessary
	// This includes things like complex optional attributes that could be nil
	// attributes that have more complex structures.
	// Good examples of this can be found in the task list and service
	// resources.

	// Create the new example model
	// The required attribute you use in the log here should be something identifying, like a name/slug/etc
	tflog.Debug(ctx, fmt.Sprintf("Create example model: %s", createRequest.RequiredAttribute), map[string]interface{}{
		"required_attribute": createRequest.RequiredAttribute,
	})
	exampleModelResponse, err := firehydrantAPIClient.ExampleModels().Create(ctx, createRequest)
	if err != nil {
		return diag.Errorf("Error creating example model %s: %v", createRequest.RequiredAttribute, err)
	}

	// Set the new example model's ID in state
	d.SetId(exampleModelResponse.ID)

	// Update state with the latest information from the API
	return readResourceFireHydrantExampleModel(ctx, d, m)
}

func updateResourceFireHydrantExampleModel(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Construct the update request
	updateRequest := firehydrant.UpdateExampleModelRequest{
		RequiredAttribute: d.Get("required_attribute").(string),
		OptionalAttribute: d.Get("optional_attribute").(string),
	}

	// Process any optional attributes and add to the update request if necessary
	// This includes things like complex optional attributes that could be nil
	// attributes that have more complex structures.
	// Good examples of this can be found in the task list and service
	// resources.

	// Update the example model
	tflog.Debug(ctx, fmt.Sprintf("Update example model: %s", d.Id()), map[string]interface{}{
		"id": d.Id(),
	})
	_, err := firehydrantAPIClient.ExampleModels().Update(ctx, d.Id(), updateRequest)
	if err != nil {
		return diag.Errorf("Error updating example model %s: %v", d.Id(), err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantExampleModel(ctx, d, m)
}

func deleteResourceFireHydrantExampleModel(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Delete the example model
	exampleModelID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Delete example model: %s", exampleModelID), map[string]interface{}{
		"id": exampleModelID,
	})
	err := firehydrantAPIClient.ExampleModels().Delete(ctx, exampleModelID)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			return nil
		}
		return diag.Errorf("Error deleting example model %s: %v", exampleModelID, err)
	}

	return diag.Diagnostics{}
}
```

## Resource test template

```go
// example_model_resource_test.go
package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// This tests the resource with a configuration that only has the required
// attributes specified.
func TestAccExampleModelResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckExampleModelResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccExampleModelResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckExampleModelResourceExistsWithAttributes_basic("firehydrant_example_model.test_example_model"),
					resource.TestCheckResourceAttrSet("firehydrant_example_model.test_example_model", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_example_model.test_example_model", "required_attribute", fmt.Sprintf("test-example-model-%s", rName)),
				),
			},
		},
	})
}

// This tests the resources ability to update and remove attributes
// with a configuration that has all attributes specified.
func TestAccExampleModelResource_update(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		CheckDestroy:      testAccCheckExampleModelResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccExampleModelResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckExampleModelResourceExistsWithAttributes_basic("firehydrant_example_model.test_example_model"),
					resource.TestCheckResourceAttrSet("firehydrant_example_model.test_example_model", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_example_model.test_example_model", "required_attribute", fmt.Sprintf("test-example-model-%s", rName)),
				),
			},
			{
				Config: testAccExampleModelResourceConfig_update(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckExampleModelResourceExistsWithAttributes_update("firehydrant_example_model.test_example_model"),
					resource.TestCheckResourceAttrSet("firehydrant_example_model.test_example_model", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_example_model.test_example_model", "required_attribute", fmt.Sprintf("test-example-model-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_example_model.test_example_model", "optional_attribute", fmt.Sprintf("test-optional-attribute-%s", rNameUpdated)),
				),
			},
			{
				Config: testAccExampleModelResourceConfig_basic(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckExampleModelResourceExistsWithAttributes_basic("firehydrant_example_model.test_example_model"),
					resource.TestCheckResourceAttrSet("firehydrant_example_model.test_example_model", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_example_model.test_example_model", "required_attribute", fmt.Sprintf("test-example-model-%s", rNameUpdated)),
				),
			},
		},
	})
}

// This tests the resource's ability to import with a configuration that
// only has the required attributes specified.
func TestAccExampleModelResourceImport_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccExampleModelResourceConfig_basic(rName),
			},
			{
				ResourceName:      "firehydrant_example_model.test_example_model",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// This tests the resource's ability to import with a configuration that
// has all attributes specified.
func TestAccExampleModelResourceImport_allAttributes(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccExampleModelResourceConfig_update(rName),
			},
			{
				ResourceName:      "firehydrant_example_model.test_example_model",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// This test does a more in-depth check to test that what was in the
// configuration matches the information we get back from the API
// with a configuration that only has the required attributes specified.
// Whenever possible, try to test every attribute as deeply as possible.
func testAccCheckExampleModelResourceExistsWithAttributes_basic(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		exampleModelResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if exampleModelResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		exampleModelResponse, err := client.ExampleModels().Get(context.TODO(), exampleModelResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := exampleModelResource.Primary.Attributes["required_attribute"], exampleModelResponse.RequiredAttribute
		if expected != got {
			return fmt.Errorf("Unexpected required_attribute. Expected: %s, got: %s", expected, got)
		}

		if exampleModelResponse.OptionalAttribute != "" {
			return fmt.Errorf("Unexpected optional_attribute. Expected no optional_attribute, got: %s", exampleModelResponse.OptionalAttribute)
		}

		return nil
	}
}

// This test does a more in-depth check to test that what was in the
// configuration matches the information we get back from the API
// with a configuration that has all attributes specified.
// Whenever possible, try to test every attribute as deeply as possible.
func testAccCheckExampleModelResourceExistsWithAttributes_update(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		exampleModelResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if exampleModelResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		exampleModelResponse, err := client.ExampleModels().Get(context.TODO(), exampleModelResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := exampleModelResource.Primary.Attributes["required_attribute"], exampleModelResponse.RequiredAttribute
		if expected != got {
			return fmt.Errorf("Unexpected required_attribute. Expected: %s, got: %s", expected, got)
		}

		expected, got = exampleModelResource.Primary.Attributes["optional_attribute"], exampleModelResponse.OptionalAttribute
		if expected != got {
			return fmt.Errorf("Unexpected optional_attribute. Expected: %s, got: %s", expected, got)
		}

		return nil
	}
}

// This tests that the resource gets destroyed properly
func testAccCheckExampleModelResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_example_model" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			_, err := client.ExampleModels().Get(context.TODO(), stateResource.Primary.ID)
			if err == nil {
				return fmt.Errorf("Example model %s still exists", stateResource.Primary.ID)
			}
		}

		return nil
	}
}

// Remember to run all example Terraform configs through the formatter
// using `terraform fmt`
func testAccExampleModelResourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_example_model" "test_example_model" {
  required_attribute = "test-example-model-%s"
}`, rName)
}

func testAccExampleModelResourceConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_example_model" "test_example_model" {
  required_attribute = "test-example-model-%s"
  optional_attribute = "test-optional-attribute-%s"
}`, rName, rName)
}
```
