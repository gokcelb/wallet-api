package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoTransaction struct {
	ID        primitive.ObjectID `bson:"_id"`
	WalletID  string             `bson:"wallet_id"`
	Type      string             `bson:"type"`
	Amount    float64            `bson:"amount"`
	CreatedAt time.Time          `bson:"created_at"`
}
