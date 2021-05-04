package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/adhocore/urlsh/common"
	"github.com/adhocore/urlsh/controller"
	"github.com/adhocore/urlsh/response"
)

func tester(token string, expectStatus int, t *testing.T) *httptest.ResponseRecorder {
	var body io.Reader

	req, _ := http.NewRequest("GET", "/api/admin/urls", body)
	if token != "" {
		req.Header.Set("Authorization", token)
	}

	res := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.Handle("/", AdminAuth(http.HandlerFunc(controller.ListURLs)))
	mux.ServeHTTP(res, req)

	if res.Result().StatusCode != expectStatus {
		t.Errorf("wanted status %v, got %v", expectStatus, res.Result().Status)
	}

	return res
}

func TestAdminAuth(t *testing.T) {
	t.Run("auth - token required", func(t *testing.T) {
		var actual response.Body

		_ = os.Setenv("APP_ADMIN_TOKEN", "some_token")
		res := tester("", http.StatusUnauthorized, t)

		_ = json.Unmarshal(res.Body.Bytes(), &actual)
		if actual["message"] != common.ErrTokenRequired.Error() {
			t.Errorf("got unexpected message: %v", actual["message"])
		}
	})

	t.Run("auth - token required", func(t *testing.T) {
		var actual response.Body

		_ = os.Setenv("APP_ADMIN_TOKEN", "correct_token")
		res := tester("Bearer wrong_token", http.StatusForbidden, t)

		_ = json.Unmarshal(res.Body.Bytes(), &actual)
		if actual["message"] != common.ErrTokenInvalid.Error() {
			t.Errorf("got unexpected message: %v", actual["message"])
		}
	})

	t.Run("auth - ok", func(t *testing.T) {
		_ = os.Setenv("APP_ADMIN_TOKEN", "correct_token")
		_ = tester("Nobearer correct_token", http.StatusUnauthorized, t)
	})

	t.Run("auth - ok empty", func(t *testing.T) {
		_ = os.Setenv("APP_ADMIN_TOKEN", "")
		_ = tester("Bearer token", http.StatusOK, t)
	})

	t.Run("auth - ok", func(t *testing.T) {
		_ = os.Setenv("APP_ADMIN_TOKEN", "correct_token")
		_ = tester("Bearer correct_token", http.StatusOK, t)
	})
}
