package metrics

import (
	"errors"
	"strconv"
)

type Gauge struct {
	items map[string]float64
}

func (g *Gauge) Process(name string, data string) error {
	floatValue, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return errors.New("metric value is not float64")
	}

	g.items[name] = floatValue
	return nil
}

func (g *Gauge) GetList() map[string]float64 { return g.items }

func NewGauge() *Gauge {
	return &Gauge{
		items: make(map[string]float64),
	}
}
