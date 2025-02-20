package conn

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/mstgnz/starter-kit/api/pkg/mstgnz/gobuilder"
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

	var err error
	var database *sql.DB

	for attempts := 1; attempts <= 5; attempts++ {
		database, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Attempt %d: Failed to open DB connection: %v", attempts, err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Veritabanı ayarlarını yapılandır
		database.SetMaxOpenConns(25)
		database.SetMaxIdleConns(5)
		database.SetConnMaxLifetime(5 * time.Minute)
		database.SetConnMaxIdleTime(2 * time.Minute)

		// Bağlantıyı test et
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = database.PingContext(ctx)
		cancel()

		if err == nil {
			log.Println("DB Connected successfully")
			db.DB = database
			return
		}

		log.Printf("Attempt %d: Failed to ping DB: %v", attempts, err)
		database.Close()
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Failed to connect to DB after 5 attempts")
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
func (db *DB) QueryExec(ctx context.Context, builder *gobuilder.GoBuilder) error {
	query, params := builder.Prepare()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, params...)
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
func (db *DB) DynamicCount(ctx context.Context, builder *gobuilder.GoBuilder) (int, error) {
	rowCount := 0

	query, params := builder.Prepare()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return rowCount, err
	}

	rows, err := stmt.QueryContext(ctx, params...)
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
func (db *DB) DynamicFind(ctx context.Context, builder *gobuilder.GoBuilder, model any) error {
	query, params := builder.Prepare()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	rows, err := stmt.QueryContext(ctx, params...)
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
func (db *DB) DynamicGet(ctx context.Context, builder *gobuilder.GoBuilder, model any) ([]any, error) {
	query, params := builder.Prepare()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, params...)
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
func (db *DB) DynamicPaginate(ctx context.Context, builder *gobuilder.GoBuilder, model any) ([]any, error) {
	query, params := builder.Prepare()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, params...)
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
func (db *DB) DynamicCreate(ctx context.Context, builder *gobuilder.GoBuilder) (int, error) {
	var id int
	query, params := builder.Prepare()
	query += " RETURNING id;"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return id, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, params...).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

// DynamicUpdate: the values specified in the table are updated.
func (db *DB) DynamicUpdate(ctx context.Context, builder *gobuilder.GoBuilder) error {
	query, params := builder.Prepare()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, params...)
	if err != nil {
		return err
	}

	return nil
}

// SoftDelete: soft delete the specified id in the specified table.
func (db *DB) SoftDelete(ctx context.Context, builder *gobuilder.GoBuilder) error {
	query, params := builder.Prepare()

	deleteAndUpdate := time.Now().Format("2006-01-02 15:04:05")
	query += fmt.Sprintf("updated_at=$%d, deleted_at=$%d;", len(params)+1, len(params)+2)
	params = append(params, deleteAndUpdate)
	params = append(params, deleteAndUpdate)

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, params...)
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
func (db *DB) HardDelete(ctx context.Context, builder *gobuilder.GoBuilder) error {

	query, params := builder.Prepare()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	result, err := stmt.ExecContext(ctx, params...)
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
func (db *DB) ExistsInTable(ctx context.Context, builder *gobuilder.GoBuilder) error {
	rowCount, err := db.count(ctx, builder)
	if err != nil {
		return err
	}
	if rowCount == 0 {
		return fmt.Errorf("no record with provided conditions exists in table")
	}
	return nil
}

// NotExistsInTable: If not exists, return "nil".
func (db *DB) NotExistsInTable(ctx context.Context, builder *gobuilder.GoBuilder) error {
	rowCount, err := db.count(ctx, builder)
	if err != nil {
		return err
	}
	if rowCount > 0 {
		return fmt.Errorf("record with provided conditions already exists in table")
	}
	return nil
}

func (db *DB) count(ctx context.Context, builder *gobuilder.GoBuilder) (int, error) {
	rowCount := 0

	query, params := builder.Prepare()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return rowCount, err
	}

	rows, err := stmt.QueryContext(ctx, params...)
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
