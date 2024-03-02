package services

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"reflect"
	"runtime"
	"time"
)

type Storage interface {
	GetMetricType(metricType string) (storages.EntityMetric, error)
	AddMetric(metricType string, metric storages.EntityMetric)
	GetList() map[string]storages.EntityMetric
}

var m runtime.MemStats
var runtimeEntityArray = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse",
	"HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse",
	"MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
	"StackSys", "Sys", "TotalAlloc"}

func UpdateStoragePeriodically(ticker <-chan time.Time, storage Storage) {
	fmt.Println("Запуск обновления метрик")

	for range ticker {
		runtime.ReadMemStats(&m)

		gauge, err := storage.GetMetricType("gauge")
		if err != nil {
			fmt.Println("Ошибка получения метрики:", err)
		}

		counter, err := storage.GetMetricType("counter")
		if err != nil {
			fmt.Println("Ошибка получения метрики:", err)
		}

		if counter.Process("PollCount", "1") != nil {
			fmt.Println("Ошибка обработки метрики:", err)
		}

		if gauge.Process("RandomValue", fmt.Sprintf("%v", time.Now().UnixNano())) != nil {
			fmt.Println("Ошибка обработки метрики:", err)
		}

		for _, field := range runtimeEntityArray {
			if gauge.Process(field, fmt.Sprintf("%v", reflect.ValueOf(m).FieldByName(field))) != nil {
				fmt.Println("Ошибка обработки метрики:", err)
			}
		}

		fmt.Println("Метрика обработана")
	}
}
