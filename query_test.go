package sqlquery

import (
	"testing"

	"github.com/huandu/go-sqlbuilder"
	"github.com/stretchr/testify/assert"
)

func TestParseSelectFilter(t *testing.T) {
	var tests = []struct {
		kind         string
		key          string
		value        interface{}
		expectedSQL  string
		expectedArgs []interface{}
	}{
		{"equals", "id", 1, `SELECT * FROM test_table WHERE id = $1`, []interface{}{1}},
		{"equals nil", "id", nil, `SELECT * FROM test_table WHERE id IS NULL`, []interface{}(nil)},
		{"in", "id.in", "1,2,3", `SELECT * FROM test_table WHERE id IN ($1, $2, $3)`, []interface{}{"1", "2", "3"}},
		{"notin", "id.notin", "1,2,3", `SELECT * FROM test_table WHERE id NOT IN ($1, $2, $3)`, []interface{}{"1", "2", "3"}},
		{"not", "id.not", 1, `SELECT * FROM test_table WHERE id <> $1`, []interface{}{1}},
		{"gt", "id.gt", 1, `SELECT * FROM test_table WHERE id > $1`, []interface{}{1}},
		{"gte", "id.gte", 1, `SELECT * FROM test_table WHERE id >= $1`, []interface{}{1}},
		{"lt", "id.lt", 1, `SELECT * FROM test_table WHERE id < $1`, []interface{}{1}},
		{"lte", "id.lte", 1, `SELECT * FROM test_table WHERE id <= $1`, []interface{}{1}},
		{"like", "id.like", 1, `SELECT * FROM test_table WHERE id LIKE $1`, []interface{}{1}},
		{"null true", "id.null", true, `SELECT * FROM test_table WHERE id.null IS NULL`, []interface{}(nil)},
		{"null false", "id.null", false, `SELECT * FROM test_table WHERE id.null IS NOT NULL`, []interface{}(nil)},
	}
	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			sb := sqlbuilder.NewSelectBuilder()
			sb.SetFlavor(sqlbuilder.Flavor(PostgreSQLFlavor))
			sb.Select("*").From("test_table")
			parseSelectFilter(sb, tt.key, tt.value)
			sqlQuery, args := sb.Build()
			assert.Equal(t, tt.expectedSQL, sqlQuery)
			assert.Equal(t, tt.expectedArgs, args)
		})
	}
}

func TestParseUpdateFilter(t *testing.T) {
	var tests = []struct {
		kind         string
		key          string
		value        interface{}
		expectedSQL  string
		expectedArgs []interface{}
	}{
		{"equals", "id", 1, `UPDATE test_table SET field = $1 WHERE id = $2`, []interface{}{"field", 1}},
		{"equals nil", "id", nil, `UPDATE test_table SET field = $1 WHERE id IS NULL`, []interface{}{"field"}},
		{"in", "id.in", "1,2,3", `UPDATE test_table SET field = $1 WHERE id IN ($2, $3, $4)`, []interface{}{"field", "1", "2", "3"}},
		{"notin", "id.notin", "1,2,3", `UPDATE test_table SET field = $1 WHERE id NOT IN ($2, $3, $4)`, []interface{}{"field", "1", "2", "3"}},
		{"not", "id.not", 1, `UPDATE test_table SET field = $1 WHERE id <> $2`, []interface{}{"field", 1}},
		{"gt", "id.gt", 1, `UPDATE test_table SET field = $1 WHERE id > $2`, []interface{}{"field", 1}},
		{"gte", "id.gte", 1, `UPDATE test_table SET field = $1 WHERE id >= $2`, []interface{}{"field", 1}},
		{"lt", "id.lt", 1, `UPDATE test_table SET field = $1 WHERE id < $2`, []interface{}{"field", 1}},
		{"lte", "id.lte", 1, `UPDATE test_table SET field = $1 WHERE id <= $2`, []interface{}{"field", 1}},
		{"like", "id.like", 1, `UPDATE test_table SET field = $1 WHERE id LIKE $2`, []interface{}{"field", 1}},
		{"null true", "id.null", true, `UPDATE test_table SET field = $1 WHERE id.null IS NULL`, []interface{}{"field"}},
		{"null false", "id.null", false, `UPDATE test_table SET field = $1 WHERE id.null IS NOT NULL`, []interface{}{"field"}},
	}
	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			ub := sqlbuilder.NewUpdateBuilder()
			ub.SetFlavor(sqlbuilder.Flavor(PostgreSQLFlavor))
			ub.Update("test_table")
			ub.Set(ub.Assign("field", "field"))
			parseUpdateFilter(ub, tt.key, tt.value)
			sqlQuery, args := ub.Build()
			assert.Equal(t, tt.expectedSQL, sqlQuery)
			assert.Equal(t, tt.expectedArgs, args)
		})
	}
}

func TestParseDeleteFilter(t *testing.T) {
	var tests = []struct {
		kind         string
		key          string
		value        interface{}
		expectedSQL  string
		expectedArgs []interface{}
	}{
		{"equals", "id", 1, `DELETE FROM test_table WHERE id = $1`, []interface{}{1}},
		{"equals nil", "id", nil, `DELETE FROM test_table WHERE id IS NULL`, []interface{}(nil)},
		{"in", "id.in", "1,2,3", `DELETE FROM test_table WHERE id IN ($1, $2, $3)`, []interface{}{"1", "2", "3"}},
		{"notin", "id.notin", "1,2,3", `DELETE FROM test_table WHERE id NOT IN ($1, $2, $3)`, []interface{}{"1", "2", "3"}},
		{"not", "id.not", 1, `DELETE FROM test_table WHERE id <> $1`, []interface{}{1}},
		{"gt", "id.gt", 1, `DELETE FROM test_table WHERE id > $1`, []interface{}{1}},
		{"gte", "id.gte", 1, `DELETE FROM test_table WHERE id >= $1`, []interface{}{1}},
		{"lt", "id.lt", 1, `DELETE FROM test_table WHERE id < $1`, []interface{}{1}},
		{"lte", "id.lte", 1, `DELETE FROM test_table WHERE id <= $1`, []interface{}{1}},
		{"like", "id.like", 1, `DELETE FROM test_table WHERE id LIKE $1`, []interface{}{1}},
		{"null true", "id.null", true, `DELETE FROM test_table WHERE id.null IS NULL`, []interface{}(nil)},
		{"null false", "id.null", false, `DELETE FROM test_table WHERE id.null IS NOT NULL`, []interface{}(nil)},
	}
	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			db := sqlbuilder.NewDeleteBuilder()
			db.SetFlavor(sqlbuilder.Flavor(PostgreSQLFlavor))
			db.DeleteFrom("test_table")
			parseDeleteFilter(db, tt.key, tt.value)
			sqlQuery, args := db.Build()
			assert.Equal(t, tt.expectedSQL, sqlQuery)
			assert.Equal(t, tt.expectedArgs, args)
		})
	}
}

func TestFindQuery(t *testing.T) {
	expectedSQLQuery := `SELECT * FROM test_table WHERE id = $1 FOR UPDATE SKIP LOCKED`
	expectedArgs := []interface{}{1}
	options := NewFindOptions(PostgreSQLFlavor).WithFilter("id", 1).WithForUpdate("SKIP LOCKED")
	sqlQuery, args := FindQuery("test_table", options)
	assert.Equal(t, expectedSQLQuery, sqlQuery)
	assert.Equal(t, expectedArgs, args)
}

func TestFindAllQuery(t *testing.T) {
	expectedSQLQuery := `SELECT * FROM test_table WHERE id = $1 ORDER BY id asc LIMIT $2 OFFSET $3 FOR UPDATE SKIP LOCKED`
	expectedArgs := []interface{}{1, 50, 10}
	options := NewFindAllOptions(PostgreSQLFlavor).
		WithFilter("id", 1).
		WithLimit(50).
		WithOffset(10).
		WithOrderBy("id asc").
		WithForUpdate("SKIP LOCKED")
	sqlQuery, args := FindAllQuery("test_table", options)
	assert.Equal(t, expectedSQLQuery, sqlQuery)
	assert.Equal(t, expectedArgs, args)
}

type player struct {
	ID   int    `db:"id" fieldtag:"insert"`
	Name string `db:"name" fieldtag:"insert,update"`
}

func TestInsertQuery(t *testing.T) {
	expectedSQLQuery := `INSERT INTO players (id, name) VALUES ($1, $2)`
	expectedArgs := []interface{}{1, "Ronaldinho 10"}
	r10 := player{ID: 1, Name: "Ronaldinho 10"}
	sqlQuery, args := InsertQuery(PostgreSQLFlavor, "insert", "players", &r10)
	assert.Equal(t, expectedSQLQuery, sqlQuery)
	assert.Equal(t, expectedArgs, args)
}

func TestUpdateQuery(t *testing.T) {
	expectedSQLQuery := `UPDATE players SET name = $1 WHERE id = $2`
	expectedArgs := []interface{}{"Ronaldinho Bruxo", 1}
	r10 := player{ID: 1, Name: "Ronaldinho Bruxo"}
	sqlQuery, args := UpdateQuery(PostgreSQLFlavor, "update", "players", r10.ID, &r10)
	assert.Equal(t, expectedSQLQuery, sqlQuery)
	assert.Equal(t, expectedArgs, args)
}

func TestDeleteQuery(t *testing.T) {
	expectedSQLQuery := `DELETE FROM players WHERE id = $1`
	expectedArgs := []interface{}{1}
	sqlQuery, args := DeleteQuery(PostgreSQLFlavor, "players", 1)
	assert.Equal(t, expectedSQLQuery, sqlQuery)
	assert.Equal(t, expectedArgs, args)
}

func TestUpdateWithOptionsQuery(t *testing.T) {
	expectedSQLQuery := `UPDATE players SET age = $1, name = $2 WHERE id = $3`
	expectedArgs := []interface{}{43, "Ronaldinho Bruxo", 1}
	options := NewUpdateOptions(PostgreSQLFlavor).
		WithAssignment("name", "Ronaldinho Bruxo").
		WithAssignment("age", 43).
		WithFilter("id", 1)
	sqlQuery, args := UpdateWithOptionsQuery("players", options)
	assert.Equal(t, expectedSQLQuery, sqlQuery)
	assert.Equal(t, expectedArgs, args)
}

func TestDeleteWithOptionsQuery(t *testing.T) {
	expectedSQLQuery := `DELETE FROM players WHERE id = $1`
	expectedArgs := []interface{}{1}
	options := NewDeleteOptions(PostgreSQLFlavor).WithFilter("id", 1)
	sqlQuery, args := DeleteWithOptionsQuery("players", options)
	assert.Equal(t, expectedSQLQuery, sqlQuery)
	assert.Equal(t, expectedArgs, args)
}
