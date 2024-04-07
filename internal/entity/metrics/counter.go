package metrics

import (
	"errors"
	"strconv"
)

type Counter struct {
	Items map[string]float64
}

func (c *Counter) Process(name string, data string) error {
	intValue, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return errors.New("metric value is not int")
	}

	c.Items[name] += float64(intValue)
	return nil
}

func (c *Counter) GetList() map[string]float64 {
	return c.Items
}

func NewCounter() *Counter {
	return &Counter{
		Items: make(map[string]float64),
	}
}
