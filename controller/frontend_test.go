package controller

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestIndex(t *testing.T) {
	t.Run("index endpoint", func(t *testing.T) {
		resp := request("GET", "/status", TestBody{}, Status)

		resp.assertStatus(200, t)
		resp.assertKeyValue("message", "it works", t)
	})
}

func TestNotFound(t *testing.T) {
	t.Run("404 not found", func(t *testing.T) {
		resp := request("GET", "/not-found", TestBody{}, NotFound)

		resp.assertStatus(404, t)
		resp.assertKeyValue("message", "requested resource is not available", t)
	})
}

func TestServeShortURL(t *testing.T) {
	t.Run("serve short url", func(t *testing.T) {
		url := fmt.Sprintf("http://urlsh.lvh.me/urlsh/lvh/me/%v", rand.Intn(100000))
		resp := request("POST", "/api/urls", TestBody{"url": url}, CreateShortURL)
		shortCode := resp.assertContains("short_code", t).(string)

		t.Run("301", func(t *testing.T) {
			resp := request("GET", "/"+shortCode, TestBody{}, ServeShortURL)
			resp.assertStatus(301, t)
		})

		t.Run("404", func(t *testing.T) {
			resp := request("GET", "/n0cod3", TestBody{}, ServeShortURL)
			resp.assertStatus(404, t)
		})

		t.Run("delete - 410", func(t *testing.T) {
			resp := request("DELETE", "/api/admin/urls?short_code="+shortCode, TestBody{}, DeleteShortURL)

			t.Run("410", func(t *testing.T) {
				resp = request("GET", "/"+shortCode, TestBody{}, ServeShortURL)
				resp.assertStatus(410, t)
			})
		})
	})
}
