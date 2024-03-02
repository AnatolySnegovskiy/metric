package metrics

import (
	"errors"
	"strconv"
)

type Counter struct {
	list map[string]int
}

func (c *Counter) Process(name string, data string) error {
	intValue, err := strconv.Atoi(data)
	if err != nil {
		return errors.New("metric value is not int")
	}

	c.list[name] += intValue
	return nil
}

func NewCounter() *Counter {
	c := &Counter{
		list: make(map[string]int),
	}
	return c
}
