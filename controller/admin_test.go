package controller

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestListURL(t *testing.T) {
	t.Run("list endpoint - not found", func(t *testing.T) {
		resp := request("GET", "/api/admin/urls?short_code=z", TestBody{}, ListURLs)

		resp.assertStatus(404, t)
	})

	t.Run("list endpoint - ok", func(t *testing.T) {
		rand.Seed(time.Now().UTC().UnixNano())

		_ = os.Setenv("APP_ALLOW_DUPE_URL", "1")
		url := fmt.Sprintf("http://localhost:1000/very/long/url-%v", rand.Intn(1000000))
		body := TestBody{"url": url, "expires_on": "2030-01-01 00:00:00", "keywords": []string{"local"}}
		resp := request("POST", "/api/urls", body, CreateShortURL)

		resp.assertStatus(200, t)

		t.Run("by page", func(t *testing.T) {
			resp := request("GET", "/api/admin/urls?page=1", TestBody{}, ListURLs)
			resp.assertStatus(200, t)
		})

		t.Run("by keyword", func(t *testing.T) {
			resp := request("GET", "/api/admin/urls?keyword=local", TestBody{}, ListURLs)
			resp.assertStatus(200, t)
		})
	})
}

func TestDeleteShortURL(t *testing.T) {
	t.Run("delete endpoint", func(t *testing.T) {
		rand.Seed(time.Now().UTC().UnixNano())

		body := TestBody{"url": fmt.Sprintf("https://localhost/test/delete/short/url/%v", rand.Intn(1000000))}
		resp := request("POST", "/api/urls", body, CreateShortURL)
		resp.assertStatus(200, t)
		resp.assertContains("short_code", t)

		shortCode := resp["short_code"]

		t.Run("delete - nok", func(t *testing.T) {
			uri := fmt.Sprintf("/api/admin/urls?short_code=%v", rand.Intn(1000000))
			resp := request("DELETE", uri, TestBody{}, DeleteShortURL)
			resp.assertStatus(404, t)
		})

		t.Run("delete - ok", func(t *testing.T) {
			uri := fmt.Sprintf("/api/admin/urls?short_code=%v", shortCode)
			resp := request("DELETE", uri, TestBody{}, DeleteShortURL)
			resp.assertStatus(200, t)

			t.Run("delete - ok - nok", func(t *testing.T) {
				resp = request("DELETE", uri, TestBody{}, DeleteShortURL)
				resp.assertStatus(404, t)
			})
		})
	})
}
