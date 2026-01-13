package provider

import (
	"context"
	"errors"
	"time"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantRole,
		ReadContext:   readResourceFireHydrantRole,
		UpdateContext: updateResourceFireHydrantRole,
		DeleteContext: deleteResourceFireHydrantRole,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the role",
			},
			"slug": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true, // Slug can't be changed after creation
				Description: "The slug for the role. If not provided, will be auto-generated from name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the role",
			},
			"permissions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "List of permission slugs assigned to this role",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Computed attributes

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

func createResourceFireHydrantRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	name := d.Get("name").(string)
	tflog.Debug(ctx, "Create role", map[string]interface{}{
		"name": name,
	})

	// Build create request
	createReq := components.CreateRole{
		Name: name,
	}

	if v, ok := d.GetOk("slug"); ok {
		slug := v.(string)
		createReq.Slug = &slug
	}

	if v, ok := d.GetOk("description"); ok {
		desc := v.(string)
		createReq.Description = &desc
	}

	if v, ok := d.GetOk("permissions"); ok {
		permissions := v.(*schema.Set).List()
		permissionSlugs := make([]components.CreateRolePermission, len(permissions))
		for i, p := range permissions {
			permissionSlugs[i] = components.CreateRolePermission(p.(string))
		}
		createReq.Permissions = permissionSlugs
	}

	// Create the role
	role, err := client.Sdk.Roles.CreateRole(ctx, createReq)
	if err != nil {
		return diag.Errorf("Error creating role %s: %v", name, err)
	}

	// Set the role ID in state
	d.SetId(*role.ID)

	return readResourceFireHydrantRole(ctx, d, m)
}

func readResourceFireHydrantRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	tflog.Debug(ctx, "Read role", map[string]interface{}{
		"id": id,
	})

	role, err := client.Sdk.Roles.GetRole(ctx, id)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			tflog.Debug(ctx, "Role not found, removing from state", map[string]interface{}{
				"id": id,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading role %s: %v", id, err)
	}

	// Extract permission slugs from the full permission objects
	permissionSlugs := make([]string, len(role.Permissions))
	for i, p := range role.Permissions {
		permissionSlugs[i] = *p.Slug
	}

	createdAtString := role.CreatedAt.Format(time.RFC3339)
	updatedAtString := role.UpdatedAt.Format(time.RFC3339)

	// Update state with current values
	attributes := map[string]interface{}{
		"name":        *role.Name,
		"slug":        *role.Slug,
		"description": *role.Description,
		"permissions": schema.NewSet(schema.HashString, convertStringSliceToInterface(permissionSlugs)),

		"built_in":   *role.BuiltIn,
		"read_only":  *role.ReadOnly,
		"created_at": createdAtString,
		"updated_at": updatedAtString,
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for role %s: %v", key, id, err)
		}
	}

	return diag.Diagnostics{}
}

func updateResourceFireHydrantRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	tflog.Debug(ctx, "Update role", map[string]interface{}{
		"id": id,
	})

	updateReq := components.UpdateRole{}

	updateReq.Name = d.Get("name").(string)
	desc := d.Get("description").(string)
	updateReq.Description = &desc

	if d.HasChange("permissions") {
		if v, ok := d.GetOk("permissions"); ok {
			permissions := v.(*schema.Set).List()
			permissionSlugs := make([]components.UpdateRolePermission, len(permissions))
			for i, p := range permissions {
				permissionSlugs[i] = p.(components.UpdateRolePermission)
			}
			updateReq.Permissions = permissionSlugs
		} else {
			// If permissions set is removed, send empty array
			updateReq.Permissions = []components.UpdateRolePermission{}
		}
	}

	_, err := client.Sdk.Roles.UpdateRole(ctx, id, updateReq)
	if err != nil {
		return diag.Errorf("Error updating role %s: %v", id, err)
	}

	return readResourceFireHydrantRole(ctx, d, m)
}

func deleteResourceFireHydrantRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	tflog.Debug(ctx, "Delete role", map[string]interface{}{
		"id": id,
	})

	err := client.Sdk.Roles.DeleteRole(ctx, id)
	if err != nil {
		return diag.Errorf("Error deleting role %s: %v", id, err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}

// Helper function to convert []string to []interface{} for Terraform sets
func convertStringSliceToInterface(strings []string) []interface{} {
	interfaces := make([]interface{}, len(strings))
	for i, s := range strings {
		interfaces[i] = s
	}
	return interfaces
}
