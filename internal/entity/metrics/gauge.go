package metrics

import (
	"context"
	"errors"
	"strconv"

	"github.com/AnatolySnegovskiy/metric/internal/repositories"
)

type Gauge struct {
	Items map[string]float64
	rep   *repositories.GaugeRepo
}

func (g *Gauge) Process(ctx context.Context, name string, data string) error {
	floatValue, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return errors.New("metric value is not float64")
	}

	g.Items[name] = floatValue

	if g.rep != nil {
		return g.rep.AddMetric(ctx, name, g.Items[name])
	}

	return nil
}

func (g *Gauge) ProcessMassive(ctx context.Context, data map[string]float64) error {
	for name, value := range data {
		g.Items[name] = value
	}

	if g.rep != nil {
		return g.rep.AddMetrics(ctx, g.Items)
	}

	return nil
}

func (g *Gauge) GetList(ctx context.Context) (map[string]float64, error) {
	if g.rep != nil {
		items, err := g.rep.GetList(ctx)
		if err != nil {
			return nil, err
		}
		g.Items = items
	}
	return g.Items, nil
}

func NewGauge(rep *repositories.GaugeRepo) *Gauge {
	return &Gauge{
		Items: make(map[string]float64),
		rep:   rep,
	}
}
