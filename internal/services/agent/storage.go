package agent

import "github.com/AnatolySnegovskiy/metric/internal/storages"

type Storage interface {
	GetMetricType(metricType string) (storages.EntityMetric, error)
	AddMetric(metricType string, metric storages.EntityMetric)
	GetList() map[string]storages.EntityMetric
}
