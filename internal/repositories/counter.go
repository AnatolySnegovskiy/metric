package repositories

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/storages/clients"
	"github.com/jackc/pgx/v5"
	"strings"
)

type CounterRepo struct {
	pg *clients.Postgres
}

func NewCounterRepo(pg *clients.Postgres) (*CounterRepo, error) {
	cr := &CounterRepo{
		pg: pg,
	}

	if err := cr.makeTable(); err != nil {
		return nil, err
	}

	return cr, nil
}

func (c *CounterRepo) makeTable() error {
	_, err := c.pg.Exec("CREATE TABLE IF NOT EXISTS counter (name varchar(100) PRIMARY KEY, value int8)")
	return err
}

func (c *CounterRepo) GetItem(name string) (int, error) {
	var value int
	err := c.pg.QueryRow("SELECT value FROM counter WHERE name = $1", name).Scan(&value)
	return value, err
}

func (c *CounterRepo) GetList() (pgx.Rows, error) {
	return c.pg.Query("SELECT * FROM counter")
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
