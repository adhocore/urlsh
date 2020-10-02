package controller

import (
    "fmt"
    "net/http"

    "github.com/adhocore/urlsh/response"
    "github.com/adhocore/urlsh/service/url"
)

func CreateShortUrl(res http.ResponseWriter, req *http.Request) {
    shortCode, err := url.CreateUrlShortCodeFromRequest(req)
    shortUrl := fmt.Sprintf("%s%s%s", "http://", req.Host, "/" + shortCode)

    if err == nil {
        response.JSON(res, http.StatusOK, response.Body{"short_code": shortCode, "short_url": shortUrl})
        return
    }

    status, body := http.StatusUnprocessableEntity, response.Body{"message": err.Error()}
    if shortCode != "" {
        status = http.StatusConflict
        body["short_code"] = shortCode
    }

    response.JSON(res, status, body)
}
