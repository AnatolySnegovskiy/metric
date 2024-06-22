package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/AnatolySnegovskiy/metric/internal/storages/clients"
)

type GaugeRepo struct {
	pg *clients.Postgres
}

func NewGaugeRepo(pg *clients.Postgres) *GaugeRepo {
	cr := &GaugeRepo{
		pg: pg,
	}

	return cr
}

func (g *GaugeRepo) GetItem(ctx context.Context, name string) (float64, error) {
	var value float64
	err := g.pg.QueryRow(ctx, "SELECT value FROM gauge WHERE name = $1", name).Scan(&value)
	return value, err
}

func (g *GaugeRepo) GetList(ctx context.Context) (map[string]float64, error) {
	rows, err := g.pg.Query(ctx, "SELECT * FROM gauge")

	if err != nil {
		return nil, err
	}

	items := make(map[string]float64)
	for rows.Next() {
		var name string
		var value float64
		_ = rows.Scan(&name, &value)
		items[name] = value
	}

	return items, nil
}

func (g *GaugeRepo) AddMetric(ctx context.Context, name string, value float64) error {
	_, err := g.pg.Exec(ctx, "INSERT INTO gauge (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2", name, value)
	return err
}

func (g *GaugeRepo) AddMetrics(ctx context.Context, metrics map[string]float64) error {
	var valueStrings []string
	var valueArgs []interface{}
	i := 1
	for name, value := range metrics {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i, i+1))
		valueArgs = append(valueArgs, name, value)
		i += 2
	}
	query := fmt.Sprintf("INSERT INTO gauge (name, value) VALUES %s ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value", strings.Join(valueStrings, ","))
	_, err := g.pg.Exec(ctx, query, valueArgs...)
	return err
}
