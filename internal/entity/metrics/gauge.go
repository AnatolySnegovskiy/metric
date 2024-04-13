package metrics

import (
	"errors"
	"github.com/AnatolySnegovskiy/metric/internal/repositories"
	"strconv"
)

type Gauge struct {
	Items map[string]float64
	rep   *repositories.GaugeRepo
}

func (g *Gauge) Process(name string, data string) error {
	floatValue, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return errors.New("metric value is not float64")
	}

	g.Items[name] = floatValue

	if g.rep != nil {
		return g.rep.AddMetric(name, g.Items[name])
	}

	return nil
}

func (g *Gauge) ProcessMassive(data map[string]float64) error {
	g.Items = data

	if g.rep != nil {
		return g.rep.AddMetrics(data)
	}

	return nil
}

func (g *Gauge) GetList() (map[string]float64, error) {
	if g.rep != nil {
		rows, err := g.rep.GetList()

		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var name string
			var value float64
			_ = rows.Scan(&name, &value)
			g.Items[name] = value
		}
	}
	return g.Items, nil
}

func NewGauge(rep *repositories.GaugeRepo) *Gauge {
	return &Gauge{
		Items: make(map[string]float64),
		rep:   rep,
	}
}
