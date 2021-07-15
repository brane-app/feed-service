package main

import (
	"github.com/brane-app/database-library"
	"github.com/brane-app/types-library"

	"net/http"
)

func feedAll(request *http.Request) (code int, r_map map[string]interface{}, err error) {
	var parsed map[string]interface{} = request.Context().Value("query").(map[string]interface{})
	var size int = parsed["size"].(int)
	var before string = parsed["before"].(string)
	var content []types.Content
	if content, size, err = database.ReadManyContent(before, size); err != nil {
		return
	}

	code = 200
	r_map = map[string]interface{}{
		"content": content,
		"size":    map[string]int{"content": size},
	}
	return
}
