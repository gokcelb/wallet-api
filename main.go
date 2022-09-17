package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gokcelb/wallet-api/config"
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

	walletCollection := mongoClient.Database(conf.Mongo.Database).Collection(conf.Mongo.Collection)

	walletRepository := walletMongo.NewMongo(walletCollection)
	walletService := wallet.NewService(walletRepository, conf)
	walletHandler := wallet.NewHandler(walletService)

	walletHandler.RegisterRoutes(e)

	go func() {
		if err := e.Start(":8000"); err != nil {
			panic(err)
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

	log.Info("system is shutting down...")
}

func connectToMongo(ctx context.Context, conf config.Conf) *mongo.Client {
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.Mongo.Uri))
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
