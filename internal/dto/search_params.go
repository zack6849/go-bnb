package dto

import (
	"GoBNB/internal/domain"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRequestsFromParams(c *gin.Context) (SearchParameters, error) {
	var zero SearchParameters
	maxDistance := 30_000 //30km
	distParam, ok := c.GetQuery("distance")
	if ok {
		parsed, err := strconv.Atoi(distParam)
		if err != nil {
			return zero, fmt.Errorf("failed to parse distance, can't convert '%s' to int: %w", distParam, err)
		}
		//overwrite the default value of 30 if present in the request
		maxDistance = parsed
	}

	valLat, latExists := c.GetQuery("latitude")
	valLng, lngExists := c.GetQuery("longitude")

	if !latExists || !lngExists {
		return zero, errors.New("latitude and longitude are required to use this endpoint")
	}

	lat, err := strconv.ParseFloat(valLat, 32)
	if err != nil {
		return zero, fmt.Errorf("failed to parse latitude value '%s' to float: %w", valLat, err)
	}
	lng, err := strconv.ParseFloat(valLng, 32)
	if err != nil {
		return zero, fmt.Errorf("failed to parse longitude value '%s' to float: %w", valLng, err)
	}

	return SearchParameters{
		SearchNear:  domain.GetLocation(lat, lng),
		MaxDistance: maxDistance,
	}, nil
}

type ListSearchType int

const (
	ListSearchOr ListSearchType = iota
	ListSearchAnd
)

type SearchParameters struct {
	SearchNear  domain.Location
	MaxDistance int //max distance in meters
}

type PriceFilter struct {
	MinPrice int
	MaxPrice int
}

type AmenitiesFilter struct {
	//maybe this should take a list of uuids instead? not sure if i want to hydrate all the models for every amenity just to pass it somewhere else where it'll be converted to a where in anyways at some point... hmmm..
	DesiredAmenities []domain.Amenity
	SearchType       ListSearchType //is this a must-have all of, or an in() search? e.g.: do they have to have a Pool AND a Microwave, or is either one ok?
}

type PropertyTypeFilter struct {
}
