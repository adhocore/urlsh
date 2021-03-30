package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/adhocore/urlsh/common"
	"github.com/adhocore/urlsh/response"
)

// AdminURIPrefix is the uri to intercept by auth middleware
const AdminURIPrefix = "/api/admin"

// validateAdminToken validates request header token against env token for admin end
// It returns possible http status code and error if auth token missing/invalid.
func validateAdminToken(req *http.Request) (int, error) {
	adminToken := os.Getenv("APP_ADMIN_TOKEN")

	// Require token if only backend is configured and the uri matches admin prefix.
	if adminToken == "" || strings.Index(req.URL.Path, AdminURIPrefix) != 0 {
		return 0, nil
	}

	token := req.Header.Get("Authorization")
	if token == "" {
		return http.StatusUnauthorized, common.ErrTokenRequired
	}

	if strings.Index(token, "Bearer ") != 0 {
		return http.StatusUnauthorized, common.ErrTokenInvalid
	}

	if token[7:] != adminToken {
		return http.StatusForbidden, common.ErrTokenInvalid
	}

	return 0, nil
}

// AdminAuth intercepts request and requires token for admin endpoints
func AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if status, err := validateAdminToken(req); err != nil {
			response.JSON(res, status, response.Body{"message": err.Error()})

			return
		}

		next.ServeHTTP(res, req)
	})
}
