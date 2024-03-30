package server

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/services/dto"
	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
	"io"
	"net/http"
	"strconv"
)

func (s *Server) writeGetMetricHandler(rw http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")
	metricValue := chi.URLParam(req, "metricValue")

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

func (s *Server) writePostMetricHandler(rw http.ResponseWriter, req *http.Request) {
	someStruct := &dto.Metrics{}
	rawBytes, err := io.ReadAll(req.Body)

	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to read body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if err := easyjson.Unmarshal(rawBytes, someStruct); err != nil {
		http.Error(rw, fmt.Sprintf("failed to unmarshal body: %s", err.Error()), http.StatusBadRequest)
		return
	}
	storage := s.storage
	metric, err := storage.GetMetricType(someStruct.MType)

	var value string
	if someStruct.Delta != nil {
		value = strconv.FormatInt(*someStruct.Delta, 10)
	} else if someStruct.Value != nil {
		value = strconv.FormatFloat(*someStruct.Value, 'f', -1, 64)
	}

	if value == "" {
		http.Error(rw, fmt.Sprintf("failed to process Value and Delta is empty"), http.StatusNotFound)
		return
	}

	metric.Process(someStruct.ID, value)
}

func (s *Server) showAllMetricHandler(rw http.ResponseWriter, req *http.Request) {
	stgList := s.storage.GetList()

	if len(stgList) == 0 {
		s.notFoundHandler(rw, req)
		return
	}

	for storageType, storage := range stgList {
		fmt.Fprintf(rw, "%s:\n", storageType)
		for metricName, metric := range storage.GetList() {
			fmt.Fprintf(rw, "\t%s: %v\n", metricName, metric)
		}
	}
}

func (s *Server) showMetricTypeHandler(rw http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "metricType")

	storage, err := s.storage.GetMetricType(metricType)
	if err != nil {
		http.Error(rw, fmt.Sprintf("metric type %s not found", metricType), http.StatusNotFound)
		return
	}
	fmt.Fprintf(rw, "%s:\n", metricType)
	for metricName, metric := range storage.GetList() {
		fmt.Fprintf(rw, "\t%s: %v\n", metricName, metric)
	}
}

func (s *Server) showMetricNameHandlers(rw http.ResponseWriter, req *http.Request) {
	metricType := chi.URLParam(req, "metricType")
	metricName := chi.URLParam(req, "metricName")

	storage, err := s.storage.GetMetricType(metricType)
	if err != nil {
		http.Error(rw, fmt.Sprintf("metric type %s not found", metricType), http.StatusNotFound)
		return
	}

	metric := storage.GetList()[metricName]

	if metric == 0 {
		s.notFoundHandler(rw, req)
		return
	}

	fmt.Fprintf(rw, "%v", storage.GetList()[metricName])
}

func (s *Server) notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
