package server

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"net/http"
	"strings"
)

type Storage interface {
	GetMetricType(metricType string) (storages.EntityMetric, error)
	Log()
}

type Server struct {
	storage *storages.MemStorage
}

func New(s *storages.MemStorage) *Server {
	return &Server{
		storage: s,
	}
}

func (s *Server) Run(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, s.HandleMetrics)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) HandleMetrics(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(rw, "method not allowed", http.StatusBadRequest)
		return
	}

	metricType, metricName, metricValue := parseURL(req.URL.Path)

	if metricName == "" {
		http.Error(rw, "metric name is required", http.StatusNotFound)
		return
	}

	storage := s.storage
	metric, err := storage.GetMetricType(metricType)

	if err != nil {
		http.Error(rw, fmt.Sprintf("metric type %s not found", metricType), http.StatusNotFound)
		return
	}

	if err := metric.Process(metricName, metricValue); err != nil {
		http.Error(rw, fmt.Sprintf("failed to process metric: %s", err.Error()), http.StatusBadRequest)
		return
	}

	storage.Log()
}

func parseURL(url string) (string, string, string) {
	elements := strings.Split(url, "/")
	return elements[2], elements[3], elements[4]
}
