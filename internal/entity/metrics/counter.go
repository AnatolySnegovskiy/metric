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
		c.rep.AddMetric(name, c.Items[name])
	}

	return nil
}

func (c *Counter) GetList() map[string]float64 {
	if c.rep != nil {
		rows := c.rep.GetList()
		for rows.Next() {
			var name string
			var value float64
			_ = rows.Scan(&name, &value)
			c.Items[name] = value
		}
	}
	return c.Items
}

func NewCounter(rep *repositories.CounterRepo) *Counter {
	return &Counter{
		Items: make(map[string]float64),
		rep:   rep,
	}
}
