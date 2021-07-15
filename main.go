package main

import (
	"github.com/gastrodon/groudon"
	"git.gastrodon.io/imonke/monkebase"
	"git.gastrodon.io/imonke/monkelib/middleware"

	"log"
	"net/http"
	"os"
)

func main() {
	monkebase.Connect(os.Getenv("DATABASE_CONNECTION"))
	groudon.RegisterMiddleware(middleware.PaginationParams)
	groudon.RegisterHandler("GET", "^/all/?$", feedAll)
	http.Handle("/", http.HandlerFunc(groudon.Route))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
