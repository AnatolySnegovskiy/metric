package agent

import (
	"net/http"
)

//go:generate mockgen -source=agent.go -destination=mocks/agent.go -package=mocks
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Agent struct {
	client         HTTPClient
	storage        Storage
	pollInterval   int
	reportInterval int
	sendAddr       string
}

type Options struct {
	Client         HTTPClient
	Storage        Storage
	PollInterval   int
	ReportInterval int
	SendAddr       string
}

func New(options Options) *Agent {
	return &Agent{
		client:         options.Client,
		storage:        options.Storage,
		pollInterval:   options.PollInterval,
		reportInterval: options.ReportInterval,
		sendAddr:       options.SendAddr,
	}
}
