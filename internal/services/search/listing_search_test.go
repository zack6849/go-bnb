package search

import (
	"fmt"
	"gobnb/internal/configuration"
	"gobnb/internal/db_assertions"
	"gobnb/internal/domain"
	"gobnb/internal/dto"
	"testing"
)

const NycLat = 40.7128
const NycLng = -74.006

func GetParameters() dto.SearchParameters {
	return dto.SearchParameters{
		MaxDistance: 30,
		SearchNear: domain.Location{
			Latitude:  NycLat,
			Longitude: NycLng,
		},
	}
}

func TestPassesDistanceToDWithin(t *testing.T) {
	params := GetParameters()
	db, _ := configuration.GetDatabaseConfiguration().DryRun()
	query := GetListingQuery(db, params).First(&domain.Listing{})
	//confirm we're using st_dwithin in the where, because otherwise the query is slow
	db_assertions.AssertQueryContainsSubstring(t,
		query,
		"st_dwithin(",
		fmt.Errorf(""),
	)
	db_assertions.AssertQueryContainsSubstring(
		t,
		query, fmt.Sprintf("%d", params.MaxDistance),
		fmt.Errorf(""),
	)
}

func TestSelectsListingsAndDistance(t *testing.T) {
	params := GetParameters()
	db, _ := configuration.GetDatabaseConfiguration().DryRun()
	query := GetListingQuery(db, params).First(&domain.Listing{})
	//confirm we're using st_dwithin in the where, because otherwise the query is slow
	db_assertions.AssertQueryContainsSubstring(t,
		query,
		"listings.*,",
		fmt.Errorf("listing search must select from listings table or the repo will break"),
	)
	db_assertions.AssertQueryContainsSubstring(
		t,
		query,
		"st_distance(",
		fmt.Errorf("must calculate the distance as well"),
	)
	db_assertions.AssertQueryContainsSubstring(
		t,
		query,
		") as distance_meters",
		fmt.Errorf("must calculate the distance as meters in a named alias"),
	)
}
