package url

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/adhocore/urlsh/request"
)

func prepare() request.URLFilter {
	rand.Seed(time.Now().UTC().UnixNano())

	longURL := fmt.Sprintf("http://localhost/some/long/url/%v", rand.Intn(1000000))
	shortCode, _ := CreateURLShortCode(request.URLInput{URL: longURL, Keywords: []string{"testing"}})

	return request.URLFilter{ShortCode: shortCode, Keyword: "testing"}
}

func TestListURLsFiltered(t *testing.T) {
	t.Run("list urls", func(t *testing.T) {
		if _, err := ListURLsFiltered(prepare()); err != nil {
			t.Errorf("should not return error")
		}
	})

	t.Run("list urls", func(t *testing.T) {
		if _, err := ListURLsFiltered(request.URLFilter{ShortCode: "zyx"}); err == nil {
			t.Errorf("should return error")
		}
	})
}

func TestDeleteURLByShortCode(t *testing.T) {
	t.Run("delete", func(t *testing.T) {
		if nil == DeleteURLByShortCode("") {
			t.Errorf("should return error")
		}
	})

	t.Run("delete", func(t *testing.T) {
		if nil == DeleteURLByShortCode("xyz") {
			t.Errorf("should return error")
		}
	})

	t.Run("delete - ok", func(t *testing.T) {
		shortCode := prepare().ShortCode
		if nil != DeleteURLByShortCode(shortCode) {
			t.Errorf("should not return error for %v", shortCode)
		}
	})
}
