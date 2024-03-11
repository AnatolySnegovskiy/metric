package internal

//go:generate mockgen -source=storages/mem_storage.go -destination=storages/mocks/mem_storage_mock.go -package=mocks
//go:generate mockgen -source=services/agent/agent.go -destination=services/agent/mocks/agent_mock.go -package=mocks
//go:generate mockgen -source=services/server/server.go -destination=services/server/mocks/server_mock.go -package=mocks
//go:generate mockgen -source=services/agent/storage.go -destination=services/agent/mocks/storage_mock.go -package=mocks
