package main

import (
	"context"
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
	logger, err := zap.NewProduction()
	handleError(err)

	c, err := NewConfig()
	handleError(err)

	db, err := pgx.Connect(context.Background(), c.dataBaseDSN)
	handleError(err)

	pg, err := clients.NewPostgres(context.Background(), db)
	handleError(err)

	var gaugeRepo *repositories.GaugeRepo
	var counterRepo *repositories.CounterRepo
	if db != nil {
		gaugeRepo, err = repositories.NewGaugeRepo(pg)
		handleError(err)
		counterRepo, err = repositories.NewCounterRepo(pg)
		handleError(err)
	}

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
