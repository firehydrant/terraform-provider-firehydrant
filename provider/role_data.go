package provider

import (
	"context"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantRole,
		Schema: map[string]*schema.Schema{
			// Input parameters - exactly one required
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The ID of the role",
				ConflictsWith: []string{"slug"},
			},
			"slug": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The slug of the role",
				ConflictsWith: []string{"id"},
			},

			// Computed attributes - all role fields available as outputs
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the role",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A description of the role",
			},
			"permissions": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "List of permission slugs assigned to this role",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"built_in": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is a built-in role",
			},
			"read_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this role is read-only",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the role was created",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the role was last updated",
			},
		},
	}
}

func dataFireHydrantRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	var role *components.PublicAPIV1RoleEntity
	var err error

	// Determine lookup method
	if id, ok := d.GetOk("id"); ok {
		// Direct ID lookup
		roleID := id.(string)
		tflog.Debug(ctx, "Looking up role by ID", map[string]interface{}{
			"id": roleID,
		})

		role, err = client.Sdk.Roles.GetRole(ctx, roleID)
		if err != nil {
			return diag.Errorf("Error reading role %s: %v", roleID, err)
		}
	} else if slug, ok := d.GetOk("slug"); ok {
		// Slug lookup - need to search through all roles
		roleSlug := slug.(string)
		tflog.Debug(ctx, "Looking up role by slug", map[string]interface{}{
			"slug": roleSlug,
		})

		role, err = findRoleBySlug(ctx, client, roleSlug)
		if err != nil {
			return diag.Errorf("Error finding role with slug %s: %v", roleSlug, err)
		}
		if role == nil {
			return diag.Errorf("No role found with slug: %s", roleSlug)
		}
	} else {
		return diag.Errorf("Either 'id' or 'slug' must be specified")
	}

	// Extract permission slugs
	permissionSlugs := make([]string, len(role.Permissions))
	for i, p := range role.Permissions {
		permissionSlugs[i] = *p.Slug
	}

	// Set all computed attributes
	attributes := map[string]interface{}{
		"id":          *role.ID,
		"name":        *role.Name,
		"slug":        *role.Slug,
		"description": *role.Description,
		"permissions": schema.NewSet(schema.HashString, convertStringSliceToInterface(permissionSlugs)),

		"built_in":   *role.BuiltIn,
		"read_only":  *role.ReadOnly,
		"created_at": *role.CreatedAt,
		"updated_at": *role.UpdatedAt,
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for role %s: %v", key, *role.ID, err)
		}
	}

	// Set the data source ID
	d.SetId(*role.ID)

	return diag.Diagnostics{}
}

// findRoleBySlug searches for a role by slug since there's no direct API endpoint
func findRoleBySlug(ctx context.Context, client *firehydrant.APIClient, slug string) (*components.PublicAPIV1RoleEntity, error) {
	// List all roles and find the one with matching slug
	// This is less efficient than direct lookup but necessary given API design
	roles, err := client.Sdk.Roles.ListRoles(ctx, &slug, nil, nil)
	if err != nil {
		return nil, err
	}

	for _, role := range roles.Data {
		if *role.Slug == slug {
			return &role, nil
		}
	}

	return nil, nil // Not found
}
