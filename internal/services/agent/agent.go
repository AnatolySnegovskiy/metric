package agent

import (
	"net/http"
	"time"
)

//go:generate mockgen -source=agent.go -destination=mocks/agent.go -package=mocks
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Agent struct {
	client         HTTPClient
	storage        Storage
	pollInterval   *time.Ticker
	reportInterval *time.Ticker
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
		pollInterval:   time.NewTicker(time.Duration(options.PollInterval) * time.Second),
		reportInterval: time.NewTicker(time.Duration(options.ReportInterval) * time.Second),
		sendAddr:       options.SendAddr,
	}
}
