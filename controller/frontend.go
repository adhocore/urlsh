package controller

import (
    "net/http"

    "github.com/adhocore/urlsh/response"
    "github.com/adhocore/urlsh/service/url"
)

func Index(res http.ResponseWriter, _ *http.Request) {
    response.JSON(res, http.StatusOK, response.Body{"message": "it works"})
}

func NotFound(res http.ResponseWriter, _ *http.Request) {
    response.JSON(res, http.StatusNotFound, response.Body{"message": "requested resource is not available"})
}

func ServeShortUrl(res http.ResponseWriter, req *http.Request) {
    shortCode := req.URL.Path[1:]
    location, status := url.LookupOriginUrl(shortCode)

    if status != http.StatusFound {
        response.JSON(res, status, response.Body{"message": "requested resource is not available"})
        return
    }

    url.IncrementHits(shortCode)
    http.Redirect(res, req, location, status)
}
