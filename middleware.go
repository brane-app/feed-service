package main

import (
	"context"
	"net/http"
	"strconv"
)

var (
	limits map[string]int = map[string]int{
		"size": 200,
	}
)

func defaults() (them map[string]int) {
	them = map[string]int{
		"size":   50,
		"offset": 0,
	}

	return
}

func contextQueryStrings(request *http.Request) (modified *http.Request, ok bool, code int, r_map map[string]interface{}, err error) {
	var parsed map[string]int = defaults()
	var limit int

	var key string
	var value []string
	for key, value = range request.URL.Query() {
		if len(value) == 0 || value[0] == "" {
			continue
		}

		if parsed[key], err = strconv.Atoi(value[0]); err != nil {
			err = nil
			code = 400
			return
		}

		if limit, ok = limits[key]; !ok {
			continue
		}

		if parsed[key] > limit {
			parsed[key] = limit
		}
	}

	ok = true
	modified = request.WithContext(context.WithValue(request.Context(), "parsed_query", parsed))
	return
}
