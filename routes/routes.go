package routes

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	r.GET("/hotels") // returns all hotels
	r.POST("/hotel") //filters
}
