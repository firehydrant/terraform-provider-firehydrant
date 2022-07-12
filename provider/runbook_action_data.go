package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceRunbookAction() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataFireHydrantRunbookAction,
		Schema: map[string]*schema.Schema{
			// Required
			"integration_slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(firehydrant.RunbookActionIntegrationSlugConfluenceCloud),
						string(firehydrant.RunbookActionIntegrationSlugFireHydrant),
						string(firehydrant.RunbookActionIntegrationSlugFireHydrantNunc),
						string(firehydrant.RunbookActionIntegrationSlugGiphy),
						string(firehydrant.RunbookActionIntegrationSlugGoogleDocs),
						string(firehydrant.RunbookActionIntegrationSlugGoogleMeet),
						string(firehydrant.RunbookActionIntegrationSlugJiraCloud),
						string(firehydrant.RunbookActionIntegrationSlugJiraServer),
						string(firehydrant.RunbookActionIntegrationSlugMicrosoftTeams),
						string(firehydrant.RunbookActionIntegrationSlugOpsgenie),
						string(firehydrant.RunbookActionIntegrationSlugPagerDuty),
						string(firehydrant.RunbookActionIntegrationSlugShortcut),
						string(firehydrant.RunbookActionIntegrationSlugSlack),
						string(firehydrant.RunbookActionIntegrationSlugStatuspage),
						string(firehydrant.RunbookActionIntegrationSlugVictorOps),
						string(firehydrant.RunbookActionIntegrationSlugWebex),
						string(firehydrant.RunbookActionIntegrationSlugZoom),
					},
					false,
				),
			},
			"slug": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						string(firehydrant.RunbookActionSlugAddServicesRelatedToFunctionality),
						string(firehydrant.RunbookActionSlugAddTaskList),
						string(firehydrant.RunbookActionSlugArchiveIncidentChannel),
						string(firehydrant.RunbookActionSlugAssignARole),
						string(firehydrant.RunbookActionSlugAssignATeam),
						string(firehydrant.RunbookActionSlugAttachARunbook),
						string(firehydrant.RunbookActionSlugCreateGoogleMeetLink),
						string(firehydrant.RunbookActionSlugCreateIncidentChannel),
						string(firehydrant.RunbookActionSlugCreateIncidentIssue),
						string(firehydrant.RunbookActionSlugCreateIncidentTicket),
						string(firehydrant.RunbookActionSlugCreateMeeting),
						string(firehydrant.RunbookActionSlugCreateNewOpsgenieIncident),
						string(firehydrant.RunbookActionSlugCreateNewPagerDutyIncident),
						string(firehydrant.RunbookActionSlugCreateNunc),
						string(firehydrant.RunbookActionSlugCreateStatuspage),
						string(firehydrant.RunbookActionSlugEmailNotification),
						string(firehydrant.RunbookActionSlugExportRetrospective),
						string(firehydrant.RunbookActionSlugFreeformText),
						string(firehydrant.RunbookActionSlugIncidentChannelGif),
						string(firehydrant.RunbookActionSlugIncidentUpdate),
						string(firehydrant.RunbookActionSlugNotifyChannel),
						string(firehydrant.RunbookActionSlugNotifyChannelCustomMessage),
						string(firehydrant.RunbookActionSlugNotifyIncidentChannelCustomMessage),
						string(firehydrant.RunbookActionSlugScript),
						string(firehydrant.RunbookActionSlugSendWebhook),
						string(firehydrant.RunbookActionSlugSetLinkedAlertsStatus),
						string(firehydrant.RunbookActionSlugUpdateStatuspage),
						string(firehydrant.RunbookActionSlugVictorOpsCreateNewIncident),
					},
					false,
				),
			},
			"type": {
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

func dataFireHydrantRunbookAction(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	firehydrantAPIClient := m.(firehydrant.Client)

	// Get the runbook action
	runbookType := d.Get("type").(string)
	actionSlug := d.Get("slug").(string)
	integrationSlug := d.Get("integration_slug").(string)
	tflog.Debug(ctx, fmt.Sprintf("Read runbook action: %s:%s", integrationSlug, actionSlug), map[string]interface{}{
		"type":             runbookType,
		"slug":             actionSlug,
		"integration_slug": integrationSlug,
	})
	runbookActionResponse, err := firehydrantAPIClient.RunbookActions().Get(ctx, runbookType, integrationSlug, actionSlug)
	if err != nil {
		return diag.Errorf("Error reading runbook action %s:%s: %v", integrationSlug, actionSlug, err)
	}

	// Update the attributes in state to the values we got from the API
	attributes := map[string]string{
		"name": runbookActionResponse.Name,
		"slug": runbookActionResponse.Slug,
	}

	if runbookActionResponse.Integration != nil {
		attributes["integration_slug"] = runbookActionResponse.Integration.Slug
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for runbook action %s:%s: %v", key, integrationSlug, actionSlug, err)
		}
	}

	// Set the runbook action's ID in state
	d.SetId(runbookActionResponse.ID)

	return diag.Diagnostics{}
}
