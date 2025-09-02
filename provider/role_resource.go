package provider

import (
	"context"
	"errors"

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
			"organization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization ID this role belongs to",
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

func createResourceFireHydrantRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	name := d.Get("name").(string)
	tflog.Debug(ctx, "Create role", map[string]interface{}{
		"name": name,
	})

	// Build create request
	createReq := firehydrant.CreateRoleRequest{
		Name: name,
	}

	if v, ok := d.GetOk("slug"); ok {
		createReq.Slug = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		createReq.Description = v.(string)
	}

	if v, ok := d.GetOk("permissions"); ok {
		permissions := v.(*schema.Set).List()
		permissionSlugs := make([]string, len(permissions))
		for i, p := range permissions {
			permissionSlugs[i] = p.(string)
		}
		createReq.Permissions = permissionSlugs
	}

	// Create the role
	role, err := firehydrantAPIClient.Roles().Create(ctx, createReq)
	if err != nil {
		return diag.Errorf("Error creating role %s: %v", name, err)
	}

	// Set the role ID in state
	d.SetId(role.ID)

	return readResourceFireHydrantRole(ctx, d, m)
}

func readResourceFireHydrantRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	id := d.Id()
	tflog.Debug(ctx, "Read role", map[string]interface{}{
		"id": id,
	})

	role, err := firehydrantAPIClient.Roles().Get(ctx, id)
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
		permissionSlugs[i] = p.Slug
	}

	// Update state with current values
	attributes := map[string]interface{}{
		"name":            role.Name,
		"slug":            role.Slug,
		"description":     role.Description,
		"permissions":     schema.NewSet(schema.HashString, convertStringSliceToInterface(permissionSlugs)),
		"organization_id": role.OrganizationID,
		"built_in":        role.BuiltIn,
		"read_only":       role.ReadOnly,
		"created_at":      role.CreatedAt,
		"updated_at":      role.UpdatedAt,
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for role %s: %v", key, id, err)
		}
	}

	return diag.Diagnostics{}
}

func updateResourceFireHydrantRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	id := d.Id()
	tflog.Debug(ctx, "Update role", map[string]interface{}{
		"id": id,
	})

	updateReq := firehydrant.UpdateRoleRequest{}

	if d.HasChange("name") {
		updateReq.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		updateReq.Description = d.Get("description").(string)
	}

	if d.HasChange("permissions") {
		if v, ok := d.GetOk("permissions"); ok {
			permissions := v.(*schema.Set).List()
			permissionSlugs := make([]string, len(permissions))
			for i, p := range permissions {
				permissionSlugs[i] = p.(string)
			}
			updateReq.Permissions = permissionSlugs
		} else {
			// If permissions set is removed, send empty array
			updateReq.Permissions = []string{}
		}
	}

	_, err := firehydrantAPIClient.Roles().Update(ctx, id, updateReq)
	if err != nil {
		return diag.Errorf("Error updating role %s: %v", id, err)
	}

	return readResourceFireHydrantRole(ctx, d, m)
}

func deleteResourceFireHydrantRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	firehydrantAPIClient := m.(firehydrant.Client)

	id := d.Id()
	tflog.Debug(ctx, "Delete role", map[string]interface{}{
		"id": id,
	})

	err := firehydrantAPIClient.Roles().Delete(ctx, id)
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
