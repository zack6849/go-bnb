package search

import (
	"GoBNB/internal/configuration"
	"GoBNB/internal/db_assertions"
	"GoBNB/internal/dto"
	"testing"
)

func TestUsesDWithinInWhere(t *testing.T) {
	params := dto.SearchParameters{}
	db, _ := configuration.GetDatabaseConfiguration().DryRun()
	query := GetListingQuery(db, params)
	//confirm we're using st_dwithin in the where, because otherwise the query is slow
	db_assertions.AssertQueryContainsWhereWithSubstring(t, query, "st_dwithin")
}
