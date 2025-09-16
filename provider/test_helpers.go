package provider

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testAccCheckRoleResourceDestroy verifies that role resources are properly cleaned up
func testAccCheckRoleResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := firehydrant.NewRestClient(os.Getenv("FIREHYDRANT_API_KEY"))
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_role" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("no instance ID is set")
			}

			_, err := client.Roles().Get(context.Background(), stateResource.Primary.ID)
			if err == nil {
				return fmt.Errorf("role %s still exists", stateResource.Primary.ID)
			}

			// If we get a 404, that's what we expect after deletion
			if !errors.Is(err, firehydrant.ErrorNotFound) {
				return fmt.Errorf("unexpected error checking role deletion: %v", err)
			}
		}

		return nil
	}
}
