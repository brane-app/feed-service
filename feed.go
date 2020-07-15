package main

import (
	"github.com/imonke/monkebase"
	"github.com/imonke/monketype"

	"net/http"
)

func feedAll(request *http.Request) (code int, r_map map[string]interface{}, err error) {
	var parsed map[string]int = request.Context().Value("parsed_query").(map[string]int)

	var content []monketype.Content
	var size int
	if content, size, err = monkebase.ReadManyContent(parsed["offset"], parsed["size"]); err != nil {
		return
	}

	code = 200
	r_map = map[string]interface{}{
		"content": content,
		"size":    size,
	}
	return
}
