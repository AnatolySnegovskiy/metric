package metrics

import (
	"context"
	"errors"
	"strconv"

	"github.com/AnatolySnegovskiy/metric/internal/repositories"
)

type Counter struct {
	Items map[string]float64
	rep   *repositories.CounterRepo
}

func (c *Counter) Process(ctx context.Context, name string, data string) error {
	intValue, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return errors.New("metric value is not int")
	}

	c.Items[name] += float64(intValue)

	if c.rep != nil {
		return c.rep.AddMetric(ctx, name, int(c.Items[name]))
	}

	return nil
}

func (c *Counter) ProcessMassive(ctx context.Context, data map[string]float64) error {
	for name, value := range data {
		c.Items[name] += value
	}

	if c.rep != nil {
		return c.rep.AddMetrics(ctx, c.Items)
	}

	return nil
}

func (c *Counter) GetList(ctx context.Context) (map[string]float64, error) {
	if c.rep != nil {
		items, err := c.rep.GetList(ctx)
		if err != nil {
			return nil, err
		}
		c.Items = items
	}
	return c.Items, nil
}

func NewCounter(rep *repositories.CounterRepo) *Counter {
	return &Counter{
		Items: make(map[string]float64),
		rep:   rep,
	}
}
