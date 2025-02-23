package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(r *gin.Engine, db *mongo.Collection) {
	// returns all hotels
	r.GET("/hotels", func(ctx *gin.Context) {
		GetAll(ctx, db)
	})

	//filters
	r.POST("/hotel", func(ctx *gin.Context) {
		GetByFilter(ctx, db)
	})
}
