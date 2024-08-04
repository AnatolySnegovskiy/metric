package main

import (
	"context"
	"fmt"
	pb "github.com/AnatolySnegovskiy/metric/internal/services/grpc"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"sync"
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
	listen, err := net.Listen("tcp", conf.grpcAddress)
	if err != nil {
		handleError(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		<-quit
		serv.ShotDown()
		logger.Info("server stopped")
		logger.Info("gRPC server stopped")
		os.Exit(0)
	}()

	go func() {
		defer wg.Done()
		logger.Info("server started on " + conf.GetServerAddress())
		handleError(serv.Run())
	}()

	go func() {
		defer wg.Done()
		grpcServer := grpc.NewServer()
		pb.RegisterMetricServiceServer(grpcServer, serv.UpGrpc())
		logger.Info("gRPC server started on " + conf.grpcAddress)
		handleError(grpcServer.Serve(listen))
	}()
	logger.Info("server started on " + conf.GetServerAddress())
	wg.Wait()
}

func setDefaultValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
