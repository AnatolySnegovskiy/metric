package server

import (
	"context"
	"errors"
	pb "github.com/AnatolySnegovskiy/metric/internal/services/grpc"
	"github.com/AnatolySnegovskiy/metric/internal/services/interfase"
	"github.com/gookit/gsr"
	"strconv"
)

type GrpcServer struct {
	pb.UnimplementedMetricServiceServer
	storage interfase.Storage
	logger  gsr.GenLogger
	conf    Config
}

func NewGrpcServer(storage interfase.Storage, logger gsr.GenLogger, conf Config) *GrpcServer {
	return &GrpcServer{
		storage: storage,
		logger:  logger,
		conf:    conf,
	}
}

func (s *GrpcServer) Update(ctx context.Context, req *pb.MetricRequest) (*pb.MetricResponse, error) {
	storage := s.storage
	metric, err := storage.GetMetricType(req.Type)

	if err != nil {
		return &pb.MetricResponse{}, err
	}

	var value string
	if req.Delta != 0 {
		value = strconv.FormatInt(req.Delta, 10)
	} else if req.Value != 0 {
		value = strconv.FormatFloat(float64(req.Value), 'f', -1, 64)
	}

	if value == "" {
		return &pb.MetricResponse{}, errors.New("failed to process Value and Delta is empty")
	}

	_ = metric.Process(ctx, req.Id, value)

	return &pb.MetricResponse{
		Id:    req.Id,
		Type:  req.Type,
		Delta: req.Delta,
		Value: req.Value,
	}, nil
}

func (s *GrpcServer) UpdateMany(ctx context.Context, req *pb.MetricRequestMany) (*pb.MetricResponseMany, error) {
	list := make(map[string]map[string]float64)

	for _, metricDTO := range req.Requests {
		if list[metricDTO.Type] == nil {
			list[metricDTO.Type] = make(map[string]float64)
		}

		if metricDTO.Delta != 0 {
			list[metricDTO.Type][metricDTO.Id] += float64(metricDTO.Delta)
		} else if metricDTO.Value != 0 {
			list[metricDTO.Type][metricDTO.Id] = float64(metricDTO.Value)
		}
	}

	storage := s.storage
	var metrics []*pb.MetricResponse

	for metricType, metric := range list {
		metricEntity, err := storage.GetMetricType(metricType)

		if err != nil {
			return &pb.MetricResponseMany{}, err
		}

		err = metricEntity.ProcessMassive(ctx, metric)

		if err != nil {
			return &pb.MetricResponseMany{}, err
		}

		for metricName, value := range metric {
			metricResponse := &pb.MetricResponse{
				Id:   metricName,
				Type: metricType,
			}
			if metricType == "gauge" {
				metricResponse.Value = float32(value)
			} else {
				metricResponse.Delta = int64(value)
			}
			metrics = append(metrics, metricResponse)
		}
	}

	return &pb.MetricResponseMany{Responses: metrics}, nil
}

func (s *GrpcServer) Get(ctx context.Context, req *pb.MetricRequest) (*pb.MetricResponse, error) {
	metricType := req.Type
	metricName := req.Id

	storage, err := s.storage.GetMetricType(metricType)
	if err != nil {
		return &pb.MetricResponse{}, err
	}

	list, err := storage.GetList(ctx)
	if err != nil {
		return &pb.MetricResponse{}, err
	}

	metric, ok := list[metricName]

	if !ok {
		return &pb.MetricResponse{}, errors.New("metric not found")
	}

	metricResponse := pb.MetricResponse{
		Id:   metricName,
		Type: metricType,
	}

	if metricType == "gauge" {
		metricResponse.Value = float32(metric)
	} else {
		metricResponse.Delta = int64(metric)
	}

	return &metricResponse, nil
}

func (s *GrpcServer) GetAll(ctx context.Context, req *pb.MetricRequest) (*pb.MetricResponseMany, error) {
	metricType := req.Type
	storage, err := s.storage.GetMetricType(metricType)
	if err != nil {
		return &pb.MetricResponseMany{}, err
	}

	list, err := storage.GetList(ctx)
	if err != nil {
		return &pb.MetricResponseMany{}, err
	}

	metrics := make([]*pb.MetricResponse, 0, len(list))

	for metricName, metric := range list {
		metricResponse := &pb.MetricResponse{
			Id:   metricName,
			Type: metricType,
		}
		if metricType == "gauge" {
			metricResponse.Value = float32(metric)
		} else {
			metricResponse.Delta = int64(metric)
		}
		metrics = append(metrics, metricResponse)
	}

	return &pb.MetricResponseMany{Responses: metrics}, nil
}
