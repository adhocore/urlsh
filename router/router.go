package router

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/adhocore/goic"
	"github.com/adhocore/urlsh/common"
	"github.com/adhocore/urlsh/controller"
	"github.com/adhocore/urlsh/middleware"
)

var routes = map[string]http.HandlerFunc{
	"GET /":                  controller.Index,
	"GET /banner.png":        controller.Banner,
	"GET /favicon.ico":       controller.Favicon,
	"GET /robots.txt":        controller.Robots,
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
	g := goic.New("/auth/o8", false)

	g.AddProvider(goic.Google.WithCredential(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET")))
	g.AddProvider(goic.Microsoft.WithCredential(os.Getenv("MICROSOFT_CLIENT_ID"), os.Getenv("MICROSOFT_CLIENT_SECRET")))
	g.AddProvider(goic.Yahoo.WithCredential(os.Getenv("YAHOO_CLIENT_ID"), os.Getenv("YAHOO_CLIENT_SECRET")))

	g.UserCallback(func(t *goic.Token, u *goic.User, res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(200)

		_ = json.NewEncoder(res).Encode(u)
	})

	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		locateHandler(req.Method, req.URL.Path)(res, req)
	})

	mux := http.NewServeMux()
	mux.Handle("/", g.MiddlewareHandler(middleware.Recover(middleware.AdminAuth(handler))))

	return mux
}
