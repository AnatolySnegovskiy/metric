package server

import (
	"context"
	"encoding/json"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/repositories"
	"github.com/AnatolySnegovskiy/metric/internal/services/interfase"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/AnatolySnegovskiy/metric/internal/storages/clients"
	"github.com/go-chi/chi/v5"
	"github.com/gookit/gsr"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var pgxConnect = pgx.Connect

type Config interface {
	GetServerAddress() string
	GetStoreInterval() int
	GetFileStoragePath() string
	GetRestore() bool
	GetDataBaseDSN() string
	GetShaKey() string
}

type Server struct {
	storage  interfase.Storage
	router   *chi.Mux
	logger   gsr.GenLogger
	dbIsOpen bool
	conf     Config
}

func New(ctx context.Context, c Config, l gsr.GenLogger) (*Server, error) {
	server := &Server{
		router: chi.NewRouter(),
		logger: l,
		conf:   c,
	}

	return server.upServer(ctx)
}

func (s *Server) setupRoutes() {
	s.router.Use(s.hashCheckMiddleware, s.gzipCompressMiddleware, s.gzipDecompressMiddleware, s.logMiddleware, s.hashResponseMiddleware)
	s.router.NotFound(s.notFoundHandler)
	s.router.With(s.JSONContentTypeMiddleware).Post("/update/", s.writePostMetricHandler)
	s.router.With(s.JSONContentTypeMiddleware).Post("/updates/", s.writeMassPostMetricHandler)
	s.router.With(s.JSONContentTypeMiddleware).Post("/value/", s.showPostMetricHandler)
	s.router.Post("/update/{metricType}/{metricName}/{metricValue}", s.writeGetMetricHandler)
	s.router.Get("/", s.showAllMetricHandler)
	s.router.Get("/value/{metricType}", s.showMetricTypeHandler)
	s.router.Get("/value/{metricType}/{metricName}", s.showMetricNameHandlers)

	s.router.Get("/ping", s.postgersPingHandler)
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.conf.GetServerAddress(), s.router)
}

func (s *Server) saveMetricsPeriodically(ctx context.Context, interval int, filePath string) {
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.saveMetricsToFile(filePath)
		}
	}
}

func (s *Server) loadMetricsOnStart(filePath string) {
	savedMetrics := loadMetricsFromFile(filePath)

	for metricType, metricValues := range savedMetrics {
		metric, err := s.storage.GetMetricType(metricType)

		if err != nil {
			s.logger.Error(err)
			continue
		}

		for _, items := range metricValues {
			for key, value := range items {
				_ = metric.Process(context.Background(), key, strconv.FormatFloat(value, 'f', -1, 64))
			}
		}
	}

	s.logger.Info("Metrics loaded: " + filePath)
}

func (s *Server) saveMetricsToFile(filePath string) {
	projectDir, _ := os.Getwd()
	absoluteFilePath := filepath.Join(projectDir, filePath)

	directory := filepath.Dir(absoluteFilePath)
	_ = os.MkdirAll(directory, os.ModePerm)

	file, _ := os.Create(absoluteFilePath)
	defer file.Close()
	jsonData, _ := json.Marshal(s.storage.GetList())

	_, _ = file.Write(jsonData)
	s.logger.Info("Metrics saved: " + absoluteFilePath)
}

func (s *Server) BDConnect() *pgx.Conn {
	db, err := pgxConnect(context.Background(), s.conf.GetDataBaseDSN())

	if err != nil {
		s.logger.Error(err)
	}

	return db
}

func (s *Server) upStorage(db *pgx.Conn) error {
	var gaugeRepo *repositories.GaugeRepo
	var counterRepo *repositories.CounterRepo

	if db != nil {
		pg := clients.NewPostgres(db)
		gaugeRepo = repositories.NewGaugeRepo(pg)
		counterRepo = repositories.NewCounterRepo(pg)
	}

	stg := storages.NewMemStorage()
	stg.AddMetric("gauge", metrics.NewGauge(gaugeRepo))
	stg.AddMetric("counter", metrics.NewCounter(counterRepo))
	s.storage = stg

	return nil
}

func (s *Server) upMigrate(ctx context.Context, db *pgx.Conn) error {
	if db == nil {
		return nil
	}

	migration, _ := migrate.NewMigrator(ctx, db, "public.schema_version")
	projectDir, _ := os.Getwd()

	if err := migration.LoadMigrations(os.DirFS(projectDir + "/internal/storages/migrations")); err != nil {
		return err
	}

	if err := migration.Migrate(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Server) upServer(ctx context.Context) (*Server, error) {
	db := s.BDConnect()
	s.dbIsOpen = db != nil

	if err := s.upMigrate(ctx, db); err != nil {
		return nil, err
	}
	if err := s.upStorage(db); err != nil {
		return nil, err
	}

	fileStorage := s.conf.GetFileStoragePath()

	if s.conf.GetRestore() {
		s.loadMetricsOnStart(fileStorage)
	}

	go s.saveMetricsPeriodically(ctx, s.conf.GetStoreInterval(), s.conf.GetFileStoragePath())

	s.setupRoutes()

	return s, nil
}

func (s *Server) ShotDown() {
	s.saveMetricsToFile(s.conf.GetFileStoragePath())
}

func loadMetricsFromFile(filePath string) map[string]map[string]map[string]float64 {
	projectDir, _ := os.Getwd()
	absoluteFilePath := filepath.Join(projectDir, filePath)
	file, _ := os.Open(absoluteFilePath)
	defer file.Close()

	var metrics map[string]map[string]map[string]float64
	decoder := json.NewDecoder(file)
	_ = decoder.Decode(&metrics)

	return metrics
}
