package controller

import (
	"net/http"

	"github.com/adhocore/urlsh/response"
	"github.com/adhocore/urlsh/service/url"
)

// Index is the controller for root aka index page
// It responds to `GET /` and does not require auth token.
func Index(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "tmpl/home.html")
}

// Banner is the controller for favicon.ico
// It responds to `GET /banner.png` and does not require auth token.
func Banner(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "tmpl/urlsh-light.png")
}

// Favicon is the controller for favicon.ico
// It responds to `GET /favicon.ico` and does not require auth token.
func Favicon(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "tmpl/u.png")
}

// Robots is the controller for robots.txt
// It responds to `GET /robots.txt` and does not require auth token.
func Robots(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "tmpl/robots.txt")
}

// Status is the controller for health/status check
// It responds to `GET /status` and does not require auth token.
func Status(res http.ResponseWriter, _ *http.Request) {
	response.JSON(res, http.StatusOK, response.Body{"message": "it works"})
}

// NotFound is the controller for handling unregistered request path
// It is auto triggered if router does not find controller for request path.
func NotFound(res http.ResponseWriter, _ *http.Request) {
	response.JSON(res, http.StatusNotFound, response.Body{"message": "requested resource is not available"})
}

// ServeShortURL is the controller for serving short urls
// It responds to `GET /{shortCode}` and does not require auth token.
func ServeShortURL(res http.ResponseWriter, req *http.Request) {
	shortCode := req.URL.Path[1:]
	urlModel, status, cached := url.LookupOriginURL(shortCode)

	if cached {
		res.Header().Add("X-Cached", "true")
	}

	if status != http.StatusMovedPermanently {
		response.JSON(res, status, response.Body{"message": "requested resource is not available"})
		return
	}

	go url.IncrementHits(urlModel)
	http.Redirect(res, req, urlModel.OriginURL, status)
}
