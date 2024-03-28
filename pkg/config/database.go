package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func DBInstance() *mongo.Client {

	if Client != nil {
		return Client
	}

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading the .env file")
	}

	db_url := os.Getenv("MONGODB_URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// mongo.NewClient(options.Client())
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db_url))
	if err != nil {
		log.Fatal(err)
	}

	Client = client

	fmt.Println("Connected to MongoDB!")

	return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cluster-0").Collection(collectionName)
	return collection
}
