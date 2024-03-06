package agent

import (
	"time"
)

type Agent struct {
	storage        Storage
	pollInterval   <-chan time.Time
	reportInterval <-chan time.Time
	flagSendAddr   string
}

type Options struct {
	Storage        Storage
	PollInterval   int
	ReportInterval int
	SendAddr       string
}

func New(options Options) *Agent {
	return &Agent{
		storage:        options.Storage,
		pollInterval:   time.Tick(time.Duration(options.PollInterval) * time.Second),
		reportInterval: time.Tick(time.Duration(options.ReportInterval) * time.Second),
		flagSendAddr:   options.SendAddr,
	}
}
