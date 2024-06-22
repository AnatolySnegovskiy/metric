package internal

//go:generate mockgen -source=storages/mem_storage.go -destination=mocks/mem_storage_mock.go -package=mocks
//go:generate mockgen -source=services/agent/agent.go -destination=mocks/agent_mock.go -package=mocks
//go:generate mockgen -source=services/server/server.go -destination=mocks/server_mock.go -package=mocks
//go:generate mockgen -source=services/interfase/storage.go -destination=mocks/storage_mock.go -package=mocks
