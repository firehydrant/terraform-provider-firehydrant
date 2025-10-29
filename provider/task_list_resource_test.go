package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTaskListResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckTaskListResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaskListResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaskListResourceExistsWithAttributes_basic("firehydrant_task_list.test_task_list"),
					resource.TestCheckResourceAttrSet("firehydrant_task_list.test_task_list", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "name", fmt.Sprintf("test-task-list-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "task_list_items.#", "1"),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "task_list_items.0.summary", fmt.Sprintf("test-summary1-%s", rName)),
				),
			},
		},
	})
}

func TestAccTaskListResource_update(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckTaskListResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaskListResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaskListResourceExistsWithAttributes_basic("firehydrant_task_list.test_task_list"),
					resource.TestCheckResourceAttrSet("firehydrant_task_list.test_task_list", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "name", fmt.Sprintf("test-task-list-%s", rName)),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "task_list_items.#", "1"),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "task_list_items.0.summary", fmt.Sprintf("test-summary1-%s", rName)),
				),
			},
			{
				Config: testAccTaskListResourceConfig_update(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaskListResourceExistsWithAttributes_update("firehydrant_task_list.test_task_list"),
					resource.TestCheckResourceAttrSet("firehydrant_task_list.test_task_list", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "name", fmt.Sprintf("test-task-list-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "description", fmt.Sprintf("test-description-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "task_list_items.#", "2"),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "task_list_items.0.summary", fmt.Sprintf("test-summary1-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "task_list_items.0.description", fmt.Sprintf("test-description1-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "task_list_items.1.summary", fmt.Sprintf("test-summary2-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "task_list_items.1.description", fmt.Sprintf("test-description2-%s", rNameUpdated)),
				),
			},
			{
				Config: testAccTaskListResourceConfig_basic(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaskListResourceExistsWithAttributes_basic("firehydrant_task_list.test_task_list"),
					resource.TestCheckResourceAttrSet("firehydrant_task_list.test_task_list", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "name", fmt.Sprintf("test-task-list-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "task_list_items.#", "1"),
					resource.TestCheckResourceAttr(
						"firehydrant_task_list.test_task_list", "task_list_items.0.summary", fmt.Sprintf("test-summary1-%s", rNameUpdated)),
				),
			},
		},
	})
}

func TestAccTaskListResourceImport_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaskListResourceConfig_basic(rName),
			},
			{
				ResourceName:      "firehydrant_task_list.test_task_list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTaskListResourceImport_allAttributes(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccTaskListResourceConfig_update(rName),
			},
			{
				ResourceName:      "firehydrant_task_list.test_task_list",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTaskListResourceExistsWithAttributes_basic(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		taskListResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if taskListResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		taskListResponse, err := client.TaskLists().Get(context.TODO(), taskListResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := taskListResource.Primary.Attributes["name"], taskListResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		if taskListResponse.Description != "" {
			return fmt.Errorf("Unexpected description. Expected no description, got: %s", taskListResponse.Description)
		}

		if len(taskListResponse.TaskListItems) != 1 {
			return fmt.Errorf("Unexpected number of task list items. Expected 1 item, got: %v", len(taskListResponse.TaskListItems))
		}

		for index, taskListItem := range taskListResponse.TaskListItems {
			key := fmt.Sprintf("task_list_items.%d", index)
			if taskListResource.Primary.Attributes[key+".summary"] != taskListItem.Summary {
				return fmt.Errorf("Unexpected task list item summary. Expected %s, got: %s", taskListItem.Summary, taskListResource.Primary.Attributes[key+".summary"])
			}

			if taskListResource.Primary.Attributes[key+".description"] != taskListItem.Description {
				return fmt.Errorf("Unexpected task list item description. Expected %s, got: %s", taskListItem.Description, taskListResource.Primary.Attributes[key+".description"])
			}
		}

		return nil
	}
}

func testAccCheckTaskListResourceExistsWithAttributes_update(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		taskListResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if taskListResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		taskListResponse, err := client.TaskLists().Get(context.TODO(), taskListResource.Primary.ID)
		if err != nil {
			return err
		}

		expected, got := taskListResource.Primary.Attributes["name"], taskListResponse.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = taskListResource.Primary.Attributes["description"], taskListResponse.Description
		if expected != got {
			return fmt.Errorf("Unexpected description. Expected: %s, got: %s", expected, got)
		}

		if len(taskListResponse.TaskListItems) != 2 {
			return fmt.Errorf("Unexpected number of task list items. Expected 2 items, got: %v", len(taskListResponse.TaskListItems))
		}

		for index, taskListItem := range taskListResponse.TaskListItems {
			key := fmt.Sprintf("task_list_items.%d", index)
			if taskListResource.Primary.Attributes[key+".summary"] != taskListItem.Summary {
				return fmt.Errorf("Unexpected task list item summary. Expected %s, got: %s", taskListItem.Summary, taskListResource.Primary.Attributes[key+".summary"])
			}

			if taskListResource.Primary.Attributes[key+".description"] != taskListItem.Description {
				return fmt.Errorf("Unexpected task list item description. Expected %s, got: %s", taskListItem.Description, taskListResource.Primary.Attributes[key+".description"])
			}
		}

		return nil
	}
}

func testAccCheckTaskListResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_task_list" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			_, err := client.TaskLists().Get(context.TODO(), stateResource.Primary.ID)
			if err == nil {
				return fmt.Errorf("Task list %s still exists", stateResource.Primary.ID)
			}
		}

		return nil
	}
}

func testAccTaskListResourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_task_list" "test_task_list" {
  name = "test-task-list-%s"

  task_list_items {
    summary = "test-summary1-%s"
  }
}`, rName, rName)
}

func testAccTaskListResourceConfig_update(rName string) string {
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
}`, rName, rName, rName, rName, rName, rName)
}
