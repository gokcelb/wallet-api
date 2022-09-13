package wallet

type Wallet struct {
	Id                    string
	UserId                string
	Balance               float64
	BalanceUpperLimit     float64
	TransactionUpperLimit float64
}
