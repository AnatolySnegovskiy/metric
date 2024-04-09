package main

import (
	"context"
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/services/server"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/AnatolySnegovskiy/metric/internal/storages/clients"
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
	s := storages.NewMemStorage()
	s.AddMetric("gauge", metrics.NewGauge())
	s.AddMetric("counter", metrics.NewCounter())

	c, err := NewConfig()
	handleError(err)

	db, _ := clients.NewPostgres(context.Background(), c.dataBaseDSN)

	serv := server.New(s, db, logger.Sugar())

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
