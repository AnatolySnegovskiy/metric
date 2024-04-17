package repositories

import (
	"context"
	"github.com/AnatolySnegovskiy/metric/internal/storages/clients"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestCounterRepo_Test(t *testing.T) {
	testCases := []struct {
		name   string
		expect func(mock pgxmock.PgxPoolIface)
		check  func(mockDB *clients.Postgres)
	}{
		{
			name: "NewCounterRepo",
			expect: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta("CREATE TABLE IF NOT EXISTS counter (name varchar(100) PRIMARY KEY, value int8)")).
					WillReturnResult(pgxmock.NewResult("CREATE", 1))
			},
			check: func(mockDB *clients.Postgres) {
				r := NewCounterRepo(mockDB)
				assert.NotNil(t, r)
			},
		},
		{
			name: "GetItem",
			expect: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT value FROM counter WHERE name = $1")).
					WithArgs("test").
					WillReturnRows(pgxmock.NewRows([]string{"value"}).AddRow(100))
			},
			check: func(mockDB *clients.Postgres) {
				cr := &CounterRepo{
					pg: mockDB,
				}
				actual, err := cr.GetItem("test")
				assert.NoError(t, err, "GetItem", err)
				assert.Equal(t, 100, actual)
			},
		},
		{
			name: "GetList",
			expect: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM counter")).
					WillReturnRows(pgxmock.NewRows([]string{"name", "value"}).AddRow("test", 100))
			},
			check: func(mockDB *clients.Postgres) {
				cr := &CounterRepo{
					pg: mockDB,
				}
				actual, err := cr.GetList()
				assert.NoError(t, err, "GetList", err)

				var name string
				var value float64
				for k, v := range actual {
					name = k
					value = v
					break
				}
				assert.Equal(t, "test", name)
				assert.Equal(t, float64(100), value)
			},
		},
		{
			name: "AddMetric",
			expect: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2")).
					WithArgs("test", 100).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			check: func(mockDB *clients.Postgres) {
				cr := &CounterRepo{
					pg: mockDB,
				}
				err := cr.AddMetric("test", 100)
				assert.NoError(t, err, "AddMetric", err)
			},
		},
		{
			name: "AddMetrics",
			expect: func(mock pgxmock.PgxPoolIface) {
				var valueArgs []interface{}
				valueArgs = append(valueArgs, "test", 500)
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value")).
					WithArgs(valueArgs...).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			check: func(mockDB *clients.Postgres) {
				cr := &CounterRepo{
					pg: mockDB,
				}
				err := cr.AddMetrics(map[string]float64{"test": 500})
				assert.NoError(t, err, "AddMetrics", err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			testCase.expect(mock)
			mockDB, _ := clients.NewPostgres(context.Background(), mock)
			testCase.check(mockDB)
		})
	}
}

func TestGaugeRepo_Test(t *testing.T) {
	testCases := []struct {
		name   string
		expect func(mock pgxmock.PgxPoolIface)
		check  func(mockDB *clients.Postgres)
	}{
		{
			name: "NewGaugeRepo",
			expect: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta("CREATE TABLE IF NOT EXISTS gauge (name varchar(100) PRIMARY KEY, value DOUBLE PRECISION)")).
					WillReturnResult(pgxmock.NewResult("CREATE", 1))
			},
			check: func(mockDB *clients.Postgres) {
				r := NewGaugeRepo(mockDB)
				assert.NotNil(t, r)
			},
		},
		{
			name: "GetItem",
			expect: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT value FROM gauge WHERE name = $1")).
					WithArgs("test").
					WillReturnRows(pgxmock.NewRows([]string{"value"}).AddRow(float64(100)))
			},
			check: func(mockDB *clients.Postgres) {
				cr := &GaugeRepo{
					pg: mockDB,
				}
				actual, err := cr.GetItem("test")
				assert.NoError(t, err, "GetItem", err)
				assert.Equal(t, float64(100), actual)
			},
		},
		{
			name: "GetList",
			expect: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM gauge")).
					WillReturnRows(pgxmock.NewRows([]string{"name", "value"}).AddRow("test", float64(100)))
			},
			check: func(mockDB *clients.Postgres) {
				cr := &GaugeRepo{
					pg: mockDB,
				}
				actual, err := cr.GetList()
				assert.NoError(t, err, "GetList", err)

				var name string
				var value float64
				for k, v := range actual {
					name = k
					value = v
					break
				}
				assert.Equal(t, "test", name)
				assert.Equal(t, float64(100), value)
			},
		},
		{
			name: "AddMetric",
			expect: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO gauge (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2")).
					WithArgs("test", float64(100)).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			check: func(mockDB *clients.Postgres) {
				cr := &GaugeRepo{
					pg: mockDB,
				}
				err := cr.AddMetric("test", float64(100))
				assert.NoError(t, err, "AddMetric", err)
			},
		},
		{
			name: "AddMetrics",
			expect: func(mock pgxmock.PgxPoolIface) {
				var valueArgs []interface{}
				valueArgs = append(valueArgs, "test", 500.50)
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO gauge (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value")).
					WithArgs(valueArgs...).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			check: func(mockDB *clients.Postgres) {
				cr := &GaugeRepo{
					pg: mockDB,
				}
				err := cr.AddMetrics(map[string]float64{"test": 500.50})
				assert.NoError(t, err, "AddMetrics", err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			testCase.expect(mock)
			mockDB, _ := clients.NewPostgres(context.Background(), mock)
			testCase.check(mockDB)
		})
	}
}
