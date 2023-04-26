package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client mongo.Client

func ConnectToDatabase() {
	db, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_CONNECTION_URI")))

	if err != nil {
		fmt.Println(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = db.Connect(ctx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Database connected")

	Client = *db

}

func GetClient() mongo.Client {
	return Client
}

func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("test").Collection(collectionName)
}
