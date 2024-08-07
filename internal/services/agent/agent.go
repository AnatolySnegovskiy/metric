package agent

import (
	grpc "github.com/AnatolySnegovskiy/metric/internal/services/grpc/metric/v1"
	"net/http"

	"github.com/AnatolySnegovskiy/metric/internal/services/interfase"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Agent struct {
	client         HTTPClient
	grpcClient     grpc.MetricV1ServiceClient
	storage        interfase.Storage
	pollInterval   int
	reportInterval int
	sendAddr       string
	maxRetries     int
	shaKey         string
	cryptoKey      string
}

type Options struct {
	Grpc           grpc.MetricV1ServiceClient
	Client         HTTPClient
	Storage        interfase.Storage
	PollInterval   int
	ReportInterval int
	SendAddr       string
	MaxRetries     int
	ShaKey         string
	CryptoKey      string
}

func New(options Options) *Agent {
	return &Agent{
		client:         options.Client,
		grpcClient:     options.Grpc,
		storage:        options.Storage,
		pollInterval:   options.PollInterval,
		reportInterval: options.ReportInterval,
		sendAddr:       options.SendAddr,
		maxRetries:     options.MaxRetries,
		shaKey:         options.ShaKey,
		cryptoKey:      options.CryptoKey,
	}
}
