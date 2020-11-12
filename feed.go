package main

import (
	"github.com/imonke/monkebase"
	"github.com/imonke/monketype"

	"net/http"
)

func feedAll(request *http.Request) (code int, r_map map[string]interface{}, err error) {
	var parsed map[string]interface{} = request.Context().Value("query").(map[string]interface{})
	var size int = parsed["size"].(int)
	var before string = parsed["before"].(string)
	var content []monketype.Content
	if content, size, err = monkebase.ReadManyContent(before, size); err != nil {
		return
	}

	code = 200
	r_map = map[string]interface{}{
		"content": content,
		"size":    map[string]int{"content": size},
	}
	return
}
