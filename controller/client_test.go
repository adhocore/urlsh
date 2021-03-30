package controller

import (
	"os"
	"strings"
	"testing"

	"github.com/adhocore/urlsh/common"
)

func TestCreateShortURL(t *testing.T) {
	t.Run("create short url - invalid length", func(t *testing.T) {
		resp := request("POST", "/api/urls", TestBody{"url": ""}, CreateShortURL)

		resp.assertStatus(422, t)
		resp.assertKeyValue("message", common.ErrInvalidURLLen.Error(), t)
	})

	t.Run("create short url - invalid url", func(t *testing.T) {
		resp := request("POST", "/api/urls", TestBody{"url": "http:/localhost"}, CreateShortURL)

		resp.assertStatus(422, t)
		resp.assertKeyValue("message", common.ErrInvalidURL.Error(), t)
	})

	t.Run("create short url - blacklist url", func(t *testing.T) {
		resp := request("POST", "/api/urls", TestBody{"url": "http://localhost/xxx"}, CreateShortURL)

		resp.assertStatus(422, t)
		resp.assertKeyValue("message", "url matches blacklist pattern", t)
	})

	t.Run("create short url - past expiry", func(t *testing.T) {
		body := TestBody{"url": "http://localhost", "expires_on": "2020-01-01 00:00:00"}
		resp := request("POST", "/api/urls", body, CreateShortURL)

		resp.assertStatus(422, t)
		resp.assertKeyValue("message", "expires_on can not be date in past", t)
	})

	t.Run("create short url - OK", func(t *testing.T) {
		tester := func(status int, message string) {
			body := TestBody{"url": "http://localhost:1000/very/long/url", "expires_on": "2030-01-01 00:00:00"}
			resp := request("POST", "/api/urls", body, CreateShortURL)

			resp.assertStatus(status, t)
			resp.assertContains("short_code", t)
			resp.assertContains("short_url", t)

			if message != "" {
				resp.assertKeyValue("message", message, t)
			}

			if !strings.Contains(resp["short_url"].(string), resp["short_code"].(string)) {
				t.Errorf("short_url must contain short_code")
			}
		}

		_ = os.Setenv("APP_ALLOW_DUPE_URL", "1")
		tester(200, "")

		_ = os.Setenv("APP_ALLOW_DUPE_URL", "0")
		tester(409, "url is already shortened")
	})
}
