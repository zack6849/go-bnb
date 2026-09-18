package main

import (
	"encoding/json"
	"fmt"
	"gobnb/internal/convert"
	"gobnb/internal/domain"
	"gobnb/internal/errs"
	"strconv"
	"time"
)

type KeyValueRecord = map[string]string

func ExtractListingFromRecord(record KeyValueRecord) (domain.Listing, error) {
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
	host, err := ExtractHostFromRecord(record)
	if err != nil {
		chamber.PushError(fmt.Errorf("failed to parse host information: %w", err))
	}
	amenities, err := ExtractAmenetiesFromRecord(record)
	if err != nil {
		chamber.PushError(fmt.Errorf("failed to parse amenities %w", err))
	}
	hosts := []domain.Host{
		host,
	}
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
		ListingURL:     record["listing_url"],
		PictureURL:     record["picture_url"],
		Tagline:        record["name"],
		Description:    record["description"],
		MinNights:      minNights,
		MaxNights:      maxNights,
		Location:       domain.GetLocation(latFloat, lngFloat),
		DistanceMeters: 0,
		Hosts:          hosts,
		Amenities:      amenities,
	}, chamber.Summary()
}

func ExtractHostFromRecord(record KeyValueRecord) (domain.Host, error) {

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
		SuperHost:   ParseAirBNBBool(record["host_is_superhost"]),
		Location:    record["host_location"], //this is just a string, eg: NYC, not a literal lat & lng location
		ProfileID:   hostProfileId,
		HostSince:   hostSince,
		Slug:        record["host_profile_url"],
		PictureURL:  record["host_picture_url"],
	}, nil
}

func ExtractAmenityNamesFromRecord(
	record KeyValueRecord,
) ([]string, error) {
	var amenities []string
	err := json.Unmarshal([]byte(record["amenities"]), &amenities)
	return amenities, err
}

func ExtractAmenetiesFromRecord(record KeyValueRecord) ([]domain.Amenity, error) {
	list := make([]domain.Amenity, 0, 1)
	amenityNames, err := ExtractAmenityNamesFromRecord(record)
	if err != nil {
		return list, err
	}
	for _, amenityName := range amenityNames {
		list = append(list, domain.Amenity{
			Name: amenityName,
		})
	}
	return list, nil
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
