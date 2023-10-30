// Package sql provides a generic interface around SQL (or SQL-like) databases.
package sql

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
)

type PaginationRequest struct {
	PerPage     int    `json:"per_page"`
	CurrentPage int    `json:"current_page"`
	Search      string `json:"search"`
	Sort        string `json:"sort"`
}

type PaginationMetaMessage struct {
	TotalItems  int64  `json:"total_items"`
	PerPage     int64  `json:"per_page"`
	CurrentPage int64  `json:"current_page"`
	TotalPage   int64  `json:"total_page"`
	NextUrl     string `json:"next_url"`
	PreviousUrl string `json:"previous_url"`
	FromItem    int64  `json:"from_item"`
	ToItem      int64  `json:"to_item"`
	Sort        string `json:"sort"`
}

// In expands slice values in args, returning the modified query string and a new arg list that can
// be executed by a database. The `query` should use the `?` bindVar. The return value uses the `?`
// bindVar.
func In(query string, args ...interface{}) (string, []interface{}, error) {
	return sqlx.In(query, args...)
}

// Paginate parse pagination meta data and return string for limit and order a query
func Paginate(request map[string]interface{}, pagination *PaginationMetaMessage, path string) string {
	// parse interface to custom model

	requestData := parsePaginationModel(request)

	totalItems := pagination.TotalItems
	page := requestData.CurrentPage
	if page == 0 {
		page = 1
	}

	pageSize := requestData.PerPage

	// Fill pagination object
	pagination.PerPage = int64(pageSize)
	pagination.CurrentPage = int64(page)
	pagination.TotalItems = totalItems
	pagination.TotalPage = int64(math.Ceil(float64(totalItems) / float64(pageSize)))
	pagination.FromItem = int64(pageSize * (page - 1))
	pagination.ToItem = int64(pageSize*page) - 1
	pagination.Sort = requestData.Sort

	prev, next := composePageUrl(request, path, pagination.CurrentPage, pagination.TotalPage)

	pagination.PreviousUrl = prev
	pagination.NextUrl = next

	offset := (page - 1) * pageSize
	return fmt.Sprintf("ORDER BY %s OFFSET %d LIMIT %d", pagination.Sort, offset, pageSize)
}

func composePageUrl(request map[string]interface{}, url string, currentPage, totalPage int64) (string, string) {
	i := 0
	for k, v := range request {
		if k == "page" {
			continue
		}

		if i == 0 {
			url += "?"
		} else {
			url += "&"
		}

		url += fmt.Sprintf("%s=%v", k, strings.Replace(fmt.Sprintf("%v", v), " ", ",", -1))
		i++
	}

	url = strings.Replace(url, "[", "", -1)
	url = strings.Replace(url, "]", "", -1)

	var prev, next string
	var page float64 = 1
	if request["page"] != nil {
		page = request["page"].(float64)
	}
	if currentPage-1 != 0 {
		prev = url + fmt.Sprintf("&page=%v", page-1)
	}

	if currentPage+1 <= totalPage {
		next = url + fmt.Sprintf("&page=%v", page+1)
	}

	return prev, next
}

func parsePaginationModel(request map[string]interface{}) *PaginationRequest {
	var perPage, currentPage int
	var search, sort string

	if request["per_page"] == nil {
		perPage = 5
	} else {
		perPage = int(request["per_page"].(float64))
	}

	if request["page"] == nil {
		currentPage = 1
	} else {
		currentPage = int(request["page"].(float64))
	}

	if request["search"] == nil {
		search = ""
	} else {
		search = request["search"].(string)
	}

	if request["sort"] == nil {
		sort = "id ASC"
	} else {
		checkSort := strings.Split(request["sort"].(string), ",")
		if len(checkSort) == 2 && (checkSort[1] == "asc" || checkSort[1] == "desc") {
			isMatch, err := regexp.MatchString(`^[a-zA-Z0-9_]+$`, checkSort[0])
			if isMatch && err == nil {
				sort = fmt.Sprintf("%s %s", checkSort[0], checkSort[1])
			}
		} else if len(checkSort) == 4 && (checkSort[1] == "asc" || checkSort[1] == "desc") && checkSort[2] == "nulls" && checkSort[3] == "last" {
			isMatch, err := regexp.MatchString(`^[a-zA-Z0-9_]+$`, checkSort[0])
			if isMatch && err == nil {
				sort = fmt.Sprintf("%s %s %s %s", checkSort[0], checkSort[1], checkSort[2], checkSort[3])
			}
		}
	}

	return &PaginationRequest{
		PerPage:     perPage,
		CurrentPage: currentPage,
		Search:      search,
		Sort:        sort,
	}
}
