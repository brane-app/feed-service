package main

import (
	"github.com/brane-app/tools-library/middleware"
	"github.com/gastrodon/groudon/v2"

	"os"
)

var (
	prefix = os.Getenv("PATH_PREFIX")

	routeAll = "^" + prefix + "/all/?$"
)

func register_handlers() {
	groudon.AddMiddleware("GET", routeAll, middleware.PaginationParams)

	groudon.AddHandler("GET", routeAll, feedAll)
}
