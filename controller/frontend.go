package controller

import (
	"net/http"

	"github.com/adhocore/urlsh/response"
	"github.com/adhocore/urlsh/service/url"
)

// EmbedAssetHandler is static assets handler.
var EmbedAssetHandler http.Handler

// Index is the controller for root aka index page
// It responds to `GET /` and does not require auth token.
func Index(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "assets/home.html"
	EmbedAssetHandler.ServeHTTP(w, r)
}

// Banner is the controller for favicon.ico
// It responds to `GET /banner.png` and does not require auth token.
func Banner(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "assets/urlsh-light.png"
	EmbedAssetHandler.ServeHTTP(w, r)
}

// Favicon is the controller for favicon.ico
// It responds to `GET /favicon.ico` and does not require auth token.
func Favicon(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "assets/u.png"
	EmbedAssetHandler.ServeHTTP(w, r)
}

// Robots is the controller for robots.txt
// It responds to `GET /robots.txt` and does not require auth token.
func Robots(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = "assets/robots.txt"
	EmbedAssetHandler.ServeHTTP(w, r)
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
