package provider

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	sharedProviderInstance *firehydrant.APIClient
	sharedProviderOnce     sync.Once
	sharedProviderErr      error
)

// getAccTestClient returns a singleton API client instance for testing.
// This reduces API calls and helps prevent rate limiting during test runs.
func getAccTestClient() (*firehydrant.APIClient, error) {
	sharedProviderOnce.Do(func() {
		ctx := context.Background()
		apiKey := os.Getenv("FIREHYDRANT_API_KEY")
		if apiKey == "" {
			sharedProviderErr = fmt.Errorf("FIREHYDRANT_API_KEY not set")
			return
		}

		baseURL := os.Getenv("FIREHYDRANT_BASE_URL")
		if baseURL == "" {
			baseURL = "https://api.firehydrant.io/v1/"
		}

		// Initialize client with test user agent
		client, err := firehydrant.NewRestClient(
			apiKey,
			firehydrant.WithBaseURL(baseURL),
			firehydrant.WithUserAgentSuffix("terraform-test-shared"),
		)
		if err != nil {
			sharedProviderErr = fmt.Errorf("could not initialize shared test client: %w", err)
			return
		}

		// Single ping for the entire test suite
		_, err = client.Ping(ctx)
		if err != nil {
			sharedProviderErr = fmt.Errorf("shared test client ping failed: %w", err)
			return
		}

		// SDK ping (only if not in acceptance test mode)
		if os.Getenv("TF_ACC") != "true" {
			_, err = client.Sdk.AccountSettings.Ping(ctx)
			if err != nil {
				sharedProviderErr = fmt.Errorf("shared test SDK ping failed: %w", err)
				return
			}
		}

		sharedProviderInstance = client
	})

	return sharedProviderInstance, sharedProviderErr
}

// sharedProviderFactories returns provider factories that use a shared provider instance.
// This should be used for acceptance tests to reduce API calls and prevent rate limiting.
func sharedProviderFactories() map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"firehydrant": func() (*schema.Provider, error) {
			client, err := getAccTestClient()
			if err != nil {
				return nil, err
			}

			// Return provider with pre-configured client
			p := Provider()
			p.ConfigureContextFunc = func(ctx context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
				// Skip normal setup, return shared client
				return client, nil
			}

			return p, nil
		},
	}
}

// resetSharedProvider resets the shared provider for testing.
// This should only be used in test cleanup or when testing provider initialization itself.
func resetSharedProvider() {
	sharedProviderOnce = sync.Once{}
	sharedProviderInstance = nil
	sharedProviderErr = nil
}

// testAccCheckResourceDestroy creates a generic destroy check function for any resource type.
// This uses the shared provider's client instead of creating a new one.
func testAccCheckResourceDestroy(resourceType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != resourceType {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			// Generic check - try to get the resource and expect it to not exist
			// This is a simplified version; specific resource types may need custom logic
			switch resourceType {
			case "firehydrant_team":
				_, err := client.Sdk.Teams.GetTeam(context.TODO(), stateResource.Primary.ID, nil)
				if err == nil {
					return fmt.Errorf("%s %s still exists", resourceType, stateResource.Primary.ID)
				}
			case "firehydrant_service":
				_, err := client.Services().Get(context.TODO(), stateResource.Primary.ID)
				if err == nil {
					return fmt.Errorf("%s %s still exists", resourceType, stateResource.Primary.ID)
				}
			case "firehydrant_environment":
				_, err := client.Environments().Get(context.TODO(), stateResource.Primary.ID)
				if err == nil {
					return fmt.Errorf("%s %s still exists", resourceType, stateResource.Primary.ID)
				}
			case "firehydrant_incident_role":
				_, err := client.Sdk.IncidentSettings.GetIncidentRole(context.TODO(), stateResource.Primary.ID)
				if err == nil {
					return fmt.Errorf("%s %s still exists", resourceType, stateResource.Primary.ID)
				}
			default:
				// For unsupported resource types, return a warning but don't fail
				fmt.Printf("Warning: Generic destroy check not implemented for %s\n", resourceType)
			}
		}

		return nil
	}
}

// testAccResourceExists creates a generic existence check function for any resource type.
// This uses the shared provider's client instead of creating a new one.
func testAccResourceExists(resourceType, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		// Generic check - try to get the resource and expect it to exist
		// This is a simplified version; specific resource types may need custom logic
		switch resourceType {
		case "firehydrant_team":
			_, err := client.Sdk.Teams.GetTeam(context.TODO(), rs.Primary.ID, nil)
			if err != nil {
				return fmt.Errorf("Error getting %s: %v", resourceType, err)
			}
		case "firehydrant_service":
			_, err := client.Services().Get(context.TODO(), rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error getting %s: %v", resourceType, err)
			}
		case "firehydrant_environment":
			_, err := client.Environments().Get(context.TODO(), rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error getting %s: %v", resourceType, err)
			}
		case "firehydrant_incident_role":
			_, err := client.Sdk.IncidentSettings.GetIncidentRole(context.TODO(), rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error getting %s: %v", resourceType, err)
			}
		default:
			return fmt.Errorf("Generic existence check not implemented for %s", resourceType)
		}

		return nil
	}
}

// testAccCheckResourceExistsWithAttributes creates a generic existence check with attribute validation.
// This is a template that can be customized for specific resource types.
func testAccCheckResourceExistsWithAttributes(resourceType, resourceName string, attributeChecks map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		// Get the resource from API to verify it exists
		switch resourceType {
		case "firehydrant_team":
			_, err := client.Sdk.Teams.GetTeam(context.TODO(), rs.Primary.ID, nil)
			if err != nil {
				return fmt.Errorf("Error getting team: %v", err)
			}
		case "firehydrant_service":
			_, err := client.Services().Get(context.TODO(), rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error getting service: %v", err)
			}
		case "firehydrant_environment":
			_, err := client.Environments().Get(context.TODO(), rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error getting environment: %v", err)
			}
		default:
			return fmt.Errorf("Attribute checking not implemented for %s", resourceType)
		}

		// Validate attributes
		for attrName, expectedValue := range attributeChecks {
			actualValue := rs.Primary.Attributes[attrName]
			if actualValue != expectedValue {
				return fmt.Errorf("Attribute %s mismatch. Expected: %s, got: %s", attrName, expectedValue, actualValue)
			}
		}

		return nil
	}
}

// testAccGetSharedResourceID returns the ID of a shared resource, or fails the test if not available.
func testAccGetSharedResourceID(t *testing.T, resourceType, resourceName string) string {
	resources, err := getSharedTestResources()
	if err != nil {
		t.Fatalf("Shared test resources not available: %v", err)
	}

	switch resourceType {
	case "team":
		id, err := resources.GetTeamID(resourceName)
		if err != nil {
			t.Fatalf("Shared team %s not available: %v", resourceName, err)
		}
		return id
	case "user":
		id, err := resources.GetUserID(resourceName)
		if err != nil {
			t.Fatalf("Shared user %s not available: %v", resourceName, err)
		}
		return id
	case "incident_role":
		id, err := resources.GetIncidentRoleID(resourceName)
		if err != nil {
			t.Fatalf("Shared incident role %s not available: %v", resourceName, err)
		}
		return id
	default:
		t.Fatalf("Shared resource type %s not supported", resourceType)
		return ""
	}
}

// getSharedTeamID returns the default shared team ID for tests
// Fails the test if shared resources are not available
func getSharedTeamID(t *testing.T) string {
	resources, err := getSharedTestResources()
	if err != nil {
		t.Fatalf("Shared test resources not available: %v", err)
	}

	if id, err := resources.GetTeamID("default"); err == nil {
		return id
	}

	t.Fatalf("Shared team 'default' not available")
	return ""
}

// getSharedOnCallScheduleID returns the default shared on-call schedule ID for tests
// Fails the test if shared resources are not available
func getSharedOnCallScheduleID(t *testing.T) string {
	resources, err := getSharedTestResources()
	if err != nil {
		t.Fatalf("Shared test resources not available: %v", err)
	}

	if id, err := resources.GetOnCallScheduleID("default"); err == nil {
		return id
	}

	t.Fatalf("Shared on-call schedule 'default' not available")
	return ""
}

// getSharedIncidentRoleID returns the default shared incident role ID for tests
// Fails the test if shared resources are not available
func getSharedIncidentRoleID(t *testing.T) string {
	resources, err := getSharedTestResources()
	if err != nil {
		t.Fatalf("Shared test resources not available: %v", err)
	}

	if id, err := resources.GetIncidentRoleID("default"); err == nil {
		return id
	}

	t.Fatalf("Shared incident role 'default' not available")
	return ""
}

// getSharedServiceID returns the default shared service ID for tests
// Fails the test if shared resources are not available
func getSharedServiceID(t *testing.T) string {
	resources, err := getSharedTestResources()
	if err != nil {
		t.Fatalf("Shared test resources not available: %v", err)
	}

	if id, err := resources.GetServiceID("default"); err == nil {
		return id
	}

	t.Fatalf("Shared service 'default' not available")
	return ""
}

// getSharedServiceID2 returns the second shared service ID for dependency tests
func getSharedServiceID2(t *testing.T) string {
	resources, err := getSharedTestResources()
	if err != nil {
		t.Fatalf("Shared test resources not available: %v", err)
	}

	if id, err := resources.GetServiceID("service2"); err == nil {
		return id
	}

	t.Fatalf("Shared service 'service2' not available")
	return ""
}
