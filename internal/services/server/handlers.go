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

	if err := metric.Process(req.Context(), metricName, metricValue); err != nil {
		http.Error(rw, fmt.Sprintf("failed to process metric: %s", err.Error()), http.StatusBadRequest)
		return
	}
}

func (s *Server) writePostMetricHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	metricDTO, err := getMetricDto(req)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "%v", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	storage := s.storage
	metric, err := storage.GetMetricType(metricDTO.MType)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "%v", fmt.Sprintf(`{"error":"metric type %s not found"}`, metricDTO.MType))
		return
	}

	var value string
	if metricDTO.Delta != nil {
		value = strconv.FormatInt(*metricDTO.Delta, 10)
	} else if metricDTO.Value != nil {
		value = strconv.FormatFloat(*metricDTO.Value, 'f', -1, 64)
	}

	if value == "" {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "%v", `{"error":"failed to process Value and Delta is empty"}`)
		return
	}

	_ = metric.Process(req.Context(), metricDTO.ID, value)
	json, _ := easyjson.Marshal(metricDTO)
	fmt.Fprintf(rw, "%v", string(json))
}

func (s *Server) showAllMetricHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	stgList := s.storage.GetList()

	if len(stgList) == 0 {
		s.notFoundHandler(rw, req)
		return
	}

	for storageType, storage := range stgList {
		list, err := storage.GetList(req.Context())
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to get list of metrics: %s", err.Error()), http.StatusInternalServerError)
			continue
		}

		fmt.Fprintf(rw, "%s:\n", storageType)
		for metricName, metric := range list {
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

	list, err := storage.GetList(req.Context())
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to get list of metrics: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(rw, "%s:\n", metricType)

	for metricName, metric := range list {
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

	list, err := storage.GetList(req.Context())
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to get list of metrics: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	metric := list[metricName]

	if metric == 0 {
		s.notFoundHandler(rw, req)
		return
	}

	fmt.Fprintf(rw, "%v", metric)
}

func (s *Server) showPostMetricHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	metricDTO, err := getMetricDto(req)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "%v", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	metricType := metricDTO.MType
	metricName := metricDTO.ID

	storage, err := s.storage.GetMetricType(metricType)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rw, "%v", fmt.Sprintf(`{"error":"metric type %s not found"}`, metricType))
		return
	}

	list, err := storage.GetList(req.Context())
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to get list of metrics: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	metric, ok := list[metricName]

	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rw, "%v", fmt.Sprintf(`{"error":"metric %s not found"}`, metricName))
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
	rawBytes, _ := io.ReadAll(req.Body)

	if err := easyjson.Unmarshal(rawBytes, metricDTO); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %s", err.Error())
	}

	return metricDTO, nil
}

func (s *Server) postgersPingHandler(writer http.ResponseWriter, _ *http.Request) {
	if s.dbIsOpen {
		writer.WriteHeader(http.StatusOK)
		return
	}

	writer.WriteHeader(http.StatusInternalServerError)
}

func (s *Server) writeMassPostMetricHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	metricDTOCollection := &dto.MetricsCollection{}
	rawBytes, _ := io.ReadAll(req.Body)

	if err := easyjson.Unmarshal(rawBytes, metricDTOCollection); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "%v", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	storage := s.storage

	for _, metricDTO := range *metricDTOCollection {
		if metricDTO.Delta != nil {
			matric, err := storage.GetMetricType(metricDTO.MType)
			if err != nil {
				rw.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(rw, "%v", fmt.Sprintf(`{"error":"metric type %s not found"}`, metricDTO.MType))
				return
			}

			metricValue := *metricDTO.Delta
			if metricValue != 0 {
				if err := matric.ProcessMassive(req.Context(), map[string]float64{metricDTO.ID: float64(metricValue)}); err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(rw, "%v", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
					return
				}
			}
		} else if metricDTO.Value != nil {
			matric, err := storage.GetMetricType(metricDTO.MType)
			if err != nil {
				rw.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(rw, "%v", fmt.Sprintf(`{"error":"metric type %s not found"}`, metricDTO.MType))
				return
			}

			metricValue := *metricDTO.Value
			if metricValue != 0 {
				if err := matric.ProcessMassive(req.Context(), map[string]float64{metricDTO.ID: float64(metricValue)}); err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(rw, "%v", fmt.Sprintf(`{"error":"%s"}`, err.Error()))
					return
				}
			}
		}
	}

	json, _ := easyjson.Marshal(metricDTOCollection)
	fmt.Fprintf(rw, "%v", string(json))
}
