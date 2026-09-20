package geocode

import (
	"context"
	"fmt"
	"gobnb/internal/configuration"
	"gobnb/internal/geo"
	"gobnb/internal/services/geoapify"

	"golang.org/x/time/rate"
)

// ReverseGeocode Attempts to geocode the given location to the given result level

func ReverseGeocode(point *geo.Point, level geoapify.ResultSpecificityLevel, ctx context.Context) (geoapify.ReverseGeocodingResult, error) {
	AllowedGeocodeResultTypes := map[geoapify.ResultSpecificityLevel]bool{
		geoapify.ResultLevelAmenity:  false, //too specific
		geoapify.ResultLevelBuilding: false, //too specific
		geoapify.ResultLevelStreet:   true,
		geoapify.ResultLevelSuburb:   true,
		geoapify.ResultLevelDistrict: true,
		geoapify.ResultLevelPostCode: true,
		geoapify.ResultLevelCity:     true,
		geoapify.ResultLevelCounty:   true,
		geoapify.ResultLevelCountry:  false, //not specific enough
	}

	zero := geoapify.ReverseGeocodingResult{}
	client := GetGeocodingClient()
	res, err := client.ReverseGeocodePoint(point, level, ctx)
	if err != nil {
		return zero, err
	}
	//we asked for a specific level, but didn't get back what we wanted
	//try midpoints to see if any of those give back the requested level
	//e.g.: center of the bbox is in the water, but center north midpoint is LA, we need to return LA
	if res.ResultLevel != level {
		allowed, exists := AllowedGeocodeResultTypes[res.ResultLevel]
		if exists {
			if allowed {
				return res, nil
			}
			gridFallback, err := GridBasedFallback(point, level, client, ctx)
			if err != nil {
				return zero, err
			}
			return gridFallback, err
		}
	}
	return res, nil
}

func CitySearch(search string, level geoapify.ResultSpecificityLevel, ctx context.Context) (geoapify.GeocodeSearchResult, error) {
	zero := geoapify.GeocodeSearchResult{}
	client := GetGeocodingClient()
	result, err := client.Autocomplete(search, level, ctx)
	if err != nil {
		return zero, fmt.Errorf("failed to geocode '%s' at level '%s' %w", search, level, err)
	}
	if result.ResultLevel != geoapify.ResultLevelCity {
		return result, fmt.Errorf("received non-%s result level back %s for term %s", level, result.ResultLevel, search)
	}
	return result, err
}

func Geocode(search string, level geoapify.ResultSpecificityLevel, ctx context.Context) (geoapify.GeocodeSearchResult, error) {
	zero := geoapify.GeocodeSearchResult{}
	client := GetGeocodingClient()
	result, err := client.Geocode(search, level, ctx)
	if err != nil {
		return zero, fmt.Errorf("failed to geocode '%s' at level '%s' %w", search, level, err)
	}
	if result.ResultLevel != geoapify.ResultLevelCity {
		return result, fmt.Errorf("received non-%s result level back %s for term %s", level, result.ResultLevel, search)
	}
	return result, err
}

func GridBasedFallback(point *geo.Point, level geoapify.ResultSpecificityLevel, client geoapify.Client, ctx context.Context) (geoapify.ReverseGeocodingResult, error) {
	zero := geoapify.ReverseGeocodingResult{}
	for _, p := range point.SubdivideIntoGrid(4, geo.HashLevelCity) {
		fallbackResult, err := ReverseGeocodePoint(p.Point, level, client, ctx)
		if err == nil {
			if fallbackResult.ResultLevel == level {
				return fallbackResult, nil
			}
		}
	}
	return zero, fmt.Errorf("failed to find a fallback record inside of the grid that maps to a %s", level)
}

func ReverseGeocodePoint(point *geo.Point, level geoapify.ResultSpecificityLevel, client geoapify.Client, ctx context.Context) (geoapify.ReverseGeocodingResult, error) {
	return client.ReverseGeocodePoint(point, level, ctx)
}

func GetGeocodingClient() geoapify.Client {
	return geoapify.CreateClient(
		configuration.GetString("GEOAPIFY_API_KEY", ""),
		rate.Limit(configuration.GetInt("GEOAPIFY_RATE_LIMIT", 1)),
	)
}
