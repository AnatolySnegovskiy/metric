package server

import (
	"bytes"
	"encoding/json"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/services/dto"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/go-chi/chi/v5"
	"github.com/gookit/slog"
	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func testHandler(t *testing.T, r chi.Router, method, path string, statusCode int, response string, requestBody []byte, headers map[string]string) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Fatal(err)
	}

	if method == http.MethodPost && requestBody != nil {
		req, err = http.NewRequest(method, path, bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatal(err)
		}
	}

	for k, v := range headers {
		req.Header.Set(k, v)
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
	r.Post("/update/{metricType}/{metricName}/{metricValue}", s.writeGetMetricHandler)
	r.Get("/", s.showAllMetricHandler)
	r.Get("/value/{metricType}", s.showMetricTypeHandler)
	r.Get("/value/{metricType}/{metricName}", s.showMetricNameHandlers)

	t.Run("test clear storage", func(t *testing.T) {
		testHandler(t, r, http.MethodGet, "/", http.StatusNotFound, "", nil, nil)
	})
}

func TestServerHandlers(t *testing.T) {
	stg := storages.NewMemStorage()
	stg.AddMetric("gauge", metrics.NewGauge())
	stg.AddMetric("type1", metrics.NewCounter())
	stg.AddMetric("type100", metrics.NewCounter())
	stg.AddMetric("typePostData", metrics.NewCounter())
	stg.AddMetric("gaugeValue", metrics.NewGauge())
	stg.AddMetric("zero", metrics.NewGauge())
	s := New(stg, slog.New())

	r := chi.NewRouter()
	r.Use(s.logMiddleware, s.gzipResponseMiddleware, s.gzipRequestMiddleware)
	r.NotFound(s.notFoundHandler) // H
	r.With(s.JSONContentTypeMiddleware).Post("/update/", s.writePostMetricHandler)
	r.With(s.JSONContentTypeMiddleware).Post("/value/", s.showPostMetricHandler)
	r.Post("/update/{metricType}/{metricName}/{metricValue}", s.writeGetMetricHandler)
	r.Get("/", s.showAllMetricHandler)
	r.Get("/value/{metricType}", s.showMetricTypeHandler)
	r.Get("/value/{metricType}/{metricName}", s.showMetricNameHandlers)

	bodyMap := map[string][]byte{}

	bodyMap["unknown"], _ = easyjson.Marshal(dto.Metrics{
		MType: "unknown",
		ID:    "unknown",
		Delta: nil,
		Value: nil,
	})

	intValue := 10
	int64Ptr := new(int64)
	*int64Ptr = int64(intValue)

	floatValue := 10.10
	float64Ptr := new(float64)
	*float64Ptr = float64(floatValue)

	bodyMap["typePostData"], _ = easyjson.Marshal(dto.Metrics{
		MType: "typePostData",
		ID:    "test",
		Delta: int64Ptr,
		Value: float64Ptr,
	})

	bodyMap["typePostDataValue"], _ = easyjson.Marshal(dto.Metrics{
		MType: "gaugeValue",
		ID:    "test",
		Value: float64Ptr,
	})

	bodyMap["typePostDataZero"], _ = easyjson.Marshal(dto.Metrics{
		MType: "zero",
		ID:    "test",
	})

	bodyMap["typePostDataGauge"], _ = easyjson.Marshal(dto.Metrics{
		MType: "gauge",
		ID:    "test",
		Delta: int64Ptr,
		Value: float64Ptr,
	})

	bodyMap["getPostValue"], _ = easyjson.Marshal(dto.Metrics{
		MType: "typePostData",
		ID:    "test",
	})

	bodyMap["getPostValueGauge"], _ = easyjson.Marshal(dto.Metrics{
		MType: "gauge",
		ID:    "test",
	})

	tests := []struct {
		name        string
		router      chi.Router
		method      string
		path        string
		statusCode  int
		response    string
		requestBody []byte
		headers     map[string]string
	}{
		{"notFoundHandler", r, http.MethodPost, "/update/", http.StatusBadRequest, "{\"error\":\"metric type nonexistent not found\"}", []byte(`{"type":"nonexistent","id":"nonexistent"}`), map[string]string{"Content-Type": "application/json"}},
		{"failed to unmarshal", r, http.MethodPost, "/update/", http.StatusBadRequest, "{\"error\":\"failed to unmarshal body: parse error: expected { near offset 12 of 'metricName'\"}", []byte(`"metricName":"example_metric","timestamp":"invalid_timestamp_format"}`), map[string]string{"Content-Type": "application/json"}},
		{"failed to process", r, http.MethodPost, "/update/", http.StatusBadRequest, "{\"error\":\"failed to process Value and Delta is empty\"}", bodyMap["typePostDataZero"], map[string]string{"Content-Type": "application/json"}},
		{"writeGetMetricHandler1", r, http.MethodPost, "/update/", http.StatusOK, "skip", bodyMap["typePostData"], map[string]string{"Content-Type": "application/json"}},
		{"writeGetMetricHandler2", r, http.MethodPost, "/update/", http.StatusOK, "skip", bodyMap["typePostDataGauge"], map[string]string{"Content-Type": "application/json"}},
		{"writeGetMetricHandler3", r, http.MethodPost, "/update/", http.StatusOK, "skip", bodyMap["typePostDataValue"], map[string]string{"Content-Type": "application/json"}},

		{"writeGetMetricHandler4", r, http.MethodPost, "/value/", http.StatusOK, "{\"id\":\"test\",\"type\":\"typePostData\",\"delta\":10}", bodyMap["getPostValue"], map[string]string{"Content-Type": "application/json"}},

		{"writeGetMetricHandler5", r, http.MethodPost, "/value/", http.StatusOK, "{\"id\":\"test\",\"type\":\"gauge\",\"value\":10}", bodyMap["getPostValueGauge"], map[string]string{"Content-Type": "application/json"}},
		{"writeGetMetricHandler6", r, http.MethodPost, "/value/", http.StatusNotFound, "{\"error\":\"metric test not found\"}", bodyMap["typePostDataZero"], map[string]string{"Content-Type": "application/json"}},

		{"writeGetMetricHandler7", r, http.MethodPost, "/update/", http.StatusBadRequest, "{\"error\":\"metric type unknown not found\"}", bodyMap["unknown"], map[string]string{"Content-Type": "application/json"}},

		{"writeGetMetricHandler7", r, http.MethodPost, "/value/", http.StatusNotFound, "{\"error\":\"metric type unknown not found\"}", bodyMap["unknown"], map[string]string{"Content-Type": "application/json"}},
		{"failed to unmarshal", r, http.MethodPost, "/value/", http.StatusBadRequest, "{\"error\":\"failed to unmarshal body: parse error: expected { near offset 12 of 'metricName'\"}", []byte(`"metricName":"example_metric","timestamp":"invalid_timestamp_format"}`), map[string]string{"Content-Type": "application/json"}},

		{"writeGetMetricHandler8", r, http.MethodPost, "/update/type1/name1/10", http.StatusOK, "", nil, nil},
		{"writeGetMetricHandler9", r, http.MethodPost, "/update/type100/name1/10", http.StatusOK, "", nil, nil},

		{"writeGetMetricHandler11", r, http.MethodPost, "/update/type1/", http.StatusNotFound, "", nil, nil},
		{"writeGetMetricHandler12", r, http.MethodPost, "/update/type23/name1/10/10", http.StatusNotFound, "", nil, nil},
		{"writeGetMetricHandler13", r, http.MethodPost, "/type1/name1/10", http.StatusNotFound, "", nil, nil},

		{"showAllMetricHandler", r, http.MethodGet, "/", http.StatusOK, "skip", nil, nil},
		{"showAllMetricHandler", r, http.MethodGet, "/", http.StatusOK, "skip", nil, map[string]string{"Content-Type": "application/json", "accept": "application/json", "Accept-Encoding": "gzip, deflate", "Connection": "keep-alive"}},
		{"showMetricTypeHandler", r, http.MethodGet, "/value/type1", http.StatusOK, "type1:\n\tname1: 10\n", nil, nil},
		{"showMetricNameHandlers", r, http.MethodGet, "/value/type1/name1", http.StatusOK, "10", nil, nil},

		{"showMetricNameHandlersNotFound1", r, http.MethodGet, "/value/not/name1", http.StatusNotFound, "metric type not not found\n", nil, nil},
		{"showMetricTypeHandlersNotFound2", r, http.MethodGet, "/value/type2", http.StatusNotFound, "metric type type2 not found\n", nil, nil},
		{"showMetricNameHandlersNotFound3", r, http.MethodGet, "/value/type1/name2", http.StatusNotFound, "", nil, nil},
		{"notFoundHandler", r, http.MethodGet, "/nonexistentpath", http.StatusNotFound, "", nil, nil},
		{"showMetricTypeHandlersNotFound", r, http.MethodGet, "/value/nonexistenttype", http.StatusNotFound, "metric type nonexistenttype not found\n", nil, nil},

		{"writeMetricHandlersBadRequest", r, http.MethodPost, "/update/type1/name1/invalidValue", http.StatusBadRequest, "failed to process metric: metric value is not int\n", nil, nil},
		{"writeGetMetricHandler", r, http.MethodPost, "/update/type23/name1/10", http.StatusBadRequest, "metric type type23 not found\n", nil, nil},
		{"writeGetMetricHandler", r, http.MethodPost, "/", http.StatusMethodNotAllowed, "", nil, nil},
		{"methodNotAllowedHandler", r, http.MethodPut, "/", http.StatusMethodNotAllowed, "", nil, nil},
		{"writeGetMetricHandler", r, http.MethodConnect, "/", http.StatusMethodNotAllowed, "", nil, nil},
		{"methodNotAllowedHandler", r, http.MethodDelete, "/", http.StatusMethodNotAllowed, "", nil, nil},
		{"writeGetMetricHandler", r, http.MethodHead, "/", http.StatusMethodNotAllowed, "", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testHandler(t, tt.router, tt.method, tt.path, tt.statusCode, tt.response, tt.requestBody, tt.headers)
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

func TestLoadMetricsFromFile(t *testing.T) {
	sampleData := map[string]map[string]map[string]float64{
		"metricType1": {
			"metricName1": {
				"value1": 1.23,
			},
		},
	}
	sampleJSON, err := json.Marshal(sampleData)
	assert.NoError(t, err)

	projectDir, _ := os.Getwd()
	filePath := "tmp/metrics.json"

	absoluteFilePath := filepath.Join(projectDir, filePath)
	directory := filepath.Dir(absoluteFilePath)

	_ = os.MkdirAll(directory, os.ModePerm)
	file, _ := os.Create(absoluteFilePath)

	_, _ = file.Write(sampleJSON)
	defer os.Remove(absoluteFilePath)
	defer os.RemoveAll(directory)
	defer file.Close()

	m, err := loadMetricsFromFile(filePath)

	assert.NoError(t, err)
	assert.Equal(t, sampleData, m)
}
