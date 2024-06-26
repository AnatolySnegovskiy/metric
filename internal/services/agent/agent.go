package agent

import (
	"github.com/AnatolySnegovskiy/metric/internal/services/interfase"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Agent struct {
	client         HTTPClient
	storage        interfase.Storage
	pollInterval   int
	reportInterval int
	sendAddr       string
	maxRetries     int
	shaKey         string
}

type Options struct {
	Client         HTTPClient
	Storage        interfase.Storage
	PollInterval   int
	ReportInterval int
	SendAddr       string
	MaxRetries     int
	ShaKey         string
}

func New(options Options) *Agent {
	return &Agent{
		client:         options.Client,
		storage:        options.Storage,
		pollInterval:   options.PollInterval,
		reportInterval: options.ReportInterval,
		sendAddr:       options.SendAddr,
		maxRetries:     options.MaxRetries,
		shaKey:         options.ShaKey,
	}
}
