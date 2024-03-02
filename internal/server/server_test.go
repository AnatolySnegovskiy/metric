package server_test

import (
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AnatolySnegovskiy/metric/internal/server"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
)

func TestServer_HandleMetrics(t *testing.T) {
	mockStorage := storages.NewMemStorage()
	mockStorage.AddMetric("gauge", metrics.NewGauge())
	s := server.New(mockStorage)
	var req *http.Request
	var rr *httptest.ResponseRecorder

	tests := []struct {
		name   string
		method string
		value  string
		want   int
	}{
		{name: "ok", method: http.MethodPost, value: "/update/gauge/testName/10", want: http.StatusOK},
		{name: "not found", method: http.MethodPost, value: "/update/", want: http.StatusNotFound},
		{name: "not found", method: http.MethodPost, value: "/update/test", want: http.StatusNotFound},
		{name: "not found", method: http.MethodPost, value: "/update/gauge2/testName/10/test", want: http.StatusNotFound},
		{name: "bad", method: http.MethodPost, value: "/update/gauge2/testName/10", want: http.StatusBadRequest},
		{name: "bad", method: http.MethodGet, value: "/update/gauge2/testName/10", want: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req = httptest.NewRequest("POST", tt.value, nil)
			rr = httptest.NewRecorder()
			s.HandleMetrics(rr, req)
			if status := rr.Code; status != tt.want {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
			}
		})
	}
}
