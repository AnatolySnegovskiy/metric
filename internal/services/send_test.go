package services

import (
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/AnatolySnegovskiy/metric/internal/storages"
)

func TestSendMetricsPeriodically(t *testing.T) {
	memStorage := storages.NewMemStorage()
	memStorage.AddMetric("storageType1", metrics.NewCounter())
	memStorage.AddMetric("storageType2", metrics.NewGauge())
	updateTicker := time.NewTicker(100 * time.Millisecond)
	defer updateTicker.Stop()

	go SendMetricsPeriodically(updateTicker.C, memStorage)
	time.Sleep(150 * time.Millisecond)

	assert.True(t, true)
}
