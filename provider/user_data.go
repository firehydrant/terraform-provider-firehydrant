package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantUser,
		Schema: map[string]*schema.Schema{
			// Required
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataFireHydrantUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the user
	email := d.Get("email").(string)
	tflog.Debug(ctx, fmt.Sprintf("Fetch user: %s", email), map[string]interface{}{
		"email": email,
	})

	params := firehydrant.GetUserParams{Query: email}
	userResponse, err := firehydrantAPIClient.GetUsers(ctx, params)
	if err != nil {
		return diag.Errorf("Error fetching user '%s': %v", email, err)
	}

	if len(userResponse.Users) == 0 {
		return diag.Errorf("Did not find user matching '%s'", email)
	}
	if len(userResponse.Users) > 1 {
		return diag.Errorf("Found multiple matching users for '%s'", email)
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"id":   userResponse.Users[0].ID,
		"name": userResponse.Users[0].Name,
	}

	// Set the data source attributes to the values we got from the API
	for key, val := range attributes {
		if err := d.Set(key, val); err != nil {
			return diag.Errorf("Error setting %s for user %s: %v", key, email, err)
		}
	}

	// Set the user's ID in state
	d.SetId(userResponse.Users[0].ID)

	return diag.Diagnostics{}
}
