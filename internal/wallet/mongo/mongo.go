package mongo

import (
	"context"
	"errors"

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

func (m *Mongo) Read(ctx context.Context, id string) (wallet.Wallet, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(err)
		return wallet.Wallet{}, err
	}

	var mongoWallet mongoWallet
	err = m.collection.FindOne(ctx, primitive.M{"_id": objectId}).Decode(&mongoWallet)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return wallet.Wallet{}, wallet.ErrWalletNotFound
	} else if err != nil {
		return wallet.Wallet{}, err
	}

	return *newWalletFromMongoWallet(&mongoWallet), nil
}

func (m *Mongo) ReadByUserId(ctx context.Context, userId string) (wallet.Wallet, error) {
	var mongoWallet mongoWallet
	err := m.collection.FindOne(ctx, primitive.M{"user_id": userId}).Decode(&mongoWallet)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return wallet.Wallet{}, wallet.ErrWalletNotFound
	} else if err != nil {
		return wallet.Wallet{}, err
	}

	return *newWalletFromMongoWallet(&mongoWallet), nil
}

func (m *Mongo) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = m.collection.DeleteOne(ctx, primitive.M{"_id": objectID})
	return err
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

func newWalletFromMongoWallet(mongoWallet *mongoWallet) *wallet.Wallet {
	return &wallet.Wallet{
		Id:                    mongoWallet.Id.Hex(),
		UserId:                mongoWallet.UserId,
		Balance:               mongoWallet.Balance,
		BalanceUpperLimit:     mongoWallet.BalanceUpperLimit,
		TransactionUpperLimit: mongoWallet.TransactionUpperLimit,
	}
}
