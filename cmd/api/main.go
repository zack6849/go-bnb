package main

import (
	"gobnb/internal/bootstrap"
	"gobnb/internal/configuration"
	"gobnb/internal/services/city_search"
	"gobnb/internal/services/listing_search"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func addListingsRoute(router *gin.Engine, db *gorm.DB) *gin.Engine {
	router.GET("/api/listings/search", FetchListings(db))
	return router
}

func addCitySearchRoute(router *gin.Engine, db *gorm.DB) *gin.Engine {
	router.GET("/api/cities/search", SearchCities(db))
	return router
}

func main() {
	_ = bootstrap.Initialize()
	configuration.Load()
	db, err := configuration.GetDatabaseConfiguration().Open()
	if err != nil {
		log.Fatal(err)
		return
	}
	// Create a Gin router with default middleware (logger and recovery)
	router := setupRouter()
	router = addListingsRoute(router, db)
	router = addCitySearchRoute(router, db)
	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	err = router.Run()
	if err != nil {
		return
	}
}
func FetchListings(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		params, err := listing_search.GetRequestsFromParams(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "failed to parse request parameters",
			})
			return
		}
		relatedListings, err := listing_search.FindMatchingListings(params, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"results": relatedListings,
		})
	}

}

// SearchCities search cities takes a substring, eg: ?search=Al and returns a list of cities matching the substring
// should return an array of cities with a lat and lng and their name
// TODO: extract this shared logic out into some kind of common APIResponse function that takes a search parameter DTO and returns a slice of structs as a response
func SearchCities(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		params, err := city_search.GetRequestsFromParams(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to parse request parameters",
				"error":   err.Error(),
			})
			return
		}
		relatedListings, err := city_search.FindMatchingCities(params, db, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"results": relatedListings,
		})
	}
}
