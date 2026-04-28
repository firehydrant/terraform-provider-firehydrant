package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
)

// maxPages is a safety cap to prevent infinite loops if the API returns a
// perpetually non-nil Next pointer.
const maxPages = 1000

// fetchAllPages iterates through all pages of a paginated SDK list endpoint and
// returns the concatenated items from every page.
//
// Parameters:
//   - fetch: calls the SDK, returns (items, pagination, error) for the current request.
//   - setPage: mutates the request's Page field to the given page number.
//   - req: the initial request value; Page will be set to 1 on the first call.
func fetchAllPages[Req any, Item any](
	ctx context.Context,
	fetch func(ctx context.Context, req Req) ([]Item, *components.NullablePaginationEntity, error),
	setPage func(req *Req, page int),
	req Req,
) ([]Item, error) {
	var all []Item

	for page := 1; page <= maxPages; page++ {
		setPage(&req, page)

		items, pagination, err := fetch(ctx, req)
		if err != nil {
			return nil, err
		}

		all = append(all, items...)

		if pagination == nil || pagination.GetNext() == nil || *pagination.GetNext() == 0 {
			break
		}

		if page == maxPages {
			return nil, fmt.Errorf("exceeded maximum page limit of %d while paginating", maxPages)
		}
	}

	return all, nil
}
