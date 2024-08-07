package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AnatolySnegovskiy/metric/internal/services/server"
	"go.uber.org/zap"
)

var buildVersion string
var buildDate string
var buildCommit string

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	fmt.Printf("Build version: %s\n", setDefaultValue(buildVersion, "N/A"))
	fmt.Printf("Build date: %s\n", setDefaultValue(buildDate, "N/A"))
	fmt.Printf("Build commit: %s\n", setDefaultValue(buildCommit, "N/A"))
	logger, err := zap.NewProduction()
	handleError(err)

	conf, err := NewConfig()
	handleError(err)

	serv, err := server.New(context.Background(), conf, logger.Sugar())
	handleError(err)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-quit
		serv.ShotDown()
		logger.Info("server stopped")
		os.Exit(0)
	}()

	logger.Info("server started on " + conf.GetServerAddress())
	handleError(serv.Run())
}

func setDefaultValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
