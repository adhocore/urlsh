package response

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		res := httptest.NewRecorder()
		body := Body{"message": "test msg"}

		JSON(res, 200, body)

		var actual Body
		_ = json.Unmarshal(res.Body.Bytes(), &actual)

		if _, ok := actual["status"]; !ok {
			t.Errorf("JSON must add status code")
		}

		if val, ok := actual["message"]; !ok || val != "test msg" {
			t.Errorf("JSON must keep message as it is, got %v", val)
		}
	})
}

func TestBody_Merge(t *testing.T) {
	t.Run("body merge", func(t *testing.T) {
		res := httptest.NewRecorder()
		body := Body{"message": "test msg"}

		JSON(res, 200, body.Merge(Body{"message": "new msg", "new_key": "new val"}))

		var actual Body
		_ = json.Unmarshal(res.Body.Bytes(), &actual)

		if val, ok := actual["message"]; !ok || val != "new msg" {
			t.Errorf("Merge must override message, got %v", val)
		}

		if val, ok := actual["new_key"]; !ok || val != "new val" {
			t.Errorf("Merge must add new key, got %v", val)
		}
	})
}
