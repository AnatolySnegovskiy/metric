package metrics

import (
	"context"
	"errors"
	"github.com/AnatolySnegovskiy/metric/internal/repositories"
	"github.com/AnatolySnegovskiy/metric/internal/storages/clients"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

var counterProcessTests = []struct {
	name        string
	data        string
	expected    map[string]float64
	expectedErr error
}{
	{"test1", "10", map[string]float64{"test1": 10}, nil},
	{"test2", "20", map[string]float64{"test1": 10, "test2": 20}, nil},
	{"test3", "invalid", nil, errors.New("metric value is not int")},
}

func TestCounter_Process(t *testing.T) {
	counter := NewCounter(nil)

	for _, tt := range counterProcessTests {
		t.Run(tt.name, func(t *testing.T) {
			err := counter.Process(tt.name, tt.data)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "Expected error")
			} else {
				assert.NoError(t, err, "Expected no error")
				list, err := counter.GetList()
				assert.Equal(t, tt.expected, list, "Expected list %v, but got: %v", tt.expected, list)
				assert.NoError(t, err, "Expected no error")
			}
		})
	}
}

var gaugeProcessTests = []struct {
	name        string
	data        string
	expected    map[string]float64
	expectedErr error
}{
	{"test1", "10.5", map[string]float64{"test1": 10.5}, nil},
	{"test2", "20.5", map[string]float64{"test1": 10.5, "test2": 20.5}, nil},
	{"test3", "invalid", nil, errors.New("metric value is not float64")},
}

func TestGauge_Process(t *testing.T) {
	gauge := NewGauge(nil)

	for _, tt := range gaugeProcessTests {
		t.Run(tt.name, func(t *testing.T) {
			err := gauge.Process(tt.name, tt.data)
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "Expected error")
			} else {
				assert.NoError(t, err, "Expected no error")
				list, err := gauge.GetList()
				assert.NoError(t, err, "Expected no error")
				assert.Equal(t, tt.expected, list, "Expected list %v, but got: %v", tt.expected, list)
			}
		})
	}
}

func TestGauge_getListErrorDB(t *testing.T) {

	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectExec(regexp.QuoteMeta("CREATE TABLE IF NOT EXISTS gauge (name varchar(100) PRIMARY KEY, value DOUBLE PRECISION)")).
		WillReturnResult(pgxmock.NewResult("CREATE", 1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM gauge")).
		WillReturnError(errors.New("db error"))
	mockDB, _ := clients.NewPostgres(context.Background(), mock)
	cr, _ := repositories.NewGaugeRepo(mockDB)
	gauge := NewGauge(cr)
	_, err = gauge.GetList()
	assert.Error(t, err)
}

func TestCounter_getListErrorDB(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	mock.ExpectExec(regexp.QuoteMeta("CREATE TABLE IF NOT EXISTS counter (name varchar(100) PRIMARY KEY, value int8)")).
		WillReturnResult(pgxmock.NewResult("CREATE", 1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM counter")).
		WillReturnError(errors.New("db error"))
	mockDB, _ := clients.NewPostgres(context.Background(), mock)
	cr, _ := repositories.NewCounterRepo(mockDB)
	counter := NewCounter(cr)
	_, err = counter.GetList()
	assert.Error(t, err)
}

func TestGauge_ProcessDB(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectExec(regexp.QuoteMeta("CREATE TABLE IF NOT EXISTS gauge (name varchar(100) PRIMARY KEY, value DOUBLE PRECISION)")).
		WillReturnResult(pgxmock.NewResult("CREATE", 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO gauge (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2")).
		WithArgs("test", float64(100)).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM gauge")).
		WillReturnRows(pgxmock.NewRows([]string{"name", "value"}).AddRow("test", float64(100)))

	mockDB, _ := clients.NewPostgres(context.Background(), mock)
	cr, _ := repositories.NewGaugeRepo(mockDB)
	gauge := NewGauge(cr)
	err = gauge.Process("test", "100")
	assert.NoError(t, err)
	list, err := gauge.GetList()
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{"test": 100}, list, "Expected list %v, but got: %v", map[string]float64{"test": 100}, list)
}

func TestCounter_ProcessDB(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectExec(regexp.QuoteMeta("CREATE TABLE IF NOT EXISTS counter (name varchar(100) PRIMARY KEY, value int8)")).
		WillReturnResult(pgxmock.NewResult("CREATE", 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = $2")).
		WithArgs("test", int(100)).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM counter")).
		WillReturnRows(pgxmock.NewRows([]string{"name", "value"}).AddRow("test", 100))

	mockDB, _ := clients.NewPostgres(context.Background(), mock)
	cr, _ := repositories.NewCounterRepo(mockDB)
	counter := NewCounter(cr)
	err = counter.Process("test", "100")
	assert.NoError(t, err)
	list, err := counter.GetList()
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{"test": 100}, list, "Expected list %v, but got: %v", map[string]float64{"test": 100}, list)
}

func TestGauge_ProcessMassiveDB(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectExec(regexp.QuoteMeta("CREATE TABLE IF NOT EXISTS gauge (name varchar(100) PRIMARY KEY, value DOUBLE PRECISION)")).
		WillReturnResult(pgxmock.NewResult("CREATE", 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO gauge (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value")).
		WithArgs("test", float64(500)).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM gauge")).
		WillReturnRows(pgxmock.NewRows([]string{"name", "value"}).AddRow("test", float64(500)))

	mockDB, _ := clients.NewPostgres(context.Background(), mock)
	cr, _ := repositories.NewGaugeRepo(mockDB)
	gauge := NewGauge(cr)
	err = gauge.ProcessMassive(map[string]float64{"test": 500})
	assert.NoError(t, err)
	list, err := gauge.GetList()
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{"test": 500}, list, "Expected list %v, but got: %v", map[string]float64{"test": 500}, list)
}

func TestCounter_ProcessMassiveDB(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	mock.ExpectExec(regexp.QuoteMeta("CREATE TABLE IF NOT EXISTS counter (name varchar(100) PRIMARY KEY, value int8)")).
		WillReturnResult(pgxmock.NewResult("CREATE", 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO counter (name, value) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value")).
		WithArgs("test", int(500)).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM counter")).
		WillReturnRows(pgxmock.NewRows([]string{"name", "value"}).AddRow("test", 500))

	mockDB, _ := clients.NewPostgres(context.Background(), mock)
	cr, _ := repositories.NewCounterRepo(mockDB)
	counter := NewCounter(cr)
	err = counter.ProcessMassive(map[string]float64{"test": 500})
	assert.NoError(t, err)
	list, err := counter.GetList()
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{"test": 500}, list, "Expected list %v, but got: %v", map[string]float64{"test": 500}, list)
}

func TestGauge_ProcessMassive(t *testing.T) {
	gauge := NewGauge(nil)
	err := gauge.ProcessMassive(map[string]float64{"test": 500})
	assert.NoError(t, err)
	list, err := gauge.GetList()
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{"test": 500}, list, "Expected list %v, but got: %v", map[string]float64{"test": 500}, list)
}

func TestCounter_ProcessMassive(t *testing.T) {
	counter := NewCounter(nil)
	err := counter.ProcessMassive(map[string]float64{"test": 500})
	assert.NoError(t, err)
	list, err := counter.GetList()
	assert.NoError(t, err)
	assert.Equal(t, map[string]float64{"test": 500}, list, "Expected list %v, but got: %v", map[string]float64{"test": 500}, list)
}
