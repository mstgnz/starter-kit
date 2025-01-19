package conn

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDatabase is creating a new connection to our database
func ConnectDatabase() *gorm.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbZone := os.Getenv("DB_ZONE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s", dbHost, dbPort, dbUser, dbPass, dbName, dbZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("Database connection failed:", err)
	}

	//_ = conn.AutoMigrate(&model.User{}, &model.Blog{}, &model.Comment{})
	log.Println("Database connection successful")
	return db
}

// CloseDatabase method is closing a connection between your app and your db
func CloseDatabase(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("Failed to close connection from database")
	}
	_ = dbSQL.Close()
}
