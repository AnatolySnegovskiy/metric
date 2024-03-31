package server

import (
	"context"
	"encoding/json"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/go-chi/chi/v5"
	"github.com/gookit/gsr"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
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

func (s *Server) SaveMetricsPeriodically(ctx context.Context, interval int, filePath string) {
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := s.saveMetricsToFile(filePath)
			if err != nil {
				s.logger.Error(err)
			}
		}
	}
}

func (s *Server) LoadMetricsOnStart(filePath string) {
	savedMetrics, err := loadMetricsFromFile(filePath)

	if err != nil {
		s.logger.Error(err)
	}

	for metricType, metricValues := range savedMetrics {
		metric, err := s.storage.GetMetricType(metricType)

		if err != nil {
			s.logger.Error(err)
			continue
		}

		for _, items := range metricValues {
			for key, value := range items {
				_ = metric.Process(key, strconv.FormatFloat(value, 'f', -1, 64))
			}
		}
	}

	s.logger.Info("Metrics loaded: " + filePath)
}

func (s *Server) HandleShutdownSignal(filePath string) error {
	return s.saveMetricsToFile(filePath)
}

func (s *Server) saveMetricsToFile(filePath string) error {
	projectDir, err := os.Getwd()
	if err != nil {
		return err
	}

	absoluteFilePath := filepath.Join(projectDir, filePath)

	directory := filepath.Dir(absoluteFilePath)
	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(absoluteFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.Marshal(s.storage.GetList())
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	s.logger.Info("Metrics saved: " + absoluteFilePath)
	return nil
}

func loadMetricsFromFile(filePath string) (map[string]map[string]map[string]float64, error) {
	projectDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	absoluteFilePath := filepath.Join(projectDir, filePath)

	file, err := os.Open(absoluteFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var metrics map[string]map[string]map[string]float64
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
