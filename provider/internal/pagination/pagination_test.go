package pagination

import (
	"context"
	"slices"
	"testing"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/firehydrant/terraform-provider-firehydrant/provider/internal/ptr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type testRequest struct {
	Page *int
}

type testResponse struct {
	Items      []int
	Pagination *components.NullablePaginationEntity
}

func (r *testResponse) GetData() []int {
	return r.Items
}

func (r *testResponse) GetPagination() *components.NullablePaginationEntity {
	return r.Pagination
}

func TestPaginate_Success(t *testing.T) {
	// Assemble
	page1Items := []int{1, 2, 3, 4, 5}
	page2Items := []int{6, 7, 8, 9, 10}
	expectedItems := append(page1Items, page2Items...)

	setPageFunc := func(request *testRequest, page *int) {
		request.Page = page
	}
	requestFunc := func(ctx context.Context, client *firehydrant.APIClient, request *testRequest) (PaginateResponse[int], diag.Diagnostics) {
		if request.Page == nil || *request.Page == 1 {
			return &testResponse{Items: page1Items, Pagination: &components.NullablePaginationEntity{Next: ptr.Of(2)}}, nil
		} else if *request.Page == 2 {
			return &testResponse{Items: page2Items}, nil
		}
		return nil, diag.Errorf("unexpected page: %d", *request.Page)
	}

	// Act
	result, err := Paginate(context.Background(), PaginateRequestOptions[testRequest, int]{
		Client:             nil,
		Request:            &testRequest{},
		SetRequestPageFunc: setPageFunc,
		GetPageFunc:        requestFunc,
	})

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if len(result) != len(page1Items)+len(page2Items) {
		t.Errorf("expected %d items, got %d", len(page1Items)+len(page2Items), len(result))
	}
	for _, item := range result {
		if !slices.Contains(expectedItems, item) {
			t.Errorf("unexpected item: %d", item)
		}
	}
}

func TestPaginate_NilResponse(t *testing.T) {
	// Assemble
	requestFunc := func(ctx context.Context, client *firehydrant.APIClient, request *testRequest) (PaginateResponse[int], diag.Diagnostics) {
		return nil, nil
	}

	// Act
	result, err := Paginate(context.Background(), PaginateRequestOptions[testRequest, int]{
		Client:      nil,
		Request:     &testRequest{},
		GetPageFunc: requestFunc,
	})

	// Assert
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if len(result) != 0 {
		t.Errorf("expected 0 items, got %d", len(result))
	}
}

func TestPaginate_Error(t *testing.T) {
	// Assemble
	requestFunc := func(ctx context.Context, client *firehydrant.APIClient, request *testRequest) (PaginateResponse[int], diag.Diagnostics) {
		return nil, diag.Errorf("test error")
	}

	// Act
	result, err := Paginate(context.Background(), PaginateRequestOptions[testRequest, int]{
		Client:      nil,
		Request:     &testRequest{},
		GetPageFunc: requestFunc,
	})

	// Assert
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if len(result) != 0 {
		t.Errorf("expected 0 items, got %d", len(result))
	}
}

func TestPaginateRequest_ErrorOnSecondPage(t *testing.T) {
	// Assemble
	setPageFunc := func(request *testRequest, page *int) {
		request.Page = page
	}
	requestFunc := func(ctx context.Context, client *firehydrant.APIClient, request *testRequest) (PaginateResponse[int], diag.Diagnostics) {
		if request.Page == nil || *request.Page == 1 {
			return &testResponse{Items: []int{1, 2}, Pagination: &components.NullablePaginationEntity{Next: ptr.Of(2)}}, nil
		} else if *request.Page == 2 {
			return nil, diag.Errorf("test error")
		}
		return nil, diag.Errorf("unexpected page: %d", *request.Page)
	}

	// Act
	result, err := Paginate(context.Background(), PaginateRequestOptions[testRequest, int]{
		Client:             nil,
		Request:            &testRequest{},
		SetRequestPageFunc: setPageFunc,
		GetPageFunc:        requestFunc,
	})

	// Assert
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if len(result) != 0 {
		t.Errorf("expected 0 items, got %d", len(result))
	}
}
