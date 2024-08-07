package bunpaginate

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"

	libtime "github.com/formancehq/stack/libs/go-libs/time"
	"github.com/uptrace/bun"
)

func UsingColumn[FILTERS any, ENTITY any](ctx context.Context,
	sb *bun.SelectQuery,
	query ColumnPaginatedQuery[FILTERS]) (*Cursor[ENTITY], error) {
	ret := make([]ENTITY, 0)

	sb = sb.Model(&ret)

	sb = sb.Limit(int(query.PageSize) + 1) // Fetch one additional item to find the next token
	order := query.Order
	if query.Reverse {
		order = order.Reverse()
	}
	sb = sb.OrderExpr(fmt.Sprintf("%s %s", query.Column, order))

	if query.PaginationID != nil {
		if query.Reverse {
			switch query.Order {
			case OrderAsc:
				sb = sb.Where(fmt.Sprintf("%s < ?", query.Column), query.PaginationID)
			case OrderDesc:
				sb = sb.Where(fmt.Sprintf("%s > ?", query.Column), query.PaginationID)
			}
		} else {
			switch query.Order {
			case OrderAsc:
				sb = sb.Where(fmt.Sprintf("%s >= ?", query.Column), query.PaginationID)
			case OrderDesc:
				sb = sb.Where(fmt.Sprintf("%s <= ?", query.Column), query.PaginationID)
			}
		}
	}

	if err := sb.Scan(ctx, &ret); err != nil {
		return nil, err
	}
	var (
		paginatedColumnIndex = 0
	)
	typeOfT := reflect.TypeOf(ret).Elem()
	for ; paginatedColumnIndex < typeOfT.NumField(); paginatedColumnIndex++ {
		field := typeOfT.Field(paginatedColumnIndex)
		tag := field.Tag.Get("bun")
		column := strings.Split(tag, ",")[0]
		if column == query.Column {
			break
		}
	}

	var (
		paginationIDs = make([]*big.Int, 0)
	)
	for _, t := range ret {
		rawPaginationID := reflect.ValueOf(t).
			Field(paginatedColumnIndex).
			Interface()
		var paginationID *big.Int
		switch rawPaginationID := rawPaginationID.(type) {
		case time.Time:
			paginationID = big.NewInt(rawPaginationID.UTC().UnixMicro())
		case libtime.Time:
			paginationID = big.NewInt(rawPaginationID.UTC().UnixMicro())
		case *BigInt:
			paginationID = (*big.Int)(rawPaginationID)
		case *big.Int:
			paginationID = rawPaginationID
		case int64:
			paginationID = big.NewInt(rawPaginationID)
		default:
			panic(fmt.Sprintf("invalid paginationID, type %T not handled", rawPaginationID))
		}
		if query.Bottom == nil {
			query.Bottom = paginationID
		}
		paginationIDs = append(paginationIDs, paginationID)
	}

	hasMore := len(ret) > int(query.PageSize)
	if hasMore {
		ret = ret[:len(ret)-1]
	}
	if query.Reverse {
		for i := 0; i < len(ret)/2; i++ {
			ret[i], ret[len(ret)-i-1] = ret[len(ret)-i-1], ret[i]
		}
	}

	var previous, next *ColumnPaginatedQuery[FILTERS]

	if query.Reverse {
		cp := query
		cp.Reverse = false
		next = &cp

		if hasMore {
			cp := query
			cp.PaginationID = paginationIDs[len(paginationIDs)-2]
			previous = &cp
		}
	} else {
		if hasMore {
			cp := query
			cp.PaginationID = paginationIDs[len(paginationIDs)-1]
			next = &cp
		}
		if query.PaginationID != nil {
			if (query.Order == OrderAsc && query.PaginationID.Cmp(query.Bottom) > 0) || (query.Order == OrderDesc && query.PaginationID.Cmp(query.Bottom) < 0) {
				cp := query
				cp.Reverse = true
				previous = &cp
			}
		}
	}

	return &Cursor[ENTITY]{
		PageSize: int(query.PageSize),
		HasMore:  next != nil,
		Previous: previous.EncodeAsCursor(),
		Next:     next.EncodeAsCursor(),
		Data:     ret,
	}, nil
}
