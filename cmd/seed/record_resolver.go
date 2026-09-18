package main

import (
	"gobnb/internal/cache"
	"gobnb/internal/domain"
	"os"

	"gorm.io/gorm"
)

func resolvePropertyType(cache *cache.CacheStore, propertyType string, db *gorm.DB) (domain.PropertyType, error) {
	return cache.GetItemByFieldValue[domain.PropertyType]("name", propertyType, func() (domain.PropertyType, error) {
		result := domain.PropertyType{}
		tx := db.Model(domain.PropertyType{}).Where(domain.PropertyType{Name: propertyType}).FirstOrCreate(&result)
		return result, tx.Error
	}, true)
}

func resolveRoomType(cache *cache.CacheStore, roomType string, db *gorm.DB) (domain.RoomType, error) {
	return cache.GetItemByFieldValue[domain.RoomType]("name", roomType, func() (domain.RoomType, error) {
		result := domain.RoomType{}
		tx := db.Model(domain.RoomType{}).Where(domain.RoomType{Name: roomType}).FirstOrCreate(&result)
		return result, tx.Error
	}, true)
}

func resolveAmenities(cache *cache.CacheStore, amenityNames []string, db *gorm.DB) []domain.Amenity {
	var resolved []domain.Amenity
	for _, amenityName := range amenityNames {
		resolvedAmenity, _ := cache.GetItemByFieldValue[domain.Amenity]("name", amenityName, func() (domain.Amenity, error) {
			result := domain.Amenity{}
			tx := db.Model(domain.Amenity{}).Where(domain.Amenity{Name: amenityName}).FirstOrCreate(&result)
			return result, tx.Error
		}, true)
		resolved = append(resolved, resolvedAmenity)
	}
	return resolved
}

func resolveHost(cache *cache.CacheStore, host domain.Host, db *gorm.DB) (domain.Host, error) {
	return cache.GetItemByFieldValue[domain.Host]("profile_id", host.ProfileID, func() (domain.Host, error) {
		result := domain.Host{}
		tx := db.Model(domain.Host{}).
			Where(domain.Host{ProfileID: host.ProfileID}).
			Attrs(host).
			FirstOrCreate(&result)
		return result, tx.Error
	}, true)
}

func resolveFileByName(filename string) (*os.File, error) {
	path := "seed_data/" + filename
	return getFileHandle(path)
}
