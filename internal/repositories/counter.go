package repositories

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/storages/clients"
	"github.com/jackc/pgx/v5"
)

type CounterRepo struct {
	db *clients.Postgres
}

func NewCounterRepo(db *clients.Postgres) *CounterRepo {
	cr := &CounterRepo{
		db: db,
	}
	cr.makeTable()
	return cr
}

func (c *CounterRepo) makeTable() {
	_, _ = c.db.Exec("CREATE TABLE IF NOT EXISTS counter (name varchar(100) PRIMARY KEY, value int)")
}

func (c *CounterRepo) GetItem(name string) float64 {
	var value float64
	_ = c.db.QueryRow("SELECT value FROM counter WHERE name = $1", name).Scan(&value)
	return value
}

func (c *CounterRepo) GetList() pgx.Rows {
	rows, _ := c.db.Query("SELECT * FROM counter")
	return rows
}

func (c *CounterRepo) AddMetric(name string, value float64) {
	_, err := c.db.Exec("INSERT INTO counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2", name, value)
	if err != nil {
		fmt.Println(" CounterRepo.AddMetric " + err.Error())
	}
}

func (c *CounterRepo) AddMetrics(metrics map[string]float64) {
	for name, value := range metrics {
		_, _ = c.db.Exec("INSERT INTO counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2", name, value)
	}
}
