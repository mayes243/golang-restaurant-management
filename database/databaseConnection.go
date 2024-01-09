package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {

	mongoDbURL := os.Getenv("DB_CONNECTION_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDbURL))

	defer cancel()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to mongodb")
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	databaseName := os.Getenv("DB_NAME")
	if databaseName == "" {
		databaseName = "restaurant"
	}
	var collection *mongo.Collection = client.Database(databaseName).Collection(collectionName)

	return collection
}
