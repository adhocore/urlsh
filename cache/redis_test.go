package cache

import (
	"net/http"
	"os"
	"testing"

	"github.com/adhocore/urlsh/model"
	"github.com/adhocore/urlsh/util"
)

func TestConnection(t *testing.T) {
	key := "APP_CACHE_HOST"
	old := os.Getenv(key)

	t.Run("nil", func(t *testing.T) {
		_ = os.Setenv(key, "")

		t.Run("not nil", func(t *testing.T) {
			_ = os.Setenv(key, old)
			if nil == Connection() {
				t.Errorf("connection should not be nil")
			}
		})
	})
}

func TestDeactivateURL(t *testing.T) {
	t.Run("deactivate", func(t *testing.T) {
		urlModel := model.URL{ShortCode: "testCode", OriginURL: "http://localhost/url"}

		DeactivateURL(urlModel)
		if !hasURL(urlModel) {
			t.Errorf("deactivated url should be in cache")
		}
	})
}

func TestLookupURL(t *testing.T) {
	t.Run("save + lookup", func(t *testing.T) {
		urlModel := model.URL{ShortCode: util.RandomString(8), OriginURL: "http://localhost/abc/def/ghi"}

		t.Run("save", func(t *testing.T) {
			SavePopularURL(urlModel, true)
			if !hasURL(urlModel) {
				t.Errorf("saved url should be in cache")
			}

			t.Run("lookup", func(t *testing.T) {
				cachedModel, status := LookupURL(urlModel.ShortCode)

				if status != http.StatusGone {
					t.Errorf("inactive url should be 410 gone")
				}
				if cachedModel.OriginURL != urlModel.OriginURL {
					t.Errorf("cached origin_url should be same as input")
				}
				if cachedModel.ShortCode != urlModel.ShortCode {
					t.Errorf("cached short_code should be same as input")
				}
			})
		})
	})
}
