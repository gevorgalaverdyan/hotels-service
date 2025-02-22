package db

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Collection {
	// Find .evn
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env", err)
	}

	// Get value from .env
	MONGO_URI := os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Connected to mongo")
	}

	collection := client.Database("hotelsdb").Collection("hotels")
	if err != nil {
		log.Fatal(err)
	}

	return collection
}
