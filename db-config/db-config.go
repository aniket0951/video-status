package dbconfig

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

func GetDBURL() string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err.Error())
	}

	return os.Getenv("LOCAL_DB_URL")
}

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(GetDBURL()))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conntError := client.Connect(ctx)

	if conntError != nil {
		log.Fatal(conntError)
	}

	pigError := client.Ping(ctx, nil)

	if pigError != nil {
		log.Fatal(pigError)
	}

	fmt.Println("Connection Established....")
	return client
}

var DB *mongo.Client = ConnectDB()

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var local = "videostatus"

	collection := client.Database(local).Collection(collectionName)
	return collection
}
