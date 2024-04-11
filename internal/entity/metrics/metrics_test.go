package metrics

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var counterProcessTests = []struct {
	name        string
	data        string
	expected    map[string]float64
	expectedErr error
}{
	{"test1", "10", map[string]float64{"test1": 10}, nil},
	{"test2", "20", map[string]float64{"test1": 10, "test2": 20}, nil},
	{"test3", "invalid", nil, errors.New("metric value is not int")},
}

func TestCounter_Process(t *testing.T) {
	counter := NewCounter(nil)

	for _, tt := range counterProcessTests {
		t.Run(tt.name, func(t *testing.T) {
			err := counter.Process(tt.name, tt.data)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "Expected error")
			} else {
				assert.NoError(t, err, "Expected no error")
				list := counter.GetList()
				assert.Equal(t, tt.expected, list, "Expected list %v, but got: %v", tt.expected, list)
			}
		})
	}
}

var gaugeProcessTests = []struct {
	name        string
	data        string
	expected    map[string]float64
	expectedErr error
}{
	{"test1", "10.5", map[string]float64{"test1": 10.5}, nil},
	{"test2", "20.5", map[string]float64{"test1": 10.5, "test2": 20.5}, nil},
	{"test3", "invalid", nil, errors.New("metric value is not float64")},
}

func TestGauge_Process(t *testing.T) {
	gauge := NewGauge(nil)

	for _, tt := range gaugeProcessTests {
		t.Run(tt.name, func(t *testing.T) {
			err := gauge.Process(tt.name, tt.data)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "Expected error")
			} else {
				assert.NoError(t, err, "Expected no error")
				list := gauge.GetList()
				assert.Equal(t, tt.expected, list, "Expected list %v, but got: %v", tt.expected, list)
			}
		})
	}
}
