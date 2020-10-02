package controller

import "testing"

func TestIndex(t *testing.T) {
    t.Run("index endpoint", func(t *testing.T) {
        resp := request("GET", "/", TestBody{}, Index)

        resp.assertStatus(200, t)
        resp.assertKeyValue("message", "it works", t)
    })
}
