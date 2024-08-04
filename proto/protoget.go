package proto

//go:generate protoc --go_out=./../internal/services/grpc --go_opt=paths=source_relative --go-grpc_out=./../internal/services/grpc --go-grpc_opt=paths=source_relative metric.proto
