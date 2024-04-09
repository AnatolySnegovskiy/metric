package clients

import (
	"context"
	"github.com/AnatolySnegovskiy/metric/internal/mocks"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestQueryFunction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mocks.NewMockPgxConnInterface(ctrl)
	typeMap := &pgtype.Map{}
	resultReader := &pgconn.ResultReader{}
	expectedRows := pgx.RowsFromResultReader(typeMap, resultReader)
	mockConn.EXPECT().Query(gomock.Any(), "SELECT * FROM table").Return(expectedRows, nil)

	db := &Postgres{
		conn: mockConn,
		ctx:  context.Background(),
	}

	rows, err := db.Query("SELECT * FROM table")

	assert.NoError(t, err)
	assert.NotNil(t, rows)
}

func TestCloseFunction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mocks.NewMockPgxConnInterface(ctrl)

	// Устанавливаем ожидания для мок объекта
	mockConn.EXPECT().Close(gomock.Any()).Return(nil)

	// Создаем экземпляр Postgres с мок объектом
	db := &Postgres{
		conn: mockConn,
		ctx:  context.Background(),
	}

	// Вызываем функцию Close
	closed, err := db.Close()

	// Проверяем, что ожидаемые данные возвращены
	assert.NoError(t, err)
	assert.True(t, closed)
}
