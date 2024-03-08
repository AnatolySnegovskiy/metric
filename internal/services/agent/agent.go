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
	pollInterval   <-chan time.Time
	reportInterval <-chan time.Time
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
		pollInterval:   time.Tick(time.Duration(options.PollInterval) * time.Second),
		reportInterval: time.Tick(time.Duration(options.ReportInterval) * time.Second),
		sendAddr:       options.SendAddr,
	}
}
