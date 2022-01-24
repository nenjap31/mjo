package queryparams

import (
	"net/url"
	"strconv"
)

type CleanQueryParams struct {
	PerPage     int
	Offset      int
	QueryParams url.Values
}

func QueryParamsCleaner(queryParams url.Values) CleanQueryParams {
	perPage, err := strconv.Atoi(queryParams.Get("per_page"))
	if err != nil {
		perPage = 10
	}
	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil {
		page = 1
	}
	offset := (page - 1) * perPage
	queryParams.Del("per_page")
	queryParams.Del("page")
	return CleanQueryParams{PerPage: perPage, Offset: offset, QueryParams: queryParams}
}
