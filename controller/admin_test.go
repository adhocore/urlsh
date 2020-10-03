package controller

import (
    "fmt"
    "math/rand"
    "os"
    "testing"
)

func TestListUrl(t *testing.T) {
    t.Run("list endpoint - not found", func(t *testing.T) {
        resp := request("GET", "/api/admin/urls?short_code=z", TestBody{}, ListUrl)

        resp.assertStatus(404, t)
    })

    t.Run("list endpoint - ok", func(t *testing.T) {
        _     = os.Setenv("APP_ALLOW_DUPE_URL", "1")
        url  := fmt.Sprintf("http://localhost:1000/very/long/url-%v", rand.Intn(1000000))
        body := TestBody{"url": url, "expires_on": "2030-01-01 00:00:00", "keywords": []string{"local"}}
        resp := request("POST", "/api/urls", body, CreateShortUrl)
        resp.assertStatus(200, t)

        t.Run("by page", func(t *testing.T) {
            resp := request("GET", "/api/admin/urls?page=1", TestBody{}, ListUrl)
            resp.assertStatus(200, t)
        })

        t.Run("by keyword", func(t *testing.T) {
            resp := request("GET", "/api/admin/urls?keyword=local", TestBody{}, ListUrl)
            resp.assertStatus(200, t)
        })
    })
}
