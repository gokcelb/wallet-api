package mongo

import (
	"context"
	"time"

	"github.com/gokcelb/wallet-api/internal/transaction"
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

func (m *Mongo) Create(ctx context.Context, txn *transaction.Transaction) (string, error) {
	mongoTxn := newMongoTransactionFromTransaction(txn)
	result, err := m.collection.InsertOne(ctx, mongoTxn)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func newMongoTransactionFromTransaction(txn *transaction.Transaction) *mongoTransaction {
	return &mongoTransaction{
		ID:        primitive.NewObjectID(),
		WalletID:  txn.WalletID,
		Type:      txn.Type,
		Amount:    txn.Amount,
		CreatedAt: time.Now(),
	}
}

func newTransactionFromMongoTransaction(mongoTxn *mongoTransaction) *transaction.Transaction {
	return &transaction.Transaction{
		ID:        mongoTxn.ID.Hex(),
		WalletID:  mongoTxn.WalletID,
		Type:      mongoTxn.Type,
		Amount:    mongoTxn.Amount,
		CreatedAt: mongoTxn.CreatedAt,
	}
}
