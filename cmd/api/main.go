package main

import (
	"GoBNB/internal/bootstrap"
	"GoBNB/internal/configuration"
	"GoBNB/internal/dto"
	"GoBNB/internal/services/search"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func addListingsRoute(router *gin.Engine) *gin.Engine {
	router.GET("/api/listings", FetchListings)
	return router
}

func main() {
	_ = bootstrap.Initialize()
	configuration.Load()
	// Create a Gin router with default middleware (logger and recovery)
	router := setupRouter()
	router = addListingsRoute(router)
	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	err := router.Run()
	if err != nil {
		return
	}
}

func FetchListings(c *gin.Context) {
	params, err := dto.GetRequestsFromParams(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "failed to parse request parameters",
		})
		return
	}
	relatedListings, err := search.FindMatchingListings(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"listings": relatedListings,
	})
}
