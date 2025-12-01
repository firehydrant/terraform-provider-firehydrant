package firehydrant

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

type Rotations interface {
	Get(ctx context.Context, teamID, scheduleID, id string) (*RotationResponse, error)
	Create(ctx context.Context, teamID, scheduleID string, createReq CreateRotationRequest) (*RotationResponse, error)
	Update(ctx context.Context, teamID, scheduleID, id string, updateReq UpdateRotationRequest) (*RotationResponse, error)
	Delete(ctx context.Context, teamID, scheduleID, id string) error
}

type RESTRotationsClient struct {
	client *APIClient
}

var _ Rotations = &RESTRotationsClient{}

type RotationStrategy struct {
	Type          string `json:"type"`
	HandoffTime   string `json:"handoff_time,omitempty"`
	HandoffDay    string `json:"handoff_day,omitempty"`
	ShiftDuration string `json:"shift_duration,omitempty"`
}

type RotationMember struct {
	UserID *string `json:"user_id,omitempty"` // Used for requests, and populated from "id" in responses
	Name   *string `json:"name,omitempty"`    // Only in responses
}

// UnmarshalJSON custom unmarshaler to handle both "id" (response) and "user_id" (request)
// Requests to the rotations api accept a user_id but responses use the generic SuccicentEntity
// which has an id field, so we use a custom unmarshaller to handle both cases and keep terraform state in sync
func (m *RotationMember) UnmarshalJSON(data []byte) error {
	// Define a temporary struct to handle both field names
	var temp struct {
		UserID *string `json:"user_id"`
		ID     *string `json:"id"`
		Name   *string `json:"name"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.ID != nil {
		m.UserID = temp.ID
	}

	m.Name = temp.Name
	return nil
}

type RotationResponse struct {
	ID                              string                `json:"id"`
	Name                            string                `json:"name"`
	Description                     string                `json:"description"`
	TimeZone                        string                `json:"time_zone"`
	Color                           string                `json:"color"`
	SlackUserGroupID                string                `json:"slack_user_group_id"`
	EnableSlackChannelNotifications bool                  `json:"enable_slack_channel_notifications,omitempty"`
	PreventShiftDeletion            bool                  `json:"prevent_shift_deletion,omitempty"`
	CoverageGapNotificationInterval string                `json:"coverage_gap_notification_interval,omitempty"`
	StartTime                       string                `json:"start_time,omitempty"`
	Strategy                        RotationStrategy      `json:"strategy"`
	Members                         []RotationMember      `json:"members"`
	Restrictions                    []RotationRestriction `json:"restrictions"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RotationsQuery is the query used to search for on call schedules
type RotationsQuery struct {
	TeamID string `url:"team_id,omitempty"` // optional
	Query  string `url:"query,omitempty"`   // optional
	Page   int    `url:"page,omitempty"`
}

type CreateRotationRequest struct {
	Name                            string                `json:"name"`
	Description                     string                `json:"description"`
	TimeZone                        string                `json:"time_zone"`
	Color                           string                `json:"color,omitempty"`
	SlackUserGroupID                string                `json:"slack_user_group_id,omitempty"`
	EnableSlackChannelNotifications bool                  `json:"enable_slack_channel_notifications,omitempty"`
	PreventShiftDeletion            bool                  `json:"prevent_shift_deletion,omitempty"`
	CoverageGapNotificationInterval string                `json:"coverage_gap_notification_interval,omitempty"`
	Strategy                        RotationStrategy      `json:"strategy"`
	Restrictions                    []RotationRestriction `json:"restrictions"`
	Members                         []RotationMember      `json:"members"`

	// StartTime is only required for `custom` strategy.
	// ISO8601 / Go RFC3339 format.
	StartTime string `json:"start_time,omitempty"`
}

type UpdateRotationRequest struct {
	Name                            string                `json:"name"`
	Description                     string                `json:"description"`
	Members                         []RotationMember      `json:"members,omitempty"`
	EffectiveAt                     string                `json:"effective_at,omitempty"`
	Color                           string                `json:"color,omitempty"`
	SlackUserGroupID                string                `json:"slack_user_group_id,omitempty"`
	EnableSlackChannelNotifications bool                  `json:"enable_slack_channel_notifications,omitempty"`
	PreventShiftDeletion            bool                  `json:"prevent_shift_deletion,omitempty"`
	CoverageGapNotificationInterval string                `json:"coverage_gap_notification_interval,omitempty"`
	Strategy                        *RotationStrategy     `json:"strategy,omitempty"`
	Restrictions                    []RotationRestriction `json:"restrictions"`
}

type RotationRestriction struct {
	StartDay  string `json:"start_day"`
	StartTime string `json:"start_time"`
	EndDay    string `json:"end_day"`
	EndTime   string `json:"end_time"`
}

func (c *RESTRotationsClient) restClient() *sling.Sling {
	return c.client.client()
}

func (c *RESTRotationsClient) Create(ctx context.Context, teamID string, scheduleID string, createReq CreateRotationRequest) (*RotationResponse, error) {
	rotationResponse := &RotationResponse{}
	apiError := &APIError{}

	response, err := c.restClient().Post(fmt.Sprintf("teams/%s/on_call_schedules/%s/rotations", teamID, scheduleID)).BodyJSON(createReq).Receive(rotationResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create rotation")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return rotationResponse, nil
}

func (c *RESTRotationsClient) Get(ctx context.Context, teamID, scheduleID, id string) (*RotationResponse, error) {
	rotationResponse := &RotationResponse{}
	apiError := &APIError{}

	response, err := c.restClient().Get(fmt.Sprintf("teams/%s/on_call_schedules/%s/rotations/%s", teamID, scheduleID, id)).Receive(rotationResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get on-call schedule")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return rotationResponse, nil
}

func (c *RESTRotationsClient) Update(ctx context.Context, teamID, scheduleID, id string, updateReq UpdateRotationRequest) (*RotationResponse, error) {
	rotationResponse := &RotationResponse{}
	apiError := &APIError{}

	if updateReq.EffectiveAt == "" {
		updateReq.EffectiveAt = time.Now().Format(time.RFC3339)
	}

	response, err := c.restClient().Patch(fmt.Sprintf("teams/%s/on_call_schedules/%s/rotations/%s", teamID, scheduleID, id)).BodyJSON(updateReq).Receive(rotationResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update on-call schedule")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return rotationResponse, nil
}

func (c *RESTRotationsClient) Delete(ctx context.Context, teamID, scheduleID, id string) error {
	apiError := &APIError{}

	response, err := c.restClient().Delete(fmt.Sprintf("teams/%s/on_call_schedules/%s/rotations/%s", teamID, scheduleID, id)).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete on-call schedule")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
