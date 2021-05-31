package firehydrant

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetService(t *testing.T) {
	resp := &ServiceResponse{}
	testServiceID := "test-service-id"
	c, teardown, err := setupClient("/services/"+testServiceID, resp)
	require.NoError(t, err)
	defer teardown()

	res, err := c.Services().Get(context.TODO(), testServiceID)
	require.NoError(t, err, "error retrieving a service")
	assert.Equal(t, resp.ID, res.ID, "returned service did not match")
	assert.Equal(t, resp.Name, res.Name, "returned service did not match")
}

func TestCreateService(t *testing.T) {
	resp := &ServiceResponse{}
	c, teardown, err := setupClient("/services", resp,
		AssertRequestJSONBody(t, CreateServiceRequest{Name: "fake-service"}),
		AssertRequestMethod(t, "POST"),
	)

	require.NoError(t, err)
	defer teardown()

	_, err = c.Services().Create(context.TODO(), CreateServiceRequest{Name: "fake-service"})
	require.NoError(t, err, "error creating a service")
}
