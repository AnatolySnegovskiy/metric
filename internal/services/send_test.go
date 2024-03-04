package services

import (
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendMetricsPeriodically(t *testing.T) {
	memStorage := storages.NewMemStorage()
	memStorage.AddMetric("storageType1", metrics.NewCounter())
	memStorage.AddMetric("storageType2", metrics.NewGauge())
	err := SendMetricsPeriodically(":8080", memStorage)

	assert.NoError(t, err)
}
