package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type mongoWallet struct {
	Id                    primitive.ObjectID `bson:"_id"`
	UserId                string             `bson:"user_id"`
	Balance               float64            `bson:"balance"`
	BalanceUpperLimit     float64            `bson:"balance_upper_limit"`
	TransactionUpperLimit float64            `bson:"transaction_upper_limit"`
}
