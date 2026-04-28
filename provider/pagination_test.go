package provider

import (
	"context"
	"errors"
	"testing"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRequest struct {
	Page *int
}

func intPtr(i int) *int { return &i }

func makeFetcher(pages [][]string, failOnPage int) func(ctx context.Context, req fakeRequest) ([]string, *components.NullablePaginationEntity, error) {
	return func(ctx context.Context, req fakeRequest) ([]string, *components.NullablePaginationEntity, error) {
		page := 1
		if req.Page != nil {
			page = *req.Page
		}

		if failOnPage > 0 && page == failOnPage {
			return nil, nil, errors.New("simulated API error")
		}

		idx := page - 1
		if idx < 0 || idx >= len(pages) {
			return []string{}, &components.NullablePaginationEntity{}, nil
		}

		items := pages[idx]
		pagination := &components.NullablePaginationEntity{}
		if page < len(pages) {
			next := page + 1
			pagination.Next = &next
		}

		return items, pagination, nil
	}
}

func TestFetchAllPages_SinglePage(t *testing.T) {
	fetch := makeFetcher([][]string{{"a", "b", "c"}}, 0)
	req := fakeRequest{}

	got, err := fetchAllPages(context.Background(), fetch, func(r *fakeRequest, p int) { r.Page = intPtr(p) }, req)

	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, got)
}

func TestFetchAllPages_MultiplePages(t *testing.T) {
	fetch := makeFetcher([][]string{{"a", "b"}, {"c", "d"}, {"e"}}, 0)
	req := fakeRequest{}

	got, err := fetchAllPages(context.Background(), fetch, func(r *fakeRequest, p int) { r.Page = intPtr(p) }, req)

	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, got)
}

func TestFetchAllPages_EmptyResult(t *testing.T) {
	fetch := makeFetcher([][]string{{}}, 0)
	req := fakeRequest{}

	got, err := fetchAllPages(context.Background(), fetch, func(r *fakeRequest, p int) { r.Page = intPtr(p) }, req)

	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestFetchAllPages_ErrorOnFirstPage(t *testing.T) {
	fetch := makeFetcher([][]string{{"a"}, {"b"}}, 1)
	req := fakeRequest{}

	got, err := fetchAllPages(context.Background(), fetch, func(r *fakeRequest, p int) { r.Page = intPtr(p) }, req)

	require.Error(t, err)
	assert.Nil(t, got)
	assert.EqualError(t, err, "simulated API error")
}

func TestFetchAllPages_ErrorOnSubsequentPage(t *testing.T) {
	fetch := makeFetcher([][]string{{"a"}, {"b"}, {"c"}}, 2)
	req := fakeRequest{}

	got, err := fetchAllPages(context.Background(), fetch, func(r *fakeRequest, p int) { r.Page = intPtr(p) }, req)

	require.Error(t, err)
	assert.Nil(t, got)
	assert.EqualError(t, err, "simulated API error")
}

func TestFetchAllPages_NilPaginationStopsLoop(t *testing.T) {
	calls := 0
	fetch := func(ctx context.Context, req fakeRequest) ([]string, *components.NullablePaginationEntity, error) {
		calls++
		return []string{"x"}, nil, nil
	}
	req := fakeRequest{}

	got, err := fetchAllPages(context.Background(), fetch, func(r *fakeRequest, p int) { r.Page = intPtr(p) }, req)

	require.NoError(t, err)
	assert.Equal(t, 1, calls)
	assert.Equal(t, []string{"x"}, got)
}

func TestFetchAllPages_MaxPageSafetyCap(t *testing.T) {
	// Fetcher always reports a next page -- should hit the safety cap.
	infiniteNext := 999
	fetch := func(ctx context.Context, req fakeRequest) ([]string, *components.NullablePaginationEntity, error) {
		pagination := &components.NullablePaginationEntity{Next: &infiniteNext}
		return []string{"item"}, pagination, nil
	}
	req := fakeRequest{}

	got, err := fetchAllPages(context.Background(), fetch, func(r *fakeRequest, p int) { r.Page = intPtr(p) }, req)

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "exceeded maximum page limit")
}
