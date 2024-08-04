package server

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/repositories"
	"github.com/AnatolySnegovskiy/metric/internal/services/interfase"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/AnatolySnegovskiy/metric/internal/storages/clients"
	"github.com/go-chi/chi/v5"
	"github.com/gookit/gsr"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
)

var pgxConnect = pgx.Connect

// Config defines the configuration interface for the server.
type Config interface {
	// GetServerAddress returns the server address.
	GetServerAddress() string
	// GetStoreInterval returns the interval for storing metrics.
	GetStoreInterval() int
	// GetFileStoragePath returns the file storage path.
	GetFileStoragePath() string
	// GetRestore returns a boolean indicating whether to restore metrics on start.
	GetRestore() bool
	// GetDataBaseDSN returns the database DSN.
	GetDataBaseDSN() string
	// GetShaKey returns the SHA key.
	GetShaKey() string
	// GetMigrationsDir returns the directory path for database migrations.
	GetMigrationsDir() string
	// GetCryptoKey returns the path to the private key file.
	GetCryptoKey() string
	// GetTrustedSubnet returns the trusted subnet.
	GetTrustedSubnet() *net.IPNet
}

// Server represents the main server struct.
type Server struct {
	storage  interfase.Storage
	router   *chi.Mux
	logger   gsr.GenLogger
	dbIsOpen bool
	conf     Config
}

// New creates a new server instance with the provided configuration and logger.
func New(ctx context.Context, c Config, l gsr.GenLogger) (*Server, error) {
	server := &Server{
		router: chi.NewRouter(),
		logger: l,
		conf:   c,
	}

	return server.upServer(ctx)
}

// setupRoutes sets up the routes for handling different HTTP endpoints.
func (s *Server) setupRoutes() {
	// Middleware functions and handlers for routing in the server.
	// Middleware functions:
	// - hashCheckMiddleware checks the hash of the request.
	// - gzipCompressMiddleware compresses the response using gzip.
	// - gzipDecompressMiddleware decompresses the request body using gzip.
	// - logMiddleware logs request information.
	// - hashResponseMiddleware hashes the response before sending.

	// NotFoundHandler handles requests for routes that are not found.
	// PostMetricHandler handles POST requests to update metrics.
	// MassPostMetricHandler handles POST requests to update multiple metrics.
	// ShowPostMetricHandler handles POST requests to display metrics.
	// WriteGetMetricHandler handles GET requests to update a specific metric.
	// ShowAllMetricHandler handles GET requests to show all metrics.
	// ShowMetricTypeHandler handles GET requests to show metrics of a specific type.
	// ShowMetricNameHandlers handles GET requests to show metrics of a specific name.

	// PostgresPingHandler handles GET requests to ping the PostgreSQL database.

	// Note: The router uses JSONContentTypeMiddleware for handling JSON content type in POST requests.
	s.router.Use(s.TrustedSubnetMiddleware, s.hashCheckMiddleware, s.DecryptMessageMiddleware, s.gzipCompressMiddleware, s.gzipDecompressMiddleware, s.logMiddleware, s.hashResponseMiddleware)
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

// Run starts the server and listens on the configured server address.
func (s *Server) Run() error {
	return http.ListenAndServe(s.conf.GetServerAddress(), s.router)
}

// saveMetricsPeriodically saves metrics to a file periodically based on the interval.
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

// loadMetricsOnStart loads metrics from a file on server start.
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

// saveMetricsToFile saves metrics to a file at a specified path.
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

// BDConnect establishes a connection to the database.
func (s *Server) BDConnect() *pgx.Conn {
	db, err := pgxConnect(context.Background(), s.conf.GetDataBaseDSN())

	if err != nil {
		s.logger.Error(err)
	}

	return db
}

// upStorage sets up the storage for metrics based on the database connection.
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

// upMigrate runs database migrations for the connected database.
func (s *Server) upMigrate(ctx context.Context, db *pgx.Conn) error {
	if db == nil {
		return nil
	}

	migration, _ := migrate.NewMigrator(ctx, db, "public.schema_version")
	if err := migration.LoadMigrations(os.DirFS(s.conf.GetMigrationsDir())); err != nil {
		return err
	}

	if err := migration.Migrate(ctx); err != nil {
		return err
	}

	return nil
}

// upServer initializes the server by connecting to the database, setting up migrations, storage, and routes.
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

// ShotDown saves metrics to a file before shutting down the server.
func (s *Server) ShotDown() {
	s.saveMetricsToFile(s.conf.GetFileStoragePath())
}

func (s *Server) UpGrpc() *GrpcServer {
	grpcServer := NewGrpcServer(s.storage, s.logger, s.conf)
	return grpcServer
}

// loadMetricsFromFile loads metrics from a file at the specified path and returns them as a map.
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
