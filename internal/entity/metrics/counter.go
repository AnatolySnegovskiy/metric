package metrics

import (
	"errors"
	"strconv"
)

type Counter struct {
	list map[string]int64
}

func (c *Counter) Process(name string, data string) error {
	intValue, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return errors.New("metric value is not int")
	}

	c.list[name] += intValue
	return nil
}

func NewCounter() *Counter {
	return &Counter{
		list: make(map[string]int64),
	}
}
