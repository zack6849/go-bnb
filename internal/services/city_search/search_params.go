package city_search

import (
	"github.com/gin-gonic/gin"
)

type CitySearchParameters struct {
	search string
}

func GetRequestsFromParams(c *gin.Context) (CitySearchParameters, error) {
	search := c.Query("search")
	return CitySearchParameters{
		search: search,
	}, nil
}
