package server

import (
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/go-chi/chi/v5"
	"net/http"
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
	r := chi.NewRouter()
	r.NotFound(s.notFoundHandler) // H
	r.Post("/update/{metricType}/{metricName}/{metricValue}", s.writeMetricHandlers)
	r.Get("/", s.showAllMetricHandlers)
	r.Get("/{metricType}", s.showMetricTypeHandlers)
	r.Get("/{metricType}/{metricName}", s.showMetricNameHandlers)

	return http.ListenAndServe(addr, r)
}
