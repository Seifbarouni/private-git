package db

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db     *mongo.Database
	client *mongo.Client
)

func Init() error {
	var err error
	mongoEndpoint := os.Getenv("MONGO_URI")
	if mongoEndpoint == "" {
		return errors.New("MONGO_URI is not set")
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		return err
	}

	slog.Info("Connected to MongoDB")

	//db = client.Database("private-git")
	db = client.Database("test")
	return nil
}

func Collection(col string) *mongo.Collection {
	return db.Collection(col)
}
