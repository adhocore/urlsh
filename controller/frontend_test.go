package controller

import "testing"

func TestIndex(t *testing.T) {
    t.Run("index endpoint", func(t *testing.T) {
        resp := request("GET", "/", TestBody{}, Index)

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
