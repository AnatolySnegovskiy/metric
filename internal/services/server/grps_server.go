package server

import (
	"context"
	"errors"
	pb "github.com/AnatolySnegovskiy/metric/internal/services/grpc/metric/v1"
	"github.com/AnatolySnegovskiy/metric/internal/services/interfase"
	"github.com/gookit/gsr"
	"log"
	"strconv"
)

type GrpcServer struct {
	pb.UnimplementedMetricV1ServiceServer
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

func (s *GrpcServer) UpdateMetricV1(ctx context.Context, req *pb.UpdateMetricV1Request) (*pb.UpdateMetricV1Response, error) {
	storage := s.storage
	metric, err := storage.GetMetricType(req.Type)

	if err != nil {
		return &pb.UpdateMetricV1Response{}, err
	}

	var value string
	if req.GetDelta() != 0 {
		value = strconv.FormatInt(req.GetDelta(), 10)
	} else if req.GetValue() != 0 {
		value = strconv.FormatFloat(float64(req.GetValue()), 'f', -1, 64)
	}

	if value == "" {
		return &pb.UpdateMetricV1Response{}, errors.New("failed to process Value and Delta is empty")
	}

	_ = metric.Process(ctx, req.Id, value)

	return &pb.UpdateMetricV1Response{
		Id:    req.Id,
		Type:  req.Type,
		Delta: req.GetDelta(),
		Value: req.GetValue(),
	}, nil
}

func (s *GrpcServer) UpdateManyMetricV1(ctx context.Context, req *pb.UpdateManyMetricV1Request) (*pb.UpdateManyMetricV1Response, error) {
	list := make(map[string]map[string]float64)

	for _, metricDTO := range req.Requests {
		if list[metricDTO.Type] == nil {
			list[metricDTO.Type] = make(map[string]float64)
		}

		if metricDTO.GetDelta() != 0 {
			list[metricDTO.Type][metricDTO.Id] += float64(metricDTO.GetDelta())
		} else if metricDTO.GetValue() != 0 {
			list[metricDTO.Type][metricDTO.Id] = float64(metricDTO.GetValue())
		}
	}

	storage := s.storage
	var metrics []*pb.UpdateMetricV1Response

	for metricType, metric := range list {
		metricEntity, err := storage.GetMetricType(metricType)

		if err != nil {
			return &pb.UpdateManyMetricV1Response{}, err
		}

		err = metricEntity.ProcessMassive(ctx, metric)

		if err != nil {
			return &pb.UpdateManyMetricV1Response{}, err
		}

		for metricName, value := range metric {
			metricResponse := &pb.UpdateMetricV1Response{
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

	return &pb.UpdateManyMetricV1Response{Responses: metrics}, nil
}

func (s *GrpcServer) GetMetricV1(ctx context.Context, req *pb.GetMetricV1Request) (*pb.GetMetricV1Response, error) {
	metricType := req.Type
	metricName := req.Id

	storage, err := s.storage.GetMetricType(metricType)

	if err != nil {
		return &pb.GetMetricV1Response{}, err
	}

	if storage == nil {
		return &pb.GetMetricV1Response{}, errors.New("storage error")
	}

	list, err := storage.GetList(ctx)
	if err != nil {
		return &pb.GetMetricV1Response{}, err
	}
	metric, ok := list[metricName]

	if !ok {
		return &pb.GetMetricV1Response{}, errors.New("metric not found")
	}

	metricResponse := &pb.GetMetricV1Response{
		Id:   metricName,
		Type: metricType,
	}

	if metricType == "gauge" {
		metricResponse.Value = float32(metric)
	} else {
		metricResponse.Delta = int64(metric)
	}
	log.Println(metricResponse)
	return metricResponse, nil
}

func (s *GrpcServer) GetAllMetricV1(ctx context.Context, req *pb.GetAllMetricV1Request) (*pb.GetAllMetricV1Response, error) {
	metricType := req.Type
	storage, err := s.storage.GetMetricType(metricType)

	if err != nil {
		return &pb.GetAllMetricV1Response{}, err
	}

	if storage == nil {
		return &pb.GetAllMetricV1Response{}, errors.New("storage error")
	}

	list, err := storage.GetList(ctx)
	if err != nil {
		return &pb.GetAllMetricV1Response{}, err
	}

	metrics := make([]*pb.UpdateMetricV1Response, 0, len(list))

	for metricName, metric := range list {
		metricResponse := &pb.UpdateMetricV1Response{
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

	return &pb.GetAllMetricV1Response{Responses: metrics}, nil
}
