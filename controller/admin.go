package controller

import (
    "net/http"

    "github.com/adhocore/urlsh/response"
    "github.com/adhocore/urlsh/service/url"
)

func ListUrl(res http.ResponseWriter, req *http.Request) {
    urls, err := url.ListUrlsFilteredFromRequest(req)

    if err != nil {
        response.JSON(res, http.StatusNotFound, response.Body{"message": err.Error(), "urls": urls})
        return
    }

    response.JSON(res, http.StatusOK, response.Body{"urls": urls})
}
