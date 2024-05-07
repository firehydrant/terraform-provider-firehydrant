package firehydrant

import (
	"context"
	"fmt"
	"time"

	"github.com/dghubble/sling"
	"github.com/pkg/errors"
)

type OnCallSchedules interface {
	Get(ctx context.Context, teamID, id string) (*OnCallScheduleResponse, error)
	Create(ctx context.Context, teamID string, createReq CreateOnCallScheduleRequest) (*OnCallScheduleResponse, error)
	Update(ctx context.Context, teamID, id string, updateReq UpdateOnCallScheduleRequest) (*OnCallScheduleResponse, error)
	Delete(ctx context.Context, teamID, id string) error
}

type RESTOnCallSchedulesClient struct {
	client *APIClient
}

var _ OnCallSchedules = &RESTOnCallSchedulesClient{}

type OnCallScheduleStrategy struct {
	Type          string `json:"type"`
	HandoffTime   string `json:"handoff_time,omitempty"`
	HandoffDay    string `json:"handoff_day,omitempty"`
	ShiftDuration string `json:"shift_duration,omitempty"`
}

type OnCallScheduleMember struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type OnCallScheduleResponse struct {
	ID           string                      `json:"id"`
	Name         string                      `json:"name"`
	Description  string                      `json:"description"`
	TimeZone     string                      `json:"time_zone"`
	Strategy     OnCallScheduleStrategy      `json:"strategy"`
	Members      []OnCallScheduleMember      `json:"members"`
	Restrictions []OnCallScheduleRestriction `json:"restrictions"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateOnCallScheduleRequest struct {
	Name         string                      `json:"name"`
	Description  string                      `json:"description"`
	TimeZone     string                      `json:"time_zone"`
	Strategy     OnCallScheduleStrategy      `json:"strategy"`
	Restrictions []OnCallScheduleRestriction `json:"restrictions"`
	MemberIDs    []string                    `json:"member_ids"`

	// StartTime is only required for `custom` strategy.
	// ISO8601 / Go RFC3339 format.
	StartTime string `json:"start_time,omitempty"`
}

type UpdateOnCallScheduleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	MemberIDs   []string `json:"member_ids,omitempty"`
	EffectiveAt string   `json:"effective_at,omitempty"`
}

type OnCallScheduleRestriction struct {
	StartDay  string `json:"start_day"`
	StartTime string `json:"start_time"`
	EndDay    string `json:"end_day"`
	EndTime   string `json:"end_time"`
}

func (c *RESTOnCallSchedulesClient) restClient() *sling.Sling {
	return c.client.client()
}

func (c *RESTOnCallSchedulesClient) Create(ctx context.Context, teamID string, createReq CreateOnCallScheduleRequest) (*OnCallScheduleResponse, error) {
	onCallScheduleResponse := &OnCallScheduleResponse{}
	apiError := &APIError{}

	response, err := c.restClient().Post(fmt.Sprintf("teams/%s/on_call_schedules", teamID)).BodyJSON(createReq).Receive(onCallScheduleResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not create on-call schedule")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return onCallScheduleResponse, nil
}

func (c *RESTOnCallSchedulesClient) Get(ctx context.Context, teamID, id string) (*OnCallScheduleResponse, error) {
	onCallScheduleResponse := &OnCallScheduleResponse{}
	apiError := &APIError{}

	response, err := c.restClient().Get(fmt.Sprintf("teams/%s/on_call_schedules/%s", teamID, id)).Receive(onCallScheduleResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not get on-call schedule")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return onCallScheduleResponse, nil
}

func (c *RESTOnCallSchedulesClient) Update(ctx context.Context, teamID, id string, updateReq UpdateOnCallScheduleRequest) (*OnCallScheduleResponse, error) {
	onCallScheduleResponse := &OnCallScheduleResponse{}
	apiError := &APIError{}

	if updateReq.EffectiveAt == "" {
		updateReq.EffectiveAt = time.Now().Format(time.RFC3339)
	}

	response, err := c.restClient().Patch(fmt.Sprintf("teams/%s/on_call_schedules/%s", teamID, id)).BodyJSON(updateReq).Receive(onCallScheduleResponse, apiError)
	if err != nil {
		return nil, errors.Wrap(err, "could not update on-call schedule")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return nil, err
	}

	return onCallScheduleResponse, nil
}

func (c *RESTOnCallSchedulesClient) Delete(ctx context.Context, teamID, id string) error {
	apiError := &APIError{}

	response, err := c.restClient().Delete(fmt.Sprintf("teams/%s/on_call_schedules/%s", teamID, id)).Receive(nil, apiError)
	if err != nil {
		return errors.Wrap(err, "could not delete on-call schedule")
	}

	err = checkResponseStatusCode(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
