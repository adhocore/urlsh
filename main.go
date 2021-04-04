package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/adhocore/urlsh/router"
)

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}

	return "2000"
}

func main() {
	port := getPort()

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router.RegisterHandlers(),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	log.Printf("Server running on port %v", port)
	log.Fatal(server.ListenAndServe())
}
