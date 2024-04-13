package repositories

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/storages/clients"
	"github.com/jackc/pgx/v5"
	"strings"
)

type GaugeRepo struct {
	pg *clients.Postgres
}

func NewGaugeRepo(pg *clients.Postgres) (*GaugeRepo, error) {
	cr := &GaugeRepo{
		pg: pg,
	}

	if err := cr.makeTable(); err != nil {
		return nil, err
	}

	return cr, nil
}

func (g *GaugeRepo) makeTable() error {
	_, err := g.pg.Exec("CREATE TABLE IF NOT EXISTS gauge (name varchar(100) PRIMARY KEY, value DOUBLE PRECISION)")
	return err
}

func (g *GaugeRepo) GetItem(name string) (float64, error) {
	var value float64
	err := g.pg.QueryRow("SELECT value FROM gauge WHERE name = $1", name).Scan(&value)
	return value, err
}

func (g *GaugeRepo) GetList() (pgx.Rows, error) {
	return g.pg.Query("SELECT * FROM gauge")
}

func (g *GaugeRepo) AddMetric(name string, value float64) error {
	_, err := g.pg.Exec("INSERT INTO gauge (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2", name, value)
	return err
}

func (g *GaugeRepo) AddMetrics(metrics map[string]float64) error {
	var valueStrings []string
	var valueArgs []interface{}
	i := 1
	for name, value := range metrics {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i, i+1))
		valueArgs = append(valueArgs, name, value)
		i += 2
	}
	query := fmt.Sprintf("INSERT INTO gauge (name, value) VALUES %s ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value", strings.Join(valueStrings, ","))
	_, err := g.pg.Exec(query, valueArgs...)
	return err
}
