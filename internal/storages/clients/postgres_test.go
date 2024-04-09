package clients

import (
	"context"
	"fmt"
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

func TestCloseFailFunction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConn := mocks.NewMockPgxConnInterface(ctrl)

	// Устанавливаем ожидания для мок объекта
	mockConn.EXPECT().Close(gomock.Any()).Return(fmt.Errorf("some error"))

	// Создаем экземпляр Postgres с мок объектом
	db := &Postgres{
		conn: mockConn,
		ctx:  context.Background(),
	}

	// Вызываем функцию Close
	closed, err := db.Close()

	// Проверяем, что ожидаемые данные возвращены
	assert.Error(t, err)
	assert.False(t, closed)
}

func TestNewPostgres(t *testing.T) {
	// Подготовка тестовых данных
	ctx := context.Background()
	validConfigString := "host=localhost port=5432 user=postgres password=root sslmode=disable"
	invalidConfigString := "invalid connection string"

	// Тестируем успешное создание объекта Postgres
	postgres, err := NewPostgres(ctx, validConfigString)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	assert.NotNil(t, postgres.conn)
	assert.Equal(t, ctx, postgres.ctx)

	// Тестируем возвращение ошибки при неверной конфигурации
	_, err = NewPostgres(ctx, invalidConfigString)
	assert.Error(t, err, "Expected an error for invalid config string")
}
