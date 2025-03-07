# ONYVA hotels-service
This is an open-source microservice for www.onyva.fun

### API ROUTES 
```go
	r.GET("/hotels", func(ctx *gin.Context) {
		GetAll(ctx, db)
	})

	//filters
	r.POST("/hotel", func(ctx *gin.Context) {
		GetByFilter(ctx, db)
	})
```
