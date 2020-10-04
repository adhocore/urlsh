package controller

import (
    "net/http"

    "github.com/adhocore/urlsh/response"
    "github.com/adhocore/urlsh/service/url"
)

// Index is the controller for root aka index page
// It responds to `GET /` and does not require auth token.
func Index(res http.ResponseWriter, _ *http.Request) {
    response.JSON(res, http.StatusOK, response.Body{"message": "it works"})
}

// NotFound is the controller for handling unregistered request path
// It is auto triggered if router does not find controller for request path.
func NotFound(res http.ResponseWriter, _ *http.Request) {
    response.JSON(res, http.StatusNotFound, response.Body{"message": "requested resource is not available"})
}

// ServeShortUrl is the controller for serving short urls
// It responds to `GET /{shortCode}` and does not require auth token.
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
