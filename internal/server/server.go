package server

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"net/http"
	"strings"
)

type Storage interface {
	GetMetricType(metricType string) (storages.EntityMetric, error)
	AddMetric(metricType string, metric storages.EntityMetric)
	GetList() map[string]storages.EntityMetric
}

type Server struct {
	storage Storage
}

func New(s Storage) *Server {
	return &Server{
		storage: s,
	}
}

func (s *Server) Run(addr string) error {
	mux := http.NewServeMux()
	mux.Handle(`/update/`, postMiddleware(http.HandlerFunc(s.writeMetricHandlers)))
	mux.Handle(`/show/`, getMiddleware(http.HandlerFunc(s.showMetricHandlers)))
	return http.ListenAndServe(addr, mux)
}

func parseURL(url string) (string, string, string, error) {
	elements := strings.Split(url, "/")

	if len(elements) < 5 || len(elements) > 5 {
		return "", "", "", fmt.Errorf("invalid url")
	}

	return elements[2], elements[3], elements[4], nil
}

func postMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
