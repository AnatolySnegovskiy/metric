package server

import (
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/go-chi/chi/v5"
	"github.com/gookit/gsr"
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
	logger  gsr.GenLogger
}

func New(s Storage, l gsr.GenLogger) *Server {
	server := &Server{
		storage: s,
		router:  chi.NewRouter(),
		logger:  l,
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.router.Use(s.logMiddleware, s.gzipResponseMiddleware, s.gzipRequestMiddleware)
	s.router.NotFound(s.notFoundHandler)
	s.router.With(s.JSONContentTypeMiddleware).Post("/update/", s.writePostMetricHandler)
	s.router.With(s.JSONContentTypeMiddleware).Post("/value/", s.showPostMetricHandler)
	s.router.Post("/update/{metricType}/{metricName}/{metricValue}", s.writeGetMetricHandler)
	s.router.Get("/", s.showAllMetricHandler)
	s.router.Get("/value/{metricType}", s.showMetricTypeHandler)
	s.router.Get("/value/{metricType}/{metricName}", s.showMetricNameHandlers)
}

func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
