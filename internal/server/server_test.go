package server_test

import (
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AnatolySnegovskiy/metric/internal/server"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
)

func TestServer_HandleMetricsOK(t *testing.T) {
	mockStorage := storages.NewMemStorage()
	mockStorage.AddMetric("gauge", metrics.NewGauge())
	s := server.New(mockStorage)
	req := httptest.NewRequest("POST", "/update/gauge/testName/10", nil)
	rr := httptest.NewRecorder()
	s.HandleMetrics(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestServer_HandleMetricsNotFound(t *testing.T) {
	mockStorage := storages.NewMemStorage()

	s := server.New(mockStorage)
	req := httptest.NewRequest("POST", "/update/", nil)
	rr := httptest.NewRecorder()
	s.HandleMetrics(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	req = httptest.NewRequest("POST", "/update/test", nil)
	rr = httptest.NewRecorder()
	s.HandleMetrics(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestServer_HandleMetricsBadRequest(t *testing.T) {
	mockStorage := storages.NewMemStorage()

	s := server.New(mockStorage)
	req := httptest.NewRequest("POST", "/update/gauge2/testName/10", nil)
	rr := httptest.NewRecorder()
	s.HandleMetrics(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	req = httptest.NewRequest("GET", "/update/gauge2/testName/10", nil)
	rr = httptest.NewRecorder()
	s.HandleMetrics(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	req = httptest.NewRequest("GET", "/update/gauge2/testName/10/TEST", nil)
	rr = httptest.NewRecorder()
	s.HandleMetrics(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
