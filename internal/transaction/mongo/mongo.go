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
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (m *Mongo) ReadByWalletID(ctx context.Context, walletID string, pageNo, pageSize int) ([]*transaction.Transaction, error) {
	opts := options.Find().SetSkip(int64(pageNo * pageSize)).SetLimit(int64(pageSize))
	filter := bson.D{bson.E{Key: "wallet_id", Value: walletID}}

	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var mongoTxns []mongoTransaction
	if err = cursor.All(ctx, &mongoTxns); err != nil {
		log.Error(err)
		return nil, err
	}

	txns := []*transaction.Transaction{}
	for _, mongoTxn := range mongoTxns {
		txns = append(txns, newTransactionFromMongoTransaction(&mongoTxn))
	}

	return txns, nil
}

func (m *Mongo) ReadByWalletIDFilterByType(ctx context.Context, walletID, typeFilter string, pageNo, pageSize int) ([]*transaction.Transaction, error) {
	opts := options.Find().SetSkip(int64(pageNo * pageSize)).SetLimit(int64(pageSize))
	filter := bson.D{bson.E{Key: "wallet_id", Value: walletID}, bson.E{Key: "type", Value: typeFilter}}

	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var mongoTxns []mongoTransaction
	if err = cursor.All(ctx, &mongoTxns); err != nil {
		log.Error(err)
		return nil, err
	}

	txns := []*transaction.Transaction{}
	for _, mongoTxn := range mongoTxns {
		txns = append(txns, newTransactionFromMongoTransaction(&mongoTxn))
	}

	return txns, nil
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
