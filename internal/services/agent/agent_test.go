package agent

import (
	"github.com/AnatolySnegovskiy/metric/internal/services/server/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	options := Options{
		Storage:        mocks.NewMockStorage(gomock.NewController(t)),
		PollInterval:   10,
		ReportInterval: 20,
		SendAddr:       "example.com:1234",
	}

	agent := New(options)
	assert.NotNil(t, agent.storage, "storage should not be nil")
	assert.NotNil(t, agent.pollInterval, "pollInterval should not be nil")
	assert.NotNil(t, agent.reportInterval, "reportInterval should not be nil")
	assert.Equal(t, "example.com:1234", agent.flagSendAddr, "send address should be example.com:1234")
}
