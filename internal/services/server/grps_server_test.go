package server

import (
	"context"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/gookit/slog"
	"testing"

	pb "github.com/AnatolySnegovskiy/metric/internal/services/grpc"
	"github.com/stretchr/testify/assert"
)

func TestGrpcServer_Update(t *testing.T) {
	stg := storages.NewMemStorage()
	stg.AddMetric("gauge", metrics.NewGauge(nil))
	stg.AddMetric("counter", metrics.NewCounter(nil))
	logger := slog.New()
	conf := getMockConf(t)

	server := NewGrpcServer(stg, logger, conf)

	tests := []struct {
		name    string
		req     *pb.MetricRequest
		want    *pb.MetricResponse
		wantErr bool
	}{
		{
			name: "successful update gauge metric",
			req: &pb.MetricRequest{
				Id:    "metric1",
				Type:  "gauge",
				Value: 123.45,
			},
			want: &pb.MetricResponse{
				Id:    "metric1",
				Type:  "gauge",
				Value: 123.45,
			},
			wantErr: false,
		},
		{
			name: "successful update counter metric",
			req: &pb.MetricRequest{
				Id:    "metric2",
				Type:  "counter",
				Delta: 123,
			},
			want: &pb.MetricResponse{
				Id:    "metric2",
				Type:  "counter",
				Delta: 123,
			},
			wantErr: false,
		},
		{
			name: "failed update with empty value and delta",
			req: &pb.MetricRequest{
				Id:   "metric1",
				Type: "gauge",
			},
			want:    &pb.MetricResponse{},
			wantErr: true,
		},
		{
			name: "failed update with empty type",
			req: &pb.MetricRequest{
				Id:    "metric1",
				Delta: 123,
				Type:  "err",
			},
			want:    &pb.MetricResponse{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := server.Update(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGrpcServer_UpdateMany(t *testing.T) {
	stg := storages.NewMemStorage()
	stg.AddMetric("gauge", metrics.NewGauge(nil))
	stg.AddMetric("counter", metrics.NewCounter(nil))
	logger := slog.New()
	conf := getMockConf(t)

	server := NewGrpcServer(stg, logger, conf)

	tests := []struct {
		name    string
		req     *pb.MetricRequestMany
		want    *pb.MetricResponseMany
		wantErr bool
	}{
		{
			name: "successful update many gauge metrics",
			req: &pb.MetricRequestMany{
				Requests: []*pb.MetricRequest{
					{
						Id:    "metric1",
						Type:  "gauge",
						Value: 123.45,
					},
					{
						Id:    "metric2",
						Type:  "gauge",
						Value: 67.89,
					},
					{
						Id:    "metric3",
						Type:  "counter",
						Delta: 10,
					},
				},
			},
			want: &pb.MetricResponseMany{
				Responses: []*pb.MetricResponse{
					{
						Id:    "metric1",
						Type:  "gauge",
						Value: 123.45,
					},
					{
						Id:    "metric2",
						Type:  "gauge",
						Value: 67.89,
					},
					{
						Id:    "metric3",
						Type:  "counter",
						Delta: 10,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "failed update with empty value and delta",
			req: &pb.MetricRequestMany{
				Requests: []*pb.MetricRequest{
					{
						Id:    "metric1",
						Type:  "ga2uge",
						Value: 67.89,
					},
					{
						Id:    "metric2",
						Type:  "count1er",
						Delta: 10,
					},
				},
			},
			want:    &pb.MetricResponseMany{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := server.UpdateMany(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGrpcServer_GetAll(t *testing.T) {
	ctx := context.Background()
	stg := storages.NewMemStorage()
	metricGauge := metrics.NewGauge(nil)
	_ = metricGauge.Process(ctx, "test_metric", "42.0")
	metricCounter := metrics.NewCounter(nil)
	_ = metricCounter.Process(ctx, "test_metric_counter", "1")
	stg.AddMetric("gauge", metricGauge)
	stg.AddMetric("counter", metricCounter)
	logger := slog.New()
	conf := getMockConf(t)

	server := NewGrpcServer(stg, logger, conf)

	req := &pb.MetricRequest{Type: "gauge"}
	resp, err := server.GetAll(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Responses))
	assert.Equal(t, "test_metric", resp.Responses[0].Id)
	assert.Equal(t, "gauge", resp.Responses[0].Type)
	assert.Equal(t, float32(42.0), resp.Responses[0].Value)

	req = &pb.MetricRequest{Type: "counter"}
	resp, err = server.GetAll(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Responses))
	assert.Equal(t, "test_metric_counter", resp.Responses[0].Id)
	assert.Equal(t, "counter", resp.Responses[0].Type)
	assert.Equal(t, int64(1), resp.Responses[0].Delta)

	stg.AddMetric("gauge", nil)
	stg.AddMetric("counter", nil)
	_, err = server.GetAll(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, "storage error", err.Error())
}

func TestGrpcServer_Get(t *testing.T) {
	ctx := context.Background()
	stg := storages.NewMemStorage()
	metricGauge := metrics.NewGauge(nil)
	_ = metricGauge.Process(ctx, "test_metric", "42.0")
	metricCounter := metrics.NewCounter(nil)
	_ = metricCounter.Process(ctx, "test_metric_counter", "1")
	stg.AddMetric("gauge", metricGauge)
	stg.AddMetric("counter", metricCounter)
	logger := slog.New()
	conf := getMockConf(t)

	server := NewGrpcServer(stg, logger, conf)

	req := &pb.MetricRequest{Type: "gauge", Id: "test_metric"}
	resp, err := server.Get(ctx, req)
	assert.NoError(t, err)

	assert.Equal(t, "test_metric", resp.Id)
	assert.Equal(t, "gauge", resp.Type)
	assert.Equal(t, float32(42.0), resp.Value)

	req = &pb.MetricRequest{Type: "counter", Id: "test_metric_counter"}
	resp, err = server.Get(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "test_metric_counter", resp.Id)
	assert.Equal(t, "counter", resp.Type)
	assert.Equal(t, int64(1), resp.Delta)

	stg.AddMetric("gauge", nil)
	stg.AddMetric("counter", nil)
	_, err = server.Get(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, "storage error", err.Error())
}
