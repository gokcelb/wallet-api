package transaction

import "time"

type Transaction struct {
	ID        string
	WalletID  string
	Type      string
	Amount    float64
	CreatedAt time.Time
}
