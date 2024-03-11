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
	router  *chi.Mux
}

func New(s Storage) *Server {
	server := &Server{
		storage: s,
		router:  chi.NewRouter(),
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.router.NotFound(s.notFoundHandler)
	s.router.Post("/update/{metricType}/{metricName}/{metricValue}", s.writeMetricHandler)
	s.router.Get("/", s.showAllMetricHandler)
	s.router.Get("/value/{metricType}", s.showMetricTypeHandler)
	s.router.Get("/value/{metricType}/{metricName}", s.showMetricNameHandlers)
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
