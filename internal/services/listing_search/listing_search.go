package listing_search

import (
	"gobnb/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func FindMatchingListings(params ListingSearchParameters, db *gorm.DB) ([]domain.Listing, error) {
	offset := (params.Page - 1) * params.Limit
	results := make([]domain.Listing, 0)
	query := GetListingQuery(db, params)
	tx := query.Offset(offset).Limit(params.Limit).Find(&results)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return results, nil
}

func GetListingQuery(db *gorm.DB, params ListingSearchParameters) *gorm.DB {
	return db.Model(&domain.Listing{}).
		Select("listings.*, st_distance(location, st_makepoint(?, ?)) as distance_meters",
			params.SearchNear.Longitude,
			params.SearchNear.Latitude,
		).
		Preload(clause.Associations).
		Where("st_dwithin(location,st_makepoint(?, ?), ?)",
			params.SearchNear.Longitude,
			params.SearchNear.Latitude,
			params.MaxDistance,
		)
}
