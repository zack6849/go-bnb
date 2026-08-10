package main

import (
	"GoBNB/internal/convert"
	"GoBNB/internal/domain"
	"GoBNB/internal/errorcollector"
	"GoBNB/internal/query_cache"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type KeyValueRecord = map[string]string

func CreateListingRecord(
	record KeyValueRecord,
	c *query_cache.CacheStore,
	db *gorm.DB,
) (domain.Listing, error) {
	ec := errorcollector.NewCollector()
	latFloat := errorcollector.Field(ec, record, "latitude", convert.Float(64))
	lngFloat := errorcollector.Field(ec, record, "longitude", convert.Float(64))
	airBnbId := errorcollector.Field(ec, record, "id", strconv.Atoi)
	accommodates := errorcollector.Field(ec, record, "accommodates", strconv.Atoi)
	hostProfileId := errorcollector.Field(ec, record, "host_profile_id", convert.Int(10, 64))
	//these are blank for some listings; fall back to the schema's own DEFAULTs
	numBeds := errorcollector.OptionalField(ec, record, "beds", strconv.Atoi)
	minNights := errorcollector.FieldOr(ec, record, "minimum_nights", 1, strconv.Atoi)
	maxNights := errorcollector.FieldOr(ec, record, "maximum_nights", 7, strconv.Atoi)
	//baths are optional and can be a float, eg: 1.5
	numBaths := errorcollector.OptionalField(ec, record, "bathrooms", convert.Float(32))
	if ec.HasErrors() {
		return domain.Listing{}, ec.Summary()
	}
	//then create listing amenities records for it
	//these all resolve to real rows, so a swallowed error here becomes a
	//zero-value struct with a zero PK and fails much later at insert time.
	amenities, err := ParseAmenities(record, c, db)
	if err != nil {
		return domain.Listing{}, err
	}
	propertyType, err := GetPropertyTypeByName(c, db, record["property_type"])
	if err != nil {
		return domain.Listing{},
			fmt.Errorf("failed to resolve property type %q: %w", record["property_type"], err)
	}
	roomType, err := GetRoomTypeByName(c, db, record["room_type"])
	if err != nil {
		return domain.Listing{},
			fmt.Errorf("failed to resolve room type %q: %w", record["room_type"], err)
	}
	host, err := GetHostForListing(c, db, hostProfileId, record)
	if err != nil {
		return domain.Listing{},
			fmt.Errorf("failed to resolve host %d: %w", hostProfileId, err)
	}

	location := domain.GetLocation(latFloat, lngFloat)

	listing := domain.Listing{
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

func ParseAmenities(
	record KeyValueRecord,
	c *query_cache.CacheStore,
	db *gorm.DB,
) ([]domain.Amenity, error) {
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

func GetRecordByFieldValue[T any](
	c *query_cache.CacheStore,
	db *gorm.DB,
	field string,
	value any,
	attrs ...any,
) (T, error) {
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

func GetListItemByName[T any](
	c *query_cache.CacheStore,
	db *gorm.DB,
	name string,
) (T, error) {
	return GetRecordByFieldValue[T](c, db, "name", name)
}

func GetAmenityByName(
	c *query_cache.CacheStore,
	db *gorm.DB,
	name string,
) (domain.Amenity, error) {
	return GetListItemByName[domain.Amenity](c, db, name)
}

func GetPropertyTypeByName(
	c *query_cache.CacheStore,
	db *gorm.DB,
	name string,
) (domain.PropertyType, error) {
	return GetListItemByName[domain.PropertyType](c, db, name)
}

func GetRoomTypeByName(
	c *query_cache.CacheStore,
	db *gorm.DB,
	name string,
) (domain.RoomType, error) {
	return GetListItemByName[domain.RoomType](c, db, name)
}

func GetHostByProfileId(
	c *query_cache.CacheStore,
	db *gorm.DB,
	profileId int64,
	values map[string]any,
) (domain.Host, error) {
	return GetRecordByFieldValue[domain.Host](c, db, "profile_id", profileId, values)
}

func GetHostForListing(
	c *query_cache.CacheStore,
	db *gorm.DB,
	profileId int64,
	record KeyValueRecord,
) (domain.Host, error) {
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

func parseHostSince(
	record KeyValueRecord,
) (time.Time, error) {
	//we only get a relative offset roughly in the dataset
	ec := errorcollector.NewCollector()
	hostMonths := errorcollector.Field(ec, record, "hosts_time_as_host_months", strconv.Atoi)
	hostYears := errorcollector.Field(ec, record, "hosts_time_as_host_years", strconv.Atoi)
	if ec.HasErrors() {
		return time.Time{}, ec.Summary()
	}
	scrapedAt := record["last_scraped"]
	relativeTo, err := time.Parse("2006-01-02", scrapedAt)
	if err != nil {
		return time.Time{},
			fmt.Errorf("failed to parse timestamp '%s' as date to determine host_since: %w", scrapedAt, err)
	}
	//add negative years and months.
	return relativeTo.AddDate(-hostYears, -hostMonths, 0), nil
}

// ParseAirBNBBool the dataset stores true / false as t/f
func ParseAirBNBBool(
	value string,
) bool {
	return value == "t"
}
