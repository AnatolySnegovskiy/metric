package main

import (
	"fmt"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/services/server"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
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

func handleShutdownSignal(quit chan os.Signal, serv *server.Server, c *Config) {
	<-quit
	fmt.Println("server stopped")
	err := serv.HandleShutdownSignal(c.fileStoragePath)
	handleError(err)
	os.Exit(0)
}

func main() {
	logger, _ := zap.NewProduction()
	s := storages.NewMemStorage()
	s.AddMetric("gauge", metrics.NewGauge())
	s.AddMetric("counter", metrics.NewCounter())

	c, err := NewConfig()
	handleError(err)

	serv := server.New(s, logger.Sugar())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go handleShutdownSignal(quit, serv, c)

	if c.restore {
		serv.LoadMetricsOnStart(c.fileStoragePath)
	}

	logger.Info("server started on " + c.flagRunAddr)

	go serv.SaveMetricsPeriodically(c.storeInterval, c.fileStoragePath)
	handleError(serv.Run(c.flagRunAddr))
}
