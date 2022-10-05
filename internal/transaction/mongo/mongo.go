package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/gokcelb/wallet-api/internal/transaction"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
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
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return "", transaction.ErrTransactionNotFound
	} else if err != nil {
		log.Error(err)
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (m *Mongo) Read(ctx context.Context, id string) (*transaction.Transaction, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var mongoTxn mongoTransaction
	err = m.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&mongoTxn)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, transaction.ErrTransactionNotFound
	} else if err != nil {
		return nil, err
	}

	return newTransactionFromMongoTransaction(&mongoTxn), nil
}

func (m *Mongo) ReadByType(ctx context.Context, id string, typeFilter string) (*transaction.Transaction, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	filter := bson.D{bson.E{Key: "_id", Value: objectID}, bson.E{Key: "type", Value: typeFilter}}

	var mongoTxn mongoTransaction
	err = m.collection.FindOne(ctx, filter).Decode(&mongoTxn)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, transaction.ErrTransactionNotFound
	} else if err != nil {
		return nil, err
	}

	return newTransactionFromMongoTransaction(&mongoTxn), nil
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
