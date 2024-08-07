package server

import (
	"context"
	"github.com/AnatolySnegovskiy/metric/internal/entity/metrics"
	pb "github.com/AnatolySnegovskiy/metric/internal/services/grpc/metric/v1"
	"github.com/AnatolySnegovskiy/metric/internal/storages"
	"github.com/gookit/slog"
	"github.com/stretchr/testify/assert"
	"testing"
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
		req     *pb.UpdateMetricV1Request
		want    *pb.UpdateMetricV1Response
		wantErr bool
	}{
		{
			name: "successful update gauge metric",
			req: &pb.UpdateMetricV1Request{
				Id:   "metric1",
				Type: "gauge",
				RequestValue: &pb.UpdateMetricV1Request_Value{
					Value: 123.45,
				},
			},
			want: &pb.UpdateMetricV1Response{
				Id:    "metric1",
				Type:  "gauge",
				Value: 123.45,
			},
			wantErr: false,
		},
		{
			name: "successful update counter metric",
			req: &pb.UpdateMetricV1Request{
				Id:   "metric2",
				Type: "counter",
				RequestValue: &pb.UpdateMetricV1Request_Delta{
					Delta: 123,
				},
			},
			want: &pb.UpdateMetricV1Response{
				Id:    "metric2",
				Type:  "counter",
				Delta: 123,
			},
			wantErr: false,
		},
		{
			name: "failed update with empty value and delta",
			req: &pb.UpdateMetricV1Request{
				Id:   "metric1",
				Type: "gauge",
			},
			want:    &pb.UpdateMetricV1Response{},
			wantErr: true,
		},
		{
			name: "failed update with empty type",
			req: &pb.UpdateMetricV1Request{
				Id: "metric1",
				RequestValue: &pb.UpdateMetricV1Request_Delta{
					Delta: 123,
				},
				Type: "err",
			},
			want:    &pb.UpdateMetricV1Response{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := server.UpdateMetricV1(context.Background(), tt.req)
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
		req     *pb.UpdateManyMetricV1Request
		want    *pb.UpdateManyMetricV1Response
		wantErr bool
	}{
		{
			name: "successful update many gauge metrics",
			req: &pb.UpdateManyMetricV1Request{
				Requests: []*pb.UpdateMetricV1Request{
					{
						Id:   "metric1",
						Type: "gauge",
						RequestValue: &pb.UpdateMetricV1Request_Value{
							Value: 123.45,
						},
					},
					{
						Id:   "metric2",
						Type: "gauge",
						RequestValue: &pb.UpdateMetricV1Request_Value{
							Value: 67.89,
						},
					},
					{
						Id:   "metric3",
						Type: "counter",
						RequestValue: &pb.UpdateMetricV1Request_Delta{
							Delta: 10,
						},
					},
				},
			},
			want: &pb.UpdateManyMetricV1Response{
				Responses: []*pb.UpdateMetricV1Response{
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
			req: &pb.UpdateManyMetricV1Request{
				Requests: []*pb.UpdateMetricV1Request{
					{
						Id:   "metric1",
						Type: "ga2uge",
						RequestValue: &pb.UpdateMetricV1Request_Value{
							Value: 67.89,
						},
					},
					{
						Id:   "metric2",
						Type: "count1er",
						RequestValue: &pb.UpdateMetricV1Request_Delta{
							Delta: 10,
						},
					},
				},
			},
			want:    &pb.UpdateManyMetricV1Response{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := server.UpdateManyMetricV1(context.Background(), tt.req)
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

	req := &pb.GetAllMetricV1Request{Type: "gauge"}
	resp, err := server.GetAllMetricV1(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Responses))
	assert.Equal(t, "test_metric", resp.Responses[0].Id)
	assert.Equal(t, "gauge", resp.Responses[0].Type)
	assert.Equal(t, float32(42.0), resp.Responses[0].Value)

	req = &pb.GetAllMetricV1Request{Type: "counter"}
	resp, err = server.GetAllMetricV1(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resp.Responses))
	assert.Equal(t, "test_metric_counter", resp.Responses[0].Id)
	assert.Equal(t, "counter", resp.Responses[0].Type)
	assert.Equal(t, int64(1), resp.Responses[0].Delta)

	stg.AddMetric("gauge", nil)
	stg.AddMetric("counter", nil)
	_, err = server.GetAllMetricV1(ctx, req)
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

	req := &pb.GetMetricV1Request{Type: "gauge", Id: "test_metric"}
	resp, err := server.GetMetricV1(ctx, req)
	assert.NoError(t, err)

	assert.Equal(t, "test_metric", resp.Id)
	assert.Equal(t, "gauge", resp.Type)
	assert.Equal(t, float32(42.0), resp.Value)

	req = &pb.GetMetricV1Request{Type: "counter", Id: "test_metric_counter"}
	resp, err = server.GetMetricV1(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "test_metric_counter", resp.Id)
	assert.Equal(t, "counter", resp.Type)
	assert.Equal(t, int64(1), resp.Delta)

	stg.AddMetric("gauge", nil)
	stg.AddMetric("counter", nil)
	_, err = server.GetMetricV1(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, "storage error", err.Error())
}
