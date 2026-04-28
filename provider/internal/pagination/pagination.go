package pagination

import (
	"context"
	"time"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// PaginateRequestOptions is the options for the Paginate function
type PaginateRequestOptions[TRequest any, TEntity any] struct {
	// Client is the API client to use for the pagination
	Client *firehydrant.APIClient
	// Request is the request to use for the pagination
	Request *TRequest
	// SetRequestPageFunc is the function to use to set the page on the request
	SetRequestPageFunc func(request *TRequest, page *int)
	// GetPageFunc is the function to use to get the page from the API
	GetPageFunc func(ctx context.Context, client *firehydrant.APIClient, request *TRequest) (PaginateResponse[TEntity], diag.Diagnostics)
	// GetPageDelay is an optional duration to sleep between pages to avoid rate limits
	GetPageDelay time.Duration
}

// PaginateResponse is an interface that the response from the API must implement.
// The SDK will return this on every response that is paginated
type PaginateResponse[TEntity any] interface {
	GetData() []TEntity
	GetPagination() *components.NullablePaginationEntity
}

// Paginate is the function that will paginate the API response for SDK-based calls to the API
func Paginate[TRequest any, TEntity any](ctx context.Context, options PaginateRequestOptions[TRequest, TEntity]) ([]TEntity, diag.Diagnostics) {
	results := []TEntity{}
	page := toIntPointer(1)

	for page != nil {
		if options.GetPageDelay > 0 {
			time.Sleep(options.GetPageDelay)
		}
		if options.SetRequestPageFunc != nil {
			options.SetRequestPageFunc(options.Request, page)
		}
		response, err := options.GetPageFunc(ctx, options.Client, options.Request)
		if err != nil {
			return nil, err
		}

		if response == nil {
			return nil, diag.Errorf("unexpected respsonse")
		}

		results = append(results, response.GetData()...)

		if response.GetPagination() != nil {
			page = response.GetPagination().GetNext()
		} else {
			page = nil
		}
	}
	return results, nil

}

func toIntPointer(v int) *int {
	return &v
}
