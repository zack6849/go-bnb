package db_assertions

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func AssertQueryContainsWhereWithSubstring(t *testing.T, query *gorm.DB, search string) {
	var wheres clause.Expression
	wheres = query.Statement.Clauses["WHERE"].Expression
	where, ok := wheres.(clause.Where)
	if !ok {
		assert.Fail(t, "expected WHERE expression to be a where clause")
	}
	found := false
	for _, e := range where.Exprs {
		expr, ok := e.(clause.Expr)
		if !ok {
			assert.Fail(t, "encountered non-clause in where clause")
		}
		found = strings.Contains(strings.ToLower(expr.SQL), search)
	}
	if !found {
		assert.Fail(t, fmt.Sprintf("failed to find substring '%s' in where", search))
	}
}
