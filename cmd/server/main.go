package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/repositories"
	"github.com/AnatolySnegovskiy/metric/internal/services/server"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/AnatolySnegovskiy/metric/internal/storages/clients"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	logger, _ := zap.NewProduction()
	c, err := NewConfig()
	db, err := pgx.Connect(context.Background(), c.dataBaseDSN)

	if err != nil {
		err = errors.New("can't connect to database")
	}

	pg, err := clients.NewPostgres(context.Background(), db)

	if err != nil {
		err = errors.New("can't connect to postgres")
	}

	var gaugeRepo *repositories.GaugeRepo
	var counterRepo *repositories.CounterRepo
	if db != nil {
		gaugeRepo, err = repositories.NewGaugeRepo(pg)
		counterRepo, err = repositories.NewCounterRepo(pg)
	}

	if err != nil {
		err = errors.New("entity repository error")
	}

	handleError(err)

	s := storages.NewMemStorage()
	s.AddMetric("gauge", metrics.NewGauge(gaugeRepo))
	s.AddMetric("counter", metrics.NewCounter(counterRepo))

	serv := server.New(s, logger.Sugar(), db != nil)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		logger.Info("server stopped")
		serv.SaveMetricsToFile(c.fileStoragePath)
		os.Exit(0)
	}()
	if c.restore {
		serv.LoadMetricsOnStart(c.fileStoragePath)
	}

	logger.Info("server started on " + c.flagRunAddr)

	go serv.SaveMetricsPeriodically(context.Background(), c.storeInterval, c.fileStoragePath)
	handleError(serv.Run(c.flagRunAddr))
}
