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
    body := response.Body{"short_code": shortCode, "short_url": shortUrl}

    if err == nil {
        response.JSON(res, http.StatusOK, body)
        return
    }

    status, errBody := http.StatusUnprocessableEntity, response.Body{"message": err.Error()}
    if shortCode != "" {
        status = http.StatusConflict
        errBody.Merge(body)
    }

    response.JSON(res, status, errBody)
}
