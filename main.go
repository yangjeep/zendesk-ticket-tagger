package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
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
