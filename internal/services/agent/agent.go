package agent

import (
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Agent struct {
	client         HTTPClient
	storage        Storage
	pollInterval   int
	reportInterval int
	sendAddr       string
	maxRetries     int
}

type Options struct {
	Client         HTTPClient
	Storage        Storage
	PollInterval   int
	ReportInterval int
	SendAddr       string
	MaxRetries     int
}

func New(options Options) *Agent {
	return &Agent{
		client:         options.Client,
		storage:        options.Storage,
		pollInterval:   options.PollInterval,
		reportInterval: options.ReportInterval,
		sendAddr:       options.SendAddr,
		maxRetries:     options.MaxRetries,
	}
}
