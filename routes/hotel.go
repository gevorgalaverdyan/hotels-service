package routes

import (
	"net/http"

	"github.com/gevorgalaverdyan/hotels-service/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Filter struct {
	Filter string `json:"filter"`
}

func GetAll(ctx *gin.Context, db *mongo.Collection) {
	hotels, err := models.GetAll(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, hotels)
}

func GetByFilter(ctx *gin.Context, db *mongo.Collection) {
	var filter Filter

	if err := ctx.ShouldBindJSON(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hotels, err := models.GetByFilter(filter.Filter, db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, hotels)
}
