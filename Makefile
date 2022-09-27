test:
	go test ./...

lint:
	golangci-lint run

mockgen:
	mockgen -destination=mocks/wallet/mock_service.go -mock_names WalletService=MockWalletService -package=mock_wallet -source=./internal/wallet/handler.go
	mockgen -destination=mocks/wallet/mock_repository.go -mock_names WalletRepository=MockWalletRepository -package=mock_wallet -source=./internal/wallet/service.go
