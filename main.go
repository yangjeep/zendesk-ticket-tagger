package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yangjeep/zendesk-ticket-tagger/config"
)

func main() {
	cfg := config.Load()

	// Now you can use the configuration
	fmt.Printf("Server running on port %d\n", cfg.ServerPort)

	s := http.Server{
		Addr:         "0.0.0.0:1234",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("hello world"))
		}),
	}

	log.Printf("Listening at http://%s", s.Addr)
	log.Fatal(s.ListenAndServe())
}
