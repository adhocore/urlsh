package url

import (
    "strings"
    "testing"
)

func TestLookupOriginUrl(t *testing.T) {
    t.Run("lookup origin url", func(t *testing.T) {
        shortCode := prepare().ShortCode

        t.Run("lookup - 302", func(t *testing.T) {
            location, status := LookupOriginUrl(shortCode)
            if status != 302 {
                t.Errorf("unexpected status: wanted 302, go %v", status)
            }
            if !strings.Contains(location, "http://localhost/some/long/url/") {
                t.Errorf("got unexpected origin url %v", location)
            }
        })

        t.Run("lookup - 410", func(t *testing.T) {
            if nil != DeleteUrlByShortCode(shortCode) {
                t.Errorf("should not return error")
            }

            if _, status := LookupOriginUrl(shortCode); status != 410 {
                t.Errorf("unexpected status: wanted 410, go %v", status)
            }
        })

        t.Run("lookup - 410", func(t *testing.T) {
            if _, status := LookupOriginUrl("n0cod3"); status != 404 {
                t.Errorf("unexpected status: wanted 404, go %v", status)
            }
        })
    })
}
