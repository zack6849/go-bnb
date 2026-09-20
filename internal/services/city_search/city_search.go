package city_search

import (
	"context"
	"fmt"
	"gobnb/internal/domain"
	"gobnb/internal/services/geoapify"
	"gobnb/internal/services/geocode"
	"log"

	"gorm.io/gorm"
)

func FindMatchingCities(params CitySearchParameters, db *gorm.DB, ctx context.Context) ([]domain.City, error) {
	results := make([]domain.City, 0)
	query := GetCityQuery(db, params)
	tx := query.Find(&results)
	if len(results) == 0 {
		log.Printf("no results for term %s, fetching from upstream API", params.search)
		_, err := FetchCityFromApi(params, db, ctx)
		if err != nil {
			return results, fmt.Errorf("no local or remote results for term %s: %w", params.search, err)
		}
		//re-set the results entirely so we ensure the ordering is correct
		results = make([]domain.City, 0)
		query.Find(&results)
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return results, nil
}

func FetchCityFromApi(params CitySearchParameters, db *gorm.DB, ctx context.Context) (domain.City, error) {
	zero := domain.City{}
	result, err := geocode.CitySearch(params.search, geoapify.ResultLevelCity, ctx)
	if err != nil {
		return zero, err
	}

	location := domain.Location{
		Latitude:  result.Latitude,
		Longitude: result.Longitude,
	}

	city := domain.City{
		Hash:        location.ToPoint().Hash(),
		Location:    location,
		Name:        result.Name,
		FullName:    result.FormattedName,
		Timezone:    result.Timezone.Name,
		SubDivision: result.State,
	}
	//save the city
	db.Create(&city)
	return city, nil
}

func GetCityQuery(db *gorm.DB, params CitySearchParameters) *gorm.DB {
	query := db.Model(&domain.City{}).Order("name asc")
	if len(params.search) > 0 {
		query = query.Where("name ilike ?", fmt.Sprintf("%%%s%%", params.search))
	}
	return query
}
