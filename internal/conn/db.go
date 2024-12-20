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
	"github.com/mstgnz/starter-kit/pkg/mstgnz"
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
func (db *DB) QueryExec(builder *mstgnz.GoBuilder) error {
	query, params := builder.Prepare()

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

// DynamicCount: returns the number of data according to the conditions
func (db *DB) DynamicCount(builder *mstgnz.GoBuilder) (int, error) {
	rowCount := 0

	query, params := builder.Prepare()

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

// DynamicFind: only renders the first matching record to the p.Model object based on the conditions
func (db *DB) DynamicFind(builder *mstgnz.GoBuilder, model any) error {
	query, params := builder.Prepare()

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

// DynamicGet: returns all records it finds
func (db *DB) DynamicGet(builder *mstgnz.GoBuilder, model any) ([]any, error) {
	query, params := builder.Prepare()

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

// DynamicPaginate: returns all records according to the conditions
func (db *DB) DynamicPaginate(builder *mstgnz.GoBuilder, model any) ([]any, error) {
	query, params := builder.Prepare()

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = stmt.Close()
		_ = rows.Close()
	}()

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

// DynamicCreate: the specified values are recorded in the specified table.
func (db *DB) DynamicCreate(builder *mstgnz.GoBuilder) (int, error) {
	var id int
	query, params := builder.Prepare()
	query += " RETURNING id;"
	stmt, err := db.Prepare(query)
	if err != nil {
		return id, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(params...).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

// DynamicUpdate: the values specified in the table are updated.
func (db *DB) DynamicUpdate(builder *mstgnz.GoBuilder) error {
	query, params := builder.Prepare()

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

// SoftDelete: soft delete the specified id in the specified table.
func (db *DB) SoftDelete(builder *mstgnz.GoBuilder) error {
	query, params := builder.Prepare()

	deleteAndUpdate := time.Now().Format("2006-01-02 15:04:05")
	query += fmt.Sprintf("updated_at=$%d, deleted_at=$%d;", len(params)+1, len(params)+2)
	params = append(params, deleteAndUpdate)
	params = append(params, deleteAndUpdate)

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
		return errors.New("not soft deleted")
	}

	return nil
}

// HardDelete: hard delete the specified id in the specified table.
func (db *DB) HardDelete(builder *mstgnz.GoBuilder) error {

	query, params := builder.Prepare()

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
		return errors.New("not hard deleted")
	}

	return nil
}

// ExistsInTable: If exists, return "nil".
func (db *DB) ExistsInTable(builder *mstgnz.GoBuilder) error {
	rowCount, err := db.count(builder)
	if err != nil {
		return err
	}
	if rowCount == 0 {
		return fmt.Errorf("no record with provided conditions exists in table")
	}
	return nil
}

// NotExistsInTable: If not exists, return "nil".
func (db *DB) NotExistsInTable(builder *mstgnz.GoBuilder) error {
	rowCount, err := db.count(builder)
	if err != nil {
		return err
	}
	if rowCount > 0 {
		return fmt.Errorf("record with provided conditions already exists in table")
	}
	return nil
}

func (db *DB) count(builder *mstgnz.GoBuilder) (int, error) {
	rowCount := 0

	query, params := builder.Prepare()

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
