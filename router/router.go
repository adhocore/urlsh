package router

import (
    "net/http"

    "github.com/adhocore/urlsh/common"
    "github.com/adhocore/urlsh/controller"
    "github.com/adhocore/urlsh/middleware"
)

var routes = map[string]http.HandlerFunc{
    "GET /": controller.Index,
    "POST /api/urls": controller.CreateShortUrl,
    "GET /api/admin/urls": controller.ListUrl,
    "DELETE /api/admin/urls" : controller.DeleteShortUrl,
}

func locateHandler(method string, path string) http.HandlerFunc {
    if handlerFunc, ok := routes[method + " " + path]; ok {
        return handlerFunc
    }

    if method == "GET" && common.ShortCodeRegex.MatchString(path[1:]) {
        return controller.ServeShortUrl
    }

    return controller.NotFound
}

func RegisterHandlers() *http.ServeMux {
    handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        locateHandler(req.Method, req.URL.Path)(res, req)
    })

    mux := http.NewServeMux()
    mux.Handle("/", middleware.AdminAuth(handler))

    return mux
}
