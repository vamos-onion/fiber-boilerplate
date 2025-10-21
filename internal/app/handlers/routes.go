package handlers

import (
	"net/http"

	"fiber-boilerplate/internal/app/router"
	api "fiber-boilerplate/internal/generated/serviceapi"

	"github.com/gofiber/fiber/v2"
)

func withoutContent(code int) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		return ctx.Status(code).SendString("")
	}
}

// Router is for routing requests
var Router = &router.Router{
	Routes: []router.Route{
		{
			Name:        "Root",
			Method:      http.MethodGet,
			Pattern:     "/",
			HandlerFunc: withoutContent(http.StatusForbidden),
		},
		{
			Name:        "FavIcon",
			Method:      http.MethodGet,
			Pattern:     "/favicon.ico",
			HandlerFunc: withoutContent(http.StatusNoContent),
		},
	},
}

// Register routes with echo
func Register(f *fiber.App) {
	for _, route := range Router.Routes {
		f.Add(route.Method, route.Pattern, route.HandlerFunc)
	}

	h := new(APIHandlerBlock)
	r := f.Group("/api")
	api.RegisterHandlersWithOptions(r, h, api.FiberServerOptions{
		BaseURL:     "",
		Middlewares: []api.MiddlewareFunc{},
	})
}
