package main

import (
	"github.com/AnatolySnegovskiy/metric/internal/server"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"log"
)

func main() {
	s := server.New(storages.NewMemStorage())
	err := s.Run(`:8080`)

	if err != nil {
		log.Fatalf("start server: %v", err)
	}
}
