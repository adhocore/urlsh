package middleware

import (
	"net/http"

	"github.com/adhocore/urlsh/common"
	"github.com/adhocore/urlsh/response"
)

// Recover recovers from panic
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		defer func() {
			err := common.ErrServerError
			if rec := recover(); rec != nil {
				switch typ := rec.(type) {
				case error:
					err = typ
				}

				response.JSON(res, http.StatusInternalServerError, response.Body{"message": err.Error()})
			}
		}()

		next.ServeHTTP(res, req)
	})
}
