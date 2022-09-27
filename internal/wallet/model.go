package wallet

type Wallet struct {
	ID                    string
	UserID                string
	Balance               float64
	BalanceUpperLimit     float64
	TransactionUpperLimit float64
}
