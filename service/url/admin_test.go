package url

import (
    "fmt"
    "math/rand"
    "testing"

    "github.com/adhocore/urlsh/request"
)

func prepare() request.UrlFilter {
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
