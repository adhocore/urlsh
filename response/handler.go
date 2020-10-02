package response

import (
    "encoding/json"
    "net/http"
)

type Body map[string]interface{}

func JSON(res http.ResponseWriter, status int, body Body) {
    body["status"] = status

    res.Header().Add("Content-Type", "application/json")
    res.WriteHeader(status)

    _ = json.NewEncoder(res).Encode(body)
}
