package router

import "github.com/gofiber/fiber/v2"

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string

	// Method is the string for the HTTP method. integration) GET, POST etc..
	Method string

	// Pattern is the pattern of the URI.
	Pattern string

	// HandlerFunc is the handler function of this route.
	HandlerFunc fiber.Handler
}

// Router :
type Router struct {
	Routes []Route
}

// FindRouteByName :
func (r Router) FindRouteByName(name string) *Route {
	for _, r := range r.Routes {
		if r.Name == name {
			return &r
		}
	}

	return nil
}
