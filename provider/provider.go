package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	apiKeyName             = "api_key"
	firehydrantBaseURLName = "firehydrant_base_url"
)

const (
	// MajorVersion is the major version
	MajorVersion = 0
	// MinorVersion is the minor version
	MinorVersion = 1
	// PatchVersion is the patch version
	PatchVersion = 0

	// UserAgentPrefix is the prefix of the User-Agent header that all terraform REST calls perform
	UserAgentPrefix = "firehydrant-terraform-provider"
)

// Version is the semver of this provider
var Version = fmt.Sprintf("%d.%d.%d", MajorVersion, MinorVersion, PatchVersion)

// Provider returns a terraform provider for the FireHydrant API
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			apiKeyName: {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("FIREHYDRANT_API_KEY", nil),
			},
			firehydrantBaseURLName: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FIREHYDRANT_BASE_URL", "https://api.firehydrant.io/v1/"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"firehydrant_environment":        resourceEnvironment(),
			"firehydrant_functionality":      resourceFunctionality(),
			"firehydrant_incident_role":      resourceIncidentRole(),
			"firehydrant_priority":           resourcePriority(),
			"firehydrant_runbook":            resourceRunbook(),
			"firehydrant_service_dependency": resourceServiceDependency(),
			"firehydrant_service":            resourceService(),
			"firehydrant_severity":           resourceSeverity(),
			"firehydrant_task_list":          resourceTaskList(),
			"firehydrant_team":               resourceTeam(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"firehydrant_environment":    dataSourceEnvironment(),
			"firehydrant_functionality":  dataSourceFunctionality(),
			"firehydrant_incident_role":  dataSourceIncidentRole(),
			"firehydrant_priority":       dataSourcePriority(),
			"firehydrant_runbook":        dataSourceRunbook(),
			"firehydrant_runbook_action": dataSourceRunbookAction(),
			"firehydrant_service":        dataSourceService(),
			"firehydrant_services":       dataSourceServices(),
			"firehydrant_severity":       dataSourceSeverity(),
			"firehydrant_task_list":      dataSourceTaskList(),
			"firehydrant_schedule":       dataSourceSchedule(),
			"firehydrant_user":           dataSourceUser(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, rd *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := provider.TerraformVersion

		if terraformVersion == "" {
			terraformVersion = "0.11+compatible"
		}

		return setupFireHydrantContext(ctx, rd, terraformVersion)
	}

	return provider
}

func setupFireHydrantContext(ctx context.Context, rd *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	apiKey := rd.Get(apiKeyName).(string)
	fireHydrantBaseURL := rd.Get(firehydrantBaseURLName).(string)

	ac, err := firehydrant.NewRestClient(apiKey, firehydrant.WithBaseURL(fireHydrantBaseURL), firehydrant.WithUserAgentSuffix(fmt.Sprintf("terraform-%s", terraformVersion)))
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("could not initialize API client: %w", err))
	}

	_, err = ac.Ping(ctx)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return ac, nil
}

func convertStringMap(sm map[string]interface{}) map[string]string {
	m := map[string]string{}
	for k, v := range sm {
		m[k] = v.(string)
	}

	return m
}
