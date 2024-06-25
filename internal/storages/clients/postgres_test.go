package clients

import (
	"context"
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestPostgres_Test(t *testing.T) {
	testCases := []struct {
		name   string
		expect func(mock pgxmock.PgxPoolIface)
		check  func(mockDB *Postgres)
	}{
		{
			name: "Exec",
			expect: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectExec(regexp.QuoteMeta("CREATE TABLE IF NOT EXISTS gauge (name varchar(100) PRIMARY KEY, value DOUBLE PRECISION)")).
					WillReturnResult(pgxmock.NewResult("CREATE", 1))
			},
			check: func(mockDB *Postgres) {
				_, err := mockDB.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS gauge (name varchar(100) PRIMARY KEY, value DOUBLE PRECISION)")
				assert.NoError(t, err)
			},
		},
		{
			name: "Query",
			expect: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT value FROM gauge WHERE name = $1")).
					WithArgs("test").
					WillReturnRows(pgxmock.NewRows([]string{"value"}).AddRow(100))
			},
			check: func(mockDB *Postgres) {
				_, err := mockDB.Query(context.Background(), "SELECT value FROM gauge WHERE name = $1", "test")
				assert.NoError(t, err)
			},
		},
		{
			name: "QueryRow",
			expect: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT value FROM gauge WHERE name = $1")).
					WithArgs("test").
					WillReturnRows(pgxmock.NewRows([]string{"value"}).AddRow(100))
			},
			check: func(mockDB *Postgres) {
				rows := mockDB.QueryRow(context.Background(), "SELECT value FROM gauge WHERE name = $1", "test")
				var value int
				err := rows.Scan(&value)
				assert.NoError(t, err)
				assert.Equal(t, 100, value)
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
			mockDB := NewPostgres(mock)
			testCase.check(mockDB)
		})
	}
}
