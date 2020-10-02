package router

import (
    "github.com/adhocore/urlsh/controller"
    "net/http"
)

type handler func(res http.ResponseWriter, req *http.Request)

var routes = map[string]handler{
    "GET /": controller.Index,
    "POST /api/urls": controller.CreateShortUrl,
}

func locateHandler(route string) handler {
    handlerFunc, ok := routes[route]

    if !ok {
        return controller.NotFound
    }

    return handlerFunc
}

func RegisterHandlers() {
    http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
        route := req.Method + " " + req.URL.Path

        locateHandler(route)(res, req)
    })
}
