package server

import (
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/go-chi/chi/v5"
	"net/http"
)

//go:generate mockgen -source=server.go -destination=mocks/server_mock.go -package=mocks
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
	s.router.Post("/update/{metricType}/{metricName}/{metricValue}", s.writeMetricHandlers)
	s.router.Get("/", s.showAllMetricHandlers)
	s.router.Get("/value/{metricType}", s.showMetricTypeHandlers)
	s.router.Get("/value/{metricType}/{metricName}", s.showMetricNameHandlers)
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
