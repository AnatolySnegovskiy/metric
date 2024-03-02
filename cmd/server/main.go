package main

import (
	"github.com/AnatolySnegovskiy/metric/internal/server"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
)

func main() {
	s := server.New(storages.NewMemStorage())
	s.Run()
}
