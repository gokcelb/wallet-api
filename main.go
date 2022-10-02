package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gokcelb/wallet-api/config"
	"github.com/gokcelb/wallet-api/internal/transaction"
	transactionMongo "github.com/gokcelb/wallet-api/internal/transaction/mongo"
	"github.com/gokcelb/wallet-api/internal/wallet"
	walletMongo "github.com/gokcelb/wallet-api/internal/wallet/mongo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	conf, err := config.Read(fmt.Sprintf(".config/%s.json", config.GetEnvOrDefault()))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	ctx := context.Background()
	mongoClient := connectToMongo(ctx, conf)
	defer disconnectFromMongo(ctx, mongoClient)

	transactionCollection := mongoClient.
		Database(conf.Mongo.Database).
		Collection(conf.Mongo.Collection.Transaction)
	transactionRepository := transactionMongo.NewMongo(transactionCollection)
	transactionService := transaction.NewService(transactionRepository)
	transactionHandler := transaction.NewHandler(transactionService)

	walletCollection := mongoClient.
		Database(conf.Mongo.Database).
		Collection(conf.Mongo.Collection.Wallet)
	walletRepository := walletMongo.NewMongo(walletCollection)
	walletService := wallet.NewService(walletRepository, transactionService, conf)
	walletHandler := wallet.NewHandler(walletService)

	walletHandler.RegisterRoutes(e)
	transactionHandler.RegisterRoutes(e)

	go func() {
		if err := e.Start(":8000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down...")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func connectToMongo(ctx context.Context, conf config.Conf) *mongo.Client {
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.Mongo.URI))
	if err != nil {
		panic(err)
	}

	return mongoClient
}

func disconnectFromMongo(ctx context.Context, mongoClient *mongo.Client) {
	if err := mongoClient.Disconnect(ctx); err != nil {
		log.Error(err)
	}
	log.Info("disconnected from mongo")
}
