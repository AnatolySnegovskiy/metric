package metrics_test

import (
	"testing"

	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
)

func TestGauge_Process(t *testing.T) {
	gauge := metrics.NewGauge()

	// Test case 1: valid data
	err := gauge.Process("test", "10.5")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	// Test case 2: invalid data
	err = gauge.Process("test", "invalid")
	if err == nil {
		t.Error("expected error, got no error")
	}
}

func TestCounter_Process(t *testing.T) {
	counter := metrics.NewCounter()

	// Test case 1: valid data
	err := counter.Process("test", "5")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	// Test case 2: invalid data
	err = counter.Process("test", "invalid")
	if err == nil {
		t.Error("expected error, got no error")
	}
}
