test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...      

lint:
	golangci-lint run

mockgen:
# wallet
	mockgen -destination=internal/wallet/mock/wallet_repository.go -package mock github.com/gokcelb/wallet-api/internal/wallet WalletRepository
	mockgen -destination=internal/wallet/mock/transaction_service.go -package mock github.com/gokcelb/wallet-api/internal/wallet TransactionService
	mockgen -destination=internal/wallet/mock/wallet_service.go -package mock github.com/gokcelb/wallet-api/internal/wallet WalletService

# transaction
	mockgen -destination=internal/transaction/mock/transaction_repository.go -package mock github.com/gokcelb/wallet-api/internal/transaction TransactionRepository
	mockgen -destination=internal/transaction/mock/transaction_service.go -package mock github.com/gokcelb/wallet-api/internal/transaction TransactionService
