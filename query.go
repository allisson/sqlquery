package sqlquery

import (
	"sort"
	"strings"

	"github.com/huandu/go-sqlbuilder"
)

func parseIn(value string) []interface{} {
	values := strings.Split(value, ",")
	result := make([]interface{}, len(values))
	for i := range values {
		result[i] = values[i]
	}
	return result
}

func parseSelectFilter(sb *sqlbuilder.SelectBuilder, key string, value interface{}) {
	if strings.Contains(key, ".") {
		split := strings.Split(key, ".")
		parsedKey := split[0]
		compare := split[1]
		switch compare {
		case "in":
			valueStr, ok := value.(string)
			if ok {
				values := parseIn(valueStr)
				sb.Where(sb.In(parsedKey, values...))
			}
		case "notin":
			valueStr, ok := value.(string)
			if ok {
				values := parseIn(valueStr)
				sb.Where(sb.NotIn(parsedKey, values...))
			}
		case "not":
			sb.Where(sb.NotEqual(parsedKey, value))
		case "gt":
			sb.Where(sb.GreaterThan(parsedKey, value))
		case "gte":
			sb.Where(sb.GreaterEqualThan(parsedKey, value))
		case "lt":
			sb.Where(sb.LessThan(parsedKey, value))
		case "lte":
			sb.Where(sb.LessEqualThan(parsedKey, value))
		case "like":
			sb.Where(sb.Like(parsedKey, value))
		case "null":
			valueBool, ok := value.(bool)
			if ok {
				if valueBool {
					sb.Where(sb.IsNull(key))
				} else {
					sb.Where(sb.IsNotNull(key))
				}
			}
		}
	} else {
		switch value.(type) {
		case nil:
			sb.Where(sb.IsNull(key))
		default:
			sb.Where(sb.Equal(key, value))
		}
	}
}

func parseUpdateFilter(ub *sqlbuilder.UpdateBuilder, key string, value interface{}) {
	if strings.Contains(key, ".") {
		split := strings.Split(key, ".")
		parsedKey := split[0]
		compare := split[1]
		switch compare {
		case "in":
			valueStr, ok := value.(string)
			if ok {
				values := parseIn(valueStr)
				ub.Where(ub.In(parsedKey, values...))
			}
		case "notin":
			valueStr, ok := value.(string)
			if ok {
				values := parseIn(valueStr)
				ub.Where(ub.NotIn(parsedKey, values...))
			}
		case "not":
			ub.Where(ub.NotEqual(parsedKey, value))
		case "gt":
			ub.Where(ub.GreaterThan(parsedKey, value))
		case "gte":
			ub.Where(ub.GreaterEqualThan(parsedKey, value))
		case "lt":
			ub.Where(ub.LessThan(parsedKey, value))
		case "lte":
			ub.Where(ub.LessEqualThan(parsedKey, value))
		case "like":
			ub.Where(ub.Like(parsedKey, value))
		case "null":
			valueBool, ok := value.(bool)
			if ok {
				if valueBool {
					ub.Where(ub.IsNull(key))
				} else {
					ub.Where(ub.IsNotNull(key))
				}
			}
		}
	} else {
		switch value.(type) {
		case nil:
			ub.Where(ub.IsNull(key))
		default:
			ub.Where(ub.Equal(key, value))
		}
	}
}

func parseDeleteFilter(db *sqlbuilder.DeleteBuilder, key string, value interface{}) {
	if strings.Contains(key, ".") {
		split := strings.Split(key, ".")
		parsedKey := split[0]
		compare := split[1]
		switch compare {
		case "in":
			valueStr, ok := value.(string)
			if ok {
				values := parseIn(valueStr)
				db.Where(db.In(parsedKey, values...))
			}
		case "notin":
			valueStr, ok := value.(string)
			if ok {
				values := parseIn(valueStr)
				db.Where(db.NotIn(parsedKey, values...))
			}
		case "not":
			db.Where(db.NotEqual(parsedKey, value))
		case "gt":
			db.Where(db.GreaterThan(parsedKey, value))
		case "gte":
			db.Where(db.GreaterEqualThan(parsedKey, value))
		case "lt":
			db.Where(db.LessThan(parsedKey, value))
		case "lte":
			db.Where(db.LessEqualThan(parsedKey, value))
		case "like":
			db.Where(db.Like(parsedKey, value))
		case "null":
			valueBool, ok := value.(bool)
			if ok {
				if valueBool {
					db.Where(db.IsNull(key))
				} else {
					db.Where(db.IsNotNull(key))
				}
			}
		}
	} else {
		switch value.(type) {
		case nil:
			db.Where(db.IsNull(key))
		default:
			db.Where(db.Equal(key, value))
		}
	}
}

// FindQuery returns compiled SELECT string and args.
func FindQuery(tableName string, options *FindOptions) (string, []interface{}) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.Flavor(options.Flavor))
	sb.Select(options.Fields...).From(tableName)
	for key, value := range options.Filters {
		parseSelectFilter(sb, key, value)
	}
	if options.ForUpdate {
		sb.ForUpdate()
		if options.ForUpdateMode != "" {
			sb.SQL(options.ForUpdateMode)
		}
	}
	return sb.Build()
}

// FindAllQuery returns compiled SELECT string and args.
func FindAllQuery(tableName string, options *FindAllOptions) (string, []interface{}) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.Flavor(options.Flavor))
	sb.Select(options.Fields...).From(tableName).Limit(options.Limit).Offset(options.Offset)
	for key, value := range options.Filters {
		parseSelectFilter(sb, key, value)
	}
	if options.OrderBy != "" {
		sb.OrderBy(options.OrderBy)
	}
	if options.ForUpdate {
		sb.ForUpdate()
		if options.ForUpdateMode != "" {
			sb.SQL(options.ForUpdateMode)
		}
	}
	return sb.Build()
}

// InsertQuery returns compiled INSERT string and args.
func InsertQuery(flavor Flavor, tag, tableName string, structValue interface{}) (string, []interface{}) {
	theStruct := sqlbuilder.NewStruct(structValue).For(sqlbuilder.Flavor(flavor))
	ib := theStruct.WithTag(tag).InsertInto(tableName, structValue)
	return ib.Build()
}

// UpdateQuery returns compiled UPDATE string and args.
func UpdateQuery(flavor Flavor, tag, tableName string, id interface{}, structValue interface{}) (string, []interface{}) {
	theStruct := sqlbuilder.NewStruct(structValue).For(sqlbuilder.Flavor(flavor))
	ub := theStruct.WithTag(tag).Update(tableName, structValue)
	ub.Where(ub.Equal("id", id))
	return ub.Build()
}

// DeleteQuery returns compiled DELETE string and args.
func DeleteQuery(flavor Flavor, tableName string, id interface{}) (string, []interface{}) {
	db := sqlbuilder.NewDeleteBuilder()
	db.SetFlavor(sqlbuilder.Flavor(flavor))
	db.DeleteFrom(tableName)
	db.Where(db.Equal("id", id))
	return db.Build()
}

// UpdateWithOptionsQuery returns compiled UPDATE string and args from UpdateOptions.
func UpdateWithOptionsQuery(tableName string, options *UpdateOptions) (string, []interface{}) {
	ub := sqlbuilder.NewUpdateBuilder()
	ub.SetFlavor(sqlbuilder.Flavor(options.Flavor))
	ub.Update(tableName)
	var assignments []string
	for key, value := range options.Assignments {
		assignments = append(assignments, ub.Assign(key, value))
	}
	sort.Strings(assignments)
	ub = ub.Set(assignments...)
	for key, value := range options.Filters {
		parseUpdateFilter(ub, key, value)
	}
	return ub.Build()
}

// DeleteWithOptionsQuery returns compiled DELETE string and args from DeleteOptions.
func DeleteWithOptionsQuery(tableName string, options *DeleteOptions) (string, []interface{}) {
	db := sqlbuilder.NewDeleteBuilder()
	db.SetFlavor(sqlbuilder.Flavor(options.Flavor))
	db.DeleteFrom(tableName)
	for key, value := range options.Filters {
		parseDeleteFilter(db, key, value)
	}
	return db.Build()
}
