package metrics

import (
	"errors"
	"github.com/AnatolySnegovskiy/metric/internal/repositories"
	"strconv"
)

type Counter struct {
	Items map[string]float64
	rep   *repositories.CounterRepo
}

func (c *Counter) Process(name string, data string) error {
	intValue, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return errors.New("metric value is not int")
	}

	c.Items[name] += float64(intValue)

	if c.rep != nil {
		return c.rep.AddMetric(name, int(c.Items[name]))
	}

	return nil
}

func (c *Counter) ProcessMassive(data map[string]float64) error {
	c.Items = data

	if c.rep != nil {
		return c.rep.AddMetrics(data)
	}

	return nil
}

func (c *Counter) GetList() (map[string]float64, error) {
	if c.rep != nil {
		rows, err := c.rep.GetList()

		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var name string
			var value int
			_ = rows.Scan(&name, &value)
			c.Items[name] = float64(value)
		}
	}
	return c.Items, nil
}

func NewCounter(rep *repositories.CounterRepo) *Counter {
	return &Counter{
		Items: make(map[string]float64),
		rep:   rep,
	}
}
