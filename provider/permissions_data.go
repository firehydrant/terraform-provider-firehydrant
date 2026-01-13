package provider

import (
	"context"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source for listing all available permissions
func dataSourcePermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantPermissions,
		Schema: map[string]*schema.Schema{
			"permissions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of all available permissions",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"slug": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The permission slug",
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Human-readable name of the permission",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of what this permission allows",
						},
						"category_display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name of the permission category",
						},
						"category_slug": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Slug of the permission category",
						},
						"parent_slug": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parent permission slug if this permission has dependencies",
						},
						"available": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this permission is available in the current context",
						},
						"dependency_slugs": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of permission slugs this permission depends on",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataFireHydrantPermissions(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	tflog.Debug(ctx, "Reading all permissions")

	permissionsResp, err := client.Sdk.Permissions.ListPermissions(ctx)
	if err != nil {
		return diag.Errorf("Error reading permissions: %v", err)
	}

	// Convert permissions to Terraform format
	permissions := make([]map[string]interface{}, len(permissionsResp.Data))
	for i, permission := range permissionsResp.Data {
		permissions[i] = map[string]interface{}{
			"slug":                  *permission.Slug,
			"display_name":          *permission.DisplayName,
			"description":           *permission.Description,
			"category_display_name": *permission.CategoryDisplayName,
			"category_slug":         *permission.CategorySlug,
			"parent_slug":           *permission.ParentSlug,
			"available":             *permission.Available,
			"dependency_slugs":      schema.NewSet(schema.HashString, convertStringSliceToInterface(permission.DependencySlugs)),
		}
	}

	if err := d.Set("permissions", permissions); err != nil {
		return diag.Errorf("Error setting permissions: %v", err)
	}

	// Use a static ID since this is a singleton data source
	d.SetId("permissions")

	return diag.Diagnostics{}
}
