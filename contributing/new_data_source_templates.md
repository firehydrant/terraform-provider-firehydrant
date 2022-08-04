# Templates for adding new data sources

The following are basic example templates for creating new data sources and the tests
that go with them.

For more complex examples, see the task list or service data sources in this provider. 

## Data source template

```go
// example_model_data_source.go
package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceExampleModel() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantExampleModel,
		Schema: map[string]*schema.Schema{
			// Required
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Optional
			"optional_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed
			"computed_attribute": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"other_computed_attribute": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataFireHydrantExampleModel(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the example model
	exampleModelID := d.Get("id").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read example model: %s", exampleModelID), map[string]interface{}{
		"id": exampleModelID,
	})
	exampleModelResponse, err := firehydrantAPIClient.ExampleModels().Get(ctx, exampleModelID)
	if err != nil {
		return diag.Errorf("Error reading example model %s: %v", exampleModelID, err)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"computed_attribute":       exampleModelResponse.ComputedAttribute,
		"other_computed_attribute": exampleModelResponse.OtherComputedAttribute,
	}

	// Process any data that could be nil or has a more complex structure
	// and add it to the attributes block
	// Good examples of this can be found in the task list and service
	// data sources.

	// Set the resource attributes to the values we got from the API
	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for example model %s: %v", key, exampleModelID, err)
		}
	}

	// Set the example model's ID in state
	d.SetId(exampleModelResponse.ID)

	return diag.Diagnostics{}
}
```

## Data source test template

```go
// example_model_data_source_test.go
package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// This tests the data source against a configuration that only has the required
// attributes specified.
func TestAccExampleModelDataSource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccExampleModelDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_example_model.test_example_model", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_example_model.test_example_model", "computed_attribute", fmt.Sprintf("test-example-model-%s", rName)),
				),
			},
		},
	})
}

// This tests the data source against a configuration that has all possible
// attributes specified.
func TestAccExampleModelDataSource_allAttributes(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: defaultProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccExampleModelDataSourceConfig_allAttributes(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_example_model.test_example_model", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_example_model.test_example_model", "computed_attribute", fmt.Sprintf("test-example-model-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_example_model.test_example_model", "other_computed_attribute", fmt.Sprintf("test-other-computed-attribute-%s", rName)),
				),
			},
		},
	})
}

// Remember to run all example Terraform configs through the formatter
// using `terraform fmt`
func testAccExampleModelDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_example_model" "test_example_model" {
  computed_attribute = "test-example-model-%s"
}

data "firehydrant_example_model" "test_example_model" {
  id = firehydrant_example_model.test_example_model.id
}`, rName)
}

func testAccExampleModelDataSourceConfig_allAttributes(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_example_model" "test_example_model" {
  computed_attribute       = "test-example-model-%s"
  other_computed_attribute = "test-other-computed-attribute-%s"
}

data "firehydrant_example_model" "test_example_model" {
  id = firehydrant_example_model.test_example_model.id
}`, rName, rName)
}
```