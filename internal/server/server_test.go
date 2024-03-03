package server

import (
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
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
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Body.String(), response)
	}
}

func TestServerHandlers(t *testing.T) {
	stg := storages.NewMemStorage()
	stg.AddMetric("type1", metrics.NewCounter())
	stg.AddMetric("type100", metrics.NewCounter())
	s := New(stg)

	r := chi.NewRouter()
	r.NotFound(s.notFoundHandler) // H
	r.Post("/update/{metricType}/{metricName}/{metricValue}", s.writeMetricHandlers)
	r.Get("/", s.showAllMetricHandlers)
	r.Get("/value/{metricType}", s.showMetricTypeHandlers)
	r.Get("/value/{metricType}/{metricName}", s.showMetricNameHandlers)

	tests := []struct {
		name       string
		router     chi.Router
		method     string
		path       string
		statusCode int
		response   string
	}{

		{"writeMetricHandlers", r, http.MethodPost, "/update/type1/name1/10", http.StatusOK, ""},
		{"writeMetricHandlers", r, http.MethodPost, "/update/type100/name1/10", http.StatusOK, ""},
		{"showAllMetricHandlers", r, http.MethodGet, "/", http.StatusOK, "skip"},
		{"showMetricTypeHandlers", r, http.MethodGet, "/value/type1", http.StatusOK, "type1:\n\tname1: 10\n"},
		{"showMetricNameHandlers", r, http.MethodGet, "/value/type1/name1", http.StatusOK, "10"},
		{"showMetricTypeHandlersNotFound", r, http.MethodGet, "/value/type2", http.StatusNotFound, "metric type type2 not found\n"},
		{"showMetricNameHandlersNotFound", r, http.MethodGet, "/value/type1/name2", http.StatusNotFound, ""},
		{"writeMetricHandlers", r, http.MethodPost, "/update/type23/name1/10", http.StatusBadRequest, "metric type type23 not found\n"},
		{"writeMetricHandlers", r, http.MethodPost, "/update/type23/name1/10/10", http.StatusNotFound, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHandler(t, tt.router, tt.method, tt.path, tt.statusCode, tt.response)
		})
	}
}
