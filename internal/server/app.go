package server

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"net/http"
	"strings"
)

type App struct {
	storage *memstorage.MemStorage
}

func New() *App {
	return &App{
		storage: memstorage.New(),
	}
}

func (app *App) Metric(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(res, "Method not allowed", http.StatusBadRequest)
		return
	}

	elements := strings.Split(req.URL.Path, "/")

	if len(elements) < 5 || len(elements) > 5 {
		http.Error(res, "Invalid request", http.StatusNotFound)
		return
	}

	metricType := elements[2]
	metricName := elements[3]
	metricValue := elements[4]

	if metricName == "" {
		http.Error(res, "Metric name is required", http.StatusNotFound)
		return
	}

	storage := app.storage
	metric, err := storage.GetMetricType(metricType)

	if err != nil {
		http.Error(res, fmt.Sprintf("Metric type %s not found", metricType), http.StatusBadRequest)
		return
	}

	err = metric.Process(metricName, metricValue)

	if err != nil {
		http.Error(res, fmt.Sprintf("Failed to process metric: %s", err.Error()), http.StatusBadRequest)
		return
	} else {
		storage.Log()
	}
}
