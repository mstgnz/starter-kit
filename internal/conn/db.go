package conn

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

// ConnectDatabase is creating a new connection to our database
func (db *DB) ConnectDatabase() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbZone := os.Getenv("DB_ZONE")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s", dbHost, dbPort, dbUser, dbPass, dbName, dbZone)
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		panic("Failed DB Connection")
	}
	if err = database.Ping(); err != nil {
		panic("Failed DB Ping")
	}
	log.Println("DB Connected")
	db.DB = database
}

// CloseDatabase method is closing a connection between your app and your db
func (db *DB) CloseDatabase() {
	if err := db.DB.Close(); err != nil {
		log.Println("Failed to close connection from the database:", err.Error())
	} else {
		log.Println("DB Connection Closed")
	}
}

// QueryExec: returns nil if the query is executed successfully
func (db *DB) QueryExec(query string, params []any) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(params...)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New("zero affected")
	}

	return nil
}

// GetObject: returns only the first record it finds
// use config.App().DB.GetObject(query, []any{id}, &model.Address{})
func (db *DB) GetObject(query string, params []any, model any) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}

	rows, err := stmt.Query(params...)
	if err != nil {
		return err
	}

	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// Create new model instance
	modelInstance := reflect.ValueOf(model).Elem()

	// Slice to map field addresses for columns returned in the query
	fieldPointers := make([]any, len(columns))

	// Match column names with model fields
	fieldMap := map[string]reflect.Value{}
	for i := 0; i < modelInstance.NumField(); i++ {
		fieldMap[strings.ToLower(modelInstance.Type().Field(i).Name)] = modelInstance.Field(i)
	}

	// Map columns to model fields
	for i, columnName := range columns {
		fieldName := strings.Join(strings.Split(columnName, "_"), "")
		if field, ok := fieldMap[fieldName]; ok {
			fieldPointers[i] = field.Addr().Interface()
		} else {
			// If the model does not have this column, assign it to a dummy value.
			var dummy sql.NullString
			fieldPointers[i] = &dummy
		}
	}

	if rows.Next() {
		if err := rows.Scan(fieldPointers...); err != nil {
			return err
		}
	} else {
		return sql.ErrNoRows
	}

	return nil
}

// ListObject: returns all records it finds
// use: config.App().DB.ListObject(query, params, &model.Address{})
func (db *DB) ListObject(query string, params []any, model any) ([]any, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var objects []any

	for rows.Next() {
		// Create new model instance
		modelInstance := reflect.ValueOf(model).Elem()

		// Slice to map field addresses for columns returned in the query
		fieldPointers := make([]any, len(columns))

		// Match column names with model fields
		fieldMap := map[string]reflect.Value{}
		for i := 0; i < modelInstance.NumField(); i++ {
			fieldMap[strings.ToLower(modelInstance.Type().Field(i).Name)] = modelInstance.Field(i)
		}

		// Map columns to model fields
		for i, columnName := range columns {
			fieldName := strings.Join(strings.Split(columnName, "_"), "")
			if field, ok := fieldMap[fieldName]; ok {
				fieldPointers[i] = field.Addr().Interface()
			} else {
				// If the model does not have this column, assign it to a dummy value.
				var dummy sql.NullString
				fieldPointers[i] = &dummy
			}
		}

		if err := rows.Scan(fieldPointers...); err != nil {
			return nil, err
		}

		objects = append(objects, modelInstance.Interface())
	}

	return objects, nil
}

// SoftDelete: soft delete the specified id in the specified table.
func (db *DB) SoftDelete(id int, table string) error {
	stmt, err := db.Prepare(fmt.Sprintf("UPDATE %s SET active=$1,deleted_at=$2,updated_at=$3 WHERE id=$4;", table))
	if err != nil {
		return err
	}

	deleteAndUpdate := time.Now().Format("2006-01-02 15:04:05")

	result, err := stmt.Exec(false, deleteAndUpdate, deleteAndUpdate, id)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return errors.New(table + " not deleted")
	}

	return nil
}

// ExistsInTable: If exists, return "nil".
func (db *DB) ExistsInTable(table string, conditions map[string]any) error {
	rowCount, err := db.count(table, conditions)
	if err != nil {
		return err
	}
	if rowCount == 0 {
		return fmt.Errorf("no record with provided conditions exists in table %s", table)
	}
	return nil
}

// NotExistsInTable: If not exists, return "nil".
func (db *DB) NotExistsInTable(table string, conditions map[string]any) error {
	rowCount, err := db.count(table, conditions)
	if err != nil {
		return err
	}
	if rowCount > 0 {
		return fmt.Errorf("record with provided conditions already exists in table %s", table)
	}
	return nil
}

// DynamicCount: returns the number of data according to the conditions
func (db *DB) DynamicCount(table string, conditions map[string]any) (int, error) {
	rowCount := 0

	query := fmt.Sprintf("SELECT count(*) FROM %s", table)
	params := []any{}
	clauses := []string{}

	if len(conditions) > 0 {
		query += " WHERE "
		paramIndex := 1
		for column, value := range conditions {
			clauses = append(clauses, fmt.Sprintf("%s=$%d", column, paramIndex))
			params = append(params, value)
			paramIndex++
		}
		query += strings.Join(clauses, " AND ")
		query += " AND deleted_at IS NULL"
	} else {
		query += " WHERE deleted_at IS NULL"
	}

	query += ";"

	stmt, err := db.Prepare(query)
	if err != nil {
		return rowCount, err
	}

	rows, err := stmt.Query(params...)
	if err != nil {
		return rowCount, err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()

	if rows.Next() {
		if err := rows.Scan(&rowCount); err != nil {
			return rowCount, err
		}
	}

	return rowCount, nil
}

// DynamicCreate: the specified values are recorded in the specified table.
func (db *DB) DynamicCreate(table string, fields map[string]any) error {
	if len(fields) == 0 {
		return fmt.Errorf("no fields provided")
	}

	query := fmt.Sprintf("INSERT INTO %s (", table)

	columns := []string{}
	values := []string{}
	params := []any{}
	paramIndex := 1

	for column, value := range fields {
		columns = append(columns, column)
		values = append(values, fmt.Sprintf("$%d", paramIndex))
		params = append(params, value)
		paramIndex++
	}

	query += strings.Join(columns, ", ") + ") VALUES (" + strings.Join(values, ", ") + ");"

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		return err
	}

	return nil
}

// DynamicUpdate: the values specified in the table are updated.
func (db *DB) DynamicUpdate(table string, updates map[string]any, conditions map[string]any) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}
	if len(conditions) == 0 {
		return fmt.Errorf("no conditions provided")
	}

	query := fmt.Sprintf("UPDATE %s SET ", table)

	setClauses := []string{}
	params := []any{}
	paramIndex := 1

	for column, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s=$%d", column, paramIndex))
		params = append(params, value)
		paramIndex++
	}

	query += strings.Join(setClauses, ", ")

	whereClauses := []string{}
	for column, value := range conditions {
		whereClauses = append(whereClauses, fmt.Sprintf("%s=$%d", column, paramIndex))
		params = append(params, value)
		paramIndex++
	}

	query += " WHERE " + strings.Join(whereClauses, " AND ")
	query += " AND deleted_at IS NULL;"

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(params...)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) RunPrepare(exec any, query string, args ...any) (*sql.Stmt, error) {
	switch e := exec.(type) {
	case *sql.DB:
		return e.Prepare(query)
	case *sql.Tx:
		return e.Prepare(query)
	default:
		return nil, fmt.Errorf("unsupported exec type")
	}
}

func (db *DB) count(table string, conditions map[string]any) (int, error) {
	rowCount := 0
	if len(conditions) == 0 {
		return rowCount, fmt.Errorf("no conditions provided")
	}

	query := fmt.Sprintf("SELECT count(*) FROM %s WHERE ", table)
	params := []any{}
	clauses := []string{}

	// create conditions
	paramIndex := 1
	for column, value := range conditions {
		clauses = append(clauses, fmt.Sprintf("%s=$%d", column, paramIndex))
		params = append(params, value)
		paramIndex++
	}

	// append conditions
	query += strings.Join(clauses, " AND ")

	query += " AND deleted_at IS NULL;"

	stmt, err := db.Prepare(query)
	if err != nil {
		return rowCount, err
	}

	rows, err := stmt.Query(params...)
	if err != nil {
		return rowCount, err
	}
	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()

	for rows.Next() {
		if err := rows.Scan(&rowCount); err != nil {
			return rowCount, err
		}
	}

	return rowCount, nil
}
