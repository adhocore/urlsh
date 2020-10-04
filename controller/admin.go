package controller

import (
    "net/http"

    "github.com/adhocore/urlsh/response"
    "github.com/adhocore/urlsh/service/url"
)

// ListUrls is the controller for url listing endpoint using filters from http.Request
// It responds to `GET /api/admin/urls` and requires auth token.
func ListUrls(res http.ResponseWriter, req *http.Request) {
    urls, err := url.ListUrlsFilteredFromRequest(req)

    if err != nil {
        response.JSON(res, http.StatusNotFound, response.Body{"message": err.Error(), "urls": urls})
        return
    }

    response.JSON(res, http.StatusOK, response.Body{"urls": urls})
}

// DeleteShortUrl is the controller for deleting short url
// It responds to `DELETE /api/admin/urls` and requires auth token.
func DeleteShortUrl(res http.ResponseWriter, req *http.Request) {
    if err := url.DeleteUrlFromRequest(req); err != nil {
        response.JSON(res, http.StatusNotFound, response.Body{"message": err.Error()})
        return
    }

    response.JSON(res, http.StatusOK, response.Body{"deleted": true})
}
