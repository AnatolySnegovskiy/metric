package server

import (
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHandler(t *testing.T, r chi.Router, method, path string, statusCode int) {
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
}

func TestServerHandlers(t *testing.T) {
	stg := storages.NewMemStorage()
	stg.AddMetric("type1", metrics.NewCounter())
	s := New(stg)

	r := chi.NewRouter()
	r.NotFound(s.notFoundHandler) // H
	r.Post("/update/{metricType}/{metricName}/{metricValue}", s.writeMetricHandlers)
	r.Get("/", s.showAllMetricHandlers)
	r.Get("/{metricType}", s.showMetricTypeHandlers)
	r.Get("/{metricType}/{metricName}", s.showMetricNameHandlers)

	tests := []struct {
		name       string
		router     chi.Router
		method     string
		path       string
		statusCode int
	}{
		{"writeMetricHandlers", r, http.MethodPost, "/update/type1/name1/10", http.StatusOK},
		{"showAllMetricHandlers", r, http.MethodGet, "/", http.StatusOK},
		{"showMetricTypeHandlers", r, http.MethodGet, "/type1", http.StatusOK},
		{"showMetricNameHandlers", r, http.MethodGet, "/type1/name1", http.StatusOK},
		{"showMetricTypeHandlersNotFound", r, http.MethodGet, "/type2", http.StatusNotFound},
		{"showMetricNameHandlersNotFound", r, http.MethodGet, "/type1/name2", http.StatusNotFound},
		{"writeMetricHandlers", r, http.MethodPost, "/update/type23/name1/10", http.StatusBadRequest},
		{"writeMetricHandlers", r, http.MethodPost, "/update/type23/name1/10/10", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHandler(t, tt.router, tt.method, tt.path, tt.statusCode)
		})
	}
}
