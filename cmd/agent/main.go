package main

import (
	"context"
	"fmt"
	pb "github.com/AnatolySnegovskiy/metric/internal/services/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/services/agent"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
)

var buildVersion string
var buildDate string
var buildCommit string

func handleError(err error) {
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

func handleShutdownSignal(quit chan os.Signal) {
	<-quit
	fmt.Println("Agent stopped")
	os.Exit(0)
}

func main() {
	fmt.Printf("Build version: %s\n", setDefaultValue(buildVersion, "N/A"))
	fmt.Printf("Build date: %s\n", setDefaultValue(buildDate, "N/A"))
	fmt.Printf("Build commit: %s\n", setDefaultValue(buildCommit, "N/A"))

	s := storages.NewMemStorage()
	s.AddMetric("gauge", metrics.NewGauge(nil))
	s.AddMetric("counter", metrics.NewCounter(nil))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go handleShutdownSignal(quit)

	fmt.Println("Agent started")
	c, err := NewConfig()
	handleError(err)

	conn, err := grpc.NewClient(c.grpcSendAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	handleError(err)
	defer func(conn *grpc.ClientConn) {
		handleError(conn.Close())
	}(conn)
	grpcClient := pb.NewMetricServiceClient(conn)

	handleError(
		agent.New(
			agent.Options{
				Storage:        s,
				Grpc:           grpcClient,
				Client:         &http.Client{},
				PollInterval:   c.PollInterval,
				ReportInterval: c.ReportInterval,
				SendAddr:       c.FlagSendAddr,
				MaxRetries:     c.maxRetries,
				ShaKey:         c.shaKey,
				CryptoKey:      c.CryptoKey,
			},
		).Run(context.Background()))
}

func setDefaultValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
