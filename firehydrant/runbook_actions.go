package firehydrant

import (
	"context"
	"fmt"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

// RunbookActionIntegrationSlug represents the integration slug of the integration
// associated with a runbook action
type RunbookActionIntegrationSlug string

// List of valid integration slugs
const (
	RunbookActionIntegrationSlugConfluenceCloud RunbookActionIntegrationSlug = "confluence_cloud"
	RunbookActionIntegrationSlugFireHydrant     RunbookActionIntegrationSlug = "patchy"
	RunbookActionIntegrationSlugFireHydrantNunc RunbookActionIntegrationSlug = "nunc"
	RunbookActionIntegrationSlugGiphy           RunbookActionIntegrationSlug = "giphy"
	RunbookActionIntegrationSlugGoogleDocs      RunbookActionIntegrationSlug = "google_docs"
	RunbookActionIntegrationSlugGoogleMeet      RunbookActionIntegrationSlug = "google_meet"
	RunbookActionIntegrationSlugJiraCloud       RunbookActionIntegrationSlug = "jira_cloud"
	RunbookActionIntegrationSlugJiraServer      RunbookActionIntegrationSlug = "jira_server"
	RunbookActionIntegrationSlugMicrosoftTeams  RunbookActionIntegrationSlug = "microsoft_teams"
	RunbookActionIntegrationSlugOpsgenie        RunbookActionIntegrationSlug = "opsgenie"
	RunbookActionIntegrationSlugPagerDuty       RunbookActionIntegrationSlug = "pager_duty"
	RunbookActionIntegrationSlugShortcut        RunbookActionIntegrationSlug = "shortcut"
	RunbookActionIntegrationSlugSlack           RunbookActionIntegrationSlug = "slack"
	RunbookActionIntegrationSlugStatuspage      RunbookActionIntegrationSlug = "statuspage"
	RunbookActionIntegrationSlugVictorOps       RunbookActionIntegrationSlug = "victorops"
	RunbookActionIntegrationSlugWebex           RunbookActionIntegrationSlug = "webex"
	RunbookActionIntegrationSlugZoom            RunbookActionIntegrationSlug = "zoom"
)

// RunbookActionSlug represents the runbook action's slug
type RunbookActionSlug string

// List of valid action slugs
const (
	RunbookActionSlugAddServicesRelatedToFunctionality  RunbookActionSlug = "add_services_related_to_functionality"
	RunbookActionSlugAddTaskList                        RunbookActionSlug = "add_task_list"
	RunbookActionSlugArchiveIncidentChannel             RunbookActionSlug = "archive_incident_channel"
	RunbookActionSlugAssignARole                        RunbookActionSlug = "assign_a_role"
	RunbookActionSlugAssignATeam                        RunbookActionSlug = "assign_a_team"
	RunbookActionSlugAttachARunbook                     RunbookActionSlug = "attach_a_runbook"
	RunbookActionSlugCreateGoogleMeetLink               RunbookActionSlug = "create_google_meet_link"
	RunbookActionSlugCreateIncidentChannel              RunbookActionSlug = "create_incident_channel"
	RunbookActionSlugCreateIncidentIssue                RunbookActionSlug = "create_incident_issue"
	RunbookActionSlugCreateIncidentTicket               RunbookActionSlug = "create_incident_ticket"
	RunbookActionSlugCreateMeeting                      RunbookActionSlug = "create_meeting"
	RunbookActionSlugCreateNewOpsgenieIncident          RunbookActionSlug = "create_new_opsgenie_incident"
	RunbookActionSlugCreateNewPagerDutyIncident         RunbookActionSlug = "create_new_pager_duty_incident"
	RunbookActionSlugCreateNunc                         RunbookActionSlug = "create_nunc"
	RunbookActionSlugCreateStatuspage                   RunbookActionSlug = "create_statuspage"
	RunbookActionSlugEmailNotification                  RunbookActionSlug = "email_notification"
	RunbookActionSlugExportRetrospective                RunbookActionSlug = "export_retrospective"
	RunbookActionSlugFreeformText                       RunbookActionSlug = "freeform_text"
	RunbookActionSlugIncidentChannelGif                 RunbookActionSlug = "incident_channel_gif"
	RunbookActionSlugIncidentUpdate                     RunbookActionSlug = "incident_update"
	RunbookActionSlugNotifyChannel                      RunbookActionSlug = "notify_channel"
	RunbookActionSlugNotifyChannelCustomMessage         RunbookActionSlug = "notify_channel_custom_message"
	RunbookActionSlugNotifyIncidentChannelCustomMessage RunbookActionSlug = "notify_incident_channel_custom_message"
	RunbookActionSlugScript                             RunbookActionSlug = "script"
	RunbookActionSlugSendWebhook                        RunbookActionSlug = "send_webhook"
	RunbookActionSlugSetLinkedAlertsStatus              RunbookActionSlug = "set_linked_alerts_status"
	RunbookActionSlugUpdateStatuspage                   RunbookActionSlug = "update_statuspage"
	RunbookActionSlugVictorOpsCreateNewIncident         RunbookActionSlug = "victorops_create_new_incident"
)

type RunbookActionsResponse struct {
	Actions []RunbookAction `json:"data"`
}

// RunbookResponse is the payload for retrieving a service
// URL: GET https://api.firehydrant.io/v1/runbooks/{id}
type RunbookAction struct {
	ID          string       `json:"id"`
	Integration *Integration `json:"integration"`
	Name        string       `json:"name"`
	Slug        string       `json:"slug"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Integration struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
}

type RunbookActionsQuery struct {
	Type  string `url:"type,omitempty"`
	Items uint   `url:"per_page,omitempty"`
}

// RunbooksClient is an interface for interacting with runbooks on FireHydrant
type RunbookActionsClient interface {
	Get(ctx context.Context, runbookType string, integrationSlug string, actionSlug string) (*RunbookAction, error)
}

// RESTRunbooksClient implements the RunbooksClient interface
type RESTRunbookActionsClient struct {
	client *APIClient
}

var _ RunbookActionsClient = &RESTRunbookActionsClient{}

func (c *RESTRunbookActionsClient) restClient() *sling.Sling {
	return c.client.client()
}

// Get returns a runbook action from the FireHydrant API
func (c *RESTRunbookActionsClient) Get(ctx context.Context, runbookType string, integrationSlug string, actionSlug string) (*RunbookAction, error) {
	runbookActionResponse := &RunbookActionsResponse{}
	apiError := &APIError{}
	query := RunbookActionsQuery{Type: runbookType, Items: 100}
	response, err := c.restClient().Get("runbooks/actions").QueryStruct(query).Receive(runbookActionResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get runbook")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	for _, action := range runbookActionResponse.Actions {
		if action.Slug == actionSlug && action.Integration.Slug == integrationSlug {
			return &action, nil
		}
	}

	return nil, fmt.Errorf("could not find runbook action")
}
