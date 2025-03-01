package db

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"syscall"
)

var client *mongo.Client

func Connect() *mongo.Collection {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Warn("No .env file found, using system environment variables")
	}

	MONGO_URI := os.Getenv("MONGO_URI")
	if MONGO_URI == "" {
		log.Fatal("MONGO_URI is not set!")
	}

	clientOptions := options.Client().ApplyURI(MONGO_URI)

	// Retry logic for MongoDB connection
	var err error
	for i := 0; i < 5; i++ {
		client, err = mongo.Connect(context.Background(), clientOptions)
		if err == nil && client.Ping(context.Background(), nil) == nil {
			log.Info("Connected to MongoDB")
			break
		}
		log.Warnf("Failed to connect to MongoDB (attempt %d/5), retrying in 2s...", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to MongoDB after retries:", err)
	}

	go handleShutdown()

	return client.Database("hotelsdb").Collection("hotels")
}

func handleShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Warn("Shutting down... Closing MongoDB connection.")

	if client != nil {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Errorf("Error closing MongoDB connection: %v", err)
		} else {
			log.Info("MongoDB connection closed successfully.")
		}
	}
	os.Exit(0)
}
