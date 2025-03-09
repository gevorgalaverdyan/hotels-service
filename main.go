package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gevorgalaverdyan/hotels-service/db"
	"github.com/gevorgalaverdyan/hotels-service/models"
	"github.com/gevorgalaverdyan/hotels-service/routes"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	//log.SetLevel(log.WarnLevel)
}

func RateLimitMiddleware(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the client's IP
		clientIP := c.ClientIP()

		// Check rate limit
		httpError := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if httpError != nil {
			log.WithFields(log.Fields{
				"ip":         clientIP,
				"endpoint":   c.Request.URL.Path,
				"method":     c.Request.Method,
				"user-agent": c.Request.UserAgent(),
			}).Warn("Rate limit exceeded")

			c.AbortWithStatusJSON(httpError.StatusCode, gin.H{"error": httpError.Message})
			return
		}

		// Log the request
		log.WithFields(log.Fields{
			"ip":         clientIP,
			"endpoint":   c.Request.URL.Path,
			"method":     c.Request.Method,
			"user-agent": c.Request.UserAgent(),
		}).Info("Incoming request")

		c.Next()
	}
}

func main() {
	file, err := os.ReadFile("./crawler/hotels.json")
	if err != nil {
		log.WithError(err).Fatal("Error reading JSON file")
	}

	var hotels []models.Hotel
	err = json.Unmarshal(file, &hotels)
	if err != nil {
		log.WithError(err).Fatal("Error unmarshalling JSON")
	}

	collection := db.Connect()

	// Check if the database is empty before inserting
	if isDatabaseEmpty(collection) {
		insertHotels(collection, hotels)
	} else {
		log.Info("Database is not empty, skipping insertion")
	}

	r := gin.Default()
	log.Info("Server Started")

	lmt := tollbooth.NewLimiter(10, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Minute})
	lmt.SetMessage("Too many requests. Please try again later.")

	// Apply rate limiting middleware with logging
	r.Use(RateLimitMiddleware(lmt))

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	routes.RegisterRoutes(r, collection)

	r.Run(":5555")
}

func isDatabaseEmpty(collection *mongo.Collection) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.WithError(err).Fatal("Failed to count documents in the database")
	}

	return count == 0
}

func insertHotels(collection *mongo.Collection, hotels []models.Hotel) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var documents []interface{}
	for _, hotel := range hotels {
		documents = append(documents, bson.M{
			"name":        hotel.Name,
			"wikiLink":    hotel.WikiLink,
			"city":        hotel.City,
			"province":    hotel.Province,
			"image":       hotel.Image,
			"coordinates": hotel.Coordinates,
			"website":     hotel.Website,
			"address":     hotel.Address,
		})
	}

	result, err := collection.InsertMany(ctx, documents)
	if err != nil {
		log.WithError(err).Fatal("Error inserting documents")
	}

	log.WithFields(log.Fields{
		"inserted_count": len(result.InsertedIDs),
	}).Info("Documents inserted successfully")
}
