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
		g.rep.AddMetric(name, g.Items[name])
	}

	return nil
}

func (g *Gauge) GetList() map[string]float64 {
	if g.rep != nil {
		rows := g.rep.GetList()
		for rows.Next() {
			var name string
			var value float64
			_ = rows.Scan(&name, &value)
			g.Items[name] = value
		}
	}
	return g.Items
}

func NewGauge(rep *repositories.GaugeRepo) *Gauge {
	return &Gauge{
		Items: make(map[string]float64),
		rep:   rep,
	}
}
