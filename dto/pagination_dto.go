package dto

import (
	"strings"
)

type Pagination struct {
	Page       int         `json:"page" query:"page"`
	Limit      int         `json:"limit" query:"limit"`
	SortBy     string      `json:"-" query:"sortBy"`
	SortOrder  string      `json:"-" query:"sortOrder"`
	Keyword    string      `json:"-" query:"keyword"`
	TotalData  int64       `json:"totalData"`
	TotalPages int         `json:"totalPages"`
	Rows       interface{} `json:"-"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	sortDir := strings.ToUpper(p.SortOrder)

	if p.SortBy == "" && (sortDir == "") {
		return ""
	}

	if sortDir == "" {
		p.SortOrder = "ASC"
	}

	return p.SortBy + " " + sortDir
}
