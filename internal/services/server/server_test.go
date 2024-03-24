package server

import (
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/go-chi/chi/v5"
	"github.com/gookit/slog"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func testHandler(t *testing.T, r chi.Router, method, path string, statusCode int, response string) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != statusCode {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, statusCode)
	}

	if response != "skip" && rr.Body.String() != response {
		t.Errorf("handler returned wrong response: got %v want %v",
			rr.Body.String(), response)
	}
}

func TestClearStorage(t *testing.T) {
	stg := storages.NewMemStorage()
	s := New(stg, slog.New())
	r := chi.NewRouter()
	r.NotFound(s.notFoundHandler)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", s.writeMetricHandler)
	r.Get("/", s.showAllMetricHandler)
	r.Get("/value/{metricType}", s.showMetricTypeHandler)
	r.Get("/value/{metricType}/{metricName}", s.showMetricNameHandlers)

	t.Run("test clear storage", func(t *testing.T) {
		testHandler(t, r, http.MethodGet, "/", http.StatusNotFound, "")
	})
}

func TestServerHandlers(t *testing.T) {
	stg := storages.NewMemStorage()
	stg.AddMetric("type1", metrics.NewCounter())
	stg.AddMetric("type100", metrics.NewCounter())
	s := New(stg, slog.New())

	r := chi.NewRouter()
	r.NotFound(s.notFoundHandler) // H
	r.Post("/update/{metricType}/{metricName}/{metricValue}", s.writeMetricHandler)
	r.Get("/", s.showAllMetricHandler)
	r.Get("/value/{metricType}", s.showMetricTypeHandler)
	r.Get("/value/{metricType}/{metricName}", s.showMetricNameHandlers)

	tests := []struct {
		name       string
		router     chi.Router
		method     string
		path       string
		statusCode int
		response   string
	}{

		{"writeMetricHandler", r, http.MethodPost, "/update/type1/name1/10", http.StatusOK, ""},
		{"writeMetricHandler", r, http.MethodPost, "/update/type100/name1/10", http.StatusOK, ""},

		{"writeMetricHandler", r, http.MethodPost, "/update/type1/", http.StatusNotFound, ""},
		{"writeMetricHandler", r, http.MethodPost, "/update/type23/name1/10/10", http.StatusNotFound, ""},
		{"writeMetricHandler", r, http.MethodPost, "/type1/name1/10", http.StatusNotFound, ""},
		{"writeMetricHandler", r, http.MethodPost, "/update/", http.StatusNotFound, ""},

		{"showAllMetricHandler", r, http.MethodGet, "/", http.StatusOK, "skip"},
		{"showMetricTypeHandler", r, http.MethodGet, "/value/type1", http.StatusOK, "type1:\n\tname1: 10\n"},
		{"showMetricNameHandlers", r, http.MethodGet, "/value/type1/name1", http.StatusOK, "10"},

		{"showMetricNameHandlersNotFound", r, http.MethodGet, "/value/not/name1", http.StatusNotFound, "metric type not not found\n"},
		{"showMetricTypeHandlersNotFound", r, http.MethodGet, "/value/type2", http.StatusNotFound, "metric type type2 not found\n"},
		{"showMetricNameHandlersNotFound", r, http.MethodGet, "/value/type1/name2", http.StatusNotFound, ""},
		{"notFoundHandler", r, http.MethodGet, "/nonexistentpath", http.StatusNotFound, ""},
		{"showMetricTypeHandlersNotFound", r, http.MethodGet, "/value/nonexistenttype", http.StatusNotFound, "metric type nonexistenttype not found\n"},

		{"writeMetricHandlersBadRequest", r, http.MethodPost, "/update/type1/name1/invalidValue", http.StatusBadRequest, "failed to process metric: metric value is not int\n"},
		{"writeMetricHandler", r, http.MethodPost, "/update/type23/name1/10", http.StatusBadRequest, "metric type type23 not found\n"},
		{"writeMetricHandler", r, http.MethodPost, "/", http.StatusMethodNotAllowed, ""},
		{"methodNotAllowedHandler", r, http.MethodPut, "/", http.StatusMethodNotAllowed, ""},
		{"writeMetricHandler", r, http.MethodConnect, "/", http.StatusMethodNotAllowed, ""},
		{"methodNotAllowedHandler", r, http.MethodDelete, "/", http.StatusMethodNotAllowed, ""},
		{"writeMetricHandler", r, http.MethodHead, "/", http.StatusMethodNotAllowed, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHandler(t, tt.router, tt.method, tt.path, tt.statusCode, tt.response)
		})
	}
}

func TestServer_Run(t *testing.T) {
	s := &Server{
		router: chi.NewRouter(),
	}
	quit := make(chan struct{})
	go func() {
		defer close(quit)
		err := s.Run(":8080")
		time.Sleep(1 * time.Millisecond)
		assert.NoError(t, err, "unexpected error")
	}()
}
