package firehydrant

import (
	"context"
	"fmt"

	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type IngestURLParams struct {
	EscalationPolicyID string
	OnCallScheduleID   string
	TeamID             string
	UserID             string
}

// IngestURLClient is an interface for interacting with ingest URLs
type IngestURLClient interface {
	Get(ctx context.Context, params IngestURLParams) (*IngestURLResponse, error)
}

// RESTIngestURLClient implements the IngestURLClient interface
type RESTIngestURLClient struct {
	client *APIClient
}

var _ IngestURLClient = &RESTIngestURLClient{}

func (c *RESTIngestURLClient) restClient() *sling.Sling {
	return c.client.client()
}

// Get retrieves an ingest URL from FireHydrant.  See below for query params.
func (c *RESTIngestURLClient) Get(ctx context.Context, params IngestURLParams) (*IngestURLResponse, error) {
	ingestURL := &IngestURLResponse{}
	apiError := &APIError{}

	// The query here should be one of:
	// - nothing (to get the default URL),
	// - a user_id only,
	// - a team_id only,
	// - an escalation_policy_id with corresponding team_id,
	// - a schedule_id with corresponding team_id
	if params.TeamID == "" && params.EscalationPolicyID != "" {
		return nil, fmt.Errorf("missing team_id for escalation_policy_id %s", params.EscalationPolicyID)
	}
	if params.TeamID == "" && params.OnCallScheduleID != "" {
		return nil, fmt.Errorf("missing team_id for on_call_schedule_id %s", params.OnCallScheduleID)
	}

	query := ""
	if params.UserID != "" {
		query = fmt.Sprintf("?user_id=%s", params.UserID)
	} else if params.TeamID != "" {
		query = fmt.Sprintf("?team_id=%s", params.TeamID)
		if params.EscalationPolicyID != "" {
			query += fmt.Sprintf("&escalation_policy_id=%s", params.EscalationPolicyID)
		} else if params.OnCallScheduleID != "" {
			query += fmt.Sprintf("&on_call_schedule_id=%s", params.OnCallScheduleID)
		}
	}

	response, err := c.restClient().Get("signals/ingest_url"+query).Receive(ingestURL, apiError)
	if err != nil {
		return nil, fmt.Errorf("could not get ingest url: %w", err)
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	if ingestURL.URL == "" {
		return nil, fmt.Errorf("no ingest URL found with options %#v", params)
	}

	tflog.Info(ctx, "found ingest URL", map[string]interface{}{
		"url": ingestURL.URL,
	})

	return ingestURL, nil
}
