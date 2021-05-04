package response

import (
	"encoding/json"
	"net/http"
)

// Body is simple map for response
type Body map[string]interface{}

// Merge merges another Body with current body
// It returns the current Body which is result of merge.
func (body Body) Merge(other Body) Body {
	for key, val := range other {
		body[key] = val
	}

	return body
}

// JSON is handy shortcut for serving json response
// It writes header, status and Body to the given http.ResponseWriter.
func JSON(res http.ResponseWriter, status int, body Body) {
	body["status"] = status

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(status)

	_ = json.NewEncoder(res).Encode(body)
}
