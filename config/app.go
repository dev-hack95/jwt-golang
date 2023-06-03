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

func DBinstance() *mongo.Client {
	// Load data from .env file
	err := godotenv.Load(".env")
	if err != nil {

		log.Fatal("Error in loding .env file")
	}

	MongoDb := os.Getenv("MONGO_URL")

	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))

	if err != nil {
		log.Fatal(err)
	}

	// Context with timeout function with Background as root if any error occured
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to mongo db")

	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	return collection
}
