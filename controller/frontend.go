package controller

import (
    "net/http"

    "github.com/adhocore/urlsh/response"
)

func Index(res http.ResponseWriter, req *http.Request) {
    response.JSON(res, http.StatusOK, response.Body{"message": "it works"})
}

func NotFound(res http.ResponseWriter, req *http.Request) {
    response.JSON(res, http.StatusNotFound, response.Body{"message": "requested resource is not available"})
}
