package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTaskListDataSource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaskListDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_task_list.test_task_list", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_task_list.test_task_list", "name", fmt.Sprintf("test-task-list-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_task_list.test_task_list", "task_list_items.#", "1"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_task_list.test_task_list", "task_list_items.0.summary", fmt.Sprintf("test-summary1-%s", rName)),
				),
			},
		},
	})
}

func TestAccTaskListDataSource_allAttributes(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaskListDataSourceConfig_allAttributes(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.firehydrant_task_list.test_task_list", "id"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_task_list.test_task_list", "name", fmt.Sprintf("test-task-list-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_task_list.test_task_list", "description", fmt.Sprintf("test-description-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_task_list.test_task_list", "task_list_items.#", "2"),
					resource.TestCheckResourceAttr(
						"data.firehydrant_task_list.test_task_list", "task_list_items.0.summary", fmt.Sprintf("test-summary1-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_task_list.test_task_list", "task_list_items.0.description", fmt.Sprintf("test-description1-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_task_list.test_task_list", "task_list_items.1.summary", fmt.Sprintf("test-summary2-%s", rName)),
					resource.TestCheckResourceAttr(
						"data.firehydrant_task_list.test_task_list", "task_list_items.1.description", fmt.Sprintf("test-description2-%s", rName)),
				),
			},
		},
	})
}

func testAccTaskListDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_task_list" "test_task_list" {
  name = "test-task-list-%s"

  task_list_items {
    summary = "test-summary1-%s"
  }
}

data "firehydrant_task_list" "test_task_list" {
  id = firehydrant_task_list.test_task_list.id
}`, rName, rName)
}

func testAccTaskListDataSourceConfig_allAttributes(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_task_list" "test_task_list" {
  name        = "test-task-list-%s"
  description = "test-description-%s"

  task_list_items {
    summary     = "test-summary1-%s"
    description = "test-description1-%s"
  }

  task_list_items {
    summary     = "test-summary2-%s"
    description = "test-description2-%s"
  }
}

data "firehydrant_task_list" "test_task_list" {
  id = firehydrant_task_list.test_task_list.id
}`, rName, rName, rName, rName, rName, rName)
}
