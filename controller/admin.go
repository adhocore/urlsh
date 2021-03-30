package controller

import (
	"net/http"

	"github.com/adhocore/urlsh/response"
	"github.com/adhocore/urlsh/service/url"
)

// ListURLs is the controller for url listing endpoint using filters from http.Request
// It responds to `GET /api/admin/urls` and requires auth token.
func ListURLs(res http.ResponseWriter, req *http.Request) {
	urls, err := url.ListURLsFilteredFromRequest(req)

	if err != nil {
		response.JSON(res, http.StatusNotFound, response.Body{"message": err.Error(), "urls": urls})
		return
	}

	response.JSON(res, http.StatusOK, response.Body{"urls": urls})
}

// DeleteShortURL is the controller for deleting short url
// It responds to `DELETE /api/admin/urls` and requires auth token.
func DeleteShortURL(res http.ResponseWriter, req *http.Request) {
	if err := url.DeleteURLFromRequest(req); err != nil {
		response.JSON(res, http.StatusNotFound, response.Body{"message": err.Error()})
		return
	}

	response.JSON(res, http.StatusOK, response.Body{"deleted": true})
}
