package url

import (
	"strings"
	"testing"
)

func TestLookupOriginURL(t *testing.T) {
	t.Run("lookup origin url", func(t *testing.T) {
		shortCode := prepare().ShortCode

		t.Run("lookup - 301", func(t *testing.T) {
			urlModel, status, _ := LookupOriginURL(shortCode)
			if status != 301 {
				t.Errorf("unexpected status: wanted 301, go %v", status)
			}
			if !strings.Contains(urlModel.OriginURL, "http://localhost/some/long/url/") {
				t.Errorf("got unexpected origin url %v", urlModel.OriginURL)
			}
		})

		t.Run("lookup - 410", func(t *testing.T) {
			if nil != DeleteURLByShortCode(shortCode) {
				t.Errorf("should not return error")
			}

			if _, status, _ := LookupOriginURL(shortCode); status != 410 {
				t.Errorf("unexpected status: wanted 410, go %v", status)
			}
		})

		t.Run("lookup - 410", func(t *testing.T) {
			if _, status, _ := LookupOriginURL("n0cod3"); status != 404 {
				t.Errorf("unexpected status: wanted 404, go %v", status)
			}
		})
	})
}
