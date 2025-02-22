package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(r *gin.Engine, db *mongo.Collection) {
	r.GET("/hotels") // returns all hotels
	r.POST("/hotel") //filters
}
