package repositories

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/storages/clients"
	"github.com/jackc/pgx/v5"
	"strings"
)

type GaugeRepo struct {
	db *clients.Postgres
}

func NewGaugeRepo(db *clients.Postgres) *GaugeRepo {
	gr := &GaugeRepo{
		db: db,
	}
	gr.makeTable()
	return gr
}

func (g *GaugeRepo) makeTable() {
	_, _ = g.db.Exec("CREATE TABLE IF NOT EXISTS guage (name varchar(100) PRIMARY KEY, value DOUBLE PRECISION)")
}

func (g *GaugeRepo) getItem(name string) float64 {
	var value float64
	_ = g.db.QueryRow("SELECT value FROM guage WHERE name = $1", name).Scan(&value)
	return value
}

func (g *GaugeRepo) GetList() pgx.Rows {
	rows, _ := g.db.Query("SELECT * FROM guage")
	return rows
}

func (g *GaugeRepo) AddMetric(name string, value float64) {
	_, err := g.db.Exec("INSERT INTO guage (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2", name, value)
	if err != nil {
		fmt.Println("GaugeRepo.AddMetric " + err.Error())
	}
}

func (g *GaugeRepo) AddMetrics(metrics map[string]float64) {
	var valueStrings []string
	var valueArgs []interface{}
	i := 1
	for name, value := range metrics {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i, i+1))
		valueArgs = append(valueArgs, name, value)
		i += 2
	}
	query := fmt.Sprintf("INSERT INTO gauge (name, value) VALUES %s ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value", strings.Join(valueStrings, ","))
	_, _ = g.db.Exec(query, valueArgs...)
}
