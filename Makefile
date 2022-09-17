mockgen:
	mockgen -destination=mocks/wallet/mock_service.go -source=./internal/wallet/handler.go Service
	mockgen -destination=mocks/wallet/mock_repository.go -source=./internal/wallet/service.go Repository
