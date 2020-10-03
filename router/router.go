package router

import (
    "net/http"

    "github.com/adhocore/urlsh/controller"
    "github.com/adhocore/urlsh/middleware"
)

var routes = map[string]http.HandlerFunc{
    "GET /": controller.Index,
    "POST /api/urls": controller.CreateShortUrl,
    "GET /api/admin/urls": controller.ListUrl,
    "DELETE /api/admin/urls" : controller.DeleteShortUrl,
}

func locateHandler(route string) http.HandlerFunc {
    handlerFunc, ok := routes[route]

    if !ok {
        return controller.NotFound
    }

    return handlerFunc
}

func RegisterHandlers() *http.ServeMux {
    handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        route := req.Method + " " + req.URL.Path

        locateHandler(route)(res, req)
    })

    mux := http.NewServeMux()
    mux.Handle("/", middleware.AdminAuth(handler))

    return mux
}
