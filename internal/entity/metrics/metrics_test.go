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
	mock.ExpectExec(regexp.QuoteMeta("CREATE TABLE IF NOT EXISTS counter (name varchar(100) PRIMARY KEY, value int)")).
		WillReturnResult(pgxmock.NewResult("CREATE", 1))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM counter")).
		WillReturnError(errors.New("db error"))
	mockDB, _ := clients.NewPostgres(context.Background(), mock)
	cr, _ := repositories.NewCounterRepo(mockDB)
	counter := NewCounter(cr)
	_, err = counter.GetList()
	assert.Error(t, err)
}
