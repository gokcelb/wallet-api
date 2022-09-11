package wallet

type Wallet struct {
	id                    string
	userId                string
	balance               float64
	balanceUpperLimit     float64
	transactionUpperLimit float64
}
