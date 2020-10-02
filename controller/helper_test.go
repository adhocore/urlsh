package controller

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"
)

type TestBody map[string]interface{}

func request(method string, uri string, data TestBody, handler http.HandlerFunc) TestBody {
    var body io.Reader

    if len(data) > 0 {
        buf, _ := json.Marshal(data)
        body = bytes.NewBuffer(buf)
    }

    req, _ := http.NewRequest(method, uri, body)
    res := httptest.NewRecorder()

    handler(res, req)

    var actual TestBody
    _ = json.Unmarshal(res.Body.Bytes(), &actual)

    return actual
}

func (body TestBody) assertContains(key string, t *testing.T) interface{} {
    if val, ok := body[key]; ok {
        return val
    }

    t.Errorf("response does not contain %s", key)
    return nil
}

func (body TestBody) assertStatus(status int, t *testing.T)  {
    actual := body.assertContains("status", t)

    if status != int(actual.(float64)) {
        t.Errorf("response status does not match: wanted %v, got %v", status, actual)
    }
}

func (body TestBody) assertKeyValue(key string, value interface{}, t *testing.T) {
    actual := body.assertContains(key, t)

    if actual != value {
        t.Errorf("response value does not match: wanted %v, got %v", value, actual)
    }
}
