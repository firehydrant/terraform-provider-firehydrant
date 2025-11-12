package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/firehydrant-go-sdk/models/operations"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
)

// SharedTestResources holds references to pre-existing production resources
// that can be reused across tests to avoid creating/destroying common resources
type SharedTestResources struct {
	Teams                 map[string]string `json:"teams"`                   // name -> ID
	Users                 map[string]string `json:"users"`                   // email -> ID
	IncidentRoles         map[string]string `json:"incident_roles"`          // name -> ID
	Environments          map[string]string `json:"environments"`            // name -> ID
	Services              map[string]string `json:"services"`                // name -> ID
	Priorities            map[string]string `json:"priorities"`              // name -> ID
	Severities            map[string]string `json:"severities"`              // name -> ID
	OnCallSchedules       map[string]string `json:"on_call_schedules"`       // name -> ID
	EscalationPolicies    map[string]string `json:"escalation_policies"`     // name -> ID
	Functionalities       map[string]string `json:"functionalities"`         // name -> ID
	Roles                 map[string]string `json:"roles"`                   // name -> ID
	Rotations             map[string]string `json:"rotations"`               // name -> ID
	Runbooks              map[string]string `json:"runbooks"`                // name -> ID
	TaskLists             map[string]string `json:"task_lists"`              // name -> ID
	IncidentTypes         map[string]string `json:"incident_types"`          // name -> ID
	SignalRules           map[string]string `json:"signal_rules"`            // name -> ID
	StatusUpdateTemplates map[string]string `json:"status_update_templates"` // name -> ID
	InboundEmails         map[string]string `json:"inbound_emails"`          // name -> ID
	CustomEventSources    map[string]string `json:"custom_event_sources"`    // name -> ID

	// Track which resources are created vs pre-existing
	CreatedResources CreatedResources
}

// CreatedResources tracks resources that were created during test initialization
// These will be cleaned up at the end of the test run
type CreatedResources struct {
	TeamIDs           []string `json:"team_ids"`
	OnCallScheduleIDs []string `json:"on_call_schedule_ids"`
	IncidentRoleIDs   []string `json:"incident_role_ids"`
	ServiceIDs        []string `json:"service_ids"`
}

var (
	sharedTestResourcesInstance *SharedTestResources
	sharedTestResourcesOnce     sync.Once
	sharedTestResourcesErr      error
)

// getSharedTestResources returns a singleton instance of shared test resources.
// This loads pre-existing production resources that can be reused across tests.
func getSharedTestResources() (*SharedTestResources, error) {
	sharedTestResourcesOnce.Do(func() {
		resources := &SharedTestResources{
			Teams:                 make(map[string]string),
			Users:                 make(map[string]string),
			IncidentRoles:         make(map[string]string),
			Environments:          make(map[string]string),
			Services:              make(map[string]string),
			Priorities:            make(map[string]string),
			Severities:            make(map[string]string),
			OnCallSchedules:       make(map[string]string),
			EscalationPolicies:    make(map[string]string),
			Functionalities:       make(map[string]string),
			Roles:                 make(map[string]string),
			Rotations:             make(map[string]string),
			Runbooks:              make(map[string]string),
			TaskLists:             make(map[string]string),
			IncidentTypes:         make(map[string]string),
			SignalRules:           make(map[string]string),
			StatusUpdateTemplates: make(map[string]string),
			InboundEmails:         make(map[string]string),
			CustomEventSources:    make(map[string]string),
		}

		// Try to load from environment variable first
		if err := resources.LoadFromEnvironment(); err != nil {
			// Fallback to API discovery
			if err := resources.LoadFromAPI(); err != nil {
				sharedTestResourcesErr = fmt.Errorf("could not load shared test resources: %w", err)
				return
			}
		}

		sharedTestResourcesInstance = resources
	})

	return sharedTestResourcesInstance, sharedTestResourcesErr
}

// InitializeSharedResources creates shared test resources if they don't exist
// Called once at the start of test run via TestMain
func (r *SharedTestResources) InitializeSharedResources(ctx context.Context, client *firehydrant.APIClient) error {
	// Create default team if not in env
	if len(r.Teams) == 0 {
		teamID, err := r.createSharedTeam(ctx, client, fmt.Sprintf("tf-test-shared-default-%d", time.Now().Unix()))
		if err != nil {
			return fmt.Errorf("could not create shared team: %w", err)
		}
		r.Teams["default"] = teamID
		r.CreatedResources.TeamIDs = append(r.CreatedResources.TeamIDs, teamID)
	}

	// Create default on-call schedule if not in env
	if len(r.OnCallSchedules) == 0 {
		scheduleID, err := r.createSharedOnCallSchedule(ctx, client, fmt.Sprintf("tf-test-shared-schedule-%d", time.Now().Unix()), r.Teams["default"])
		if err != nil {
			return fmt.Errorf("could not create shared on-call schedule: %w", err)
		}
		r.OnCallSchedules["default"] = scheduleID
		r.CreatedResources.OnCallScheduleIDs = append(r.CreatedResources.OnCallScheduleIDs, scheduleID)
	}

	// Create default incident role if not in env
	if len(r.IncidentRoles) == 0 {
		roleID, err := r.createSharedIncidentRole(ctx, client, fmt.Sprintf("tf-test-shared-role-%d", time.Now().Unix()))
		if err != nil {
			return fmt.Errorf("could not create shared incident role: %w", err)
		}
		r.IncidentRoles["default"] = roleID
		r.CreatedResources.IncidentRoleIDs = append(r.CreatedResources.IncidentRoleIDs, roleID)
	}

	// Create default service if not in env
	if len(r.Services) == 0 {
		serviceID, err := r.createSharedService(ctx, client, fmt.Sprintf("tf-test-shared-service-%d", time.Now().Unix()))
		if err != nil {
			return fmt.Errorf("could not create shared service: %w", err)
		}
		r.Services["default"] = serviceID
		r.CreatedResources.ServiceIDs = append(r.CreatedResources.ServiceIDs, serviceID)
	}

	// Create second shared service for dependency tests
	if _, exists := r.Services["service2"]; !exists {
		serviceID, err := r.createSharedService(ctx, client, fmt.Sprintf("tf-test-shared-service-2-%d", time.Now().Unix()))
		if err != nil {
			return fmt.Errorf("could not create shared service 2: %w", err)
		}
		r.Services["service2"] = serviceID
		r.CreatedResources.ServiceIDs = append(r.CreatedResources.ServiceIDs, serviceID)
	}

	// Create second shared team for tests that need multiple teams
	if _, exists := r.Teams["team2"]; !exists {
		teamID, err := r.createSharedTeam(ctx, client, fmt.Sprintf("tf-test-shared-team2-%d", time.Now().Unix()))
		if err != nil {
			return fmt.Errorf("could not create shared team 2: %w", err)
		}
		r.Teams["team2"] = teamID
		r.CreatedResources.TeamIDs = append(r.CreatedResources.TeamIDs, teamID)
	}

	return nil
}

// stringPtr returns a pointer to the given string
func stringPtr(s string) *string {
	return &s
}

// createSharedTeam creates a team with the given name
func (r *SharedTestResources) createSharedTeam(ctx context.Context, client *firehydrant.APIClient, name string) (string, error) {
	createRequest := components.CreateTeam{
		Name:        name,
		Description: &name, // Use name as description for test resources
	}

	teamResponse, err := client.Sdk.Teams.CreateTeam(ctx, createRequest)
	if err != nil {
		return "", fmt.Errorf("failed to create team %s: %w", name, err)
	}

	return *teamResponse.GetID(), nil
}

// createSharedOnCallSchedule creates an on-call schedule with the given name and team ID
func (r *SharedTestResources) createSharedOnCallSchedule(ctx context.Context, client *firehydrant.APIClient, name, teamID string) (string, error) {
	// Use current time as start time to ensure it's within the allowed range
	startTime := time.Now().UTC().Format(time.RFC3339)

	// Create a basic on-call schedule with the team
	createRequest := components.CreateTeamOnCallSchedule{
		Name:        name,
		Description: &name, // Use name as description for test resources
		TimeZone:    stringPtr("UTC"),
		Strategy: &components.CreateTeamOnCallScheduleStrategy{
			Type:          components.CreateTeamOnCallScheduleTypeWeekly,
			HandoffTime:   stringPtr("09:00"),
			HandoffDay:    (*components.CreateTeamOnCallScheduleHandoffDay)(stringPtr("monday")),
			ShiftDuration: stringPtr("7"),
		},
		StartTime: stringPtr(startTime),
		MemberIds: []string{}, // Empty for now, can be populated later if needed
	}

	scheduleResponse, err := client.Sdk.Signals.CreateTeamOnCallSchedule(ctx, teamID, createRequest)
	if err != nil {
		return "", fmt.Errorf("failed to create on-call schedule %s: %w", name, err)
	}

	return *scheduleResponse.GetID(), nil
}

// createSharedIncidentRole creates an incident role with the given name
func (r *SharedTestResources) createSharedIncidentRole(ctx context.Context, client *firehydrant.APIClient, name string) (string, error) {
	createRequest := components.CreateIncidentRole{
		Name:    name,
		Summary: fmt.Sprintf("Test incident role: %s", name),
	}

	incidentRole, err := client.Sdk.IncidentSettings.CreateIncidentRole(ctx, createRequest)
	if err != nil {
		return "", fmt.Errorf("failed to create incident role %s: %w", name, err)
	}

	return *incidentRole.GetID(), nil
}

// createSharedService creates a service with the given name
func (r *SharedTestResources) createSharedService(ctx context.Context, client *firehydrant.APIClient, name string) (string, error) {
	createRequest := components.CreateService{
		Name: name,
		Labels: map[string]string{
			"test": "shared",
		},
	}

	serviceResponse, err := client.Sdk.CatalogEntries.CreateService(ctx, createRequest)
	if err != nil {
		return "", fmt.Errorf("failed to create service %s: %w", name, err)
	}

	return *serviceResponse.ID, nil
}

// LoadFromEnvironment loads shared resources from FIREHYDRANT_TEST_RESOURCES JSON environment variable
func (r *SharedTestResources) LoadFromEnvironment() error {
	envResources := os.Getenv("FIREHYDRANT_TEST_RESOURCES")
	if envResources == "" {
		return fmt.Errorf("FIREHYDRANT_TEST_RESOURCES not set")
	}

	if err := json.Unmarshal([]byte(envResources), r); err != nil {
		return fmt.Errorf("could not parse FIREHYDRANT_TEST_RESOURCES JSON: %w", err)
	}

	return nil
}

// LoadFromAPI discovers shared resources by querying the API for resources with test naming convention
func (r *SharedTestResources) LoadFromAPI() error {
	client, err := getAccTestClient()
	if err != nil {
		return fmt.Errorf("could not get shared provider for API discovery: %w", err)
	}

	ctx := context.Background()

	// Load teams with test naming convention
	if err := r.loadTeamsFromAPI(ctx, client); err != nil {
		return fmt.Errorf("could not load teams from API: %w", err)
	}

	// Load users with test naming convention
	if err := r.loadUsersFromAPI(ctx, client); err != nil {
		return fmt.Errorf("could not load users from API: %w", err)
	}

	// Load incident roles with test naming convention
	if err := r.loadIncidentRolesFromAPI(ctx, client); err != nil {
		return fmt.Errorf("could not load incident roles from API: %w", err)
	}

	// Add more resource types as needed...

	return nil
}

// loadTeamsFromAPI discovers teams with tf-test-shared- prefix
func (r *SharedTestResources) loadTeamsFromAPI(ctx context.Context, client *firehydrant.APIClient) error {
	request := operations.ListTeamsRequest{}
	teamsResponse, err := client.Sdk.Teams.ListTeams(ctx, request)
	if err != nil {
		return err
	}

	for _, team := range teamsResponse.GetData() {
		if team.GetName() != nil && len(*team.GetName()) > 12 && (*team.GetName())[:12] == "tf-test-shared" {
			r.Teams[*team.GetName()] = *team.GetID()
		}
	}

	return nil
}

// loadUsersFromAPI discovers users with test email patterns
func (r *SharedTestResources) loadUsersFromAPI(ctx context.Context, client *firehydrant.APIClient) error {
	// Note: This would need to be implemented based on available user API endpoints
	// For now, we'll rely on environment variable configuration
	return nil
}

// loadIncidentRolesFromAPI discovers incident roles with tf-test-shared- prefix
func (r *SharedTestResources) loadIncidentRolesFromAPI(ctx context.Context, client *firehydrant.APIClient) error {
	// Note: This would need to be implemented based on available incident role API endpoints
	// For now, we'll rely on environment variable configuration
	return nil
}

// Helper methods for accessing shared resources

// GetTeamID returns the ID for a team by name, or error if not found
func (r *SharedTestResources) GetTeamID(name string) (string, error) {
	if id, exists := r.Teams[name]; exists {
		return id, nil
	}
	return "", fmt.Errorf("shared team %s not found", name)
}

// GetUserID returns the ID for a user by email, or error if not found
func (r *SharedTestResources) GetUserID(email string) (string, error) {
	if id, exists := r.Users[email]; exists {
		return id, nil
	}
	return "", fmt.Errorf("shared user %s not found", email)
}

// GetIncidentRoleID returns the ID for an incident role by name, or error if not found
func (r *SharedTestResources) GetIncidentRoleID(name string) (string, error) {
	if id, exists := r.IncidentRoles[name]; exists {
		return id, nil
	}
	return "", fmt.Errorf("shared incident role %s not found", name)
}

// GetEnvironmentID returns the ID for an environment by name, or error if not found
func (r *SharedTestResources) GetEnvironmentID(name string) (string, error) {
	if id, exists := r.Environments[name]; exists {
		return id, nil
	}
	return "", fmt.Errorf("shared environment %s not found", name)
}

// GetServiceID returns the ID for a service by name, or error if not found
func (r *SharedTestResources) GetServiceID(name string) (string, error) {
	if id, exists := r.Services[name]; exists {
		return id, nil
	}
	return "", fmt.Errorf("shared service %s not found", name)
}

// GetPriorityID returns the ID for a priority by name, or error if not found
func (r *SharedTestResources) GetPriorityID(name string) (string, error) {
	if id, exists := r.Priorities[name]; exists {
		return id, nil
	}
	return "", fmt.Errorf("shared priority %s not found", name)
}

// GetSeverityID returns the ID for a severity by name, or error if not found
func (r *SharedTestResources) GetSeverityID(name string) (string, error) {
	if id, exists := r.Severities[name]; exists {
		return id, nil
	}
	return "", fmt.Errorf("shared severity %s not found", name)
}

// GetOnCallScheduleID returns the ID for an on-call schedule by name, or error if not found
func (r *SharedTestResources) GetOnCallScheduleID(name string) (string, error) {
	if id, exists := r.OnCallSchedules[name]; exists {
		return id, nil
	}
	return "", fmt.Errorf("shared on-call schedule %s not found", name)
}

// GetEscalationPolicyID returns the ID for an escalation policy by name, or error if not found
func (r *SharedTestResources) GetEscalationPolicyID(name string) (string, error) {
	if id, exists := r.EscalationPolicies[name]; exists {
		return id, nil
	}
	return "", fmt.Errorf("shared escalation policy %s not found", name)
}

// DestroyCreatedResources cleans up resources that were created during test initialization
// Called once at the end of test run via TestMain
func (r *SharedTestResources) DestroyCreatedResources(ctx context.Context, client *firehydrant.APIClient) error {
	// Destroy in reverse order of creation to handle dependencies
	for _, roleID := range r.CreatedResources.IncidentRoleIDs {
		if err := client.Sdk.IncidentSettings.DeleteIncidentRole(ctx, roleID); err != nil {
			fmt.Printf("Warning: Failed to delete incident role %s: %v\n", roleID, err)
		}
	}

	for _, scheduleID := range r.CreatedResources.OnCallScheduleIDs {
		// Find the team ID for this schedule (we'll need to track this or use a default)
		// For now, we'll use the first team ID as a fallback
		teamID := ""
		if len(r.CreatedResources.TeamIDs) > 0 {
			teamID = r.CreatedResources.TeamIDs[0]
		}
		if teamID != "" {
			if err := client.Sdk.Signals.DeleteTeamOnCallSchedule(ctx, teamID, scheduleID); err != nil {
				fmt.Printf("Warning: Failed to delete on-call schedule %s: %v\n", scheduleID, err)
			}
		}
	}

	for _, teamID := range r.CreatedResources.TeamIDs {
		if err := client.Sdk.Teams.DeleteTeam(ctx, teamID); err != nil {
			fmt.Printf("Warning: Failed to delete team %s: %v\n", teamID, err)
		}
	}

	// Destroy services
	for _, serviceID := range r.CreatedResources.ServiceIDs {
		if err := client.Sdk.CatalogEntries.DeleteService(ctx, serviceID); err != nil {
			fmt.Printf("Warning: Failed to delete service %s: %v\n", serviceID, err)
		}
	}

	return nil
}

// resetSharedTestResources resets the shared test resources for testing.
// This should only be used in test cleanup or when testing resource loading itself.
func resetSharedTestResources() {
	sharedTestResourcesOnce = sync.Once{}
	sharedTestResourcesInstance = nil
	sharedTestResourcesErr = nil
}
