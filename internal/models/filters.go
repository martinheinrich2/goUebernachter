package models

import (
	"github.com/martinheinrich2/goUebernachter/internal/validator"
	"strings"
)

// The Filters struct holds the supported values for the filters.
type Filters struct {
	Page          int
	PageSize      int
	Sort          string
	SortDirection string
	SortSafeList  []string
}

// Check that the client-provided Sort field matches one of the entries in our safelist and if it does,
// extract the column name from the Sort field by stripping the leading hyphen character (if one exists).
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

// Return the sort direction ("ASC" or "DESC") depending on the prefix character of the Sort field.
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	} else if strings.HasPrefix(f.Sort, "") {
		return "ASC"
	}
	return "DESC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

// ValidateFilters checks if the filters are within limits and safe to use.
func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(validator.PermittedValue(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

// Metadata is a  struct for holding the pagination metadata.
type Metadata struct {
	CurrentPage   int
	PreviousPage  int
	NextPage      int
	PlusFive      int
	PlusTen       int
	MinusFive     int
	MinusTen      int
	PageSize      int
	FirstPage     int
	LastPage      int
	TotalRecords  int
	SortDirection string
}

// The calculateMetadata() function calculates the appropriate pagination metadata values
// given the total number of records, current page, and page size values. The last page value
// is calculated by dividing two int values, resulting in an int value and dropping the modulus.
func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		// Return an empty metadata struct if there are no records.
		return Metadata{}
	}
	firstPage := 1
	lastPage := (totalRecords + pageSize - 1) / pageSize

	return Metadata{
		CurrentPage:  page,
		PreviousPage: page - 1,
		NextPage:     page + 1,
		PlusFive:     page + 5,
		PlusTen:      page + 10,
		MinusFive:    page - 5,
		MinusTen:     page - 10,
		PageSize:     pageSize,
		FirstPage:    firstPage,
		LastPage:     lastPage,
		TotalRecords: totalRecords,
	}

}
