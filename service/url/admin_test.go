package url

import (
    "fmt"
    "math/rand"
    "testing"
    "time"

    "github.com/adhocore/urlsh/request"
)

func prepare() request.UrlFilter {
    rand.Seed(time.Now().UTC().UnixNano())

    longUrl := fmt.Sprintf("http://localhost/some/long/url/%v", rand.Intn(1000000))
    shortCode, _:= CreateUrlShortCode(request.UrlInput{Url: longUrl, Keywords: []string{"testing"}})

    return request.UrlFilter{ShortCode: shortCode, Keyword: "testing"}
}

func TestListUrlsFiltered(t *testing.T) {
    t.Run("list urls", func(t *testing.T) {
        if _, err := ListUrlsFiltered(prepare()); err != nil {
            t.Errorf("should not return error")
        }
    })

    t.Run("list urls", func(t *testing.T) {
        if _, err := ListUrlsFiltered(request.UrlFilter{ShortCode: "zyx"}); err == nil {
            t.Errorf("should return error")
        }
    })
}

func TestDeleteUrlByShortCode(t *testing.T) {
    t.Run("delete", func(t *testing.T) {
        if nil == DeleteUrlByShortCode("xyz") {
            t.Errorf("should return error")
        }
    })

    t.Run("delete - ok", func(t *testing.T) {
        shortCode := prepare().ShortCode
        if nil != DeleteUrlByShortCode(shortCode) {
            t.Errorf("should not return error for %v", shortCode)
        }
    })
}
