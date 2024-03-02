package main

import (
	"github.com/AnatolySnegovskiy/metric/internal/server"
	"log"
	"net/http"
)

func main() {
	app := server.New()
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, app.Metric)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
