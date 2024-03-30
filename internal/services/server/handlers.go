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
	metricDTO, err := getMetricDto(req)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	storage := s.storage
	metric, err := storage.GetMetricType(metricDTO.MType)

	if err != nil {
		http.Error(rw, fmt.Sprintf("metric type %s not found", metricDTO.MType), http.StatusNotFound)
		return
	}

	var value string
	if metricDTO.Delta != nil {
		value = strconv.FormatInt(*metricDTO.Delta, 10)
	} else if metricDTO.Value != nil {
		value = strconv.FormatFloat(*metricDTO.Value, 'f', -1, 64)
	}

	if value == "" {
		http.Error(rw, fmt.Sprintf("failed to process Value and Delta is empty"), http.StatusNotFound)
		return
	}

	if err := metric.Process(metricDTO.ID, value); err != nil {
		http.Error(rw, fmt.Sprintf("failed to process metric: %s", err.Error()), http.StatusBadRequest)
		return
	}
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

func (s *Server) showPostMetricHandler(rw http.ResponseWriter, req *http.Request) {
	metricDTO, err := getMetricDto(req)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	metricType := metricDTO.MType
	metricName := metricDTO.ID

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

	if metricDTO.MType == "gauge" {
		metricDTO.Value = &metric
	} else {
		val := int64(metric)
		metricDTO.Delta = &val
	}

	json, _ := easyjson.Marshal(metricDTO)
	fmt.Fprintf(rw, "%v", string(json))
}

func (s *Server) notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

func getMetricDto(req *http.Request) (*dto.Metrics, error) {
	metricDTO := &dto.Metrics{}
	rawBytes, err := io.ReadAll(req.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read body: %s", err.Error())
	}

	if err := easyjson.Unmarshal(rawBytes, metricDTO); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %s", err.Error())
	}

	return metricDTO, nil
}
