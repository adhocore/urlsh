package router

import (
    "net/http"

    "github.com/adhocore/urlsh/controller"
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
