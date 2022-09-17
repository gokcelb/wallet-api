package mongo

import (
	"context"

	"github.com/gokcelb/wallet-api/internal/wallet"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct {
	collection *mongo.Collection
}

func NewMongo(collection *mongo.Collection) *Mongo {
	return &Mongo{collection}
}

func (m *Mongo) Create(ctx context.Context, wallet wallet.Wallet) (string, error) {
	mongoWallet := newMongoWalletFromWallet(wallet)
	result, err := m.collection.InsertOne(ctx, mongoWallet)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func newMongoWalletFromWallet(wallet wallet.Wallet) *mongoWallet {
	return &mongoWallet{
		Id:                    primitive.NewObjectID(),
		UserId:                wallet.UserId,
		Balance:               wallet.Balance,
		BalanceUpperLimit:     wallet.BalanceUpperLimit,
		TransactionUpperLimit: wallet.TransactionUpperLimit,
	}
}
