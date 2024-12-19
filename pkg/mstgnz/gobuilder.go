package mstgnz

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const (
	InnerJoin = "INNER"
	LeftJoin  = "LEFT"
	RightJoin = "RIGHT"
	FullJoin  = "FULL"
)

// https://github.com/mstgnz/gobuilder
type GoBuilder struct {
	TableClause   string
	SelectClause  string
	WhereClause   string
	JoinClause    string
	GroupByClause string
	HavingClause  string
	OrderByClause string
	LimitClause   string
	UnionClause   string
}

// NewGoBuilder initializes a new instance of GoBuilder
func NewGoBuilder() *GoBuilder {
	return &GoBuilder{}
}

// Table specifies the table name for the query
func (gb *GoBuilder) Table(table string) *GoBuilder {
	gb.TableClause = table
	return gb
}

// Select defines the columns to be selected in the query
func (gb *GoBuilder) Select(columns ...string) *GoBuilder {
	if len(columns) == 0 {
		columns = append(columns, "*")
	}
	gb.SelectClause = fmt.Sprintf("SELECT %s FROM %s", strings.Join(columns, ", "), gb.TableClause)
	return gb
}

// SelectDistinct defines the distinct columns to be selected
func (gb *GoBuilder) SelectDistinct(columns ...string) *GoBuilder {
	if len(columns) == 0 {
		columns = append(columns, "*")
	}
	gb.SelectClause = fmt.Sprintf("SELECT DISTINCT %s FROM %s", strings.Join(columns, ", "), gb.TableClause)
	return gb
}

// Insert adds an INSERT INTO statement to the query with bind parameters
func (gb *GoBuilder) Insert(args map[string]any) *GoBuilder {
	if len(args) != 0 {
		keys := make([]string, 0, len(args))
		values := make([]string, 0, len(args))

		for key, value := range args {
			cleanedValue := cleanValue(value)
			keys = append(keys, key)
			values = append(values, cleanedValue)
		}

		gb.SelectClause = fmt.Sprintf("INSERT INTO %s (%v) VALUES (%v)", gb.TableClause, strings.Join(keys, ", "), strings.Join(values, ", "))
	}
	return gb
}

// Update builds an UPDATE statement with bind parameters
func (gb *GoBuilder) Update(args map[string]any) *GoBuilder {
	if len(args) != 0 {
		setClauses := make([]string, 0, len(args))

		for key, value := range args {
			setClauses = append(setClauses, fmt.Sprintf("%s = %v", key, cleanValue(value)))
		}

		gb.SelectClause = fmt.Sprintf("UPDATE %s SET %v", gb.TableClause, strings.Join(setClauses, ", "))
	}
	return gb
}

// Delete builds a DELETE statement
func (gb *GoBuilder) Delete() *GoBuilder {
	gb.SelectClause = fmt.Sprintf("DELETE FROM %v", gb.TableClause)
	return gb
}

// Where adds a WHERE clause with bind parameters
func (gb *GoBuilder) Where(key, opt string, val any) *GoBuilder {
	return gb.addWhereClause("AND", key, opt, val)
}

// OrWhere adds an OR WHERE clause with bind parameters
func (gb *GoBuilder) OrWhere(key, opt string, val any) *GoBuilder {
	return gb.addWhereClause("OR", key, opt, val)
}

// In adds an IN clause with bind parameters
func (gb *GoBuilder) In(column string, args ...any) *GoBuilder {
	return gb.addInClause("AND", column, args...)
}

// OrIn adds an OR IN clause with bind parameters
func (gb *GoBuilder) OrIn(column string, args ...any) *GoBuilder {
	return gb.addInClause("OR", column, args...)
}

// Between adds a BETWEEN clause with bind parameters
func (gb *GoBuilder) Between(column string, args ...any) *GoBuilder {
	return gb.between("AND", column, args...)
}

// OrBetween adds an OR BETWEEN clause with bind parameters
func (gb *GoBuilder) OrBetween(column string, args ...any) *GoBuilder {
	return gb.between("OR", column, args...)
}

// IsNull adds an IS NULL clause
func (gb *GoBuilder) IsNull(column string) *GoBuilder {
	clause := fmt.Sprintf("%s IS NULL", column)
	gb.addClause("AND", clause)
	return gb
}

// OrIsNull adds an OR IS NULL clause
func (gb *GoBuilder) OrIsNull(column string) *GoBuilder {
	clause := fmt.Sprintf("%s IS NULL", column)
	gb.addClause("OR", clause)
	return gb
}

// IsNotNull adds an IS NOT NULL clause
func (gb *GoBuilder) IsNotNull(column string) *GoBuilder {
	clause := fmt.Sprintf("%s IS NOT NULL", column)
	gb.addClause("AND", clause)
	return gb
}

// OrIsNotNull adds an OR IS NOT NULL clause
func (gb *GoBuilder) OrIsNotNull(column string) *GoBuilder {
	clause := fmt.Sprintf("%s IS NOT NULL", column)
	gb.addClause("OR", clause)
	return gb
}

// Having adds a HAVING clause
func (gb *GoBuilder) Having(condition string) *GoBuilder {
	gb.HavingClause = fmt.Sprintf("HAVING %s", condition)
	return gb
}

// OrHaving adds an OR HAVING clause
func (gb *GoBuilder) OrHaving(condition string) *GoBuilder {
	if gb.HavingClause != "" {
		gb.HavingClause = fmt.Sprintf("%s OR %s", gb.HavingClause, condition)
	} else {
		gb.HavingClause = fmt.Sprintf("HAVING %s", condition)
	}
	return gb
}

// Join adds a JOIN clause
func (gb *GoBuilder) Join(joinType, table, condition string) *GoBuilder {
	gb.JoinClause = fmt.Sprintf("%s JOIN %s ON %s", joinType, table, condition)
	return gb
}

// Limit adds a LIMIT clause
func (gb *GoBuilder) Limit(start, limit int) *GoBuilder {
	gb.LimitClause = fmt.Sprintf("LIMIT %d, %d", start, limit)
	return gb
}

// GroupBy adds a GROUP BY clause
func (gb *GoBuilder) GroupBy(columns ...string) *GoBuilder {
	gb.GroupByClause = fmt.Sprintf("GROUP BY %v", strings.Join(columns, ", "))
	return gb
}

// OrderBy adds an ORDER BY ASC clause
func (gb *GoBuilder) OrderBy(columns ...string) *GoBuilder {
	gb.OrderByClause = fmt.Sprintf("ORDER BY %v ASC", strings.Join(columns, ", "))
	return gb
}

// OrderByDesc adds an ORDER BY DESC clause
func (gb *GoBuilder) OrderByDesc(columns ...string) *GoBuilder {
	gb.OrderByClause = fmt.Sprintf("ORDER BY %v DESC", strings.Join(columns, ", "))
	return gb
}

// Union adds a UNION clause
func (gb *GoBuilder) Union(sql string) *GoBuilder {
	if gb.UnionClause == "" {
		gb.UnionClause = fmt.Sprintf("UNION %v", sql)
	} else {
		gb.UnionClause = fmt.Sprintf("%v UNION %v", gb.UnionClause, sql)
	}
	return gb
}

// ToSql returns the final SQL query and the associated bind parameters
func (gb *GoBuilder) ToSql() string {
	clauses := []string{
		gb.SelectClause,
		gb.JoinClause,
		gb.WhereClause,
		gb.GroupByClause,
		gb.HavingClause,
		gb.OrderByClause,
		gb.LimitClause,
		gb.UnionClause,
	}
	query := strings.Join(clauses, " ")
	re := regexp.MustCompile(`\s+`)
	gb.reset()
	return strings.TrimSpace(re.ReplaceAllString(query, " "))
}

// Private method to RESET builder
func (gb *GoBuilder) reset() {
	v := reflect.ValueOf(gb).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String {
			field.SetString("")
		}
	}
}

// Private method to add WHERE or OR WHERE clauses with bind parameters
func (gb *GoBuilder) addWhereClause(OP, key, opt string, val any) *GoBuilder {
	clause := ""
	switch v := val.(type) {
	case *GoBuilder:
		clause = fmt.Sprintf("%v %v (%v)", key, opt, v.ToSql())
	default:
		clause = fmt.Sprintf("%v %v %v", key, opt, cleanValue(val))
	}
	gb.addClause(OP, clause)
	return gb
}

// Private method to add clauses with logical operators
func (gb *GoBuilder) addClause(OP, clause string) {
	if gb.WhereClause != "" {
		gb.WhereClause = fmt.Sprintf("%v %v %v", gb.WhereClause, OP, clause)
	} else {
		gb.WhereClause = fmt.Sprintf("WHERE %v", clause)
	}
}

// Private method to add IN clauses with values directly
func (gb *GoBuilder) addInClause(OP, column string, args ...any) *GoBuilder {
	if len(args) > 0 {
		values := make([]string, len(args))
		for i, arg := range args {
			switch v := arg.(type) {
			case string:
				values[i] = fmt.Sprintf("%v", cleanValue(v))
			default:
				values[i] = fmt.Sprintf("%d", v)
			}
		}
		clause := fmt.Sprintf("%v IN (%v)", column, strings.Join(values, ", "))
		gb.addClause(OP, clause)
	}
	return gb
}

// Private method to add BETWEEN clauses with values directly
func (gb *GoBuilder) between(OP, column string, args ...any) *GoBuilder {
	if len(args) == 2 {
		var clause string
		switch args[0].(type) {
		case string:
			clause = fmt.Sprintf("%v BETWEEN %s AND %s", column, cleanValue(fmt.Sprintf("%v", args[0])), cleanValue(fmt.Sprintf("%v", args[1])))
		default:
			clause = fmt.Sprintf("%v BETWEEN %v AND %v", column, args[0], args[1])
		}
		gb.addClause(OP, clause)
	}
	return gb
}

// cleanValue trims and escapes potentially harmful characters from the value
func cleanValue(value any) string {
	switch v := value.(type) {
	case string:
		cleaned := strings.ReplaceAll(v, "'", "''")
		return "'" + cleaned + "'"
	default:
		return fmt.Sprintf("%v", v)
	}
}
