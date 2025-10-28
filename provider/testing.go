package provider

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	sharedProviderInstance *firehydrant.APIClient
	sharedProviderOnce     sync.Once
	sharedProviderErr      error
)

// getSharedProvider returns a singleton provider instance for testing.
// This reduces API calls and helps prevent rate limiting during test runs.
func getSharedProvider() (*firehydrant.APIClient, error) {
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
			client, err := getSharedProvider()
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
