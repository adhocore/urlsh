package router

import (
	"net/http"

	"github.com/adhocore/urlsh/common"
	"github.com/adhocore/urlsh/controller"
	"github.com/adhocore/urlsh/middleware"
)

var routes = map[string]http.HandlerFunc{
	"GET /":                  controller.Index,
	"GET /status":            controller.Status,
	"POST /api/urls":         controller.CreateShortURL,
	"GET /api/admin/urls":    controller.ListURLs,
	"DELETE /api/admin/urls": controller.DeleteShortURL,
}

// locateHandler locates controller for given http request method and path
// It also handles not found case and short code redirection.
func locateHandler(method string, path string) http.HandlerFunc {
	if handlerFunc, ok := routes[method+" "+path]; ok {
		return handlerFunc
	}

	if method == "GET" && common.ShortCodeRegex.MatchString(path[1:]) {
		return controller.ServeShortURL
	}

	return controller.NotFound
}

// RegisterHandlers registers middlewares, handlers and route locators
// It returns server mux which then can be attached to a http server.
func RegisterHandlers() *http.ServeMux {
	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		locateHandler(req.Method, req.URL.Path)(res, req)
	})

	mux := http.NewServeMux()
	mux.Handle("/", middleware.AdminAuth(handler))

	return mux
}
