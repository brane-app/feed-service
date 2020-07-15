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

func contextQueryStrings(request *http.Request) (modified *http.Request, ok bool, code int, _ map[string]interface{}, err error) {
	var parsed map[string]int = defaults()
	var limited bool
	var limit, it int

	var key string
	var value []string
	for key, value = range request.URL.Query() {
		if len(value) == 0 || value[0] == "" {
			continue
		}

		limit, limited = limits[key]
		it, err = strconv.Atoi(value[0])
		switch {
		case err != nil:
			err = nil
			code = 400
			return
		case it > limit && limited:
			parsed[key] = limits[key]
		case it > 0 && (!limited || it <= limit):
			parsed[key] = it
		}
	}

	ok = true
	modified = request.WithContext(context.WithValue(request.Context(), "parsed_query", parsed))
	return
}
