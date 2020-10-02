package main

import (
    "log"
    "net/http"
    "os"
    "time"

    "github.com/adhocore/urlsh/controller"
)

func main() {
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "1000"
    }

    server := &http.Server{
        Addr:         ":" + port,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    http.HandleFunc("/", controller.Index)

    log.Fatal(server.ListenAndServe())
}
