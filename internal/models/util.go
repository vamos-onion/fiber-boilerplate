package models

import (
	"strconv"

	"fiber-boilerplate/internal/pkg/util"

	"gopkg.in/guregu/null.v4"
)

// ListOptions :
type ListOptions struct {
	Sorting    SortingBlock
	Pagination PaginationBlock
}

// SortingBlock :
type SortingBlock struct {
	Provided bool
	Orders   []string
}

// PaginationBlock :
type PaginationBlock struct {
	Provided bool
	Limit    int
	Offset   int
}

func MakeListOptions(sorting *SortingBlock, pagination *PaginationBlock) (options *ListOptions) {
	if options == nil {
		options = &ListOptions{}
	}
	if sorting != nil && sorting.Provided {
		options.Sorting = *sorting
	}
	if pagination != nil && pagination.Provided {
		options.Pagination = *pagination
	}
	return
}

// Parameterize :
func (l *ListOptions) Parameterize() (option null.String) {
	orderBy := func() null.String {
		if l.Sorting.Provided {
			return l.Sorting.Parameterize()
		}
		return null.String{}
	}()
	pagination := func() null.String {
		if l.Pagination.Provided {
			return l.Pagination.Parameterize()
		}
		return null.String{}
	}()

	if !orderBy.Valid && !pagination.Valid {
		return
	}

	option.String = util.String.Join(" ", orderBy.String, pagination.String)
	option.Valid = true

	return
}

// Parameterize :
func (s *SortingBlock) Parameterize() (orderBy null.String) {
	// Provided : true, Orders : ["osType", "ASC", "version", "DESC"]
	if s.Provided {
		for i := range s.Orders {
			if i == 0 {
				orderBy.String = util.String.Words("ORDER BY", s.Orders[i])
				continue
			}
			if i%2 == 0 {
				orderBy.String = util.String.Join(",", orderBy.String, s.Orders[i])
			} else {
				orderBy.String += " " + s.Orders[i]
			}
		}
		orderBy.Valid = true
	}

	return
}

// Parameterize :
func (p *PaginationBlock) Parameterize() (pagination null.String) {
	if p.Provided {
		pagination = null.StringFrom(util.String.Words("LIMIT", strconv.Itoa(p.Limit), "OFFSET", strconv.Itoa(p.Offset)))
		pagination.Valid = true
	}

	return

}

// NullableInt :
func NullableInt(raw null.Int) (ptr *int) {
	if raw.Valid {
		data := int(raw.Int64)
		ptr = &data
	}

	return
}

// NullableInt64 :
func NullableInt64(raw null.Int) (ptr *int64) {
	if raw.Valid {
		data := int64(raw.Int64)
		ptr = &data
	}

	return
}

// NullableFloat :
func NullableFloat(raw null.Float) (ptr *float32) {
	if raw.Valid {
		data := float32(raw.Float64)
		ptr = &data
	}

	return
}

// NullableFloat64 :
func NullableFloat64(raw null.Float) (ptr *float64) {
	if raw.Valid {
		data := raw.Float64
		ptr = &data
	}

	return
}

// NullableBool :
func NullableBool(raw null.Bool) (ptr *bool) {
	if raw.Valid {
		ptr = &raw.Bool
	}

	return
}

// NullableTS :
func NullableTS(raw null.Time) (ptr *int64) {
	if raw.Valid {
		data := util.Time.UnixMilli(raw.Time)
		ptr = &data
	}

	return
}

// NullableTM :
func NullableTM(raw null.Time) (ptr *string) {
	if raw.Valid {
		data := util.Time.Sprintf("%0R", raw.Time)
		ptr = &data
	}

	return
}

// NullableString :
func NullableString(raw null.String) (ptr *string) {
	if raw.Valid {
		ptr = &raw.String
	}

	return
}

// GenderToNullString :
func GenderToNullString(gender string) null.String {
	if gender == "" {
		return null.StringFrom(string(NullEnumGender{}.EnumGender))
	}

	return null.StringFrom(string(NullEnumGender{
		EnumGender: EnumGender(gender),
		Valid:      true,
	}.EnumGender))
}
