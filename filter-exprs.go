package finder

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
)

type FilterKind int

const (
	FilterInt FilterKind = iota
	FilterFloat
	FilterTime
)

type FilterExpr struct {
	Expr string
	Kind FilterKind
}

func ParseFilterExprs(
	filters string,
	exprs map[string]FilterExpr,
) ([]squirrel.Sqlizer, error) {
	filters = strings.TrimSpace(filters)
	if filters == "" {
		return nil, nil
	}

	var wheres []squirrel.Sqlizer
	for _, part := range strings.Split(filters, ",") {
		args := strings.SplitN(strings.TrimSpace(part), ":", 3)
		if len(args) != 3 {
			return nil, fmt.Errorf("invalid filter %q", part)
		}

		field, ok := exprs[args[0]]
		if !ok {
			return nil, fmt.Errorf("unknown filter field %q", args[0])
		}

		raws := strings.Split(args[2], "|")
		switch args[1] {
		case "bw":
			if len(raws) != 2 {
				return nil, fmt.Errorf("filter %q needs min|max", args[0])
			}
		case "eq", "nq":
		case "gt", "gte", "lt", "lte":
			if len(raws) != 1 {
				return nil, fmt.Errorf("filter %q takes one value", args[0])
			}
		default:
			return nil, fmt.Errorf("unknown filter operator %q", args[1])
		}

		values := make([]any, 0, len(raws))
		for _, raw := range raws {
			value, err := parseFilterValue(field.Kind, strings.TrimSpace(raw))
			if err != nil {
				return nil, fmt.Errorf("filter %q: %w", args[0], err)
			}
			values = append(values, value)
		}

		switch args[1] {
		case "eq":
			wheres = append(wheres, squirrel.Eq{field.Expr: values})
		case "nq":
			wheres = append(wheres, squirrel.NotEq{field.Expr: values})
		case "gt":
			wheres = append(wheres, squirrel.Gt{field.Expr: values[0]})
		case "gte":
			wheres = append(wheres, squirrel.GtOrEq{field.Expr: values[0]})
		case "lt":
			wheres = append(wheres, squirrel.Lt{field.Expr: values[0]})
		case "lte":
			wheres = append(wheres, squirrel.LtOrEq{field.Expr: values[0]})
		case "bw":
			wheres = append(wheres, squirrel.And{
				squirrel.GtOrEq{field.Expr: values[0]},
				squirrel.LtOrEq{field.Expr: values[1]},
			})
		}
	}

	return wheres, nil
}

func parseFilterValue(kind FilterKind, raw string) (any, error) {
	switch kind {
	case FilterInt:
		value, err := strconv.Atoi(raw)
		if err != nil {
			return nil, fmt.Errorf("%q is not a valid integer", raw)
		}
		return value, nil
	case FilterFloat:
		value, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return nil, fmt.Errorf("%q is not a valid number", raw)
		}
		return value, nil
	default:
		if value, err := time.Parse(time.RFC3339, raw); err == nil {
			return value, nil
		}
		value, err := time.Parse(time.DateOnly, raw)
		if err != nil {
			return nil, fmt.Errorf("%q is not a valid date", raw)
		}
		return value, nil
	}
}
