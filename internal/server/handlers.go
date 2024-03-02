package server

import (
	"fmt"
	"net/http"
)

func (s *Server) writeMetricHandlers(rw http.ResponseWriter, req *http.Request) {
	metricType, metricName, metricValue, err := parseURL(req.URL.Path)

	if err != nil || metricType == "" || metricName == "" || metricValue == "" {
		http.Error(rw, "metric name is required", http.StatusNotFound)
		return
	}

	storage := s.storage
	metric, err := storage.GetMetricType(metricType)

	if err != nil {
		http.Error(rw, fmt.Sprintf("metric type %s not found", metricType), http.StatusBadRequest)
		return
	}

	if err := metric.Process(metricName, metricValue); err != nil {
		http.Error(rw, fmt.Sprintf("failed to process metric: %s", err.Error()), http.StatusBadRequest)
		return
	}
}

func (s *Server) showMetricHandlers(rw http.ResponseWriter, _ *http.Request) {
	for storageType, storage := range s.storage.GetList() {
		fmt.Fprintf(rw, "%s:\n", storageType)
		for metricName, metric := range storage.GetList() {
			fmt.Fprintf(rw, "\t%s: %v\n", metricName, metric)
		}
	}
}
