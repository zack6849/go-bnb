package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func GetExampleRecord() (record KeyValueRecord) {
	record = KeyValueRecord{
		"id":                             "14316232",
		"listing_url":                    "https://www.airbnb.com/rooms/14316232",
		"scrape_id":                      "20260616211424",
		"last_scraped":                   "2026-06-16",
		"source":                         "city scrape",
		"name":                           "Unique Modern Room, Perfect Location Albany",
		"description":                    "Single family home with driveway close to Albany Med, St. Peter’s, St. Rose, U Albany, CDTA public transport, restaurants, downtown, and nightlife. <br /><br />Recently renovated room packed with basic amenities.",
		"neighborhood_overview":          "",
		"picture_url":                    "https://a0.muscache.com/pictures/0d7f6d7c-0faf-43d6-a63f-b935103cf056.jpg",
		"host_id":                        "36438637",
		"host_url":                       "https://www.airbnb.com/users/show/36438637",
		"host_profile_id":                "1462672002391360660",
		"host_profile_url":               "https://www.airbnb.com/users/profile/1462672002391360660",
		"host_name":                      "Chad",
		"host_since":                     "",
		"hosts_time_as_user_years":       "10",
		"hosts_time_as_user_months":      "11",
		"hosts_time_as_host_years":       "10",
		"hosts_time_as_host_months":      "2",
		"host_location":                  "Albany, NY",
		"host_about":                     "",
		"host_response_time":             "",
		"host_response_rate":             "",
		"host_acceptance_rate":           "",
		"host_is_superhost":              "f",
		"host_thumbnail_url":             "",
		"host_picture_url":               "https://a0.muscache.com/im/pictures/user/b3382f39-e130-4af4-b2b6-b755d0a79db8.jpg?aki_policy=profile_x_medium",
		"host_neighbourhood":             "",
		"host_listings_count":            "13",
		"host_total_listings_count":      "",
		"host_verifications":             "None",
		"host_has_profile_pic":           "t",
		"host_identity_verified":         "t",
		"neighbourhood":                  "",
		"neighbourhood_cleansed":         "TENTH WARD",
		"neighbourhood_group_cleansed":   "",
		"latitude":                       "42.65726",
		"longitude":                      "-73.78219",
		"property_type":                  "Private room in home",
		"room_type":                      "Private room",
		"accommodates":                   "1",
		"bathrooms":                      "1.0",
		"bathrooms_text":                 "1 shared bath",
		"bedrooms":                       "",
		"beds":                           "1",
		"amenities":                      "[\"Free washer \\u2013 In unit\", \"Microwave\", \"Stove\", \"Blender\", \"Long term stays allowed\", \"TV\", \"Wifi\", \"Dishwasher\", \"Smoke alarm\", \"Shared backyard \\u2013 Fully fenced\", \"Lock on bedroom door\", \"Coffee maker: drip coffee maker, Keurig coffee machine\", \"Fire extinguisher\", \"Clothing storage: closet\", \"Essentials\", \"Cooking basics\", \"Keypad\", \"Hangers\", \"Hot water\", \"Carbon monoxide alarm\", \"Freezer\", \"Iron\", \"Dining table\", \"Books and reading material\", \"Kitchen\", \"Central heating\", \"Refrigerator\", \"Hair dryer\", \"Baking sheet\", \"Free street parking\", \"Dedicated workspace\", \"Window AC unit\", \"Cleaning products\", \"Toaster\", \"Self check-in\", \"Free driveway parking on premises\", \"First aid kit\", \"Bathtub\", \"Bed linens\", \"Oven\", \"Free dryer \\u2013 In unit\", \"Luggage dropoff allowed\", \"Dishes and silverware\"]",
		"price":                          "$51.89",
		"price_quote_checkin_date":       "2026-08-10",
		"price_quote_checkout_date":      "2026-08-28",
		"price_quote_total_price":        "934.00",
		"price_quote_price_per_night":    "51.89",
		"price_quote_raw":                "{\"quote\": {\"taxes\": null, \"currency\": \"USD\", \"date_match\": null, \"service_fee\": null, \"total_price\": \"934\", \"cleaning_fee\": null, \"is_available\": true, \"discount_amount\": null, \"price_per_night\": \"51.88888888888888888888888889\", \"nightly_subtotal\": \"933.48\", \"discounted_subtotal\": \"933.48\", \"raw_price_line_items\": [{\"amount\": \"933.48\", \"item_type\": \"nightly_subtotal\", \"description\": \"18 nights x $51.86\", \"price_string\": \"$933.48\"}, {\"amount\": \"933.48\", \"item_type\": \"discounted_subtotal\", \"description\": \"Price after discount\", \"price_string\": \"$933.48\"}], \"returned_checkin_date\": null, \"requested_checkin_date\": \"2026-08-10\", \"returned_checkout_date\": null, \"requested_checkout_date\": \"2026-08-28\"}}",
		"minimum_nights":                 "18",
		"maximum_nights":                 "1125",
		"minimum_minimum_nights":         "18",
		"maximum_minimum_nights":         "18",
		"minimum_maximum_nights":         "1125",
		"maximum_maximum_nights":         "1125",
		"minimum_nights_avg_ntm":         "18.0",
		"maximum_nights_avg_ntm":         "1125.0",
		"calendar_updated":               "",
		"has_availability":               "t",
		"availability_30":                "0",
		"availability_60":                "5",
		"availability_90":                "35",
		"availability_365":               "276",
		"calendar_last_scraped":          "2026-06-16",
		"number_of_reviews":              "55",
		"number_of_reviews_ltm":          "1",
		"number_of_reviews_l30d":         "0",
		"availability_eoy":               "110",
		"number_of_reviews_ly":           "3",
		"estimated_occupancy_l365d":      "36",
		"estimated_revenue_l365d":        "1868",
		"first_review":                   "2016-08-07",
		"last_review":                    "2026-02-28",
		"review_scores_rating":           "4.67",
		"review_scores_accuracy":         "4.8",
		"review_scores_cleanliness":      "4.6",
		"review_scores_checkin":          "4.94",
		"review_scores_communication":    "4.76",
		"review_scores_location":         "4.67",
		"review_scores_value":            "4.76",
		"license":                        "",
		"instant_bookable":               "",
		"calculated_host_listings_count": "2",
		"calculated_host_listings_count_entire_homes":  "0",
		"calculated_host_listings_count_private_rooms": "2",
		"calculated_host_listings_count_shared_rooms":  "0",
		"reviews_per_month":                            "0.46",
	}
	return
}

func TestParsesListingCorrectly(t *testing.T) {
	record := GetExampleRecord()
	listing, err := ParseListingRecord(record)
	require.NoError(t, err)
	assert.Equal(t, listing.PropertyType.Name, record["property_type"])
	assert.Equal(t, listing.RoomType.Name, record["room_type"])
	assert.Equal(t, 14316232, listing.AirbnbID)
	assert.Equal(t, 1, listing.Accommodates)
	assert.Equal(t, float32(1.0), listing.NumBaths)
	assert.Equal(t, 1, listing.NumBeds)
	assert.Equal(t, 18, listing.MinNights)
	assert.Equal(t, 1125, listing.MaxNights)
	assert.Equal(t, record["listing_url"], listing.ListingUrl)
	assert.Equal(t, record["description"], listing.Description)
	assert.Equal(t, record["name"], listing.Tagline)
	assert.Equal(t, 42.65726, listing.Location.Latitude)
	assert.Equal(t, -73.78219, listing.Location.Longitude)
	assert.Len(t, listing.Hosts, 1)
	assert.Equal(t, record["host_name"], listing.Hosts[0].Name)
}

func TestParsesListingRecordFallsBackToDefaultsWhenFieldsAreMissing(t *testing.T) {
	record := GetExampleRecord()
	delete(record, "beds")
	delete(record, "minimum_nights")
	delete(record, "maximum_nights")
	delete(record, "bathrooms")

	listing, err := ParseListingRecord(record)
	if err != nil {
		t.Fail()
		return
	}
	assert.Equal(t, 0, listing.NumBeds)
	assert.Equal(t, 1, listing.MinNights)
	assert.Equal(t, 7, listing.MaxNights)
	assert.Equal(t, float32(0), listing.NumBaths)
}

func TestParseHostInformation(t *testing.T) {
	record := GetExampleRecord()
	host, err := ParseHostInformation(record)
	if err != nil {
		t.Fail()
		return
	}
	assert.Equal(t, record["host_name"], host.Name)
	assert.Equal(t, record["host_about"], host.Description)
	assert.Equal(t, record["host_location"], host.Location)
	assert.Equal(t, false, host.Superhost)
	assert.Equal(t, int64(1462672002391360660), host.ProfileID)
	assert.Equal(t, time.Date(2016, time.April, 16, 0, 0, 0, 0, time.UTC), host.HostSince)
}

func TestParseHostInformationErrorsOnUnparseableProfileId(t *testing.T) {
	record := GetExampleRecord()
	record["host_profile_id"] = "not-a-number"

	_, err := ParseHostInformation(record)
	assert.Error(t, err)
}

func TestParseHostInformationErrorsOnUnparseableHostSince(t *testing.T) {
	record := GetExampleRecord()
	record["last_scraped"] = "not-a-date"

	_, err := ParseHostInformation(record)
	assert.Error(t, err)
}

func TestParseAirBNBBool(t *testing.T) {
	tests := map[string]bool{
		"t":     true,
		"f":     false,
		"":      false,
		"true":  false,
		"false": false,
	}
	for input, expected := range tests {
		assert.Equal(t, expected, ParseAirBNBBool(input), "input %q", input)
	}
}
