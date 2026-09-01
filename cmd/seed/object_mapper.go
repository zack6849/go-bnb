package main

import (
	"GoBNB/internal/convert"
	"GoBNB/internal/domain"
	"GoBNB/internal/errs"
	"GoBNB/internal/query_cache"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type KeyValueRecord = map[string]string

func ParseListingRecord(record KeyValueRecord) (domain.Listing, error) {
	chamber := errs.CreateEchoChamber()
	latFloat := chamber.ConvertValue(record, "latitude", convert.Float(64))
	lngFloat := chamber.ConvertValue(record, "longitude", convert.Float(64))
	airBnbId := chamber.ConvertValue(record, "id", strconv.Atoi)
	accommodates := chamber.ConvertValue(record, "accommodates", strconv.Atoi)
	//these are blank for some listings; fall back to the schema's own DEFAULTs
	numBeds := chamber.ConvertFieldOrZeroValue(record, "beds", strconv.Atoi)
	minNights := chamber.ConvertFieldOrDefaultValue(record, "minimum_nights", 1, strconv.Atoi)
	maxNights := chamber.ConvertFieldOrDefaultValue(record, "maximum_nights", 7, strconv.Atoi)
	//baths are optional and can be a float, eg: 1.5
	numBaths := chamber.ConvertFieldOrZeroValue(record, "bathrooms", convert.Float(32))
	host, err := ParseHostInformation(record)
	if err != nil {
		chamber.PushError(fmt.Errorf("failed to parse host information: %w", err))
	}
	hosts := []domain.Host{
		host,
	}
	//then create listing amenities records for it
	//these all resolve to real rows, so a swallowed error here becomes a
	//zero-value struct with a zero PK and fails much later at insert time.
	return domain.Listing{
		PropertyType: domain.PropertyType{
			Name: record["property_type"],
		},
		RoomType: domain.RoomType{
			Name: record["room_type"],
		},
		AirbnbID:       airBnbId,
		Accommodates:   accommodates,
		NumBaths:       float32(numBaths),
		NumBeds:        numBeds,
		ListingUrl:     record["listing_url"],
		Tagline:        record["name"],
		Description:    record["description"],
		MinNights:      minNights,
		MaxNights:      maxNights,
		Location:       domain.GetLocation(latFloat, lngFloat),
		DistanceMeters: 0,
		Hosts:          hosts,
		Amenities:      nil,
	}, chamber.Summary()
}

func ParseHostInformation(record KeyValueRecord) (domain.Host, error) {

	hostProfileId, err := convert.Int(10, 64)(record["host_profile_id"])
	if err != nil {
		return domain.Host{}, fmt.Errorf("failed to parse host information, profile ID cannot be parsed: %w", err)
	}
	hostSince, err := parseHostSince(record)
	if err != nil {
		return domain.Host{}, fmt.Errorf("failed to parse host-since information: %w", err)
	}

	return domain.Host{
		Name:        record["host_name"],
		Description: record["host_about"],
		Superhost:   ParseAirBNBBool(record["host_is_superhost"]),
		Location:    record["host_location"], //this is just a string, eg: NYC, not a literal lat & lng location
		ProfileID:   hostProfileId,
		HostSince:   hostSince,
	}, nil
}

func CreateListingRecord(
	record KeyValueRecord,
	c *query_cache.CacheStore,
	db *gorm.DB,
) (domain.Listing, error) {
	eb := errs.CreateEchoChamber()
	listing, err := ParseListingRecord(record)
	eb.PushError(err)
	amenities, err := ParseAmenities(record, c, db)
	eb.PushError(err)
	listing.Amenities = amenities
	propertyType, err := ResolvePropertyType(c, db, listing.PropertyType.Name)
	eb.PushError(err)
	listing.PropertyType = propertyType
	roomType, err := ResolveRoomType(c, db, listing.RoomType.Name)
	eb.PushError(err)
	listing.RoomType = roomType
	host, err := ResolveHostByProfileId(c, db, listing.Hosts[0].ProfileID, record)
	eb.PushError(err)
	listing.Hosts[0] = host
	return listing, eb.Summary()
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
	return c.GetItemByFieldValue[T](field, value, func() (T, error) {
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
		//if res.RowsAffected == 1 {
		//log.Printf("Created new %T: %v\n", result, value)
		//}
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

func ResolvePropertyType(
	c *query_cache.CacheStore,
	db *gorm.DB,
	name string,
) (domain.PropertyType, error) {
	return GetListItemByName[domain.PropertyType](c, db, name)
}

func ResolveRoomType(
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

func ResolveHostByProfileId(
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
	errorBag := errs.CreateEchoChamber()
	hostMonths := errorBag.ConvertValue(record, "hosts_time_as_host_months", strconv.Atoi)
	hostYears := errorBag.ConvertValue(record, "hosts_time_as_host_years", strconv.Atoi)
	if errorBag.HasErrors() {
		return time.Time{}, errorBag.Summary()
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
