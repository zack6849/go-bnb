package db_assertions

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func AssertQueryContainsSubstring(t *testing.T, query *gorm.DB, search string, error any) {
	rawSql := GetSQLFromQuery(query)
	normalizedSql := strings.ToLower(rawSql)
	assert.Contains(t,
		normalizedSql,
		search,
		error,
	)
}

// GetSQLFromQuery this function only works if the query has been built already
func GetSQLFromQuery(query *gorm.DB) string {
	//this looks weird, but it's more readable than
	//return query.Dialector.Explain(query.Statement.SQL.String(), query.Statement.Vars...)
	return query.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return query
	})
}
