package repositories

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/storages/clients"
	"strings"
)

type CounterRepo struct {
	pg *clients.Postgres
}

func NewCounterRepo(pg *clients.Postgres) *CounterRepo {
	cr := &CounterRepo{
		pg: pg,
	}

	return cr
}

func (c *CounterRepo) GetItem(name string) (int, error) {
	var value int
	err := c.pg.QueryRow("SELECT value FROM counter WHERE name = $1", name).Scan(&value)
	return value, err
}

func (c *CounterRepo) GetList() (map[string]float64, error) {
	rows, err := c.pg.Query("SELECT * FROM counter")

	if err != nil {
		return nil, err
	}
	items := make(map[string]float64)
	for rows.Next() {
		var name string
		var value int
		_ = rows.Scan(&name, &value)
		items[name] = float64(value)
	}

	return items, nil
}

func (c *CounterRepo) AddMetric(name string, value int) error {
	_, err := c.pg.Exec("INSERT INTO counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2", name, value)
	return err
}

func (c *CounterRepo) AddMetrics(metrics map[string]float64) error {
	var valueStrings []string
	var valueArgs []interface{}
	i := 1
	for name, value := range metrics {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i, i+1))
		valueArgs = append(valueArgs, name, int(value))
		i += 2
	}
	query := fmt.Sprintf("INSERT INTO counter (name, value) VALUES %s ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value", strings.Join(valueStrings, ","))
	_, err := c.pg.Exec(query, valueArgs...)
	return err
}
