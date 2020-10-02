package main

import (
    "log"
    "net/http"
    "os"
    "time"

    "github.com/adhocore/urlsh/router"
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

    router.RegisterHandlers()

    log.Println("Server running on port " + port)
    log.Fatal(server.ListenAndServe())
}
