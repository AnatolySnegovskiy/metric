package server

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"log"
	"net/http"
	"strings"
)

type Server struct {
	storage *storages.MemStorage
}

func New(storage *storages.MemStorage) *Server {
	return &Server{
		storage: storage,
	}
}

func (server *Server) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, server.HandleMetrics)
	err := http.ListenAndServe(`:8080`, mux)

	if err != nil {
		log.Fatalf("start server: %v", err)
	}
}

func (server *Server) HandleMetrics(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(rw, "method not allowed", http.StatusBadRequest)
		return
	}

	metricType, metricName, metricValue := server.parseURL(req.URL.Path)

	if metricName == "" {
		http.Error(rw, "metric name is required", http.StatusNotFound)
		return
	}

	storage := server.storage
	metric, err := storage.GetMetricType(metricType)

	if err != nil {
		http.Error(rw, fmt.Sprintf("metric type %s not found", metricType), http.StatusBadRequest)
		return
	}

	err = metric.Process(metricName, metricValue)

	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to process metric: %s", err.Error()), http.StatusBadRequest)
		return
	}

	storage.Log()
}

func (server *Server) parseURL(url string) (string, string, string) {
	elements := strings.Split(url, "/")
	return elements[2], elements[3], elements[4]
}
