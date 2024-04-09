package internal

//go:generate mockgen -source=storages/mem_storage.go -destination=mocks/mem_storage_mock.go -package=mocks
//go:generate mockgen -source=services/agent/agent.go -destination=mocks/agent_mock.go -package=mocks
//go:generate mockgen -source=services/agent/storage.go -destination=mocks/storage_mock.go -package=mocks
//go:generate mockgen -source=pgxconn.go -destination=mocks/mock_pgxconn.go -package=mocks
//go:generate mockgen -source=context.go -destination=mocks/mock_context.go -package=mocks
