package httpx

import (
	"net/http"
	"strconv"
)

type Page struct {
	Limit  int
	Offset int
}

func ParsePage(r *http.Request, defLimit int) Page {
	limit := defLimit
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}
	return Page{Limit: limit, Offset: offset}
}
