package services

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"reflect"
	"runtime"
	"time"
)

//go:generate mockgen -source=update.go -destination=mocks/update_mock.go -package=mocks

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

func UpdateStoragePeriodically(storage Storage) error {
	runtime.ReadMemStats(&m)

	gauge, err := storage.GetMetricType("gauge")
	if err != nil {
		return fmt.Errorf("error getting gauge: %w", err)
	}

	counter, err := storage.GetMetricType("counter")
	if err != nil {
		return fmt.Errorf("error getting counter: %w", err)
	}

	if counter.Process("PollCount", "1") != nil {
		return fmt.Errorf("error while processing field: PollCount")
	}
	if gauge.Process("RandomValue", fmt.Sprintf("%v", time.Now().UnixNano())) != nil {
		return fmt.Errorf("error while processing field: RandomValue")
	}

	for _, field := range runtimeEntityArray {
		if gauge.Process(field, fmt.Sprintf("%v", reflect.ValueOf(m).FieldByName(field))) != nil {
			return fmt.Errorf("error while processing field: %s", field)
		}
	}

	return nil
}
