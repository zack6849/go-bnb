package main

import (
	"GoBNB/internal/domain"
	"GoBNB/internal/query_cache"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type KeyValueRecord = map[string]string

func CreateListingRecord(record KeyValueRecord, c *query_cache.CacheStore, db *gorm.DB) (domain.Listing, error) {
	latFloat, _ := strconv.ParseFloat(record["latitude"], 64)
	lngFloat, _ := strconv.ParseFloat(record["longitude"], 64)
	airBnbId, _ := strconv.Atoi(record["id"])
	accommodates, _ := strconv.Atoi(record["accommodates"])
	numBeds, _ := strconv.Atoi(record["beds"])
	minNights, _ := strconv.Atoi(record["minimum_nights"])
	maxNights, _ := strconv.Atoi(record["maximum_nights"])
	hostProfileId, _ := strconv.ParseInt(record["host_profile_id"], 10, 64)
	//then create listing amenities records for it
	//these all resolve to real rows, so a swallowed error here becomes a
	//zero-value struct with a zero PK and fails much later at insert time.
	amenities, err := ParseAmenities(record, c, db)
	if err != nil {
		return domain.Listing{}, err
	}
	propertyType, err := GetPropertyTypeByName(c, db, record["property_type"])
	if err != nil {
		return domain.Listing{}, fmt.Errorf("failed to resolve property type %q: %w", record["property_type"], err)
	}
	roomType, err := GetRoomTypeByName(c, db, record["room_type"])
	if err != nil {
		return domain.Listing{}, fmt.Errorf("failed to resolve room type %q: %w", record["room_type"], err)
	}
	host, err := GetHostForListing(c, db, hostProfileId, record)
	if err != nil {
		return domain.Listing{}, fmt.Errorf("failed to resolve host %d: %w", hostProfileId, err)
	}

	numBaths := 0.0
	_, hasBaths := record["bathrooms"]
	if hasBaths {
		numBaths, _ = strconv.ParseFloat(record["bathrooms"], 32) //baths can be a float, eg: 1.5
	}

	location := domain.Location{
		Latitude:  latFloat,
		Longitude: lngFloat,
	}

	listing := domain.Listing{
		//I can't wait for go 1.27, so I can access promoted fields directly
		PropertyType: propertyType,
		RoomType:     roomType,
		Amenities:    amenities,
		AirbnbID:     airBnbId,
		NumBaths:     float32(numBaths),
		NumBeds:      numBeds,
		ListingUrl:   record["listing_url"],
		Tagline:      record["name"],
		Description:  record["description"],
		Accommodates: accommodates,
		MinNights:    minNights,
		MaxNights:    maxNights,
		Location:     location,
		Hosts:        []domain.Host{host},
	}

	return listing, nil
}

func ParseAmenities(record KeyValueRecord, c *query_cache.CacheStore, db *gorm.DB) ([]domain.Amenity, error) {
	list := make([]domain.Amenity, 0, 1)
	var amenities []string
	if err := json.Unmarshal([]byte(record["amenities"]), &amenities); err != nil {
		return nil, fmt.Errorf("failed to read amenities json for record: %w", err)
	}
	for _, amenityName := range amenities {
		amenity, err := GetAmenityByName(c, db, amenityName)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve amenity by name %s: %w", amenityName, err)
		}
		list = append(list, amenity)
	}
	return list, nil
}

func GetRecordByFieldValue[T any](c *query_cache.CacheStore, db *gorm.DB, field string, value any, attrs ...any) (T, error) {
	return query_cache.GetItemByFieldValue[T](c, field, value, func() (T, error) {
		var result T
		q := db.Where(map[string]any{field: value})
		if len(attrs) > 0 {
			q = q.Attrs(attrs...)
		}
		res := q.FirstOrCreate(&result)
		if res.Error != nil {
			var zero T
			return zero, res.Error
		}
		if res.RowsAffected == 1 {
			fmt.Printf("Created new %T: %v\n", result, value)
		}
		return result, nil
	}, true)
}

func GetListItemByName[T any](c *query_cache.CacheStore, db *gorm.DB, name string) (T, error) {
	return GetRecordByFieldValue[T](c, db, "name", name)
}

func GetAmenityByName(c *query_cache.CacheStore, db *gorm.DB, name string) (domain.Amenity, error) {
	return GetListItemByName[domain.Amenity](c, db, name)
}

func GetPropertyTypeByName(c *query_cache.CacheStore, db *gorm.DB, name string) (domain.PropertyType, error) {
	return GetListItemByName[domain.PropertyType](c, db, name)
}

func GetRoomTypeByName(c *query_cache.CacheStore, db *gorm.DB, name string) (domain.RoomType, error) {
	return GetListItemByName[domain.RoomType](c, db, name)
}

func GetHostByProfileId(c *query_cache.CacheStore, db *gorm.DB, profileId int64, values map[string]interface{}) (domain.Host, error) {
	return GetRecordByFieldValue[domain.Host](c, db, "profile_id", profileId, values)
}

func GetHostForListing(c *query_cache.CacheStore, db *gorm.DB, profileId int64, record KeyValueRecord) (domain.Host, error) {
	hostSince, err := parseHostSince(record)
	if err != nil {
		return domain.Host{}, err
	}
	return GetHostByProfileId(c, db, profileId, map[string]any{
		"name":        record["host_name"],
		"description": record["host_about"],
		"superhost":   ParseAirBNBBool(record["host_is_superhost"]),
		"location":    record["host_location"],
		"slug":        record["host_profile_url"],
		"host_since":  hostSince,
	})
}

func parseHostSince(record KeyValueRecord) (time.Time, error) {
	hostMonths, err := strconv.Atoi(record["hosts_time_as_host_months"])
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse host since months for record: %w", err)
	}
	hostYears, err := strconv.Atoi(record["hosts_time_as_host_years"])
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse host since years for record: %w", err)
	}
	//add negative years and months.
	return time.Now().AddDate(-hostYears, -hostMonths, 0), nil
}

// the dataset stores true / false as t/f
func ParseAirBNBBool(value string) bool {
	return value == "t"
}
