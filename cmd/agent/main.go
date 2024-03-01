package main

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"net/http"
	"strings"
)

func mainPage(res http.ResponseWriter, req *http.Request) {
	availableTypes := []string{"gauge", "counter"}

	elements := strings.Split(req.URL.Path, "/")
	metricType := elements[2]
	metricName := elements[3]
	metricValue := elements[4]

	if stringInArray(availableTypes, metricType) == false {
		http.Error(res, fmt.Sprintf("metric type %s Not Found", metricType), http.StatusBadRequest)
		return
	}

	if metricName == "" {
		http.Error(res, "need metric name", http.StatusNotFound)
		return
	}

	storage := storages.NewMetricStorage()
	storage.Metrics[metricType].Process(metricValue)

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, mainPage)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func stringInArray(arr []string, target string) bool {
	for _, element := range arr {
		if element == target {
			return true
		}
	}
	return false
}
