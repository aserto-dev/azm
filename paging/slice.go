package paging

import (
	dsc2 "github.com/aserto-dev/go-directory/aserto/directory/common/v2"
	"github.com/aserto-dev/go-directory/pkg/derr"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type KeyComparer[T any] func([]string, T) bool
type KeyMapper[T any] func(T) []string

type Result[T any] struct {
	Items    []T
	PageInfo *dsc2.PaginationResponse
}

func PaginateSlice[T any](
	s []T,
	page *dsc2.PaginationRequest,
	keyCount int,
	cmp KeyComparer[T],
	mapper KeyMapper[T],
) (*Result[T], error) {
	result := &Result[T]{}

	start := 0
	if page != nil && page.Token != "" {
		cursor, err := DecodeCursor(page.Token)
		if err != nil {
			return result, err
		}

		if len(cursor.Keys) != keyCount {
			return result, derr.ErrInvalidCursor
		}

		_, start, _ = lo.FindIndexOf(s, func(o T) bool {
			return cmp(cursor.Keys, o)
		})

		if start == -1 {
			return result, derr.ErrInvalidCursor
		}
	}

	pageSize := lo.Min([]int32{page.Size, int32(len(s) - start)})
	end := start + int(pageSize)

	var next *string
	if end < len(s) {
		cursor := &Cursor{Keys: mapper(s[end])}
		n, err := cursor.Encode()
		if err != nil {
			return result, errors.Wrap(err, "failed to encode cursor")
		}
		next = &n
	}

	result.Items = s[start:end]
	result.PageInfo = &dsc2.PaginationResponse{
		NextToken: lo.FromPtr(next),
	}

	return result, nil
}
