package search

import (
	"GoBNB/internal/configuration"
	"GoBNB/internal/domain"
	"GoBNB/internal/dto"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func FindMatchingListings(params dto.SearchParameters) ([]domain.Listing, error) {
	results := make([]domain.Listing, 0)
	db, err := configuration.GetDatabaseConfiguration().Open()
	if err != nil {
		return results, err
	}
	query := GetListingQuery(db, params)
	tx := query.Find(&results)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return results, nil
}

func GetListingQuery(db *gorm.DB, params dto.SearchParameters) *gorm.DB {
	return db.Model(&domain.Listing{}).
		Select("listings.*, st_distance(location, st_makepoint(?, ?)) as distance_meters",
			params.SearchNear.Latitude,
			params.SearchNear.Longitude,
		).
		Preload(clause.Associations).
		Where("st_dwithin(location,st_makepoint(?, ?), ?)",
			params.SearchNear.Latitude,
			params.SearchNear.Longitude,
			params.MaxDistance,
		)
}
